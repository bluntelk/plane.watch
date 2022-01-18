package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"math"
	"net/url"
	"os"
	"plane.watch/lib/dedupe"
	"plane.watch/lib/logging"
	"plane.watch/lib/monitoring"
	"plane.watch/lib/producer"
	"plane.watch/lib/sink"
	"plane.watch/lib/tracker"
	"strconv"
	"strings"
)

var (
	prometheusInputBeastFrames = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pw_ingest_input_beast_total",
		Help: "The total number of beast frames processed.",
	})
	prometheusInputAvrFrames = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pw_ingest_input_avr_total",
		Help: "The total number of AVR frames processed.",
	})
	prometheusInputSbs1Frames = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pw_ingest_input_sbs1_total",
		Help: "The total number of SBS1 frames processed.",
	})

	prometheusOutputFrame = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pw_ingest_output_frame_total",
		Help: "The total number of raw frames output. (no dedupe)",
	})
	prometheusOutputFrameDedupe = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pw_ingest_output_frame_dedupe_total",
		Help: "The total number of deduped frames output.",
	})
	prometheusOutputPlaneLocation = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pw_ingest_output_location_update_total",
		Help: "The total number of plane location events output.",
	})
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

	app.Flags = []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "fetch",
			Usage:   "The Source in URL Form. [avr|beast|sbs1]://host:port?tag=MYTAG&refLat=-31.0&refLon=115.0",
			EnvVars: []string{"SOURCE"},
		},
		&cli.StringSliceFlag{
			Name:    "listen",
			Usage:   "The Source in URL Form. [avr|beast|sbs1]://host:port?tag=MYTAG&refLat=-31.0&refLon=115.0",
			EnvVars: []string{"LISTEN"},
		},
		&cli.StringSliceFlag{
			Name:    "file",
			Usage:   "The Source in URL Form. [avr|beast|sbs1]:///path/to/file?tag=MYTAG&refLat=-31.0&refLon=115.0",
			EnvVars: []string{"FILE"},
		},
		&cli.StringSliceFlag{
			Name:    "sink",
			Usage:   "The place to send decoded JSON in URL Form. [redis|amqp]://user:pass@host:port/vhost?ttl=60",
			EnvVars: []string{"SINK"},
		},
		&cli.StringSliceFlag{
			Name:    "publish-types",
			Usage:   fmt.Sprintf("The types of output we want to publish from this binary. Default: All Types. Valid options are %v", sink.AllQueues),
			EnvVars: []string{"PUBLISH"},
		},
		&cli.BoolFlag{
			Name:  "rabbitmq-test-queues",
			Usage: fmt.Sprintf("Create a queue (named after the publishing routing key) and bind it. This allows you to see the messages being published."),
		},

		&cli.StringFlag{
			Name:    "tag",
			Usage:   "A value that is included in the payloads output to the Sinks. Useful for knowing where something came from",
			EnvVars: []string{"TAG"},
		},
		&cli.IntFlag{
			Name:  "sink-message-ttl",
			Value: 60,
			Usage: "Instruct our sinks to hold onto generated messages this long. In Seconds",
		},
		&cli.Float64Flag{
			Name:  "ref-lat",
			Usage: "The reference latitude for decoding messages. Needs to be within 45nm of where the messages are generated.",
		},
		&cli.Float64Flag{
			Name:  "ref-lon",
			Usage: "The reference longitude for decoding messages. Needs to be within 45nm of where the messages are generated.",
		},
	}
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
	}

	app.Before = func(c *cli.Context) error {
		logging.SetLoggingLevel(c)
		return nil
	}

	if err := app.Run(os.Args); nil != err {
		log.Error().Err(err).Msg("Finishing with an error")
	}
}

