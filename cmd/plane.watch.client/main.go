package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"net/url"
	"os"
	"plane.watch/lib/dedupe"
	"plane.watch/lib/producer"
	"plane.watch/lib/sink"
	"plane.watch/lib/tracker"
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
		`example: plane.watch.client --source=beast://crawled.mapwithlove.com:3004 --sink=amqp://guest:guest@localhost:5672/pw --tag="cool-stuff" --quiet simple`


	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "source",
			Usage:   "The Source in URL Form. [avr|beast|sbs1]://host:port",
			EnvVars: []string{"SOURCE"},
		},
		&cli.StringFlag{
			Name:    "sink",
			Usage:   "The place to send decoded JSON in URL Form. [redis|amqp]://user:pass@host:port/vhost",
			EnvVars: []string{"SINK"},
		},
		&cli.StringFlag{
			Name:    "tag",
			Usage:   "A value that is included in the payloads output to the Sinks. Useful for knowing where something came from",
			EnvVars: []string{"TAG"},
		},
		&cli.StringSliceFlag{
			Name:    "rabbit-queue",
			Usage:   fmt.Sprintf("The types of output we want from this binary. Valid options are %v", sink.AllQueues),
			EnvVars: []string{"QUEUES"},
		},

		&cli.StringFlag{
			Name:    "rabbit-host",
			Usage:   "the rabbitmq host to talk to",
			EnvVars: []string{"RABBITMQ_HOST"},
		},
		&cli.IntFlag{
			Name:    "rabbit-port",
			Value:   5672,
			Usage:   "The rabbitmq port to talk to",
			EnvVars: []string{"RABBITMQ_PORT"},
		},
		&cli.StringFlag{
			Name:    "rabbit-user",
			Value:   "guest",
			Usage:   "user for rabbitmq",
			EnvVars: []string{"RABBITMQ_USER"},
		},
		&cli.StringFlag{
			Name:    "rabbit-pass",
			Value:   "guest",
			Usage:   "rabbitmq password",
			EnvVars: []string{"RABBITMQ_PASS"},
		},
		&cli.StringFlag{
			Name:    "rabbit-vhost",
			Value:   "plane.watch",
			Usage:   "the virtual host on the rabbit server to use",
			EnvVars: []string{"RABBITMQ_VHOST"},
		},
		&cli.StringFlag{
			Name:    "dump1090-host",
			Value:   "",
			Usage:   "The host to read dump1090 from",
			EnvVars: []string{"DUMP1090_HOST"},
		},
		&cli.StringFlag{
			Name:    "dump1090-port",
			Value:   "30005",
			Usage:   "The port on the dump 1090 host to read from",
			EnvVars: []string{"DUMP1090_PORT"},
		},
		&cli.StringFlag{
			Name:  "feed-type",
			Value: "",
			Usage: "if not on a standard port, specify the type of feed (avr, sbs1, beast)",
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
	}

	if err := app.Run(os.Args); nil != err {
		fmt.Println(err)
	}
}

