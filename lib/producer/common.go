package producer

import (
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"plane.watch/lib/tracker"
	"strings"
	"sync"
	"time"
)

const (
	cmdExit = 1
)

type (
	producer struct {
		label, ident string
		out          chan tracker.Event
		outClosed    bool
		outLocker    sync.Mutex

		cmdChan chan int
	}
)

func NewProducer(label string) *producer {
	p := &producer{
		label:     label,
		out:       make(chan tracker.Event, 100),
		outClosed: false,
		cmdChan:   make(chan int),
	}

	return p
}

func newSource(label, ident string) tracker.Source {
	return tracker.Source{
		OriginIdentifier: ident,
		Name:             label,
	}
}

func (p *producer) String() string {
	return p.label
}

func (p *producer) Listen() chan tracker.Event {
	return p.out
}

func (p *producer) addFrame(f tracker.Frame, s tracker.Source) {
	p.AddEvent(tracker.NewFrameEvent(f, s))
}

func (p *producer) addDebug(sfmt string, v ...interface{}) {
	p.AddEvent(tracker.NewLogEvent(tracker.LogLevelDebug, p.label, fmt.Sprintf(sfmt, v...)))
}

func (p *producer) addInfo(sfmt string, v ...interface{}) {
	p.AddEvent(tracker.NewLogEvent(tracker.LogLevelInfo, p.label, fmt.Sprintf(sfmt, v...)))
}

func (p *producer) addError(err error) {
	p.AddEvent(tracker.NewLogEvent(tracker.LogLevelError, p.label, fmt.Sprint(err)))
}

func (p *producer) Stop() {
	p.cmdChan <- cmdExit
}

func (p *producer) AddEvent(e tracker.Event) {
	p.outLocker.Lock()
	defer p.outLocker.Unlock()
	if !p.outClosed {
		p.out <- e
	}
}

func (p *producer) Cleanup() {
	p.outLocker.Lock()
	defer p.outLocker.Unlock()
	if p.outClosed {
		return
	}
	p.outClosed = true
	close(p.out)
}

func (p *producer) readFiles(dataFiles []string, read func(io.Reader, string) error) {
	var err error
	var inFile *os.File
	var gzipFile *gzip.Reader
	go func() {
		for _, inFileName := range dataFiles {
			p.addDebug("Loading lines from %s", inFileName)
			inFile, err = os.Open(inFileName)
			if err != nil {
				p.addError(fmt.Errorf("failed to open file {%s}: %s", inFileName, err))
				continue
			}

			isGzip := strings.ToLower(inFileName[len(inFileName)-2:]) == "gz"
			isBzip2 := strings.ToLower(inFileName[len(inFileName)-3:]) == "bz2"
			p.addDebug("Is Gzip? %t Is Bzip2? %t", isGzip, isBzip2)

			if isGzip {
				gzipFile, err = gzip.NewReader(inFile)
				if nil != err {
					err = read(gzipFile, inFileName)
				}
				err = read(gzipFile, inFileName)
			} else if isBzip2 {
				bzip2File := bzip2.NewReader(inFile)
				err = read(bzip2File, inFileName)
			} else {
				err = read(inFile, inFileName)
			}
			if nil != err {
				p.addError(err)
			}
			_ = inFile.Close()
			p.addDebug("Finished with file %s", inFileName)
		}
		p.addDebug("Done loading lines")
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
}

func (p *producer) fetcher(host, port string, read func(net.Conn) error) {
	var conn net.Conn
	var wLock sync.RWMutex
	working := true

	isWorking := func() bool {
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
				backOff = backOff*2 + ((time.Duration(rand.Intn(20)) * time.Millisecond * 100) - time.Second)
				if backOff > time.Minute {
					backOff = time.Minute
				}
				continue
			}
			p.addDebug("Connected!")
			backOff = time.Second

			if err = read(conn); nil != err {
				p.addError(err)
			}
		}
		p.addDebug("Done with Producer %s", p)
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

}
