package tracker

import (
	"fmt"
	"time"
)

const LogEventType = "log-event"
const PlaneLocationEventType = "plane-location-event"

type (
	// an Event is something that we want to know about. This is the base of our sending of data
	Event interface {
		Type() string
		String() string
	}

	//LogEvent allows us to send out logs in a structured manner
	LogEvent struct {
		When    time.Time
		Level   int
		Section string
		Message string
	}

	// a PlaneLocationEvent is send whenever a planes information has been updated
	PlaneLocationEvent struct {
		p *Plane
	}

	// FrameEvent is for whenever we get a frame of data from our producers
	FrameEvent struct {
		frame Frame
	}
)

func (t *Tracker) AddEvent(e Event) {
	t.events <- e
}

func (t *Tracker) processEvents() {
	for e := range t.events {
		for _, sink := range t.sinks {
			sink.OnEvent(e)
		}
	}
}

func NewLogEvent(level int, section, msg string) *LogEvent {
	return &LogEvent{
		When:    time.Now(),
		Section: section,
		Level:   level,
		Message: msg,
	}
}

func (l *LogEvent) Type() string {
	return LogEventType
}
func (l *LogEvent) String() string {
	return fmt.Sprintf("%s - %s - %s", l.When.Format(time.Stamp), l.Section, l.Message)
}

func newPlaneLocationEvent(p *Plane) *PlaneLocationEvent {
	return &PlaneLocationEvent{p: p}
}

func (p *PlaneLocationEvent) Type() string {
	return PlaneLocationEventType
}
func (p *PlaneLocationEvent) String() string {
	return p.p.String()
}

func (f *FrameEvent) Type() string {
	return PlaneLocationEventType
}
func (f *FrameEvent) String() string {
	return f.frame.IcaoStr()
}
