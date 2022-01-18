package rabbitmq

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"net"
	"net/url"
	"strings"
	"time"
)

// back off logic
const rabbitmqRetryInterval = 2
const rabbitmqRetryIntervalMax = 120
const PlaneWatchExchange = "plane.watch.data"

type RabbitMQ struct {
	uri          string
	conn         *amqp.Connection
	channel      *amqp.Channel
	disconnected chan *amqp.Error
	connected    bool
	log          zerolog.Logger
}

type ConfigSSL struct {
	PrivateKeyFile string `json:"private_key_file"`
	CertChainFile  string `json:"cert_chain_file"`
}

type Config struct {
	Host     string    `json:"host"`
	Port     string    `json:"port"`
	Vhost    string    `json:"vhost"`
	User     string    `json:"user"`
	Password string    `json:"password"`
	Ssl      ConfigSSL `json:"ssl"`
}

var (
	ErrNilChannel = errors.New("trying to use a nil channel")
)

func NewConfigFromUrl(connectUrl string) (*Config, error) {
	rabbitUrl, err := url.Parse(connectUrl)
	if err != nil {
		return nil, err
	}

	rabbitPassword, _ := rabbitUrl.User.Password()

	rabbitConfig := Config{
		Host:     rabbitUrl.Hostname(),
		Port:     rabbitUrl.Port(),
		User:     rabbitUrl.User.Username(),
		Password: rabbitPassword,
		Vhost:    rabbitUrl.Path,
		Ssl:      ConfigSSL{},
	}

	return &rabbitConfig, nil
}

func (cfg Config) String() string {
	u := url.URL{
		Scheme: "amqp",
		Host:   net.JoinHostPort(cfg.Host, cfg.Port),
		Path:   strings.TrimLeft(fmt.Sprintf("/%s", cfg.Vhost), "/"),
		User:   url.UserPassword(cfg.User, cfg.Password),
	}
	return u.String()
}

func New(cfg *Config) *RabbitMQ {
	if nil == cfg {
		return nil
	}
	return &RabbitMQ{
		uri: cfg.String(),
		log: log.With().Str("Service", "RabbitMQ").Logger(),
	}
}

func (r *RabbitMQ) Connect(connected chan bool) {
	reset := make(chan bool)
	done := make(chan bool)
	timer := time.AfterFunc(0, func() {
		r.connect(r.uri, done)
		reset <- true
	})
	defer timer.Stop()

	var backoffIntervalCounter, backoffInterval int64

	for {
		select {
		case <-done:
			r.log.Debug().Msg("RabbitMQ connected and channel established")
			r.connected = true
			connected <- true
			backoffIntervalCounter = 0
			backoffInterval = 0
			return
		case <-reset:
			r.connected = false
			backoffIntervalCounter++
			if 0 == backoffInterval {
				backoffInterval = rabbitmqRetryInterval
			} else {
				backoffInterval = backoffInterval * rabbitmqRetryInterval
			}

			if backoffInterval > rabbitmqRetryIntervalMax {
				backoffInterval = rabbitmqRetryIntervalMax
			}

			r.log.Error().
				Int64("attempt", backoffIntervalCounter).
				Msgf("Failed to connect, attempt %d, Retrying in %d seconds", backoffIntervalCounter, backoffInterval)

			timer.Reset(time.Duration(backoffInterval) * time.Second)
		}
	}
}
func (r *RabbitMQ) ConnectAndWait(timeout time.Duration) error {
	r.log.Info().Str("Url", r.uri).Dur("Timeout MS", timeout).Msg("Connecting")
	connected := make(chan bool)
	go r.Connect(connected)
	select {
	case <-connected:
		r.log.Info().Str("Service", "RabbitMQ").Msg("Connected")
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("failed to connect to rabbit in a timely manner")
	}
}

func (r *RabbitMQ) Disconnect() {
	r.log.Info().Msg("Disconnecting")
	if r.connected {
		_ = r.conn.Close()
	}
	r.connected = false
}

func (r *RabbitMQ) Disconnected() chan *amqp.Error {
	return r.disconnected
}

