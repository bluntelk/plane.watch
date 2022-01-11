package main

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"plane.watch/lib/export"
	"plane.watch/lib/rabbitmq"
	"plane.watch/lib/randstr"
	"syscall"
)

type (
	PwWsBroker struct {
		PwWsBrokerRabbit
		PwWsBrokerWeb
	}
)

func NewPlaneWatchWebSocketBroker(rabbitUrl, routeLow, routeHigh, httpAddr string, serveTestWeb bool) (*PwWsBroker, error) {
	rabbitCfg, err := rabbitmq.NewConfigFromUrl(rabbitUrl)
	if nil != err {
		return nil, err
	}
	prefix := "broker_" + randstr.RandString(10) + "_"
	return &PwWsBroker{
		PwWsBrokerRabbit{
			rabbit:      rabbitmq.New(rabbitCfg),
			queuePrefix: prefix,
			queueLow:    prefix + "low",
			queueHigh:   prefix + "high",
			routeLow:    routeLow,
			routeHigh:   routeHigh,
		},
		PwWsBrokerWeb{
			Addr:      httpAddr,
			ServeTest: serveTestWeb,
		},
	}, nil
}

func (b *PwWsBroker) Setup() error {
	if err := b.configureRabbitMq(); nil != err {
		return err
	}
	if err := b.configureWeb(); nil != err {
		return err
	}

	b.processMessage = func(highLow string, loc *export.EnrichedPlaneLocation) {
		tile := loc.PlaneLocation.TileLocation + highLow
		b.clients.SendLocationUpdate(highLow, tile, loc)
	}
	return nil
}

func (b *PwWsBroker) Run() {
	go b.consumeAll()

	log.Info().Str("HttpAddr", b.httpServer.Addr).Msg("HTTP Listening on")
	if err := b.httpServer.ListenAndServe(); nil != err {
		log.Error().Err(err).Send()
	}
}

func (b *PwWsBroker) Wait() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func (b *PwWsBroker) Close() {
	if err := b.httpServer.Close(); nil != err {
		log.Error().Err(err).Msg("Failed to close web server cleanly")
	}
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