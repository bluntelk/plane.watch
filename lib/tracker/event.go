package tracker

import (
	"fmt"
	"time"
)

const LogEventType = "log-event"
const PlaneLocationEventType = "plane-location-event"
const InfoEventType = "info-event"

type (
	// Event is something that we want to know about. This is the base of our sending of data
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

	//PlaneLocationEvent is send whenever a planes information has been updated
	PlaneLocationEvent struct {
		new, removed bool
		p            *Plane
	}

	// FrameEvent is for whenever we get a frame of data from our producers
	FrameEvent struct {
		frame  Frame
		source *FrameSource
	}
	DedupedFrameEvent struct {
		FrameEvent
	}

	FrameSource struct {
		OriginIdentifier string
		Name, Tag        string
		RefLat, RefLon   *float64
	}

	// InfoEvent periodically sends out some interesting stats
	InfoEvent struct {
		receivedFrames uint64
		numReceivers   int
		uptime         float64
	}
)

func (t *Tracker) AddEvent(e Event) {
	t.eventSync.RLock()
	defer t.eventSync.RUnlock()
	if t.eventsOpen {
		t.events <- e
	}
}

func (t *Tracker) processEvents() {
	t.eventsWaiter.Add(1)
	for e := range t.events {
		for _, sink := range t.sinks {
			sink.OnEvent(e)
		}
	}
	t.eventsWaiter.Done()
}

func (l *LogEvent) Type() string {
	return LogEventType
}
func (l *LogEvent) String() string {
	return l.Message
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

func NewFrameEvent(f Frame, s *FrameSource) *FrameEvent {
	return &FrameEvent{frame: f, source: s}
}

func NewDedupedFrameEvent(f Frame, s *FrameSource) *DedupedFrameEvent {
	return &DedupedFrameEvent{FrameEvent{frame: f, source: s}}
}

func (f *FrameEvent) Type() string {
	return PlaneLocationEventType
}

func (f *FrameEvent) String() string {
	return f.frame.IcaoStr()
}

func (f *FrameEvent) Frame() Frame {
	return f.frame
}

func (f *FrameEvent) Source() *FrameSource {
	return f.source
}

func (i *InfoEvent) Type() string {
	return InfoEventType
}

func (i *InfoEvent) String() string {
	return fmt.Sprintf("Info: #feeders=%d, #frames=%d. uptime(s)=%0.2f", i.numReceivers, i.receivedFrames, i.uptime)
}

func (i *InfoEvent) NumReceivers() int {
	return i.numReceivers
}

func (i *InfoEvent) NumFrames() uint64 {
	return i.receivedFrames
}

func (i *InfoEvent) Uptime() float64 {
	return i.uptime
}
