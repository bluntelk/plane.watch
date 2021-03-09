package tracker

import (
	"errors"
	"os"
	"plane.watch/lib/tracker/mode_s"
	"plane.watch/lib/tracker/sbs1"
	"sync"
	"time"
)

type (
	InputHandlerOption func(*InputHandler)

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
		Stop()
	}

	// a Middleware has a chance to modify a frame before we send it to the plane Tracker
	Middleware func(Frame) Frame

	InputHandler struct {
		producers   []Producer
		middlewares []Middleware

		producerWaiter sync.WaitGroup

		decodeWorkerCount   uint
		decodingQueue       chan Frame
		decodingQueueWaiter sync.WaitGroup
		Tracker             *Tracker
	}
)

func (ih *InputHandler) setVerbosity(logLevel int) {
	ih.Tracker.logLevel = logLevel
	ih.Tracker.SetLoggerOutput(os.Stderr)
}

func NewInputHandler(opts ...InputHandlerOption) *InputHandler {
	ph := &InputHandler{
		producers:     []Producer{},
		middlewares:   []Middleware{},
		decodingQueue: make(chan Frame, 1000), // a nice deep buffer
		Tracker:       NewTracker(),
	}

	for _, opt := range opts {
		opt(ph)
	}

	for i := 0; i < 5; i++ {
		go ph.decodeQueue()
	}

	return ph
}

func WithVerboseOutput() InputHandlerOption {
	return func(ih *InputHandler) {
		ih.setVerbosity(logLevelDebug)
	}
}
func WithInfoOutput() InputHandlerOption {
	return func(ih *InputHandler) {
		ih.setVerbosity(logLevelInfo)
	}
}
func WithQuietOutput() InputHandlerOption {
	return func(ih *InputHandler) {
		ih.setVerbosity(logLevelQuiet)
	}
}
func WithDecodeWorkerCount(numDecodeWorkers uint) InputHandlerOption {
	return func(ih *InputHandler) {
		ih.decodeWorkerCount = numDecodeWorkers
	}
}

func (ih *InputHandler) Finish() {
	close(ih.decodingQueue)
}

func (ih *InputHandler) AddProducer(p Producer) {
	ih.producers = append(ih.producers, p)
	ih.producerWaiter.Add(1)
	go func() {
		for f := range p.Listen() {
			ih.decodingQueue <- f
		}
		ih.producerWaiter.Done()
	}()
}

func (ih *InputHandler) AddMiddleware(m Middleware) {
	ih.middlewares = append(ih.middlewares, m)
}

// Wait waits for all producers to stop producing input and then return
func (ih *InputHandler) Wait() {
	ih.producerWaiter.Wait()
	close(ih.decodingQueue)
	ih.decodingQueueWaiter.Wait()
}

func (ih *InputHandler) handleError(err error) {
	if nil != err {
		ih.Tracker.errorMessage("%s", err)
	}
}

func (ih *InputHandler) decodeQueue() {
	ih.decodingQueueWaiter.Add(1)
	for f := range ih.decodingQueue {
		ok, err := f.Decode()
		if nil != err {
			// the decode operation failed to produce valid output, and we tell someone about it
			ih.handleError(err)
			continue
		}
		if !ok {
			// the decode operation did not produce a valid frame, but this is not an error
			// example: NoOp heartbeat
			continue
		}

		for _, m := range ih.middlewares {
			f = m(f)
		}

		switch f.(type) {
		case *mode_s.Frame:
			ih.Tracker.HandleModeSFrame(f.(*mode_s.Frame))
		case *sbs1.Frame:
			ih.Tracker.HandleSbs1Frame(f.(*sbs1.Frame))
		default:
			ih.handleError(errors.New("unknown frame type, cannot track"))
		}
	}
	ih.decodingQueueWaiter.Done()
}
