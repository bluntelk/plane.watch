package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"plane.watch/lib/rabbitmq"
	"plane.watch/lib/tile_grid"
	"time"
)

type (
	rabbitMqRouter struct {
		rabbitmq.RabbitMQ
		doneChan chan bool
	}
)

func NewRabbitMqRouter(url string) *rabbitMqRouter {
	if "" == url {
		return nil
	}
	rr := rabbitMqRouter{
		doneChan: make(chan bool),
	}
	rr.SetUrl(url)
	rr.Setup()
	return &rr
}

func (rr *rabbitMqRouter) connect() error {
	var err error
	// connect to Rabbit
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = rr.ConnectAndWait(ctx)

	if nil != err {
		log.Error().
			Err(err).
			Str("MQ", "rabbitmq").
			Msg("Unable to determine configuration from URL")
	}
	return err
}

func (rr *rabbitMqRouter) listen(subject string, incomingMessages chan []byte) error {
	var err error
	if err = rr.rabbitMqMakeQueue("reducer-in", subject); nil != err {
		return err
	}

	ch, err := rr.Consume("reducer-in", "pw-router")
	if nil != err {
		log.Info().Msg("Failed to consume reducer-in")
		return err
	}
	go func() {
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					log.Error().Msg("failed to get message from rabbit nicely")
					return
				}

				incomingMessages <- msg.Body
			case <-rr.doneChan:
				return
			}
		}
	}()

	return nil
}

func (rr *rabbitMqRouter) publish(subject string, msg []byte) error {
	err := rr.Publish(rabbitmq.PlaneWatchExchange, subject, amqp.Publishing{
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		Timestamp:       time.Now(),
		Body:            msg,
	})
	return err
}

func (rr *rabbitMqRouter) close() {
	rr.doneChan <- true
}

func (rr *rabbitMqRouter) rabbitMqSetupTestQueues() error {
	log.Info().Msg("Setting up test queues")
	// we need a _low and a _high for each tile
	suffixes := []string{qSuffixLow, qSuffixHigh}
	for _, name := range tile_grid.GridLocationNames() {
		for _, suffix := range suffixes {
			if err := rr.rabbitMqMakeQueue(name+suffix, name+suffix); nil != err {
				return err
			}
		}
	}
	return nil
}
func (rr *rabbitMqRouter) rabbitMqMakeQueue(name, bindRouteKey string) error {
	_, err := rr.QueueDeclare(name, 60000) // 60sec TTL
	if nil != err {
		log.Error().Err(err).Msgf("Failed to create queue '%s'", name)
		return err
	}

	if err = rr.QueueBind(name, bindRouteKey, rabbitmq.PlaneWatchExchange); nil != err {
		log.Error().Err(err).Msgf("Failed to QueueBind to route-key:%s to queue %s", bindRouteKey, name)
		return err
	}
	log.Debug().Str("queue", name).Str("route-key", bindRouteKey).Msg("Setup Queue")

	return nil
}
