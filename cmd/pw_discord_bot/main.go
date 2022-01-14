package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"plane.watch/lib/discord_bot"
	"plane.watch/lib/discord_bot/config"
	"plane.watch/lib/mapping"
)

func main() {
	app := cli.NewApp()

	app.Name = "PlaneWatch Discord Bot"

	conf := config.Load()

	pwBot, err := discord_bot.NewBot(conf.Token)
	if nil != err {
		log.Fatalln(err)
	}
	mapping.SetHereMapsApiKey(conf.HereMapsApiKey)

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "nuke-commands",
			Usage: "Use the command to re-register all bot commands",
			Value: false,
		},
	}

	app.Action = pwBot.Run

	if err = app.Run(os.Args); nil != err {
		log.Println(err)
	}
}
