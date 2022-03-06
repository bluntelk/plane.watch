package setup

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"math"
	"net/url"
	"plane.watch/lib/sink"
	"plane.watch/lib/tracker"
	"strconv"
	"strings"
)

var (
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
)

func IncludeSinkFlags(app *cli.App) {
	app.Flags = append(app.Flags, []cli.Flag{
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

		&cli.IntFlag{
			Name:  "sink-message-ttl",
			Value: 60,
			Usage: "Instruct our sinks to hold onto generated messages this long. In Seconds",
		},
	}...)
}

func HandleSinkFlags(c *cli.Context) ([]tracker.Sink, error) {
	defaultTTl := c.Int("sink-message-ttl")
	defaultTag := c.String("tag")
	defaultQueues := c.StringSlice("publish-types")
	sinks := make([]tracker.Sink, 0)

	for _, sinkUrl := range c.StringSlice("sink") {
		log.Debug().Str("sink-url", sinkUrl).Msg("With Sink")
		s, err := handleSink(sinkUrl, defaultTag, defaultTTl, defaultQueues, c.Bool("rabbitmq-test-queues"))
		if nil != err {
			log.Error().Err(err).Str("url", sinkUrl).Str("what", "sink").Msg("Failed setup sink")
			return nil, err
		} else {
			sinks = append(sinks, s)
		}
	}
	return sinks, nil
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
