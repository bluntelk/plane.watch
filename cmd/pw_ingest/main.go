package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
	"plane.watch/lib/dedupe"
	"plane.watch/lib/example_finder"
	"plane.watch/lib/logging"
	"plane.watch/lib/monitoring"
	"plane.watch/lib/setup"
	"plane.watch/lib/sink"
	"plane.watch/lib/tracker"
)

var (
	prometheusGaugeCurrentPlanes = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pw_ingest_current_tracked_planes_count",
		Help: "The number of planes this instance is currently tracking",
	})
)

func main() {
	app := cli.NewApp()

	app.Version = "1.0.0"
	app.Name = "Plane Watch Client"
	app.Usage = "Reads from dump1090 and sends it to https://plane.watch/"

	app.Description = `This program takes a stream of plane tracking info (beast/avr/sbs1), tracks the planes and ` +
		`outputs all sorts if interesting information to the configured sink, including decoded and tracked planes in JSON format.` +
		"\n\n" +
		`example: pw_ingest --fetch=beast://crawled.mapwithlove.com:3004 --sink=amqp://guest:guest@localhost:5672/pw?queues=location-updates --tag="cool-stuff" --quiet simple`

	setup.IncludeSourceFlags(app)
	setup.IncludeSinkFlags(app)
	logging.IncludeVerbosityFlags(app)
	monitoring.IncludeMonitoringFlags(app, 9602)

	app.Commands = []*cli.Command{
		{
			Name:   "run",
			Usage:  "Gather ADSB data and sends it to the configured output. has a simple TUI",
			Action: run,
		},
		{
			Name:      "simple",
			Usage:     "Gather ADSB data and sends it to the configured output. just a log of info",
			Action:    runSimple,
			ArgsUsage: "[app.log - A file name to output to or stdout if not specified]",
		},
		{
			Name:   "daemon",
			Usage:  "Docker Daemon Mode",
			Action: runDaemon,
		},
		{
			Name:   "filter",
			Usage:  "Find examples from input",
			Action: runDfFilter,
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:  "icao",
					Usage: "Plane ICAO to filter on. e,g, --icao=E48DF6 --icao=123ABC",
				},
				&cli.BoolFlag{
					Name:  "locations-only",
					Usage: "Filter location updates only",
				},
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		logging.SetLoggingLevel(c)
		return nil
	}

	if err := app.Run(os.Args); nil != err {
		log.Error().Err(err).Msg("Finishing with an error")
		os.Exit(1)
	}
}

func commonSetup(c *cli.Context) (*tracker.Tracker, error) {
	monitoring.RunWebServer(c)

	// let's parse our URL forms

	trackerOpts := make([]tracker.Option, 0)
	trackerOpts = append(trackerOpts, tracker.WithPrometheusCounters(prometheusGaugeCurrentPlanes))
	trk := tracker.NewTracker(trackerOpts...)

	trk.AddMiddleware(dedupe.NewFilter())
	sinks, err := setup.HandleSinkFlags(c)
	if nil != err {
		return nil, err
	}
	for _, s := range sinks {
		trk.AddSink(s)
	}

	producers, err := setup.HandleSourceFlags(c)
	if nil != err {
		return nil, err
	}
	for _, p := range producers {
		trk.AddProducer(p)
	}

	return trk, nil
}

func runSimple(c *cli.Context) error {
	logging.ConfigureForCli()

	trk, err := commonSetup(c)

	if nil != err {
		return err
	}
	var opts []sink.Option
	if c.Bool("quiet") {
		opts = append(opts, sink.WithoutLoggingLocation())
	}
	trk.AddSink(sink.NewLoggerSink(opts...))

	go trk.StopOnCancel()
	trk.Wait()
	return nil
}

// runDfFilter is a special mode for hunting down DF examples from live inputs
func runDfFilter(c *cli.Context) error {
	logging.ConfigureForCli()

	trk, err := commonSetup(c)
	if nil != err {
		return err
	}
	var opts []sink.Option
	if c.Bool("quiet") {
		opts = append(opts, sink.WithoutLoggingLocation())
	}
	opts = append(opts, sink.WithoutLoggingLocation())
	trk.AddSink(sink.NewLoggerSink(opts...))

	var filterOpts []example_finder.Option
	if c.Bool("locations-only") {
		filterOpts = append(filterOpts, example_finder.WithDF17MessageTypeLocation())
	} else {
		filterOpts = append(filterOpts, example_finder.WithDownlinkFormatType(17))
	}
	for _, icao := range c.StringSlice("icao") {
		filterOpts = append(filterOpts, example_finder.WithPlaneIcaoStr(icao))
	}
	trk.AddMiddleware(example_finder.NewFilter(filterOpts...))

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

// runDaemon does not have pretty cli output (just JSON from logging)
func runDaemon(c *cli.Context) error {
	trk, err := commonSetup(c)
	if nil != err {
		return err
	}
	var opts []sink.Option
	opts = append(opts, sink.WithoutLoggingLocation())
	trk.AddSink(sink.NewLoggerSink(opts...))

	go trk.StopOnCancel()
	trk.Wait()
	return nil
}
