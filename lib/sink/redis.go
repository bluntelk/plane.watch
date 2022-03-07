package sink

import (
	"plane.watch/lib/tracker"
)

type (
	RedisSink struct {
		Config
	}
)

func NewRedisSink(opts ...Option) (tracker.Sink, error) {
	r := &RedisSink{}
	r.setupConfig(opts)

	// TODO: Connect to redis

	return NewSink(&r.Config, r), nil
}

func (r *RedisSink) PublishJson(queue string, msg []byte) error {
	panic("IMPLEMENT ME")
}
func (r *RedisSink) PublishText(queue string, msg []byte) error {
	panic("IMPLEMENT ME")
}

func (r *RedisSink) HealthCheck() bool {
	return false
}
func (r *RedisSink) Stop() {}

func (r *RedisSink) HealthCheckName() string {
	return "Redis"
}
