package producer

import (
	"bufio"
	"net"
	"plane.watch/lib/tracker"
)

func NewBeastFile(host, port string) tracker.Producer {
	p := NewProducer("Beast Listener for: " + net.JoinHostPort(host, port))

	go func() {
		// todo: Listen for incoming connection
	}()

	return p
}

func NewBeastListener(host, port string) tracker.Producer {
	p := NewProducer("Beast Listener for: " + net.JoinHostPort(host, port))

	go func() {
		// todo: Listen for incoming connection
	}()

	return p
}

func ScanBeast(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	return 1, data[0:1], nil
}


func NewBeastFetcher(host, port string) tracker.Producer {
	p := NewProducer("Beast Fetcher for: " + net.JoinHostPort(host, port))

	p.fetcher(host, port, func(conn net.Conn) error {
		scan := bufio.NewScanner(conn)
		scan.Split(ScanBeast)

		return nil
	})

	return p
}