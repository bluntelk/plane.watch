package producer

import (
	"fmt"
	"plane.watch/lib/tracker"
)

type Producer struct {
	label string
	out  chan tracker.Frame
	errs chan tracker.LogItem
}

func NewProducer(label string) *Producer {
	return &Producer{
		label: label,
		out:  make(chan tracker.Frame),
		errs: make(chan tracker.LogItem),
	}
}

func (p *Producer) Listen() chan tracker.Frame {
	return p.out
}

func (p *Producer) addFrame(f tracker.Frame) {
	p.out <- f
}

func (p *Producer) addDebug(sfmt string, v ...interface{}) {
	p.errs <- tracker.LogItem{
		Level:   tracker.LogLevelDebug,
		Section: p.label,
		Message: fmt.Sprintf("Debug: "+ sfmt, v...),
	}
}

func (p *Producer) addInfo(sfmt string, v ...interface{}) {
	p.errs <- tracker.LogItem{
		Level:   tracker.LogLevelInfo,
		Section: p.label,
		Message: fmt.Sprintf("Info : "+ sfmt, v...),
	}
}

func (p *Producer) addError(err error) {
	p.errs <- tracker.LogItem{
		Level:   tracker.LogLevelError,
		Section: p.label,
		Message: fmt.Sprintf("Error: %s", err),
	}
}

func (p *Producer) Logs() chan tracker.LogItem {
	return p.errs
}

func (p *Producer) Stop() {

}
