package producer

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"os"
	"plane.watch/lib/tracker"
	"strings"
	"sync"
)

const (
	cmdExit = 1
)

type producer struct {
	label string
	out   chan tracker.Event
	outClosed bool
	outLocker sync.Mutex

	cmdChan chan int
}

func NewProducer(label string) *producer {
	return &producer{
		label:   label,
		out:     make(chan tracker.Event, 100),
		outClosed: false,
		cmdChan: make(chan int),
	}
}

func (p *producer) String() string {
	return p.label
}

func (p *producer) Listen() chan tracker.Event {
	return p.out
}

func (p *producer) addFrame(f tracker.Frame) {
	p.AddEvent(tracker.NewFrameEvent(f))
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
	p.outClosed=true
	close(p.out)
}

func (p *producer) readFiles(dataFiles []string) chan string {
	fileLines := make(chan string)
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

			var scanner *bufio.Scanner
			if isGzip {
				gzipFile, err = gzip.NewReader(inFile)
				if nil != err {
					p.addError(err)
					continue
				}
				scanner = bufio.NewScanner(gzipFile)
			} else if isBzip2 {
				bzip2File := bzip2.NewReader(inFile)
				scanner = bufio.NewScanner(bzip2File)
			} else {
				scanner = bufio.NewScanner(inFile)
			}
			for scanner.Scan() {
				fileLines <- scanner.Text()
			}
			_ = inFile.Close()
			p.addDebug("Finished with file %s", inFileName)
		}
		close(fileLines)
		p.addDebug("Done loading lines")
	}()
	return fileLines
}
