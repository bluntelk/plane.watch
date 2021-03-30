package producer

import (
	"bufio"
	"io"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/sbs1"
)

func NewSbs1File(filePaths []string) tracker.Producer {
	p := NewProducer("AVR File")

	p.readFiles(filePaths, func(reader io.Reader) error {
		var count uint64
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			count++
			p.AddEvent(tracker.NewFrameEvent(sbs1.NewFrame(scanner.Text())))
		}
		p.addInfo("We processed %d lines", count)
		return scanner.Err()

	})

	return p
}
