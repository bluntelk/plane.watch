package sink

import (
	"github.com/rs/zerolog/log"
	"net"
	"net/url"
	"plane.watch/lib/nats_io"
	"plane.watch/lib/tracker"
)

type (
	NatsSink struct {
		Config
		server *nats_io.Server
	}
)

func NewNatsSink(opts ...Option) (tracker.Sink, error) {
	n := &NatsSink{}
	n.setupConfig(opts)
	if err := n.connect(); nil != err {
		log.Error().Err(err).Msg("Unable to setup nats sink")
		return nil, err
	}
	return NewSink(&n.Config, n), nil
}

func (n *NatsSink) connect() error {
	var err error
	serverUrl := url.URL{
		Scheme:  "nats", // tls for secure
		User:    url.UserPassword(n.user, n.pass),
		Host:    net.JoinHostPort(n.host, n.port),
		Path:    "",
		RawPath: "",
	}
	n.server, err = nats_io.NewServer(serverUrl.String())
	return err
}

func (n *NatsSink) PublishJson(subject string, msg []byte) error {
	return n.server.Publish(subject, msg)
}
func (n *NatsSink) PublishText(subject string, msg []byte) error {
	return n.server.Publish(subject, msg)
}

func (n *NatsSink) Stop() {
	n.server.Close()
}

func (n *NatsSink) HealthCheck() bool {
	return false
}

func (n *NatsSink) HealthCheckName() string {
	return "nats.io"
}
