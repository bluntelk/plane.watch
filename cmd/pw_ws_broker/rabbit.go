package main

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"plane.watch/lib/export"
	"plane.watch/lib/rabbitmq"
	"time"
)

type (
	PwWsBrokerRabbit struct {
		rabbit      *rabbitmq.RabbitMQ
		queuePrefix string
		queueLow    string
		queueHigh   string
		routeLow    string
		routeHigh   string
	}
)

func (br *PwWsBrokerRabbit) configureRabbitMq() error {
	if nil == br.rabbit {
		return errors.New("you need to configure the rabbit client")
	}
	if err := br.rabbit.ConnectAndWait(5 * time.Second); nil != err {
		return err
	}

	if _, err := br.rabbit.QueueDeclare(br.queuePrefix+"high", 60000); nil != err {
		return err
	}

	if _, err := br.rabbit.QueueDeclare(br.queuePrefix+"low", 60000); nil != err {
		return err
	}

	// bind routing keys to our queues
	if err := br.rabbit.QueueBind(br.queueLow, br.routeLow, rabbitmq.PlaneWatchExchange); nil != err {
		return err
	}

	if err := br.rabbit.QueueBind(br.queueHigh, br.routeHigh, rabbitmq.PlaneWatchExchange); nil != err {
		return err
	}
	return nil
}

func (br *PwWsBrokerRabbit) processMessage(what string, msg []byte) {
	planeLoc := export.PlaneLocation{}
	err := json.Unmarshal(msg, &planeLoc)
	if nil != err {
		log.Debug().Err(err).Msg("did not understand msg")
	}

}

func (br *PwWsBrokerRabbit) consume(queue, what string) {
	ch, err := br.rabbit.Consume(queue, "pw_ws_broker"+what)
	if nil != err {
		log.Error().Err(err).Msg("Failed to consume")
		return
	}

	for msg := range ch {
		br.processMessage(what, msg.Body)
	}
}
func (br *PwWsBrokerRabbit) consumeAll() {
	go br.consume(br.queueLow, "_low")
	go br.consume(br.queueHigh, "_high")
}