func getTag(parsedUrl *url.URL, defaultTag string) string {
	if nil == parsedUrl {
		return ""
	}
	if parsedUrl.Query().Has("tag") {
		return parsedUrl.Query().Get("tag")
	}
	return defaultTag
}
func getRef(parsedUrl *url.URL, what string, defaultRef float64) float64 {
	if nil == parsedUrl {
		return 0
	}
	if parsedUrl.Query().Has(what) {
		f, err := strconv.ParseFloat(parsedUrl.Query().Get(what), 64)
		if nil == err {
			return f
		}
		log.Error().Err(err).Str("query_param", what).Msg("Could not determine reference value")
	}
	return defaultRef
}

func handleSink(urlSink, defaultTag string, defaultTtl int, defaultQueues []string, rabbitmqTestQueues bool) (tracker.Sink, error) {
	parsedUrl, err := url.Parse(urlSink)
	if nil != err {
		return nil, err
	}
	switch strings.ToLower(parsedUrl.Scheme) {
	case "redis":
		redisHost := parsedUrl.Hostname()
		redisPort := parsedUrl.Port()
		return sink.NewRedisSink(sink.WithHost(redisHost, redisPort), sink.WithSourceTag(getTag(parsedUrl, defaultTag))), nil
	case "amqp", "rabbitmq":
		rabbitPass, _ := parsedUrl.User.Password()
		messageTtl := defaultTtl
		if parsedUrl.Query().Has("ttl") {
			var requestedTtl int64
			requestedTtl, err = strconv.ParseInt(parsedUrl.Query().Get("ttl"), 10, 32)
			if requestedTtl > 0 && requestedTtl < math.MaxInt32 {
				messageTtl = int(requestedTtl)
			}
		}

		rabbitQueues := defaultQueues
		if parsedUrl.Query().Has("queues") {
			rabbitQueues = strings.Split(parsedUrl.Query().Get("queues"), ",")
		}

		return sink.NewRabbitMqSink(
			sink.WithHost(parsedUrl.Hostname(), parsedUrl.Port()),
			sink.WithUserPass(parsedUrl.User.Username(), rabbitPass),
			sink.WithRabbitVhost(parsedUrl.Path),
			sink.WithRabbitQueues(rabbitQueues),
			sink.WithSourceTag(getTag(parsedUrl, defaultTag)),
			sink.WithMessageTtl(messageTtl),
			sink.WithRabbitTestQueues(rabbitmqTestQueues),
			sink.WithPrometheusCounters(prometheusOutputFrame, prometheusOutputFrameDedupe, prometheusOutputPlaneLocation),
		)

	default:
		return nil, fmt.Errorf("unknown scheme: %s, expected one of [redis|amqp|rabbitmq]", parsedUrl.Scheme)
	}

}

func handleSource(urlSource, defaultTag string, defaultRefLat, defaultRefLon float64, listen bool) (tracker.Producer, error) {
	parsedUrl, err := url.Parse(urlSource)
	if nil != err {
		return nil, err
	}

	producerOpts := make([]producer.Option, 3)
	producerOpts[0] = producer.WithSourceTag(getTag(parsedUrl, defaultTag))

	switch strings.ToLower(parsedUrl.Scheme) {
	case "avr":
		producerOpts[1] = producer.WithType(producer.Avr)
	case "beast":
		producerOpts[1] = producer.WithType(producer.Beast)
	case "sbs1":
		producerOpts[1] = producer.WithType(producer.Sbs1)
	default:
		return nil, fmt.Errorf("unknown scheme: %s, expected one of [avr|beast|sbs1]", parsedUrl.Scheme)
	}
	producerOpts[2] = producer.WithPrometheusCounters(prometheusInputAvrFrames, prometheusInputBeastFrames, prometheusInputSbs1Frames)

	refLat := getRef(parsedUrl, "refLat", defaultRefLat)
	refLon := getRef(parsedUrl, "refLon", defaultRefLon)
	if refLat != 0 && refLon != 0 {
		producerOpts = append(producerOpts, producer.WithReferenceLatLon(refLat, refLon))
	}

	if listen {
		producerOpts = append(producerOpts, producer.WithListener(parsedUrl.Hostname(), parsedUrl.Port()))
	} else {
		producerOpts = append(producerOpts, producer.WithFetcher(parsedUrl.Hostname(), parsedUrl.Port()))
	}

	return producer.New(producerOpts...), nil
}

