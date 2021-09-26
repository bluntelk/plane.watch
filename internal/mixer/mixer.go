package mixer

import (
	"github.com/urfave/cli/v2"
)

type (
	// incoming can be both push/pull with auth
	incoming struct {
		host, port string
	}

	// default sink -> plane.watch
	outgoing struct {
	}

	mixer struct {
		inputs  []incoming
		outputs []outgoing
	}
)

func Run(c *cli.Context) error {
	return nil
}
