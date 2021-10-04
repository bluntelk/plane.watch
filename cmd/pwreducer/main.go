package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"github.com/urfave/cli/v2"

	"plane.watch/lib/logging"
	"plane.watch/lib/rabbitmq"
)

type (
	rabbit struct {
		rmq  *rabbitmq.RabbitMQ
		conf *rabbitmq.Config

		queues map[string]*amqp.Queue

		sync_samples sync.Map
	}

	planeLocationLast struct {
		lastSignificantUpdate planeLocation
		candidateUpdate       planeLocation
	}

	// straight up copy from lib/sink/rabbitmq.go
	planeLocation struct {
		original          []byte
		New, Removed      bool
		Icao              string
		Lat, Lon, Heading float64
		Velocity          float64
		Altitude          int
		VerticalRate      int
		AltitudeUnits     string
		FlightNumber      string
		FlightStatus      string
		OnGround          bool
		Airframe          string
		AirframeType      string
		HasLocation       bool
		HasHeading        bool
		HasVerticalRate   bool
		HasVelocity       bool
		SourceTag         string
		Squawk            string
		Special           string
		TrackedSince      time.Time
		LastMsg           time.Time
	}
)

var (
	updatesProccessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pwreducer_updates_processed_total",
		Help: "The total number of messages processed.",
	})
	updatesSignificant = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pwreducer_updates_significant_total",
		Help: "The total number of messages determined to be significant.",
	})
	updatesIgnored = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pwreducer_updates_ignored_total",
		Help: "The total number of messages determined to be insignificant and thus ignored.",
	})
	updatesPublished = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pwreducer_updates_published_total",
		Help: "The total number of messages published to the output queue.",
	})
	updatesError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pwreducer_updates_error_total",
		Help: "The total number of messages that could not be processed due to an error.",
	})
)

func main() {
	app := cli.NewApp()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":9601", nil)

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "rabbitmq",
			Usage:   "The place to send decoded JSON in URL Form. amqp://user:pass@host:port/vhost?ttl=60",
			EnvVars: []string{"RABBITMQ"},
		},
		&cli.StringFlag{
			Name:    "source-queue-name",
			Usage:   "Name of the queue to read location updates from. Default: location-updates",
			EnvVars: []string{"SOURCE_QUEUE_NAME"},
		},
		&cli.IntFlag{
			Name:    "num-workers",
			Usage:   "Number of workers to process updates. Default: 4",
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
	r.rmq = rabbitmq.New(config)
	connected := make(chan bool)
	go r.rmq.Connect(connected)
	select {
	case <-connected:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("failed to connect to rabbit in a timely manner")
	}
}

func (r *rabbit) makeQueue(name string) error {
	q, err := r.rmq.QueueDeclare(name, 60000) // 60sec TTL
	if nil != err {
		return err
	}
	r.queues[name] = &q
	return nil
}

func run(c *cli.Context) error {
	var err error
	// connect to rabbitmq, create ourselves 2 queues
	r := rabbit{
		queues:       map[string]*amqp.Queue{},
		sync_samples: sync.Map{},
	}

	url, err := url.Parse(c.String("rabbitmq"))

	if err != nil {
		return err
	}

	rabbitPassword, _ := url.User.Password()

	rabbitConfig := rabbitmq.Config{
		Host:     url.Hostname(),
		Port:     url.Port(),
		User:     url.User.Username(),
		Password: rabbitPassword,
		Ssl:      rabbitmq.ConfigSSL{},
	}

	// connect to Rabbit
	if err = r.connect(rabbitConfig, time.Second*5); nil != err {
		return err
	}

	if err = r.makeQueue("reducer-in"); nil != err {
		log.Info().Msg("Failed to makeQueue reducer-in")
		return err
	}

	var queueName string

	if c.String("source-queue-name") == "" {
		queueName = c.String("source-queue-name")
	} else {
		queueName = "location-updates" //default name
	}

	if err = r.rmq.QueueBind("reducer-in", queueName, "plane.watch.data"); nil != err {
		log.Info().Msg("Failed to QueueBind to input queue")
		return err
	}

	if err = r.makeQueue("reducer-out"); nil != err {
		log.Info().Msg("Failed to makeQueue reducer-out")
		return err
	}

	if err = r.rmq.QueueBind("reducer-out", "location-updates-reduced", "plane.watch.data"); nil != err {
		log.Info().Msg("Failed to QueueBind to output queue")
		return err
	}

	// open a channel to the reducer-in queue - I think "pw-reducer" is the name?? of the client.
	ch, err := r.rmq.Consume("reducer-in", "pw-reducer")
	if nil != err {
		log.Info().Msg("Failed to consume reducer-in")
		return err
	}

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	var workerCount int

	if c.Int("num-workers") > 0 {
		workerCount = c.Int("num-workers")
	} else {
		workerCount = 4
	}

	log.Info().Msgf("Starting with %d workers...", workerCount)
	for i := 0; i < workerCount; i++ {
		worker := worker{
			rabbit: &r,
		}
		wg.Add(1)

		go worker.run(ctx, ch, &wg)
	}

	wg.Wait()

	return nil
}
