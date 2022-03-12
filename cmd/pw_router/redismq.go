package main

import (
	"github.com/rs/zerolog/log"
	"plane.watch/lib/redismq"
)

type (
	redisRouter struct {
		redismq.Server
		doneChan chan bool
	}
)

func NewRedisRouter(url string) *redisRouter {
	if "" == url {
		return nil
	}
	rr := &redisRouter{
		doneChan: make(chan bool),
	}
	rr.SetUrl(url)
	return rr
}

func (rr *redisRouter) connect() error {
	var err error
	err = rr.Connect()

	if nil != err {
		log.Error().
			Err(err).
			Str("MQ", "redis").
			Msg("Unable to determine configuration from URL")
	}
	return err
}

func (rr *redisRouter) listen(sourceRouteKey string, incomingMessages chan []byte) error {
	ch, err := rr.Subscribe(sourceRouteKey)
	if nil != err {
		return err
	}

	go func() {
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					log.Error().Msg("failed to get message from redis nicely")
					return
				}

				incomingMessages <- []byte(msg.Payload)
			case <-rr.doneChan:
				return
			}
		}
	}()
	return nil
}

func (rr *redisRouter) publish(subject string, msg []byte) error {
	err := rr.Publish(subject, msg)
	if nil != err {
		log.Warn().Err(err).Str("mq", "redis").Msg("Failed to send update")
		return err
	}
	return nil
}

func (rr *redisRouter) close() {
	rr.doneChan <- true
}
