package sink

import "plane.watch/lib/tracker"

type (
	RedisSink struct {
		Config
	}
)

func NewRedisSink(opts ...Option) *RedisSink {
	r := &RedisSink{}
	for _, opt := range opts {
		opt(&r.Config)
	}
	return r
}

func (r *RedisSink) OnEvent(e tracker.Event) {

}
