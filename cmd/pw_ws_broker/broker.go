package main

import (
	"errors"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"plane.watch/lib/rabbitmq"
	"plane.watch/lib/randstr"
	"syscall"
	"time"
)

type (
	PwWsBroker struct {
		rabbit      *rabbitmq.RabbitMQ
		queuePrefix string
		queueLow    string
		queueHigh   string
	}
)

func NewPlaneWatchWebSocketBroker(rabbitUrl string) (*PwWsBroker, error) {
	rabbitCfg, err := rabbitmq.NewConfigFromUrl(rabbitUrl)
	if nil != err {
		return nil, err
	}
	prefix := "broker_" + randstr.RandString(10) + "_"
	return &PwWsBroker{
		rabbit:      rabbitmq.New(rabbitCfg),
		queuePrefix: prefix,
		queueLow:    prefix + "low",
		queueHigh:   prefix + "high",
	}, nil
}

func (b *PwWsBroker) Start(routeLow, routeHigh string) error {
	if nil == b.rabbit {
		return errors.New("you need to configure the rabbit client")
	}
	if err := b.rabbit.ConnectAndWait(5 * time.Second); nil != err {
		return err
	}

	if _, err := b.rabbit.QueueDeclare(b.queuePrefix+"high", 60000); nil != err {
		return err
	}

	if _, err := b.rabbit.QueueDeclare(b.queuePrefix+"low", 60000); nil != err {
		return err
	}

	// bind routing keys to our queues
	if err := b.rabbit.QueueBind(b.queueLow, routeLow, rabbitmq.PlaneWatchExchange); nil != err {
		return err
	}

	if err := b.rabbit.QueueBind(b.queueHigh, routeHigh, rabbitmq.PlaneWatchExchange); nil != err {
		return err
	}
	return nil
}

func (b *PwWsBroker) Run() {
	log.Debug().Msg("RUN")

}

func (b *PwWsBroker) Wait() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func (b *PwWsBroker) Close() {
	if err := b.rabbit.QueueRemove(b.queueLow); nil != err {
		log.Error().Err(err).Str("Queue", b.queueLow).Msg("Removing Queue")
	}
	if err := b.rabbit.QueueRemove(b.queueHigh); nil != err {
		log.Error().Err(err).Str("Queue", b.queueHigh).Msg("Removing Queue")
	}

	if nil != b.rabbit {
		b.rabbit.Disconnect()
	}
	return
}
