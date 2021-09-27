package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"math"
	"net/url"
	"os"
	"plane.watch/lib/dedupe"
	"plane.watch/lib/logging"
	"plane.watch/lib/producer"
	"plane.watch/lib/sink"
	"plane.watch/lib/tracker"
	"strconv"
	"strings"
)

func main() {
	app := cli.NewApp()

	app.Version = "1.0.0"
	app.Name = "Plane Watch Client"
	app.Usage = "Reads from dump1090 and sends it to https://plane.watch/"

	app.Description = `This program takes a stream of plane tracking info (beast/avr/sbs1), tracks the planes and ` +
		`outputs all sorts if interesting information to the configured sink, including decoded and tracked planes in JSON format.` +
		"\n\n" +
		`example: plane.watch.client --fetch=beast://crawled.mapwithlove.com:3004 --sink=amqp://guest:guest@localhost:5672/pw?queues=location-updates --tag="cool-stuff" --quiet simple`

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
			Name:    "sink",
			Usage:   "The place to send decoded JSON in URL Form. [redis|amqp]://user:pass@host:port/vhost?ttl=60",
			EnvVars: []string{"SINK"},
		},
		&cli.StringSliceFlag{
			Name:    "file",
			Usage:   "The Source in URL Form. [avr|beast|sbs1]:///path/to/file?tag=MYTAG&refLat=-31.0&refLon=115.0",
			EnvVars: []string{"FILE"},
		},
		&cli.StringSliceFlag{
			Name:    "rabbit-queue",
			Usage:   fmt.Sprintf("The types of output we want from this binary. Valid options are %v", sink.AllQueues),
			EnvVars: []string{"QUEUES"},
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
		logging.SetVerboseOrQuiet(c.Bool("debug"), c.Bool("quiet"))
		return nil
	}

	if err := app.Run(os.Args); nil != err {
		log.Error().Err(err).Msg("Finishing with an error")
	}
}

func getTag(parsedUrl *url.URL, defaultTag string) string {
	if parsedUrl.Query().Has("tag") {
		return parsedUrl.Query().Get("tag")
	}
	return defaultTag
}
func getRef(parsedUrl *url.URL, what string, defaultRef float64) float64 {
	if parsedUrl.Query().Has(what) {
		f, err := strconv.ParseFloat(parsedUrl.Query().Get(what), 64)
		if nil == err {
			return f
		}
		log.Error().Err(err).Str("query_param", what).Msg("Could not determine reference value")
	}
	return defaultRef
}

func handleSink(urlSink, defaultTag string, defaultTtl int, defaultQueues []string) (tracker.Sink, error) {
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

	producerOpts := make([]producer.Option, 2)
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
	refLat := c.Float64("refLat")
	refLon := c.Float64("refLon")

	// let's parse our URL forms
	defaultTTl := c.Int("sink-message-ttl")
	defaultTag := c.String("tag")
	defaultQueues := c.StringSlice("rabbit-queue")

	trackerOpts := make([]tracker.Option, 0)
	trk := tracker.NewTracker(trackerOpts...)

	trk.AddMiddleware(dedupe.NewFilter())

	for _, sinkUrl := range c.StringSlice("sink") {
		log.Debug().Str("sink-url", sinkUrl).Send()
		p, err := handleSink(sinkUrl, defaultTag, defaultTTl, defaultQueues)
		if nil != err {
			log.Error().Err(err).Str("url", sinkUrl).Msgf("Failed to understand URL: %s", err)
		} else {
			trk.AddSink(p)
		}
	}
	for _, fetchUrl := range c.StringSlice("fetch") {
		p, err := handleSource(fetchUrl, defaultTag, refLat, refLon, false)
		if nil != err {
			log.Error().Err(err).Str("url", fetchUrl).Msgf("Failed to understand URL: %s", err)
		} else {
			trk.AddProducer(p)
		}
	}
	for _, listenUrl := range c.StringSlice("listen") {
		p, err := handleSource(listenUrl, defaultTag, refLat, refLon, true)
		if nil != err {
			log.Error().Err(err).Str("url", listenUrl).Msgf("Failed to understand URL: %s", err)
		} else {
			trk.AddProducer(p)
		}
	}
	for _, fileUrl := range c.StringSlice("file") {
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

// run is our method for running things
func runDaemon(c *cli.Context) error {
	logging.SetVerboseOrQuiet(false, true)
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
