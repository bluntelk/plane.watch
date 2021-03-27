package producer

import (
	"bufio"
	"math/rand"
	"net"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/mode_s"
	"sync"
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
	var conn net.Conn
	var wLock sync.RWMutex
	working := true

	isWorking := func () bool {
		wLock.RLock()
		defer wLock.RUnlock()
		return working
	}

	go func() {
		var backOff = time.Second
		var err error
		for isWorking() {
			p.addDebug("We are working!")
			wLock.Lock()
			conn, err = net.Dial("tcp", net.JoinHostPort(host, port))
			wLock.Unlock()
			if nil != err {
				p.addError(err)
				time.Sleep(backOff)
				backOff = backOff * 2 + ((time.Duration(rand.Intn(20)) * time.Millisecond * 100) - time.Second)
				if backOff > time.Minute {
					backOff = time.Minute
				}
				continue
			}
			p.addDebug("Connected!")
			backOff = time.Second
			scan := bufio.NewScanner(conn)
			for scan.Scan() {
				line := scan.Text()
				p.addFrame(mode_s.NewFrame(line, time.Now()))
				p.addDebug("AVR Frame: %s", line)
			}
			if err = scan.Err(); nil != err {
				p.addError(err)
			}
		}
		p.addDebug("Done with Producer", p)
		p.Cleanup()
	}()

	go func() {
		for cmd := range p.cmdChan {
			switch cmd {
			case cmdExit:
				p.addDebug("Got Cmd Exit")
				wLock.Lock()
				working = false
				if nil != conn {
					_ = conn.Close()
				}
				wLock.Unlock()
				return
			}
		}
	}()

	return p
}

func NewAvrFile(filePaths []string) tracker.Producer {
	p := NewProducer("AVR File")

	go func() {
		for line := range p.readFiles(filePaths) {
			p.AddEvent(tracker.NewFrameEvent(mode_s.NewFrame(line, time.Now())))
		}
		p.addDebug("Done with reading")
		p.Cleanup()
	}()

	go func() {
		for cmd := range p.cmdChan {
			switch cmd {
			case cmdExit:
				return
			}
		}
	}()

	return p
}
