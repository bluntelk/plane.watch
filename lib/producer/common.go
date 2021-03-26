package producer

import (
	"fmt"
	"plane.watch/lib/tracker"
)

const (
	cmdExit = 1
)

type producer struct {
	label string
	out   chan tracker.Event

	cmdChan chan int
}

func NewProducer(label string) *producer {
	return &producer{
		label:   label,
		out:     make(chan tracker.Event, 100),
		cmdChan: make(chan int),
	}
}

func (p *producer) Listen() chan tracker.Event {
	return p.out
}

func (p *producer) addFrame(f tracker.Frame) {
	p.out <- tracker.NewFrameEvent(f)
}

func (p *producer) addDebug(sfmt string, v ...interface{}) {
	p.out <- tracker.NewLogEvent(tracker.LogLevelDebug, p.label, fmt.Sprintf(sfmt, v...))
}

func (p *producer) addInfo(sfmt string, v ...interface{}) {
	p.out <- tracker.NewLogEvent(tracker.LogLevelInfo, p.label, fmt.Sprintf(sfmt, v...))
}

func (p *producer) addError(err error) {
	p.out <- tracker.NewLogEvent(tracker.LogLevelError, p.label, fmt.Sprint(err))
}

func (p *producer) Stop() {
	p.cmdChan <- cmdExit
}

func (p *producer) Cleanup() {
	close(p.out)
}
