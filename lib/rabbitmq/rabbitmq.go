package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"time"
)

// back off logic
const rabbitmqRetryInterval = 2
const rabbitmqRetryIntervalMax = 120
const PlaneWatchExchange = "plane.watch.data"

type (
	rabbitMqConnection struct {
		conn         *amqp.Connection
		channel      *amqp.Channel
		disconnected chan *amqp.Error
		connected    bool
		log          zerolog.Logger
	}

	RabbitMQ struct {
		uri           string
		receive, send rabbitMqConnection
		log           zerolog.Logger
	}
)

var (
	ErrNilChannel = errors.New("trying to use a nil channel")
)

func New(cfg *Config) *RabbitMQ {
	if nil == cfg {
		return nil
	}
	rmq := &RabbitMQ{}
	rmq.SetUrl(cfg.String())
	rmq.Setup()

	return rmq
}

func (r *RabbitMQ) Setup() {
	r.log = log.With().Str("Service", "RabbitMQ").Logger()
	r.send.log = r.log.With().Str("conn", "send").Logger()
	r.receive.log = r.log.With().Str("conn", "receive").Logger()

}

func (r *RabbitMQ) SetUrl(serverUrl string) {
	r.uri = serverUrl
}

// ConnectAndWait connects to our RabbitMQ server and sets up our connections.
// it waits until the timeout for the connection to come up and returns an error if anything went wrong
func (r *RabbitMQ) ConnectAndWait(ctx context.Context) error {
	t, ok := ctx.Deadline()
	r.log.Info().Str("Url", r.uri).Time("Timeout", t).Bool("reached", ok).Msg("Connecting")
	connectedSend := make(chan bool)
	connectedReceive := make(chan bool)
	var isConnectedSend, isConnectedReceive bool
	go r.send.handleConnect(r.uri, connectedSend)
	go r.receive.handleConnect(r.uri, connectedReceive)
	for {
		select {
		case <-connectedSend:
			r.log.Info().Msg("Connected")
			isConnectedSend = true
			if isConnectedSend && isConnectedReceive {
				return nil
			}
		case <-connectedReceive:
			isConnectedReceive = true
			if isConnectedSend && isConnectedReceive {
				return nil
			}
		case <-ctx.Done():
			r.log.Error().Err(ctx.Err()).Msg("Connect Timeout")
			return fmt.Errorf("failed to connect to rabbit in a timely manner")
		}
	}
}

func (r *RabbitMQ) Disconnect() {
	r.log.Info().Msg("Disconnecting")

	if r.send.connected {
		_ = r.send.conn.Close()
	}
	r.send.connected = false

	if r.receive.connected {
		_ = r.receive.conn.Close()
	}
	r.receive.connected = false
}

func (r *RabbitMQ) ExchangeDeclare(name, kind string) error {
	r.log.Debug().Str("Exchange", name).Str("type", kind).Msg("Declaring Exchange")
	if r.send.isHealthy() {
		return r.send.channel.ExchangeDeclare(
			name,
			kind,
			true,  // All exchanges are not declared durable
			false, // auto-deleted
			false, // internal
			false, // no-wait
			nil,   // arguments
		)
	}

	r.log.Warn().Err(ErrNilChannel).Str("what", "ExchangeDeclare").Msg("We do not have a channel, cannot proceed")
	return ErrNilChannel
}

// QueueDeclare makes sure we have our queue setup with a default message expiry
// ttlMs is the number of seconds we will wait before expiring a message, in MilliSeconds
func (r *RabbitMQ) QueueDeclare(name string, ttlMs int) (amqp.Queue, error) {
	r.log.Debug().Str("Queue", name).Int("TTL (ms)", ttlMs).Msg("Declaring Queue")
	if r.receive.isHealthy() {
		return r.receive.channel.QueueDeclare(
			name,
			false,
			true,
			false,
			false,
			map[string]interface{}{
				"x-message-ttl": ttlMs,
			},
		)
	}
	r.log.Warn().Err(ErrNilChannel).Str("what", "QueueDeclare").Msg("We do not have a channel, cannot proceed")
	return amqp.Queue{}, ErrNilChannel
}

func (r *RabbitMQ) QueueBind(name, routingKey, sourceExchange string) error {
	r.log.Debug().
		Str("Queue", name).
		Str("Routing Key", routingKey).
		Str("Exchange", sourceExchange).
		Msg("Binding Queue")
	if r.receive.isHealthy() {
		return r.receive.channel.QueueBind(
			name,
			routingKey,
			sourceExchange,
			true,
			nil,
		)
	}
	r.log.Warn().Err(ErrNilChannel).Str("what", "QueueBind").Msg("We do not have a channel, cannot proceed")
	return ErrNilChannel
}

