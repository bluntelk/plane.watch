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

func IncludeVerbosityFlags(app *cli.App) {
	app.Flags = append(app.Flags,
		&cli.BoolFlag{
			Name:  "very-verbose",
			Usage: "Enable trace level debugging",
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
	)
}

func SetLoggingLevel(c *cli.Context) {
	SetVerboseOrQuiet(
		c.Bool("very-verbose"),
		c.Bool("debug"),
		c.Bool("quiet"),
	)
}

func SetVerboseOrQuiet(trace, verbose, quiet bool) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if trace {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if quiet {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}
	//log.Info().Str("log-level", zerolog.GlobalLevel().String()).Msg("Logging Set")
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
