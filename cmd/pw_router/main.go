package main

import (
	"context"
	"os"
	"os/signal"
	"plane.watch/lib/nats_io"
	"plane.watch/lib/redismq"
	"sync"
	"syscall"
	"time"

	"plane.watch/lib/dedupe"
	"plane.watch/lib/monitoring"
	"plane.watch/lib/tile_grid"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"plane.watch/lib/logging"
	"plane.watch/lib/rabbitmq"
)

// queue suffixes for a low (only significant) and high (every message) tile queues
const (
	qSuffixLow  = "_low"
	qSuffixHigh = "_high"
)

type (
	pwRouter struct {
		rmq   *rabbitmq.RabbitMQ
		nats  *nats_io.Server
		redis *redismq.Server

		syncSamples *dedupe.ForgetfulSyncMap

		haveSourceSinkConnection bool

		incomingMessages chan []byte
	}
)

var (
	updatesProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pw_router_updates_processed_total",
		Help: "The total number of messages processed.",
	})
	updatesSignificant = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pw_router_updates_significant_total",
		Help: "The total number of messages determined to be significant.",
	})
	updatesInsignificant = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pw_router_updates_insignificant_total",
		Help: "The total number of messages determined to be insignificant.",
	})
	updatesIgnored = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pw_router_updates_ignored_total",
		Help: "The total number of messages determined to be insignificant and thus ignored.",
	})
	updatesPublished = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pw_router_updates_published_total",
		Help: "The total number of messages published to the output queue.",
	})
	updatesError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pw_router_updates_error_total",
		Help: "The total number of messages that could not be processed due to an error.",
	})
	cacheEntries = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pw_router_cache_planes_count",
		Help: "The number of planes in the reducer cache.",
	})
	cacheEvictions = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pw_router_cache_eviction_total",
		Help: "The number of cache evictions made from the cache.",
	})
)

func main() {
	app := cli.NewApp()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	app.Version = "1.0.0"
	app.Name = "Plane Watch Router (pw_router)"
	app.Usage = "Reads location updates from AMQP and publishes only significant updates."

	app.Description = `This program takes a stream of plane tracking data (location updates) from an AMQP message bus  ` +
		`and filters messages and only returns significant changes for each aircraft.` +
		"\n\n" +
		`example: ./pw_router --rabbitmq="amqp://guest:guest@localhost:5672" --source-route-key=location-updates --num-workers=8 --prom-metrics-port=9601`

	app.Commands = cli.Commands{
		{
			Name:        "daemon",
			Description: "For prod, Logging is JSON formatted",
			Action:      runDaemon,
		},
		{
			Name:        "cli",
			Description: "Runs in your terminal with human readable output",
			Action:      runCli,
		},
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "rabbitmq",
			Usage: "Rabbitmq URL for fetching and publishing updates.",
			//Value:   "amqp://guest:guest@rabbitmq:5672/pw",
			EnvVars: []string{"RABBITMQ"},
		},
		&cli.StringFlag{
			Name:  "nats",
			Usage: "Nats.io URL for fetching and publishing updates.",
			//Value:   "nats://guest:guest@nats:4222/",
			EnvVars: []string{"NATS"},
		},
		&cli.StringFlag{
			Name:  "redis",
			Usage: "redis server URL for fetching and publishing updates.",
			//Value:   "redis://guest:guest@redis:6379/",
			EnvVars: []string{"NATS"},
		},
		&cli.StringFlag{
			Name:    "source-route-key",
			Usage:   "Name of the routing key to read location updates from.",
			Value:   "location-updates-enriched",
			EnvVars: []string{"SOURCE_ROUTE_KEY"},
		},
		&cli.StringFlag{
			Name:    "destination-route-key",
			Usage:   "Name of the routing key to publish significant updates to.",
			Value:   "location-updates-enriched-reduced",
			EnvVars: []string{"DEST_ROUTE_KEY"},
		},
		&cli.IntFlag{
			Name:    "num-workers",
			Usage:   "Number of workers to process updates.",
			Value:   10,
			EnvVars: []string{"NUM_WORKERS"},
		},
		&cli.BoolFlag{
			Name:    "spread-updates",
			Usage:   "publish location updates to their respective tileXX_high and tileXX_low routing keys as well.",
			EnvVars: []string{"SPREAD"},
		},
		&cli.IntFlag{
			Name:    "update-age",
			Usage:   "seconds to keep an update before aging it out of the cache.",
			Value:   30,
			EnvVars: []string{"UPDATE_AGE"},
		},
		&cli.IntFlag{
			Name:    "update-age-sweep-interval",
			Usage:   "Seconds between cache age sweeps.",
			Value:   5,
			EnvVars: []string{"UPDATE_SWEEP"},
		},
		&cli.BoolFlag{
			Name:  "register-test-queues",
			Usage: "Subscribes a bunch of queues to our routing keys.",
		},
	}
	logging.IncludeVerbosityFlags(app)
	monitoring.IncludeMonitoringFlags(app, 9601)

	app.Before = func(c *cli.Context) error {
		logging.SetLoggingLevel(c)
		return nil
	}

	if err := app.Run(os.Args); nil != err {
		log.Error().Err(err).Send()
	}
}

