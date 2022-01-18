package sink

import "plane.watch/lib/tracker"

type (
	RedisSink struct {
		Config
		events chan tracker.Event
	}
)

func NewRedisSink(opts ...Option) *RedisSink {
	r := &RedisSink{
		events: make(chan tracker.Event),
	}
	for _, opt := range opts {
		opt(&r.Config)
	}
	return r
}

func (r *RedisSink) OnEvent(e tracker.Event) {
	panic("Implement REDIS")
}

func (r *RedisSink) Listen() chan tracker.Event {
	return r.events
}

func (r *RedisSink) Stop() {
	close(r.events)
}

func (r *RedisSink) HealthCheck() bool {
	return false
}

func (r *RedisSink) HealthCheckName() string {
	return "Redis"
}
