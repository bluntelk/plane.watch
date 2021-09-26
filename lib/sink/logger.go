package sink

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"plane.watch/lib/tracker"
	"time"
)

type (
	LoggerSink struct {
		Config
	}
)

func NewLoggerSink(opts ...Option) *LoggerSink {
	l := &LoggerSink{}
	l.logLocation = true

	for _, opt := range opts {
		opt(&l.Config)
	}
	return l
}

func WithCliLogger() Option {
	return func(config *Config) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.UnixDate})
	}
}

func (l *LoggerSink) Finish() {
	l.Config.Finish()
}

func (l *LoggerSink) OnEvent(e tracker.Event) {
	switch e.(type) {
	case *tracker.LogEvent:
		log.Info().Str("event","method").Msg(e.String())
	case *tracker.PlaneLocationEvent:
		if l.logLocation {
			log.Info().Msg(e.String())
		}
	case *tracker.InfoEvent:
		log.Info().Msg(e.String())
	}
}
