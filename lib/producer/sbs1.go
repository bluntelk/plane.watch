package producer

import (
	"bufio"
	"io"
	"net"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/sbs1"
)

func NewSbs1Fetcher(host, port string) tracker.Producer {
	p := NewProducer("SBS1 Fetcher for: " + net.JoinHostPort(host, port))

	source := newSource(p.String(), "sbs1://"+net.JoinHostPort(host, port))
	p.fetcher(host, port, func(conn net.Conn) error {
		scan := bufio.NewScanner(conn)
		for scan.Scan() {
			line := scan.Text()
			p.addFrame(sbs1.NewFrame(scan.Text()), source)
			p.addDebug("AVR Frame: %s", line)
		}
		return scan.Err()
	})

	return p
}

func NewSbs1File(filePaths []string) tracker.Producer {
	p := NewProducer("AVR File")

	p.readFiles(filePaths, func(reader io.Reader, fileName string) error {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			p.AddEvent(
				tracker.NewFrameEvent(
					sbs1.NewFrame(scanner.Text()),
					newSource(p.String(), "file://"+fileName),
				),
			)
		}
		return scanner.Err()

	})

	return p
}
