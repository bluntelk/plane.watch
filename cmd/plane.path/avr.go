package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"plane.watch/lib/logging"
	"plane.watch/lib/producer"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/mode_s"
	"sync"
	"time"
)

type timeFiddler struct{
	events chan tracker.Event
}

func NewTimeFiddler() *timeFiddler {
	return &timeFiddler{
		events: make(chan tracker.Event),
	}
}

func (fm *timeFiddler) Listen() chan tracker.Event {
	return fm.events
}

func (fm *timeFiddler) String() string {
	return "Time Fiddler"
}

func (fm *timeFiddler) Stop() {
	close(fm.events)
}

func parseAvr(c *cli.Context) error {
	opts := make([]tracker.Option, 0)
	var verbose bool
	logging.SetVerboseOrQuiet(c.Bool("verbose"), c.Bool("quiet"))

	out, err := getOutput(c)
	if nil != err {
		fmt.Println(err)
	}

	trk := tracker.NewTracker(opts...)
	trk.AddMiddleware(NewTimeFiddler())
	if verbose {
		logging.SetVerboseOrQuiet(verbose, false)
	}
	trk.AddProducer(producer.New(producer.WithType(producer.Avr), producer.WithFiles(getFilePaths(c))))
	trk.Wait()
	return writeResult(trk, out)
}

var lastSeenMap sync.Map

// Handle ensures we have enough time between messages for a plane to have travelled the distance it says it did
// this is because we do not have the timestamp for when it was collected when processing AVR frames
func (fm *timeFiddler) Handle(f tracker.Frame, src *tracker.FrameSource) tracker.Frame {
	switch f.(type) {
	case *mode_s.Frame:
		lastSeen, _ := lastSeenMap.LoadOrStore(f.Icao(), time.Now().Add(-24*time.Hour))
		t := lastSeen.(time.Time)
		frame := f.(*mode_s.Frame)
		if 17 == frame.DownLinkType() {
			switch frame.MessageTypeString() {
			case mode_s.DF17FrameSurfacePos, mode_s.DF17FrameAirPositionGnss, mode_s.DF17FrameAirPositionBarometric:
				if frame.IsEven() {
					t = t.Add(10 * time.Second)
				}
			}
		} else {
			t = t.Add(100 * time.Millisecond)
		}
		fp := f.(*mode_s.Frame)
		fp.SetTimeStamp(t)
		lastSeenMap.Store(f.Icao(), t)
	}

	return f
}
