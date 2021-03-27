package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"plane.watch/lib/producer"
	"plane.watch/lib/sink"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/mode_s"
	"sync"
	"time"
)

func parseAvr(c *cli.Context) error {
	opts := make([]tracker.Option,0)
	if c.GlobalBool("verbose") {
		opts = append(opts, tracker.WithVerboseOutput())
	} else {
		opts = append(opts, tracker.WithInfoOutput())
	}
	out, err := getOutput(c)
	if nil != err {
		fmt.Println(err)
	}

	ih := tracker.NewTracker(opts...)
	ih.AddProducer(producer.NewAvrFile(getFilePaths(c)))
	ih.AddMiddleware(timeFiddler)
	ih.AddSink(sink.NewLoggerSink(sink.WithLogOutput(os.Stderr)))
	ih.Wait()
	return writeResult(ih, out)
}

var lastSeenMap sync.Map

// timeFiddler ensures we have enough time between messages for a plane to have travelled the distance it says it did
// this is because we do not have the timestamp for when it was collected when processing AVR frames
func timeFiddler(f tracker.Frame) tracker.Frame {
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
