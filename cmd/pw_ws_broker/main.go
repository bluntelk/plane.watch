package main

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
	"plane.watch/lib/logging"
	"plane.watch/lib/stats"
)

var (
	prometheusNumClients = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pw_ws_broker_num_clients",
		Help: "The current number of websocket clients we are currently serving",
	})
)

func main() {
	app := cli.NewApp()
	app.Name = "Plane.Watch WebSocket Broker (pw_ws_broker)"
	app.Usage = "Websocket Broker"
	app.Description = "Acts as a go between external display elements our the data pipeline"
	app.Authors = []*cli.Author{
		{
			Name:  "Jason Playne",
			Email: "jason@jasonplayne.com",
		},
	}

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
			Name:    "source",
			Usage:   "A place to fetch data from. amqp://user:pass@host:port/vhost?ttl=60",
			Value:   "amqp://guest:guest@rabbitmq:5672/pw",
			EnvVars: []string{"SOURCE"},
		},
		&cli.StringFlag{
			Name:    "route-key-low",
			Usage:   "The routing key that has only the significant flight update events",
			Value:   "location-updates-enriched-reduced",
			EnvVars: []string{"ROUTE_KEY_LOW"},
		},
		&cli.StringFlag{
			Name:    "route-key-high",
			Usage:   "The routing key that has all of the flight update events",
			Value:   "location-updates-enriched",
			EnvVars: []string{"ROUTE_KEY_HIGH"},
		},
		&cli.StringFlag{
			Name:    "http-addr",
			Usage:   "What the HTTP server listens on",
			Value:   ":80",
			EnvVars: []string{"HTTP_ADDR"},
		},
		&cli.BoolFlag{
			Name:    "serve-test-web",
			Usage:   "Serve up a test website for websocket testing",
			EnvVars: []string{"TEST_WEB"},
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

	stats.IncludePrometheusFlags(app, 9603)

	app.Before = func(c *cli.Context) error {
		logging.SetVerboseOrQuiet(c.Bool("debug"), c.Bool("quiet"))
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

func run(c *cli.Context) error {
	stats.RunPrometheusWebServer(c)
	source := c.String("source")
	lowRoute := c.String("route-key-low")
	highRoute := c.String("route-key-high")

	isValid := true
	if "" == source {
		log.Info().Msg("Please provide rabbitmq connection details. (--source)")
		isValid = false
	}
	if "" == lowRoute {
		log.Info().Msg("Please provide the routing key for significant updates. (--route-key-low)")
		isValid = false
	}
	if "" == highRoute {
		log.Info().Msg("Please provide the routing key for all updates. (--route-key-high)")
		isValid = false
	}
	if !isValid {
		return errors.New("invalid configuration")
	}

	broker, err := NewPlaneWatchWebSocketBroker(
		source,
		lowRoute,
		highRoute,
		c.String("http-addr"),
		c.Bool("serve-test-web"),
	)
	if nil != err {
		return err
	}
	defer broker.Close()

	if err = broker.Setup(); nil != err {
		return err
	}
	go broker.Run()

	broker.Wait()

	return nil
}
