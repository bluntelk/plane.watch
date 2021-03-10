package producer

import (
	"bufio"
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
	// todo: gracefully close/stop

	go func() {
		var conn net.Conn
		var err error
		for {
			conn, err = net.Dial("tcp", net.JoinHostPort(host, port))
			if nil != err {
				p.addError(err)
				continue
			}
			scan := bufio.NewScanner(conn)
			for scan.Scan() {
				p.addFrame(mode_s.NewFrame(scan.Text(), time.Now()))
			}
		}
	}()

	return p
}
