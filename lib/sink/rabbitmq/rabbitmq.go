package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"net"
	"net/url"
	"strings"
	"time"
)

// back off logic
const rabbitmqRetryInterval = 2
const rabbitmqRetryIntervalMax = 120

type RabbitMQ struct {
	uri          string
	conn         *amqp.Connection
	channel      *amqp.Channel
	disconnected chan *amqp.Error
	connected    bool
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

func (c Config) String() string {
	return createRabbitmqUri(c)
}

func New(cfg Config) *RabbitMQ {
	uri := createRabbitmqUri(cfg)
	return &RabbitMQ{uri: uri}
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
			log.Println("RabbitMQ connected and channel established")
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

			log.Printf("Failed to connect, attempt %d, Retrying in %d seconds", backoffIntervalCounter, backoffInterval)

			timer.Reset(time.Duration(backoffInterval) * time.Second)
		}
	}
}

func (r *RabbitMQ) Disconnect() {
	if r.connected {
		_ = r.conn.Close()
	}
	r.connected = false
}

func (r *RabbitMQ) Disconnected() chan *amqp.Error {
	return r.disconnected
}

func (r *RabbitMQ) ExchangeDeclare(name, kind string) error {
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

// QueueDeclare makes sure we have our queue setup with a default message expiry
// ttlMs is the number of seconds we will wait before expiring a message, in MilliSeconds
func (r *RabbitMQ) QueueDeclare(name string, ttlMs int) (amqp.Queue, error) {
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

func (r *RabbitMQ) QueueBind(name, routingKey, sourceExchange string) error {
	return r.channel.QueueBind(
		name,
		routingKey,
		sourceExchange,
		false,
		nil,
	)
}

func (r *RabbitMQ) Consume(name, consumer string) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		name,
		consumer,
		false,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQ) Publish(exchange, key string, msg amqp.Publishing) error {
	return r.channel.Publish(
		exchange,
		key,
		false,
		false,
		msg,
	)
}

func (r *RabbitMQ) connect(uri string, done chan bool) {
	var err error

	log.Printf("Dialing %q", uri)
	r.conn, err = amqp.Dial(uri)
	if err != nil {
		log.Printf("Dial: %s", err)
		return
	}

	log.Printf("Config established, getting Channel")
	r.channel, err = r.conn.Channel()
	if err != nil {
		log.Printf("Channel: %s", err)
		return
	}

	// Notify disconnect channel when disconnected
	r.disconnected = make(chan *amqp.Error)
	r.channel.NotifyClose(r.disconnected)

	done <- true
}

func createRabbitmqUri(cfg Config) string {
	u := url.URL{
		Scheme: "amqp",
		Host:   net.JoinHostPort(cfg.Host, cfg.Port),
		Path:   strings.TrimLeft(fmt.Sprintf("/%s", cfg.Vhost), "/"),
		User:   url.UserPassword(cfg.User, cfg.Password),
	}
	return u.String()
}
