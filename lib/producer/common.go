package producer

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"errors"
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

	Avr = iota
	Beast
	Sbs1
)

type (
	producer struct {
		tracker.FrameSource
		producerType int

		out       chan tracker.Event
		outClosed bool
		outLocker sync.Mutex

		cmdChan chan int

		splitter   bufio.SplitFunc
		beastDelay bool

		run func()
	}

	Option func(*producer)
)

func New(opts ...Option) *producer {
	p := &producer{
		FrameSource: tracker.FrameSource{
			OriginIdentifier: "",
			Name:             "",
			RefLat:           nil,
			RefLon:           nil,
		},
		out:       make(chan tracker.Event, 100),
		outClosed: false,
		cmdChan:   make(chan int),
		run: func() {
			println("You did not specify any sources")
			os.Exit(1)
		},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// producer.New(WithFetcher(host, port), WithType(producer.Avr), WithRefLatLon(lat, lon))

func WithListener(host, port string) Option {
	return func(p *producer) {
		// TODO: implement a listener
	}
}

func WithFetcher(host, port string) Option {
	hp := net.JoinHostPort(host, port)
	return func(p *producer) {
		p.FrameSource.OriginIdentifier = hp
		p.run = func() {
			p.addInfo("Fetching From Host: %s:%s", host, port)
			p.fetcher(host, port, func(conn net.Conn) error {
				scan := bufio.NewScanner(conn)
				scan.Split(p.splitter)
				return p.readFromScanner(scan)
			})
		}
	}
}

func WithOriginName(name string) Option {
	return func(p *producer) {
		p.FrameSource.Name = name
	}
}

func WithFiles(filePaths []string) Option {
	return func(p *producer) {
		p.run = func() {
			p.readFiles(filePaths, func(reader io.Reader, fileName string) error {
				scanner := bufio.NewScanner(reader)
				p.FrameSource.OriginIdentifier = "file://"+fileName
				return p.readFromScanner(scanner)
			})
		}
	}
}

func WithBeastDelay(beastDelay bool) Option {
	return func(p *producer) {
		p.beastDelay = beastDelay
	}
}

func WithType(producerType int) Option {
	return func(p *producer) {
		p.producerType = producerType
		if producerType == Beast {
			p.splitter = ScanBeast
		} else {
			p.splitter = bufio.ScanLines
		}
	}
}

func (p *producer) readFromScanner(scan *bufio.Scanner) error {
	scan.Split(p.splitter)

	switch p.producerType {
	case Avr:
		return p.avrScanner(scan)
	case Sbs1:
		return p.sbsScanner(scan)
	case Beast:
		return p.beastScanner(scan)
	default:
		return errors.New("unknown producer type")
	}
}

// WithReferenceLatLon sets up the reference lat/lon for decoding surface position messages
func WithReferenceLatLon(lat, lon float64) Option {
	return func(p *producer) {
		p.RefLat = &lat
		p.RefLon = &lon
	}
}

func (p *producer) String() string {
	return p.Name
}

func (p *producer) Listen() chan tracker.Event {
	go p.run()
	return p.out
}

func (p *producer) addFrame(f tracker.Frame, s *tracker.FrameSource) {
	p.AddEvent(tracker.NewFrameEvent(f, s))
}

func (p *producer) addDebug(sfmt string, v ...interface{}) {
	p.AddEvent(tracker.NewLogEvent(tracker.LogLevelDebug, p.Name, fmt.Sprintf(sfmt, v...)))
}

func (p *producer) addInfo(sfmt string, v ...interface{}) {
	p.AddEvent(tracker.NewLogEvent(tracker.LogLevelInfo, p.Name, fmt.Sprintf(sfmt, v...)))
}

func (p *producer) addError(err error) {
	p.AddEvent(tracker.NewLogEvent(tracker.LogLevelError, p.Name, fmt.Sprint(err)))
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
			p.FrameSource.OriginIdentifier = "file://" + inFileName
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
