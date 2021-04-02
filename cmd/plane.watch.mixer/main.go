package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"plane.watch/internal/mixer"
)

func main() {
	app := cli.NewApp()

	app.Name = "Plane Watch Mixer"
	app.Version = "1"
	app.Usage = "Aggregate individual feeds and reduce bandwidth"

	app.Commands = []*cli.Command{
		{
			Name:   "env",
			Action: mixer.ShowConfig,
		},
		{
			Name:   "run",
			Action: mixer.Run,
		},
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config-file",
			Value:   "/etc/plane.watch/mixer.json",
			EnvVars: []string{"CONFIG_FILE"},
		},
	}

	if err := app.Run(os.Args); nil != err {
		fmt.Println(err)
		os.Exit(1)
	}
}
