package sink

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"plane.watch/lib/rabbitmq"
	"plane.watch/lib/tracker"
	"time"
)

type (
	RabbitMqSink struct {
		Config
		mq       *rabbitmq.RabbitMQ
		exchange string
	}

	rabbitFrameMsg struct {
		Type, RouteKey string
		Body           []byte
		Source         *tracker.FrameSource
	}
)

func NewRabbitMqSink(opts ...Option) (tracker.Sink, error) {
	r := &RabbitMqSink{
		exchange: rabbitmq.PlaneWatchExchange,
	}
	r.setupConfig(opts)

	for _, opt := range opts {
		opt(&r.Config)
	}

	var err error
	if r.mq, err = r.connect(time.Second * 5); nil != err {
		r.mq = nil
		return nil, err
	}

	if err = r.setup(); nil != err {
		return nil, err
	}

	return NewSink(&r.Config, r), nil
}

func WithRabbitVhost(vhost string) Option {
	return func(config *Config) {
		config.vhost = vhost
	}
}
func WithRabbitTestQueues(value bool) Option {
	return func(conf *Config) {
		conf.createTestQueues = value
	}
}

func (r *RabbitMqSink) Stop() {
	r.mq.Disconnect()
}
func (r *RabbitMqSink) PublishJson(queue string, msg []byte) error {
	return r.mq.Publish(r.exchange, queue, amqp.Publishing{
		//ContentType:     "text/plain",
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		Timestamp:       time.Now(),
		Body:            msg,
	})
}
func (r *RabbitMqSink) PublishText(queue string, msg []byte) error {
	return r.mq.Publish(r.exchange, queue, amqp.Publishing{
		ContentType:     "text/plain",
		ContentEncoding: "utf-8",
		Timestamp:       time.Now(),
		Body:            msg,
	})
}

func (r *RabbitMqSink) connect(timeout time.Duration) (*rabbitmq.RabbitMQ, error) {
	var rabbitConfig rabbitmq.Config
	rabbitConfig.Host = r.host
	rabbitConfig.Port = r.port
	rabbitConfig.User = r.user
	rabbitConfig.Password = r.pass
	rabbitConfig.Vhost = r.vhost

	log.Info().Str("host", rabbitConfig.String()).Msg("Connecting to RabbitMQ")
	rabbit := rabbitmq.New(&rabbitConfig)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return rabbit, rabbit.ConnectAndWait(ctx)
}

func (r *RabbitMqSink) setup() error {
	var err error

	// let's make sure all of our queues and exchanges are setup
	if err = r.mq.ExchangeDeclare(r.exchange, amqp.ExchangeDirect); nil != err {
		return err
	}
	if r.Config.createTestQueues {
		for t, q := range r.queue {
			log.Debug().Str("Queue", q).Msg("Creating Test Queue")
			if _, err = r.mq.QueueDeclare(q, r.messageTtlSeconds*1000); nil != err {
				return err
			}

			if err = r.mq.QueueBind(q, t, r.exchange); nil != err {
				return err
			}
		}
	}

	return nil
}

func (r *RabbitMqSink) HealthCheck() bool {
	log.Debug().Msg("RabbitmqSink Health Check")
	if nil == r.mq {
		return false
	}
	return r.mq.HealthCheck()
}

func (r *RabbitMqSink) HealthCheckName() string {
	return "Rabbit MQ Sink"
}
