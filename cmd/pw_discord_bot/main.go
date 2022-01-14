package main

import (
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
	"plane.watch/lib/discord_bot"
	"plane.watch/lib/discord_bot/config"
	"plane.watch/lib/logging"
	"plane.watch/lib/mapping"
)

func main() {
	app := cli.NewApp()

	app.Name = "PlaneWatch Discord Bot"

	app.Commands = cli.Commands{
		{
			Name:        "daemon",
			Description: "For prod, Logging is JSON formatted",
			Action:      runDaemon,
		},
		{
			Name:        "cli",
			Description: "Runs in your terminal with human readable output",
			Action:      runCli,
		},
	}

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "nuke-commands",
			Usage: "Use the command to re-register all bot commands",
			Value: false,
		},
		&cli.StringFlag{
			Name:  "host",
			Usage: "The default host to talk to for our websocket goodness",
			Value: "localhost:80",
		},
		&cli.BoolFlag{
			Name:  "insecure",
			Usage: "Use this if your connection is not TLS protected",
			Value: false,
		},
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   "Show Extra Debug Information",
			EnvVars: []string{"DEBUG"},
		},
		&cli.BoolFlag{
			Name:    "quiet",
			Usage:   "Only show important messages",
			EnvVars: []string{"QUIET"},
		},
	}

	app.Before = func(c *cli.Context) error {
		logging.SetVerboseOrQuiet(c.Bool("debug"), c.Bool("quiet"))
		return nil
	}

	if err := app.Run(os.Args); nil != err {
		log.Error().Err(err).Send()
	}
}

func runDaemon(c *cli.Context) error {
	return run(c)
}

func runCli(c *cli.Context) error {
	logging.ConfigureForCli()
	return run(c)
}

func run(c *cli.Context) error {

	conf := config.Load()

	pwBot, err := discord_bot.NewBot(conf.Token)
	if nil != err {
		return err
	}
	mapping.SetHereMapsApiKey(conf.HereMapsApiKey)

	return pwBot.Run(c)
}