func (r *RabbitMQ) ExchangeDeclare(name, kind string) error {
	r.log.Debug().Str("Exchange", name).Str("type", kind).Msg("Declaring Exchange")
	if nil != r.channel {
		return r.channel.ExchangeDeclare(
			name,
			kind,
			true,  // All exchanges are not declared durable
			false, // auto-deleted
			false, // internal
			false, // no-wait
			nil,   // arguments
		)
	}

	r.log.Warn().Err(ErrNilChannel).Str("what", "ExchangeDeclare").Send()
	return ErrNilChannel
}

// QueueDeclare makes sure we have our queue setup with a default message expiry
// ttlMs is the number of seconds we will wait before expiring a message, in MilliSeconds
func (r *RabbitMQ) QueueDeclare(name string, ttlMs int) (amqp.Queue, error) {
	r.log.Debug().Str("Queue", name).Int("TTL (ms)", ttlMs).Msg("Declaring Queue")
	if nil != r.channel {
		return r.channel.QueueDeclare(
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
	r.log.Warn().Err(ErrNilChannel).Str("what", "QueueDeclare").Send()
	return amqp.Queue{}, ErrNilChannel
}

func (r *RabbitMQ) QueueBind(name, routingKey, sourceExchange string) error {
	r.log.Debug().
		Str("Queue", name).
		Str("Routing Key", routingKey).
		Str("Exchange", sourceExchange).
		Msg("Binding Queue")
	if nil != r.channel {
		return r.channel.QueueBind(
			name,
			routingKey,
			sourceExchange,
			true,
			nil,
		)
	}
	r.log.Warn().Err(ErrNilChannel).Str("what", "QueueBind").Send()
	return ErrNilChannel
}

func (r *RabbitMQ) QueueRemove(name string) error {
	r.log.Debug().
		Str("Queue", name).
		Msg("Removing Queue")
	if nil != r.channel {
		n, err := r.channel.QueueDelete(
			name,
			false,
			false,
			false,
		)
		r.log.Debug().Int("Num Messages Lost", n).Msg("Queue Removed")
		return err
	}
	r.log.Warn().Err(ErrNilChannel).Str("what", "QueueRemove").Send()

	return ErrNilChannel
}

func (r *RabbitMQ) Consume(name, consumer string) (<-chan amqp.Delivery, error) {
	r.log.Debug().
		Str("Queue", name).
		Str("Consumer", consumer).
		Msg("Consuming Queue Messages")
	if nil != r.channel {
		return r.channel.Consume(
			name,
			consumer,
			true,
			false,
			false,
			false,
			nil,
		)
	}
	r.log.Warn().Err(ErrNilChannel).Str("what", "Consume").Send()
	return nil, ErrNilChannel
}

func (r *RabbitMQ) Publish(exchange, key string, msg amqp.Publishing) error {
	if nil != r.channel {
		return r.channel.Publish(
			exchange,
			key,
			false,
			false,
			msg,
		)
	}
	r.log.Warn().Err(ErrNilChannel).Str("what", "Publish").Send()
	return ErrNilChannel
}

func (r *RabbitMQ) connect(uri string, done chan bool) {
	var err error

	r.log.Debug().Str("Url", uri).Msg("Dialing")
	r.conn, err = amqp.Dial(uri)
	if err != nil {
		r.log.Error().Err(err).Msg("Cannot Connect")
		return
	}

	r.log.Debug().Msg("Config established, getting Channel")
	r.channel, err = r.conn.Channel()
	if err != nil {
		log.Error().Err(err).Msg("Channel Error")
		return
	}

	// Notify disconnect channel when disconnected
	r.disconnected = make(chan *amqp.Error)
	r.channel.NotifyClose(r.disconnected)

	done <- true
}

func (r *RabbitMQ) HealthCheck() bool {
	log.Debug().Msg("RabbitMQ Health Check")
	if nil == r.conn {
		return false
	}
	return !r.conn.IsClosed()
}

func (r *RabbitMQ) HealthCheckName() string {
	return "RabbitMQ"
}
