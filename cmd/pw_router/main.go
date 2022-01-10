package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"plane.watch/lib/tile_grid"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"github.com/urfave/cli/v2"

	"plane.watch/lib/export"
	"plane.watch/lib/logging"
	"plane.watch/lib/rabbitmq"
)

// queue suffixes for a low (only significant) and high (every message) tile queues
const (
	qSuffixLow  = "_low"
	qSuffixHigh = "_high"
)

type (
	rabbit struct {
		rmq  *rabbitmq.RabbitMQ
		conf *rabbitmq.Config

		queues map[string]*amqp.Queue

		syncSamples sync.Map
	}

	planeLocationLast struct {
		lastSignificantUpdate export.PlaneLocation
		candidateUpdate       export.PlaneLocation
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

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "rabbitmq",
			Usage:   "Rabbitmq URL for reaching and publishing updates.",
			EnvVars: []string{"RABBITMQ"},
		},
		&cli.StringFlag{
			Name:    "source-route-key",
			Usage:   "Name of the routing key to read location updates from.",
			Value:   "location-updates",
			EnvVars: []string{"SOURCE_ROUTE_KEY"},
		},
		&cli.StringFlag{
			Name:    "destination-route-key",
			Usage:   "Name of the routing key to publish significant updates to.",
			Value:   "location-updates-reduced",
			EnvVars: []string{"DEST_ROUTE_KEY"},
		},
		&cli.IntFlag{
			Name:    "num-workers",
			Usage:   "Number of workers to process updates.",
			Value:   4,
			EnvVars: []string{"NUM_WORKERS"},
		},
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   "Show Extra Debug Information",
			EnvVars: []string{"DEBUG"},
		},
		&cli.BoolFlag{
			Name:    "quiet",
			Usage:   "Only show important messages",
			EnvVars: []string{"QUIET"},
		},
		&cli.IntFlag{
			Name:    "prom-metrics-port",
			Usage:   "Port to listen on for prometheus app metrics.",
			Value:   9601,
			EnvVars: []string{"PROM_METRICS_PORT"},
		},
		&cli.BoolFlag{
			Name:  "register-test-queues",
			Usage: "Subscribes a bunch of queues to our routing keys",
		},
	}

	app.Action = run

	app.Before = func(c *cli.Context) error {
		logging.ConfigureForCli()
		logging.SetVerboseOrQuiet(c.Bool("debug"), c.Bool("quiet"))
		return nil
	}

	if err := app.Run(os.Args); nil != err {
		log.Error().Err(err).Send()
	}
}

func (r *rabbit) connect(config rabbitmq.Config, timeout time.Duration) error {
	log.Info().Str("host", config.String()).Msg("Connecting to RabbitMQ")
	r.rmq = rabbitmq.New(&config)
	return r.rmq.ConnectAndWait(timeout)
}

func (r *rabbit) makeQueue(name, bindRouteKey string) error {
	q, err := r.rmq.QueueDeclare(name, 60000) // 60sec TTL
	if nil != err {
		log.Error().Err(err).Msgf("Failed to create queue '%s'", name)
		return err
	}
	r.queues[name] = &q

	if err = r.rmq.QueueBind(name, bindRouteKey, rabbitmq.PlaneWatchExchange); nil != err {
		log.Error().Err(err).Msgf("Failed to QueueBind to route-key:%s to queue %s", bindRouteKey, name)
		return err
	}
	log.Debug().Str("queue", name).Str("route-key", bindRouteKey).Msg("Setup Queue")
	return nil
}

func (r *rabbit) setupTestQueues() error {
	log.Info().Msg("Setting up test queues")
	// we need a _low and a _high for each tile
	suffixes := []string{qSuffixLow, qSuffixHigh}
	for _, name := range tile_grid.GridLocationNames() {
		for _, suffix := range suffixes {
			if err := r.makeQueue(name+suffix, name+suffix); nil != err {
				return err
			}
		}
	}
	return nil
}

func run(c *cli.Context) error {
	// setup and start the prom exporter
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(fmt.Sprintf(":%d", c.Int("prom-metrics-port")), nil)
	}()

	var err error
	// connect to rabbitmq, create ourselves 2 queues
	r := rabbit{
		queues:      map[string]*amqp.Queue{},
		syncSamples: sync.Map{},
	}

	if "" == c.String("rabbitmq") {
		return errors.New("please specify the --rabbitmq parameter")
	}

	rabbitUrl, err := url.Parse(c.String("rabbitmq"))
	if err != nil {
		return err
	}

	rabbitPassword, _ := rabbitUrl.User.Password()

	rabbitConfig := rabbitmq.Config{
		Host:     rabbitUrl.Hostname(),
		Port:     rabbitUrl.Port(),
		User:     rabbitUrl.User.Username(),
		Password: rabbitPassword,
		Vhost:    rabbitUrl.Path,
		Ssl:      rabbitmq.ConfigSSL{},
	}

	// connect to Rabbit
	if err = r.connect(rabbitConfig, time.Second*5); nil != err {
		return err
	}

	if err = r.makeQueue("reducer-in", c.String("source-route-key")); nil != err {
		return err
	}

	if c.Bool("register-test-queues") {
		if err = r.setupTestQueues(); nil != err {
			return err
		}
	}

	ch, err := r.rmq.Consume("reducer-in", "pw-router")
	if nil != err {
		log.Info().Msg("Failed to consume reducer-in")
		return err
	}

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	log.Info().Msgf("Starting with %d workers...", c.Int("num-workers"))
	for i := 0; i < c.Int("num-workers"); i++ {
		worker := worker{
			rabbit:         &r,
			destRoutingKey: c.String("destination-route-key"),
		}
		wg.Add(1)
		go func() {
			worker.run(ctx, ch)
			wg.Done()
		}()
	}

	wg.Wait()

	return nil
}