func handleFileSource(urlFile, defaultTag string, defaultRefLat, defaultRefLon float64) (tracker.Producer, error) {
	parsedUrl, err := url.Parse(urlFile)
	if nil != err {
		return nil, err
	}
	producerOpts := make([]producer.Option, 1)
	switch strings.ToLower(parsedUrl.Scheme) {
	case "avr":
		producerOpts[0] = producer.WithType(producer.Avr)
	case "beast":
		producerOpts[0] = producer.WithType(producer.Beast)
	case "sbs1":
		producerOpts[0] = producer.WithType(producer.Sbs1)
	default:
		return nil, fmt.Errorf("unknown file Type: %s", parsedUrl.Scheme)
	}
	refLat := getRef(parsedUrl, "refLat", defaultRefLat)
	refLon := getRef(parsedUrl, "refLon", defaultRefLon)
	if refLat != 0 && refLon != 0 {
		producerOpts = append(producerOpts, producer.WithReferenceLatLon(refLat, refLon))
	}

	producerOpts = append(
		producerOpts,
		producer.WithSourceTag(getTag(parsedUrl, defaultTag)),
		producer.WithFiles([]string{parsedUrl.Path}),
	)

	return producer.New(producerOpts...), nil
}

func commonSetup(c *cli.Context) (*tracker.Tracker, error) {
	monitoring.RunWebServer(c)
	refLat := c.Float64("refLat")
	refLon := c.Float64("refLon")

	// let's parse our URL forms
	defaultTTl := c.Int("sink-message-ttl")
	defaultTag := c.String("tag")
	defaultQueues := c.StringSlice("publish-types")

	trackerOpts := make([]tracker.Option, 0)
	trackerOpts = append(trackerOpts, tracker.WithPrometheusCounters(prometheusGaugeCurrentPlanes))
	trk := tracker.NewTracker(trackerOpts...)

	trk.AddMiddleware(dedupe.NewFilter())

	for _, sinkUrl := range c.StringSlice("sink") {
		log.Debug().Str("sink-url", sinkUrl).Msg("With Sink")
		p, err := handleSink(sinkUrl, defaultTag, defaultTTl, defaultQueues, c.Bool("rabbitmq-test-queues"))
		if nil != err {
			log.Error().Err(err).Str("url", sinkUrl).Msgf("Failed to understand URL: %s", err)
		} else {
			trk.AddSink(p)
		}
	}
	for _, fetchUrl := range c.StringSlice("fetch") {
		log.Debug().Str("fetch-url", fetchUrl).Msg("With Fetch")
		p, err := handleSource(fetchUrl, defaultTag, refLat, refLon, false)
		if nil != err {
			log.Error().Err(err).Str("url", fetchUrl).Msgf("Failed to understand URL: %s", err)
		} else {
			trk.AddProducer(p)
		}
	}
	for _, listenUrl := range c.StringSlice("listen") {
		log.Debug().Str("listen-url", listenUrl).Msg("With Listen")
		p, err := handleSource(listenUrl, defaultTag, refLat, refLon, true)
		if nil != err {
			log.Error().Err(err).Str("url", listenUrl).Msgf("Failed to understand URL: %s", err)
		} else {
			trk.AddProducer(p)
		}
	}
	for _, fileUrl := range c.StringSlice("file") {
		log.Debug().Str("file-url", fileUrl).Msg("With File")
		p, err := handleFileSource(fileUrl, defaultTag, refLat, refLon)
		if nil != err {
			log.Error().Err(err).Str("url", fileUrl).Msgf("Failed to understand URL: %s", err)
		} else {
			trk.AddProducer(p)
		}
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

	trk.Wait()
	return nil
}
