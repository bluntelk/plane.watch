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
	out   chan tracker.Frame
	logs  chan tracker.LogItem

	cmdChan chan int
}

func NewProducer(label string) *producer {
	return &producer{
		label:   label,
		out:     make(chan tracker.Frame),
		logs:    make(chan tracker.LogItem),
		cmdChan: make(chan int),
	}
}

func (p *producer) Listen() chan tracker.Frame {
	return p.out
}

func (p *producer) addFrame(f tracker.Frame) {
	p.out <- f
}

func (p *producer) addDebug(sfmt string, v ...interface{}) {
	p.logs <- tracker.LogItem{
		Level:   tracker.LogLevelDebug,
		Section: p.label,
		Message: fmt.Sprintf("Debug: "+sfmt, v...),
	}
}

func (p *producer) addInfo(sfmt string, v ...interface{}) {
	p.logs <- tracker.LogItem{
		Level:   tracker.LogLevelInfo,
		Section: p.label,
		Message: fmt.Sprintf("Info : "+sfmt, v...),
	}
}

func (p *producer) addError(err error) {
	p.logs <- tracker.LogItem{
		Level:   tracker.LogLevelError,
		Section: p.label,
		Message: fmt.Sprintf("Error: %s", err),
	}
}

func (p *producer) Logs() chan tracker.LogItem {
	return p.logs
}

func (p *producer) Stop() {
	p.cmdChan <- cmdExit
}
