package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"plane.watch/lib/producer"
	"plane.watch/lib/sink"
	"plane.watch/lib/tracker"
)

var (
	pwPort                          int
	showDebug                       bool
	dump1090Host                    string
	dump1090Port                    string
)

func main() {
	app := cli.NewApp()

	app.Version = "1.0.0"
	app.Name = "Plane Watch Client"
	app.Usage = "Reads from dump1090 and sends it to http://plane.watch/"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "rabbit-host",
			Value:       "localhost",
			Usage:       "the rabbitmq host to talk to",
		},
		&cli.IntFlag{
			Name:        "rabbit-port",
			Value:       5672,
			Usage:       "The rabbitmq port to talk to",
			Destination: &pwPort,
		},
		&cli.StringFlag{
			Name:        "rabbit-user",
			Value:       "",
			Usage:       "user for rabbitmq",
		},
		&cli.StringFlag{
			Name:        "rabbit-pass",
			Value:       "",
			Usage:       "rabbitmq password",
		},
		&cli.StringFlag{
			Name:        "rabbit-vhost",
			Value:       "plane.watch",
			Usage:       "the virtual host on the rabbit server to use",
		},
		&cli.StringFlag{
			Name:        "dump1090_host",
			Value:       "",
			Usage:       "The host to read dump1090 from",
			Destination: &dump1090Host,
			EnvVars:     []string{"DUMP1090_HOST"},
		},
		&cli.StringFlag{
			Name:        "dump1090_port",
			Value:       "30002",
			Usage:       "The port to read dump1090 from",
			Destination: &dump1090Port,
			EnvVars:     []string{"DUMP1090_PORT"},
		},
		&cli.StringFlag{
			Name:  "avr-file",
			Value: "",
			Usage: "A file to read AVR frames from",
		},
		&cli.StringFlag{
			Name:  "beast-file",
			Value: "",
			Usage: "A file to read beast format AVR frames from",
		},
		&cli.Float64Flag{
			Name:  "ref-lat",
			Usage: "The reference latitude for decoding messages. Needs to be within 45nm of where the messages are generated.",
		},
		&cli.Float64Flag{
			Name:  "ref-lon",
			Usage: "The reference longitude for decoding messages. Needs to be within 45nm of where the messages are generated.",
		},
		&cli.BoolFlag{
			Name:        "debug",
			Usage:       "Show Extra Debug Information",
			Destination: &showDebug,
			EnvVars:     []string{"DEBUG"},
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:   "run",
			Usage:  "Gather ADSB data and sends it to plane.watch",
			Action: run,
		},
		{
			Name:      "simple",
			Usage:     "Gather ADSB data and sends it to plane.watch",
			Action:    runSimple,
			ArgsUsage: "[app.log - A file name to output to or stdout if not specified]",
		},
	}

	if err := app.Run(os.Args); nil != err {
		fmt.Println(err)
	}
}

func commonSetup(c *cli.Context) (*tracker.Tracker, error) {
	opts := make([]tracker.Option, 0)
	if c.Bool("debug") {
		opts = append(opts, tracker.WithVerboseOutput())
	} else {
		opts = append(opts, tracker.WithInfoOutput())
	}
	refLat := c.Float64("refLat")
	refLon := c.Float64("refLon")
	if refLat != 0 && refLon != 0 {
		opts = append(opts, tracker.WithReferenceLatLon(refLat, refLon))
	}

	trk := tracker.NewTracker(opts...)

	if "" != c.String("redis-host") {
		trk.AddSink(
			sink.NewRedisSink(
				sink.WithHost(c.String("redis-host"), c.String("redis-port")),
			),
		)
	}
	if "" != c.String("rabbit-host") {
		rabbitSink, err := sink.NewRabbitMqSink(
			sink.WithHost(c.String("rabbit-host"), c.String("rabbit-port")),
			sink.WithUserPass(c.String("rabbit-user"), c.String("rabbit-pass")),
			sink.WithRabbitVhost(c.String("rabbit-vhost")),
			sink.WithAllRabbitQueues(),
		)
		if nil != err {
			return nil, err
		}

		trk.AddSink(rabbitSink)
	}

	if "" != dump1090Host {
		switch dump1090Port {
		case "30002":
			trk.AddProducer(producer.NewAvrFetcher(dump1090Host, dump1090Port))
		case "30003":
			trk.AddProducer(producer.NewSbs1Fetcher(dump1090Host, dump1090Port))
		case "30005":
			trk.AddProducer(producer.NewBeastFetcher(dump1090Host, dump1090Port))
		default:
			return nil, errors.New("don't know how to handle port:" + dump1090Port)
		}
	}
	if file := c.String("avr-file"); "" != file {
		trk.AddProducer(producer.NewAvrFile([]string{file}))
	}
	if file := c.String("beast-file"); "" != file {
		trk.AddProducer(producer.NewBeastFile([]string{file}, true))
	}
	return trk, nil
}

func runSimple(c *cli.Context) error {
	trk, err := commonSetup(c)
	if nil != err {
		return err
	}
	trk.AddSink(sink.NewLoggerSink(sink.WithLogOutput(os.Stdout)))

	trk.Wait()
	return nil
}

// run is our method for running things
func run(c *cli.Context) error {
	app, err := newAppDisplay()
	if nil != err {
		return err
	}

	trk, err := commonSetup(c)
	if nil != err {
		return err
	}
	trk.AddSink(sink.NewLoggerSink(sink.WithLogFile("app.log")))
	trk.AddSink(app)

	err = app.Run()
	trk.Stop()
	return err
}
