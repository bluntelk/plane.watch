package main

import (
	"errors"
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
