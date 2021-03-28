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
		new, removed bool
		p *Plane
	}

	// FrameEvent is for whenever we get a frame of data from our producers
	FrameEvent struct {
		frame Frame
	}

	// InfoEvent periodically sends out some interesting stats
	InfoEvent struct {
		receivedFrames uint64
		numReceivers uint
		uptime uint64
	}

)

func (t *Tracker) AddEvent(e Event) {
	t.events <- e
}

func (t *Tracker) processEvents() {
	for e := range t.events {
		for _, sink := range t.sinks {
			go sink.OnEvent(e)
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
	if l.Level > LogLevelDebug {
		l.Level = LogLevelDebug
	}
	lvl := Levels[l.Level]
	return fmt.Sprintf("%s - %s - %s - %s", l.When.Format(time.Stamp), lvl, l.Section, l.Message)
}

func newPlaneLocationEvent(p *Plane) *PlaneLocationEvent {
	return &PlaneLocationEvent{p: p}
}

func newPlaneActionEvent(p *Plane, isNew, isRemoved bool) *PlaneLocationEvent {
	return &PlaneLocationEvent{p: p, new: isNew, removed: isRemoved}
}

func (p *PlaneLocationEvent) Type() string {
	return PlaneLocationEventType
}
func (p *PlaneLocationEvent) String() string {
	return p.p.String()
}
func (p *PlaneLocationEvent) Plane() *Plane {
	return p.p
}
func (p *PlaneLocationEvent) New() bool {
	return p.new
}
func (p *PlaneLocationEvent) Removed() bool {
	return p.removed
}

func (f *FrameEvent) Type() string {
	return PlaneLocationEventType
}
func (f *FrameEvent) String() string {
	return f.frame.IcaoStr()
}