func runDaemon(c *cli.Context) error {
	return run(c)
}

func runCli(c *cli.Context) error {
	logging.ConfigureForCli()
	return run(c)
}

func (r *pwRouter) rabbitMqMakeQueue(name, bindRouteKey string) error {
	if nil != r.rmq {
		_, err := r.rmq.QueueDeclare(name, 60000) // 60sec TTL
		if nil != err {
			log.Error().Err(err).Msgf("Failed to create queue '%s'", name)
			return err
		}

		if err = r.rmq.QueueBind(name, bindRouteKey, rabbitmq.PlaneWatchExchange); nil != err {
			log.Error().Err(err).Msgf("Failed to QueueBind to route-key:%s to queue %s", bindRouteKey, name)
			return err
		}
		log.Debug().Str("queue", name).Str("route-key", bindRouteKey).Msg("Setup Queue")
	}

	if nil != r.nats {

	}
	return nil
}
func (r *pwRouter) rabbitMqSetupTestQueues() error {
	log.Info().Msg("Setting up test queues")
	// we need a _low and a _high for each tile
	suffixes := []string{qSuffixLow, qSuffixHigh}
	for _, name := range tile_grid.GridLocationNames() {
		for _, suffix := range suffixes {
			if err := r.rabbitMqMakeQueue(name+suffix, name+suffix); nil != err {
				return err
			}
		}
	}
	return nil
}
func (r *pwRouter) connect(config rabbitmq.Config, timeout time.Duration) error {
	log.Info().Str("host", config.String()).Msg("Connecting to RabbitMQ")
	r.rmq = rabbitmq.New(&config)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return r.rmq.ConnectAndWait(ctx)
}
func (r *pwRouter) connectRabbit(c *cli.Context, done context.Context) error {
	url := c.String("rabbitmq")
	if "" == url {
		return nil
	}

	conf, err := rabbitmq.NewConfigFromUrl(url)
	if nil != err {
		log.Error().
			Err(err).
			Str("url", url).
			Str("MQ", "rabbitmq").
			Msg("Unable to determine configuration from URL")
	}

	// connect to Rabbit
	if err = r.connect(*conf, time.Second*5); nil != err {
		return err
	}

	if c.Bool("register-test-queues") {
		if err = r.rabbitMqSetupTestQueues(); nil != err {
			return err
		}
	}

	if err = r.rabbitMqMakeQueue("reducer-in", c.String("source-route-key")); nil != err {
		return err
	}

	ch, err := r.rmq.Consume("reducer-in", "pw-router")
	if nil != err {
		log.Info().Msg("Failed to consume reducer-in")
		return err
	}
	go func() {
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					log.Error().Msg("failed to get message from rabbit nicely")
					return
				}

				r.incomingMessages <- msg.Body
			case <-done.Done():
				return
			}
		}
	}()

	monitoring.AddHealthCheck(r.rmq)
	r.haveSourceSinkConnection = true
	return nil
}

