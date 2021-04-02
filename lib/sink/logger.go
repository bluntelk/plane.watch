package sink

import (
	"bufio"
	"fmt"
	"plane.watch/lib/tracker"
)

type (
	LoggerSink struct {
		Config
		bufOut *bufio.Writer
	}
)

func NewLoggerSink(opts ...Option) *LoggerSink {
	l := &LoggerSink{}
	for _, opt := range opts {
		opt(&l.Config)
	}
	l.bufOut = bufio.NewWriter(l.out)
	return l
}

func (l *LoggerSink) Finish() {
	_ = l.bufOut.Flush()
	l.Config.Finish()
}

func (l *LoggerSink) OnEvent(e tracker.Event) {
	switch e.(type) {
	case *tracker.LogEvent:
		_, _ = fmt.Fprintln(l.bufOut, e.String())
	case *tracker.PlaneLocationEvent:
		_, _ = fmt.Fprintln(l.bufOut, e.String())
	case *tracker.InfoEvent:
		_, _ = fmt.Fprintln(l.bufOut, e.String())
		_ = l.bufOut.Flush()
	}
}
