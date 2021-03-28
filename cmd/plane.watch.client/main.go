package main

import (
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
		},
	}

	if err := app.Run(os.Args); nil != err {
		fmt.Println(err)
	}
}

func runSimple(c *cli.Context) error {
	opts := make([]tracker.Option, 0)
	if c.GlobalBool("debug") {
		opts = append(opts, tracker.WithVerboseOutput())
	} else {
		opts = append(opts, tracker.WithInfoOutput())
	}
	trk := tracker.NewTracker(opts...)
	//trk.AddSink(sink.NewLoggerSink(sink.WithLogOutput(os.Stdout)))
	trk.AddSink(sink.NewLoggerSink(sink.WithLogFile("app.log")))
	if "" != dump1090Host {
		trk.AddProducer(producer.NewAvrFetcher(dump1090Host, dump1090Port))
	}
	if file := c.GlobalString("avr-file"); "" != file {
		trk.AddProducer(producer.NewAvrFile([]string{file}))
	}

	trk.Wait()
	return nil
}

// run is our method for running things
func run(c *cli.Context) error {
	app, err := newAppDisplay()
	if nil != err {
		return err
	}

	opts := make([]tracker.Option, 0)
	if c.GlobalBool("debug") {
		opts = append(opts, tracker.WithVerboseOutput())
	} else {
		opts = append(opts, tracker.WithInfoOutput())
	}
	trk := tracker.NewTracker(opts...)
	trk.AddSink(app)
	f, err := os.Create("app.log")
	if nil != err {
		return err
	}
	trk.AddSink(sink.NewLoggerSink(sink.WithLogOutput(f)))

	if "" != c.String("redis-host") {
		trk.AddSink(
			sink.NewRedisSink(
				sink.WithHost(c.String("redis-host"), c.String("redis-port")),
			),
		)
	}
	if "" != c.String("rabbit-host") {
		trk.AddSink(
			sink.NewRabbitMqSink(
				sink.WithHost(c.String("rabbit-host"), c.String("rabbit-port")),
				sink.WithQueue(c.String("rabbit-queue")),
			),
		)
	}

	if "" != dump1090Host {
		trk.AddProducer(producer.NewAvrFetcher(dump1090Host, dump1090Port))
	}
	if file := c.GlobalString("avr-file"); "" != file {
		trk.AddProducer(producer.NewAvrFile([]string{file}))
	}
	err = app.Run()
	trk.Stop()
	return err
}
