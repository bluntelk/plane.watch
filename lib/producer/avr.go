package producer

import (
	"bufio"
	"io"
	"net"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/mode_s"
	"time"
)

func NewAvrListener(host, port string) tracker.Producer {
	p := NewProducer("AVR Listener for: " + net.JoinHostPort(host, port))

	go func() {
		// todo: Listen for incoming connection
	}()

	return p
}

func NewAvrFetcher(host, port string) tracker.Producer {
	p := NewProducer("AVR Fetcher for: " + net.JoinHostPort(host, port))

	source := newSource(p.String(), "avr://"+net.JoinHostPort(host, port))
	p.fetcher(host, port, func(conn net.Conn) error {
		scan := bufio.NewScanner(conn)
		for scan.Scan() {
			line := scan.Text()
			p.addFrame(mode_s.NewFrame(line, time.Now()), source)
			p.addDebug("AVR Frame: %s", line)
		}
		return scan.Err()
	})

	return p
}

func NewAvrFile(filePaths []string) tracker.Producer {
	p := NewProducer("AVR File")

	p.readFiles(filePaths, func(reader io.Reader, fileName string) error {
		var count uint64
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			count++
			p.AddEvent(
				tracker.NewFrameEvent(
					mode_s.NewFrame(scanner.Text(), time.Now()),
					newSource(p.String(), "file://"+fileName),
				),
			)
		}
		p.addInfo("We processed %d lines", count)
		return scanner.Err()
	})
	return p
}
