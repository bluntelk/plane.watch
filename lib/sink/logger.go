package sink

import (
	"github.com/rs/zerolog/log"
	"plane.watch/lib/tracker"
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

func (l *LoggerSink) Stop() {
	l.Config.Finish()
}

func (l *LoggerSink) OnEvent(e tracker.Event) {
	switch e.(type) {
	case *tracker.LogEvent:
		log.Info().Str("event", "method").Msg(e.String())
	case *tracker.PlaneLocationEvent:
		if l.logLocation {
			log.Info().Msg(e.String())
		}
	case *tracker.InfoEvent:
		i := e.(*tracker.InfoEvent)
		log.Info().
			Int("num-receivers", i.NumReceivers()).
			Uint64("num-frames", i.NumFrames()).
			Float64("uptime", i.Uptime()).
			Msg(e.String())
	}
}
