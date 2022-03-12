package sink

import (
	"net"
	"net/url"
	"plane.watch/lib/redismq"
	"plane.watch/lib/tracker"
)

type (
	RedisSink struct {
		Config
		rc *redismq.Server
	}
)

func NewRedisSink(opts ...Option) (tracker.Sink, error) {
	r := &RedisSink{}
	r.setupConfig(opts)

	serverUrl := url.URL{
		Scheme:  "redis", // tls for secure
		User:    url.UserPassword(r.user, r.pass),
		Host:    net.JoinHostPort(r.host, r.port),
		Path:    "",
		RawPath: "",
	}

	var err error
	r.rc, err = redismq.NewServer(serverUrl.String())

	return NewSink(&r.Config, r), err
}

func (r *RedisSink) PublishJson(queue string, msg []byte) error {
	return r.rc.Publish(queue, msg)
}

func (r *RedisSink) PublishText(queue string, msg []byte) error {
	return r.rc.Publish(queue, msg)
}

func (r *RedisSink) HealthCheck() bool {
	return r.rc != nil
}
func (r *RedisSink) Stop() {
	r.rc.Close()
}

func (r *RedisSink) HealthCheckName() string {
	return "Redis"
}
