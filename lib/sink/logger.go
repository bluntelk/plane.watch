package sink

import (
	"fmt"
	"plane.watch/lib/tracker"
)

type (
	LoggerSink struct {
		Config
	}
)

func NewLoggerSink(opts ...Option) *LoggerSink {
	r := &LoggerSink{}
	for _, opt := range opts {
		opt(&r.Config)
	}
	return r
}

func (l *LoggerSink) OnEvent(e tracker.Event) {
	switch e.(type) {
	case *tracker.LogEvent:
		_, _ = fmt.Fprintln(l.out, e.String())
	case *tracker.PlaneLocationEvent:
		_, _ = fmt.Fprintln(l.out, e.String())
	}
}
