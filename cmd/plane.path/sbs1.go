package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"plane.watch/lib/logging"
	"plane.watch/lib/producer"
	"plane.watch/lib/tracker"
)

func parseSbs1(c *cli.Context) error {
	opts := make([]tracker.Option, 0)
	logging.SetVerboseOrQuiet(c.Bool("verbose"), c.Bool("quiet"))

	out, err := getOutput(c)
	if nil != err {
		fmt.Println(err)
	}

	trk := tracker.NewTracker(opts...)

	trk.AddProducer(producer.New(producer.WithType(producer.Sbs1), producer.WithFiles(getFilePaths(c))))
	trk.AddMiddleware(NewTimeFiddler())
	trk.Wait()
	return writeResult(trk, out)
}