func (r *pwRouter) connectNatsIo(c *cli.Context, done context.Context) error {
	var err error
	url := c.String("nats")
	if "" == url {
		return nil
	}

	r.nats, err = nats_io.NewServer(url)
	if nil == err {
		monitoring.AddHealthCheck(r.nats)
		r.haveSourceSinkConnection = true
	} else {
		log.Error().
			Err(err).
			Str("url", url).
			Str("MQ", "nats.io").
			Msg("Unable to determine configuration from URL")
		return err
	}

	ch, err := r.nats.Subscribe(c.String("source-route-key"))
	go func() {
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					log.Error().Msg("failed to get message from nats.io nicely")
					return
				}

				r.incomingMessages <- msg.Data
			case <-done.Done():
				close(ch)
				return
			}
		}
	}()
	return nil
}

func (r *pwRouter) connectRedis(c *cli.Context, done context.Context) error {
	url := c.String("redis")
	if "" == url {
		return nil
	}
	var err error
	r.redis, err = redismq.NewServer(url)

	if nil == err {
		monitoring.AddHealthCheck(r.redis)
		r.haveSourceSinkConnection = true
	} else {
		log.Error().
			Err(err).
			Str("url", url).
			Str("MQ", "redis").
			Msg("Unable to determine configuration from URL")
		return err
	}

	ch, err := r.redis.Subscribe(c.String("source-route-key"))
	go func() {
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					log.Error().Msg("failed to get message from redis nicely")
					return
				}

				r.incomingMessages <- []byte(msg.Payload)
			case <-done.Done():
				return
			}
		}
	}()
	return nil
}

func run(c *cli.Context) error {
	// setup and start the prom exporter
	monitoring.RunWebServer(c)

	var err error
	// connect to rabbitmq, create ourselves 2 queues
	r := pwRouter{
		syncSamples: dedupe.NewForgetfulSyncMap(time.Duration(c.Int("update-age-sweep-interval"))*time.Second, time.Duration(c.Int("update-age"))*time.Second),
	}

	r.syncSamples.SetEvictionAction(func(key interface{}, value interface{}) {
		cacheEvictions.Inc()
		cacheEntries.Dec()
		log.Debug().Msgf("Evicting cache entry Icao: %s", key)
	})

	r.incomingMessages = make(chan []byte, 300)

	doneRabbit, rabbitCancel := context.WithCancel(context.Background())
	defer rabbitCancel()
	doneNats, natsCancel := context.WithCancel(context.Background())
	defer natsCancel()
	doneRedis, redisCancel := context.WithCancel(context.Background())
	defer redisCancel()

	if err = r.connectRabbit(c, doneRabbit); nil != err {
		log.Error().Err(err).Msg("Failed to connect to rabbitmq")
	}
	if err = r.connectNatsIo(c, doneNats); nil != err {
		log.Error().Err(err).Msg("Failed to connect to nats.io")
	}
	if err = r.connectRedis(c, doneRedis); nil != err {
		log.Error().Err(err).Msg("Failed to connect to nats.io")
	}
	if !r.haveSourceSinkConnection {
		cli.ShowAppHelpAndExit(c, 1)
	}

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-chSignal // wait for our cancel signal
		log.Info().Msg("Shutting Down")
		// and then close all the things
		rabbitCancel()
		natsCancel()
		cancel()
	}()

	log.Info().Msgf("Starting with %d workers...", c.Int("num-workers"))
	for i := 0; i < c.Int("num-workers"); i++ {
		wkr := worker{
			router:         &r,
			destRoutingKey: c.String("destination-route-key"),
			spreadUpdates:  c.Bool("spread-updates"),
		}
		wg.Add(1)
		go func() {
			wkr.run(ctx, r.incomingMessages)
			wg.Done()
		}()
	}

	wg.Wait()

	return nil
}