func (r *RabbitMQ) QueueRemove(name string) error {
	r.log.Debug().
		Str("Queue", name).
		Msg("Removing Queue")
	if r.receive.isHealthy() {
		n, err := r.receive.channel.QueueDelete(
			name,
			false,
			false,
			false,
		)
		r.log.Debug().Int("Num Messages Lost", n).Msg("Queue Removed")
		return err
	}
	r.log.Warn().Err(ErrNilChannel).Str("what", "QueueRemove").Msg("We do not have a channel, cannot proceed")
	return ErrNilChannel
}

func (r *RabbitMQ) Consume(name, consumer string) (<-chan amqp.Delivery, error) {
	r.log.Debug().
		Str("Queue", name).
		Str("Consumer", consumer).
		Msg("Consuming Queue Messages")
	if r.receive.isHealthy() {
		return r.receive.channel.Consume(
			name,
			consumer,
			true,
			false,
			false,
			false,
			nil,
		)
	}
	r.log.Warn().Err(ErrNilChannel).Str("what", "Consume").Msg("We do not have a channel, cannot proceed")
	return nil, ErrNilChannel
}

func (r *RabbitMQ) Publish(exchange, key string, msg amqp.Publishing) error {
	if nil != r.send.channel {
		return r.send.channel.Publish(
			exchange,
			key,
			false,
			false,
			msg,
		)
	}
	r.log.Warn().Err(ErrNilChannel).Str("what", "Publish").Msg("We do not have a channel, cannot proceed")
	return ErrNilChannel
}

func (r *RabbitMQ) HealthCheck() bool {
	log.Debug().Msg("RabbitMQ Health Check")
	if nil == r || !r.send.isHealthy() || !r.receive.isHealthy() {
		log.Error().Msg("Rabbit Healthcheck Bad")
		return false
	}
	return true
}

func (r *RabbitMQ) HealthCheckName() string {
	return "RabbitMQ"
}

func (rc *rabbitMqConnection) handleConnect(uri string, connected chan bool) {
	reset := make(chan bool)
	connectedChan := make(chan bool)
	timer := time.AfterFunc(0, func() {
		err := rc.doConnect(uri)
		if nil == err {
			connectedChan <- true
			reset <- true
		} else {
			rc.log.Error().Err(err).Msg("Could not establish connection")
		}
	})
	defer timer.Stop()

	var backoffIntervalCounter, backoffInterval int64

	for {
		select {
		case <-connectedChan:
			rc.log.Debug().Msg("RabbitMQ connected and channel established")
			rc.connected = true
			connected <- true
			backoffIntervalCounter = 0
			backoffInterval = 0
			return
		case <-reset:
			rc.connected = false
			backoffIntervalCounter++
			if 0 == backoffInterval {
				backoffInterval = rabbitmqRetryInterval
			} else {
				backoffInterval = backoffInterval * rabbitmqRetryInterval
			}

			if backoffInterval > rabbitmqRetryIntervalMax {
				backoffInterval = rabbitmqRetryIntervalMax
			}

			rc.log.Error().
				Int64("backoff (s)", backoffInterval).
				Int64("attempt", backoffIntervalCounter).
				Msgf("Failed to connect, attempt %d, Retrying in %d seconds", backoffIntervalCounter, backoffInterval)

			timer.Reset(time.Duration(backoffInterval) * time.Second)
		}
	}
}

func (rc *rabbitMqConnection) handleDisconnect() {
	for msg := range rc.disconnected {
		rc.log.Info().
			Str("connection", "disconnected").
			Str("error", msg.Error()).
			Bool("server", msg.Server).
			Str("reason", msg.Reason).
			Int("code", msg.Code).
			Bool("recover", msg.Recover).
			Send()
		rc.connected = false
		rc.channel = nil
	}
}

func (rc *rabbitMqConnection) doConnect(uri string) error {
	var err error

	rc.log.Debug().Str("Url", uri).Msg("Dialing")
	rc.conn, err = amqp.Dial(uri)
	if err != nil {
		rc.log.Error().Err(err).Msg("Cannot Connect")
		return err
	}

	rc.log.Debug().Msg("Connection established, getting Channels")
	rc.channel, err = rc.conn.Channel()
	if err != nil {
		log.Error().Err(err).Str("which", "receive").Msg("Channel Error")
		return err
	}

	// Notify disconnect channel when disconnected
	rc.disconnected = make(chan *amqp.Error)
	rc.channel.NotifyClose(rc.disconnected)
	go rc.handleDisconnect()
	return nil
}

func (rc *rabbitMqConnection) isHealthy() bool {
	if nil == rc.conn || nil == rc.channel {
		return false
	}
	return rc.connected
}
