package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"plane.watch/lib/producer"
	"plane.watch/lib/tracker"
)

func parseSbs1(c *cli.Context) error {
	opts := make([]tracker.Option,0)
	if c.Bool("verbose") {
		opts = append(opts, tracker.WithVerboseOutput())
	} else {
		opts = append(opts, tracker.WithInfoOutput())
	}
	out, err := getOutput(c)
	if nil != err {
		fmt.Println(err)
	}

	trk := tracker.NewTracker(opts...)

	trk.AddProducer(producer.NewSbs1File(getFilePaths(c)))
	trk.AddMiddleware(timeFiddler)
	trk.Wait()
	return writeResult(trk, out)
}
