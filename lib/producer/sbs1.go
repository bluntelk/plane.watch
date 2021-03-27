package producer

import (
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/sbs1"
)

func NewSbs1File(filePaths []string) tracker.Producer {
	p := NewProducer("AVR File")

	go func() {
		for line := range p.readFiles(filePaths) {
			p.AddEvent(tracker.NewFrameEvent(sbs1.NewFrame(line)))
		}
		p.Cleanup()
	}()

	return p
}
