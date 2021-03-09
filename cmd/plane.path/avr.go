package main

import (
	"github.com/urfave/cli"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/mode_s"
	"sync"
	"time"
)


func parseAvr(c *cli.Context) error {
	newFrameFunc := func(line string) tracker.Frame {
		return mode_s.NewFrame(line, time.Now())
	}
	p, err := produceOutput(c, newFrameFunc)
	if nil != err {
		return err
	}

	ih := tracker.NewInputHandler(tracker.WithVerboseOutput())
	ih.AddProducer(p)
	ih.AddMiddleware(timeFiddler)
	ih.Wait()

	return writeResult(ih.Tracker, p.outFile)
}

var lastSeenMap sync.Map
// timeFiddler ensures we have enough time between messages for a plane to have travelled the distance it says it did
// this is because we do not have the timestamp for when it was collected when processing AVR frames
func timeFiddler(f tracker.Frame) tracker.Frame {
	switch f.(type) {
	case *mode_s.Frame:
		lastSeen, _ := lastSeenMap.LoadOrStore(f.Icao(), time.Now().Add(-24 *time.Hour))
		t := lastSeen.(time.Time)
		t = t.Add(10*time.Second)
		fp := f.(*mode_s.Frame)
		fp.SetTimeStamp(t)
		lastSeenMap.Store(f.Icao(), t)
	}

	return f
}

