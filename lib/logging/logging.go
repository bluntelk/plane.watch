package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"time"
)

var isCli bool

func IncludeDebugQuiet(app *cli.App) {
	app.Flags = append(app.Flags,
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
	)
}

func SetVerboseOrQuiet(verbose, quiet bool) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if quiet {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}
}

func cliWriter() zerolog.ConsoleWriter {
	return zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.UnixDate}
}

func ConfigureForCli() {
	isCli = true
	log.Logger = log.Output(cliWriter())
}

func AddLogDestination(newLogger io.Writer) {
	var multi zerolog.LevelWriter
	if isCli {
		multi = zerolog.MultiLevelWriter(cliWriter(), newLogger)
	} else {
		multi = zerolog.MultiLevelWriter(log.Logger, newLogger)
	}
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()
}
