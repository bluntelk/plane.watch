package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"plane.watch/lib/producer"
	"plane.watch/lib/sink"
	"plane.watch/lib/tracker"
)

var (
	pwHost, pwUser, pwPass, pwVhost string
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
		cli.StringFlag{
			Name:        "pw_host",
			Value:       "mq.plane.watch",
			Usage:       "How we connect to plane.watch",
			Destination: &pwHost,
			EnvVar:      "PW_HOST",
		},
		cli.StringFlag{
			Name:        "pw_user",
			Value:       "",
			Usage:       "user for plane.watch",
			Destination: &pwUser,
			EnvVar:      "PW_USER",
		},
		cli.StringFlag{
			Name:        "pw_pass",
			Value:       "",
			Usage:       "plane.watch password",
			Destination: &pwPass,
			EnvVar:      "PW_PASS",
		},
		cli.IntFlag{
			Name:        "pw_port",
			Value:       5672,
			Usage:       "How we connect to plane.watch",
			Destination: &pwPort,
			EnvVar:      "PW_PORT",
		},
		cli.StringFlag{
			Name:        "pw_vhost",
			Value:       "/pw_feedin",
			Usage:       "the virtual host on the plane watch rabbit server",
			Destination: &pwVhost,
			EnvVar:      "PW_VHOST",
		},
		cli.StringFlag{
			Name:        "dump1090_host",
			Value:       "",
			Usage:       "The host to read dump1090 from",
			Destination: &dump1090Host,
			EnvVar:      "DUMP1090_HOST",
		},
		cli.StringFlag{
			Name:        "dump1090_port",
			Value:       "30002",
			Usage:       "The port to read dump1090 from",
			Destination: &dump1090Port,
			EnvVar:      "DUMP1090_PORT",
		},
		cli.StringFlag{
			Name:        "avr-file",
			Value:       "",
			Usage:       "A file to read AVR frames from",
		},
		cli.StringFlag{
			Name:        "beast-file",
			Value:       "",
			Usage:       "A file to read beast format AVR frames from",
		},
		cli.Float64Flag{
			Name:        "ref-lat",
			Usage:       "The reference latitude for decoding messages. Needs to be within 45nm of where the messages are generated.",
		},
		cli.Float64Flag{
			Name:        "ref-lon",
			Usage:       "The reference longitude for decoding messages. Needs to be within 45nm of where the messages are generated.",
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Show Extra Debug Information",
			Destination: &showDebug,
			EnvVar:      "DEBUG",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "Gather ADSB data and sends it to plane.watch",
			Action: run,
		},
		{
			Name:   "simple",
			Usage:  "Gather ADSB data and sends it to plane.watch",
			Action: runSimple,
			ArgsUsage: "[app.log - A file name to output to or stdout if not specified]",
		},
	}

	if err := app.Run(os.Args); nil != err {
		fmt.Println(err)
	}
}

func commonSetup(c *cli.Context) (*tracker.Tracker, error) {
	opts := make([]tracker.Option, 0)
	if c.GlobalBool("debug") {
		opts = append(opts, tracker.WithVerboseOutput())
	} else {
		opts = append(opts, tracker.WithInfoOutput())
	}
	refLat := c.GlobalFloat64("refLat")
	refLon := c.GlobalFloat64("refLon")
	if refLat != 0 && refLon != 0 {
		opts = append(opts, tracker.WithReferenceLatLon(refLat, refLon))
	}

	trk := tracker.NewTracker(opts...)

	if "" != c.GlobalString("redis-host") {
		trk.AddSink(
			sink.NewRedisSink(
				sink.WithHost(c.String("redis-host"), c.String("redis-port")),
			),
		)
	}
	if "" != c.GlobalString("rabbit-host") {
		trk.AddSink(
			sink.NewRabbitMqSink(
				sink.WithHost(c.String("rabbit-host"), c.String("rabbit-port")),
				sink.WithQueue(c.String("rabbit-queue")),
			),
		)
	}


	if "" != dump1090Host {
		switch dump1090Port {
		case "30002":
			trk.AddProducer(producer.NewAvrFetcher(dump1090Host, dump1090Port))
		case "30005":
			trk.AddProducer(producer.NewBeastFetcher(dump1090Host, dump1090Port))
		default:
			return nil, errors.New("don't know how to handle port:" +dump1090Port)
		}
	}
	if file := c.GlobalString("avr-file"); "" != file {
		trk.AddProducer(producer.NewAvrFile([]string{file}))
	}
	if file := c.GlobalString("beast-file"); "" != file {
		trk.AddProducer(producer.NewBeastFile([]string{file}))
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
