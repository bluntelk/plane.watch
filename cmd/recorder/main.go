package main

import (
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
	"plane.watch/lib/logging"
	"plane.watch/lib/monitoring"
	"plane.watch/lib/setup"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/beast"
	"plane.watch/lib/tracker/mode_s"
	"plane.watch/lib/tracker/sbs1"
	"sync"
)

type (
	frameProcessor struct {
		events         chan tracker.Event
		logFileHandles sync.Map
	}
)

func main() {
	app := cli.NewApp()

	app.Description = "Records Beast/AVR/SBS1 source to a file"

	setup.IncludeSourceFlags(app)
	logging.IncludeVerbosityFlags(app)
	monitoring.IncludeMonitoringFlags(app, 9606)

	app.Before = func(c *cli.Context) error {
		logging.SetLoggingLevel(c)
		return nil
	}

	app.Action = runCli

	if err := app.Run(os.Args); nil != err {
		log.Error().Err(err).Msg("did not complete successfully")
		os.Exit(1)
	}
}

func runCli(c *cli.Context) error {
	logging.ConfigureForCli()
	monitoring.RunWebServer(c)

	trackerOpts := make([]tracker.Option, 0)
	trk := tracker.NewTracker(trackerOpts...)

	producers, err := setup.HandleSourceFlags(c)
	if nil != err {
		return err
	}
	for _, p := range producers {
		trk.AddProducer(p)
	}

	// have our own middlewares
	fp := newFrameLogger()
	trk.AddMiddleware(fp)

	// setup a file sync?

	trk.Wait()
	return nil
}

func newFrameLogger() *frameProcessor {
	return &frameProcessor{
		events: make(chan tracker.Event),
	}
}

func (fp *frameProcessor) Stop() {
	close(fp.events)
	fp.logFileHandles.Range(func(name, handle interface{}) bool {
		fp.logFileHandles.Delete(name)
		_ = handle.(*os.File).Close()
		return true
	})
}

func (fp *frameProcessor) String() string {
	return "Recorder Frame Processor"
}

func (fp *frameProcessor) Listen() chan tracker.Event {
	return fp.events
}
func (fp *frameProcessor) Handle(frame tracker.Frame, frameSource *tracker.FrameSource) tracker.Frame {
	if nil == frame {
		return nil
	}
	var fileName string
	switch (frame).(type) {
	case *beast.Frame:
		fileName = "frames.beast"
	case *mode_s.Frame:
		fileName = "frames.avr"
	case *sbs1.Frame:
		fileName = "frames.sbs1"
	default:
	}

	fh, ok := fp.logFileHandles.Load(fileName)
	if !ok {
		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if nil != err {
			log.Error().Err(err).Bytes("frame", frame.Raw()).Msg("Failed to log frame")
			return frame
		}
		fp.logFileHandles.Store(fileName, f)
		fh = f
	}

	if _, err := fh.(*os.File).Write(frame.Raw()); nil != err {
		log.Error().Err(err).Msg("Failed to write frame to log file")
	}

	return frame
}
