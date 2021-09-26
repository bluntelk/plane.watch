package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"plane.watch/lib/producer"
	"plane.watch/lib/tracker"
)

func parseSbs1(c *cli.Context) error {
	opts := make([]tracker.Option,0)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if c.Bool("verbose") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if c.Bool("quiet") {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}

	out, err := getOutput(c)
	if nil != err {
		fmt.Println(err)
	}

	trk := tracker.NewTracker(opts...)

	trk.AddProducer(producer.New(producer.WithType(producer.Sbs1), producer.WithFiles(getFilePaths(c))))
	trk.AddMiddleware(timeFiddler)
	trk.Wait()
	return writeResult(trk, out)
}
