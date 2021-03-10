package tracker

import (
	"errors"
	"io/ioutil"
	"os"
	"plane.watch/lib/tracker/mode_s"
	"plane.watch/lib/tracker/sbs1"
	"time"
)

type (
	// Option allows us to configure our new Tracker as we need it
	Option func(*Tracker)

	//LogItem allows us to send out logs in a structured manner
	LogItem struct {
		Level   int
		Section string
		Message string
	}

	// Frame is our general object for a tracking update, AVR, SBS1, Modes Beast Binary
	Frame interface {
		Icao() uint32
		IcaoStr() string
		Decode() (bool, error)
		TimeStamp() time.Time
	}
	// A Producer can listen or generate Frames, it provides the output via a channel that the handler then
	// processes further
	Producer interface {
		Listen() chan Frame
		Logs() chan LogItem
		Stop()
	}

	// a Middleware has a chance to modify a frame before we send it to the plane Tracker
	Middleware func(Frame) Frame
)

func (t *Tracker) setVerbosity(logLevel int) {
	t.logLevel = logLevel
	if logLevel == LogLevelQuiet {
		t.SetLoggerOutput(ioutil.Discard)
	} else {
		t.SetLoggerOutput(os.Stderr)
	}
}

func WithVerboseOutput() Option {
	return func(ih *Tracker) {
		ih.setVerbosity(LogLevelDebug)
	}
}
func WithInfoOutput() Option {
	return func(ih *Tracker) {
		ih.setVerbosity(LogLevelInfo)
	}
}
func WithQuietOutput() Option {
	return func(ih *Tracker) {
		ih.setVerbosity(LogLevelQuiet)
	}
}
func WithDecodeWorkerCount(numDecodeWorkers uint) Option {
	return func(ih *Tracker) {
		ih.decodeWorkerCount = numDecodeWorkers
	}
}

func (t *Tracker) Finish() {
	close(t.decodingQueue)
}

func (t *Tracker) AddProducer(p Producer) {
	t.debugMessage("Adding a producer")
	t.producers = append(t.producers, p)
	t.producerWaiter.Add(2)
	go func() {
		for f := range p.Listen() {
			t.decodingQueue <- f
		}
		t.producerWaiter.Done()
	}()
	go func() {
		for log := range p.Logs() {
			switch log.Level {
			case LogLevelQuiet:
				continue
			case LogLevelError:
				t.errorMessage("%s: %s", log.Section, log.Message)
			case LogLevelInfo:
				t.infoMessage("%s: %s", log.Section, log.Message)
			case LogLevelDebug:
				t.debugMessage("%s: %s", log.Section, log.Message)
			}
		}
		t.producerWaiter.Done()
	}()
	t.debugMessage("Just added a producer")
}

func (t *Tracker) AddMiddleware(m Middleware) {
	t.middlewares = append(t.middlewares, m)
}

// Wait waits for all producers to stop producing input and then return
func (t *Tracker) Wait() {
	t.debugMessage("and we are up and running...")
	t.producerWaiter.Wait()
	close(t.decodingQueue)
	t.decodingQueueWaiter.Wait()
}

func (t *Tracker) handleError(err error) {
	if nil != err {
		t.errorMessage("%s", err)
	}
}

func (t *Tracker) decodeQueue() {
	t.decodingQueueWaiter.Add(1)
	for f := range t.decodingQueue {
		ok, err := f.Decode()
		if nil != err {
			// the decode operation failed to produce valid output, and we tell someone about it
			t.handleError(err)
			continue
		}
		if !ok {
			// the decode operation did not produce a valid frame, but this is not an error
			// example: NoOp heartbeat
			continue
		}

		for _, m := range t.middlewares {
			f = m(f)
		}

		switch f.(type) {
		case *mode_s.Frame:
			t.HandleModeSFrame(f.(*mode_s.Frame))
		case *sbs1.Frame:
			t.HandleSbs1Frame(f.(*sbs1.Frame))
		default:
			t.handleError(errors.New("unknown frame type, cannot track"))
		}
	}
	t.decodingQueueWaiter.Done()
}