func commonSetup(c *cli.Context) (*tracker.Tracker, error) {
	isDebug := c.Bool("debug")
	isQuiet := c.Bool("quiet")
	refLat := c.Float64("refLat")
	refLon := c.Float64("refLon")
	redisHost := c.String("redis-host")
	redisPort := c.String("redis-port")
	rabbitHost := c.String("rabbit-host")
	rabbitPort := c.String("rabbit-port")
	rabbitUser := c.String("rabbit-user")
	rabbitPass := c.String("rabbit-pass")
	rabbitVHost := c.String("rabbit-vhost")
	rabbitQueues := c.StringSlice("rabbit-queue")

	sourceHost := c.String("dump1090-host")
	sourcePort := c.String("dump1090-port")
	sourceFormat := c.String("feed-type")

	sourceFileAvr := c.String("avr-file")
	sourceFileBeast := c.String("beast-file")

	// let's parse our URL forms
	urlSource := c.String("source")
	urlSink := c.String("sink")
	tag := c.String("tag")

	if "" != urlSource {
		parsedUrl, err := url.Parse(urlSource)
		if nil != err {
			return nil, err
		}

		switch strings.ToLower(parsedUrl.Scheme) {
		case "avr", "beast", "sbs1":
			sourceFormat = strings.ToLower(parsedUrl.Scheme)
		default:
			return nil, fmt.Errorf("unknown scheme: %s, expected one of [avr|beast|sbs1]", parsedUrl.Scheme)
		}
		sourceHost = parsedUrl.Hostname()
		sourcePort = parsedUrl.Port()
	}

	if "" != urlSink {
		parsedUrl, err := url.Parse(urlSink)
		if nil != err {
			return nil, err
		}
		switch strings.ToLower(parsedUrl.Scheme) {
		case "redis":
			redisHost = parsedUrl.Hostname()
			redisPort = parsedUrl.Port()
		case "amqp", "rabbitmq":
			rabbitHost = parsedUrl.Hostname()
			rabbitPort = parsedUrl.Port()
			rabbitUser = parsedUrl.User.Username()
			rabbitPass, _ = parsedUrl.User.Password()
			rabbitVHost = parsedUrl.Path
		default:
			return nil, fmt.Errorf("unknown scheme: %s, expected one of [redis|amqp|rabbitmq]", parsedUrl.Scheme)
		}

	}

	trackerOpts := make([]tracker.Option, 0)
	if isDebug {
		trackerOpts = append(trackerOpts, tracker.WithVerboseOutput())
	} else if isQuiet {
		trackerOpts = append(trackerOpts, tracker.WithQuietOutput())
	} else {
		trackerOpts = append(trackerOpts, tracker.WithInfoOutput())
	}
	trk := tracker.NewTracker(trackerOpts...)

	dedupeFilter := dedupe.NewFilter()
	trk.AddMiddleware(dedupeFilter.DeDupe)

	producerOpts := make([]producer.Option, 0)
	if refLat != 0 && refLon != 0 {
		producerOpts = append(producerOpts, producer.WithReferenceLatLon(refLat, refLon))
	}
	if "" != tag {
		producerOpts = append(producerOpts, producer.WithOriginName(tag))
	}

	if "" != redisHost {
		trk.AddSink(sink.NewRedisSink(sink.WithHost(redisHost, redisPort), sink.WithSourceTag(tag)))
	}
	if "" != rabbitHost {
		rabbitSink, err := sink.NewRabbitMqSink(
			sink.WithHost(rabbitHost, rabbitPort),
			sink.WithUserPass(rabbitUser, rabbitPass),
			sink.WithRabbitVhost(rabbitVHost),
			sink.WithRabbitQueues(rabbitQueues),
			sink.WithSourceTag(tag),
		)
		if nil != err {
			return nil, err
		}

		trk.AddSink(rabbitSink)
	}

	if "" != sourceHost {
		producerOpts = append(producerOpts, producer.WithFetcher(sourceHost, sourcePort))
		if "" != sourceFormat {
			switch sourceFormat {
			case "avr":
				producerOpts = append(producerOpts, producer.WithType(producer.Avr))
			case "sbs1":
				producerOpts = append(producerOpts, producer.WithType(producer.Sbs1))
			case "beast":
				producerOpts = append(producerOpts, producer.WithType(producer.Beast))
			default:
				return nil, errors.New("don't know how to handle type:" + sourceFormat)

			}
		} else {
			switch sourcePort {
			case "30002":
				producerOpts = append(producerOpts, producer.WithType(producer.Avr))
			case "30003":
				producerOpts = append(producerOpts, producer.WithType(producer.Sbs1))
			case "30005":
				producerOpts = append(producerOpts, producer.WithType(producer.Beast))
			default:
				return nil, errors.New("don't know how to handle port:" + sourcePort)
			}
		}
	} else {
		if "" != sourceFileAvr {
			producerOpts = append(
				producerOpts,
				producer.WithType(producer.Avr),
				producer.WithFiles([]string{sourceFileAvr}),
			)
		}
		if "" != sourceFileBeast {
			producerOpts = append(
				producerOpts,
				producer.WithType(producer.Beast),
				producer.WithFiles([]string{sourceFileBeast}),
			)
		}
	}

	trk.AddProducer(producer.New(producerOpts...))
	return trk, nil
}

func runSimple(c *cli.Context) error {
	trk, err := commonSetup(c)
	if nil != err {
		return err
	}
	opts := []sink.Option{sink.WithLogOutput(os.Stdout)}
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
