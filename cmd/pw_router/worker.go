package main

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"time"

	"plane.watch/lib/rabbitmq"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"plane.watch/lib/export"
)

type (
	worker struct {
		router         *pwRouter
		destRoutingKey string
		spreadUpdates  bool
	}
)

var (
	ErrUnhandledMessage = errors.New("a message wasn't handled at all")
)

const SigHeadingChange = 1.0 // at least 1.0 degrees change.

func (w *worker) isSignificant(last export.PlaneLocation, candidate export.PlaneLocation) bool {
	// check the candidate vs last, if any of the following have changed
	// - Heading, VerticalRate, Velocity, Altitude, FlightNumber, FlightStatus, OnGround, Special, Squawk

	// if any of these fields differ, indicate this update is significant
	if candidate.HasHeading && last.HasHeading && math.Abs(candidate.Heading-last.Heading) > SigHeadingChange {
		log.Debug().
			Str("aircraft", candidate.Icao).
			Float64("last", last.Heading).
			Float64("current", candidate.Heading).
			Float64("diff_value", last.Heading-candidate.Heading).
			Int64("diff_time", int64(candidate.LastMsg.Sub(last.LastMsg))).
			Msg("Significant heading change.")
		return true
	}

	if candidate.HasVelocity && last.HasVelocity && candidate.Velocity != last.Velocity {
		log.Debug().
			Str("aircraft", candidate.Icao).
			Float64("last", last.Velocity).
			Float64("current", candidate.Velocity).
			Float64("diff_value", last.Velocity-candidate.Velocity).
			Int64("diff_time", int64(candidate.LastMsg.Sub(last.LastMsg))).
			Msg("Significant velocity change.")
		return true
	}

	if candidate.HasVerticalRate && last.HasVerticalRate && candidate.VerticalRate != last.VerticalRate {
		log.Debug().
			Str("aircraft", candidate.Icao).
			Int("last", last.VerticalRate).
			Int("current", candidate.VerticalRate).
			Int("diff_value", last.VerticalRate-candidate.VerticalRate).
			Int64("diff_time", int64(candidate.LastMsg.Sub(last.LastMsg))).
			Msg("Significant vertical rate change.")
		return true
	}

	if candidate.Altitude != last.Altitude {
		log.Debug().
			Str("aircraft", candidate.Icao).
			Int("last", last.Altitude).
			Int("current", candidate.Altitude).
			Int("diff_value", last.Altitude-candidate.Altitude).
			Int64("diff_time", int64(candidate.LastMsg.Sub(last.LastMsg))).
			Msg("Significant altitude change.")
		return true
	}

	if candidate.FlightStatus != last.FlightStatus {
		log.Debug().
			Str("aircraft", candidate.Icao).
			Str("last", last.FlightStatus).
			Str("current", candidate.FlightStatus).
			Int64("diff_time", int64(candidate.LastMsg.Sub(last.LastMsg))).
			Msg("Significant FlightStatus change.")
		return true
	}

	if candidate.OnGround != last.OnGround {
		log.Debug().
			Str("aircraft", candidate.Icao).
			Bool("last", last.OnGround).
			Bool("current", candidate.OnGround).
			Int64("diff_time", int64(candidate.LastMsg.Sub(last.LastMsg))).
			Msg("Significant OnGround change.")
		return true
	}

	if candidate.Special != last.Special {
		log.Debug().
			Str("aircraft", candidate.Icao).
			Str("last", last.Special).
			Str("current", candidate.Special).
			Int64("diff_time", int64(candidate.LastMsg.Sub(last.LastMsg))).
			Msg("Significant Special change.")
		return true
	}

	if candidate.Squawk != last.Squawk {
		log.Debug().
			Str("aircraft", candidate.Icao).
			Str("last", last.Squawk).
			Str("current", candidate.Squawk).
			Int64("diff_time", int64(candidate.LastMsg.Sub(last.LastMsg))).
			Msg("Significant Squawk change.")
		return true
	}

	if candidate.TileLocation != last.TileLocation {
		log.Debug().
			Str("aircraft", candidate.Icao).
			Str("last", last.TileLocation).
			Str("current", candidate.TileLocation).
			Int64("diff_time", int64(candidate.LastMsg.Sub(last.LastMsg))).
			Msg("Significant TileLocation change.")
		return true
	}

	log.Debug().
		Str("aircraft", candidate.Icao).
		Msg("Ignoring insignificant event.")

	return false
}

func (w *worker) run(ctx context.Context, ch <-chan []byte) {
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				log.Error().Msg("Worker ending due to error.")
				return
			}

			var gErr error
			if gErr = w.handleMsg(msg); nil != gErr {
				log.Error().Err(gErr).Send()
			}
		case <-ctx.Done():
			return
		}
	}
}

func (w *worker) handleMsg(msg []byte) error {
	var err error

	// unmarshal the JSON and ensure it's valid.
	// report the error if not and skip this message.
	update := export.PlaneLocation{}
	if err = json.Unmarshal(msg, &update); nil != err {
		log.Error().Err(err).Msg("Unable to unmarshal JSON")
		updatesError.Inc()
		return err
	}

	if "" == update.Icao {
		log.Debug().Str("payload", string(msg)).Msg("empty ICAO")
		updatesError.Inc()
		return nil
	}

	// this is considered "processed" at this point as it's valid JSON
	if err == nil {
		updatesProcessed.Inc()
	}

	// lookup what we know about this plane.
	item, ok := w.router.syncSamples.Load(update.Icao)

	// if this Icao is not in the cache, it's new.
	if !ok {
		w.handleNewUpdate(update, msg)
		return nil // finish here, no significance check as we have nothing to compare.
	}

	// upstream signals that this plane has been removed / lost.
	if update.Removed {
		w.handleRemovedUpdate(update, msg)
		return nil // don't need to do anything else with this.
	}

	// is this update significant versus the previous one
	lastRecord := item.(export.PlaneLocation)
	if w.isSignificant(lastRecord, update) {
		w.handleSignificantUpdate(update, msg)
		return nil
	} else {
		w.handleInsignificantUpdate(update, msg)
		return nil
	}

	return ErrUnhandledMessage
}

func (w *worker) handleRemovedUpdate(update export.PlaneLocation, msg []byte) {
	//check if this is a removed record and purge it from the cache and emit an event
	// this ensures downstream pipeline components always know about a removed record.
	// we get the removed flag from pw_ingest - this shortcuts our cache expiry for efficiency.
	w.router.syncSamples.Delete(update.Icao)
	cacheEntries.Dec()
	cacheEvictions.Inc()

	// emit the event to both queues
	w.publishLocationUpdate(w.destRoutingKey, msg) // to the reduced full-feed queue

	if w.spreadUpdates {
		w.publishLocationUpdate(update.TileLocation+qSuffixLow, msg)  // to the low-speed tile-queue.
		w.publishLocationUpdate(update.TileLocation+qSuffixHigh, msg) // to the high-speed tile-queue.
	}
}

func (w *worker) handleSignificantUpdate(update export.PlaneLocation, msg []byte) {
	// store the new update in-place of the old one
	w.router.syncSamples.Store(update.Icao, update)
	updatesSignificant.Inc()

	// emit the new lastSignificant
	w.publishLocationUpdate(w.destRoutingKey, msg) // all low speed messages
	if w.spreadUpdates {
		w.publishLocationUpdate(update.TileLocation+qSuffixLow, msg)
		w.publishLocationUpdate(update.TileLocation+qSuffixHigh, msg)
	}
}

func (w *worker) handleNewUpdate(update export.PlaneLocation, msg []byte) {
	// store the new update
	w.router.syncSamples.Store(update.Icao, update)
	cacheEntries.Inc()

	log.Debug().
		Str("aircraft", update.Icao).
		Msg("First time seeing aircraft.")

	// always publish to the main output queue
	w.publishLocationUpdate(w.destRoutingKey, msg)

	// if spreading updates is enabled, output to spread queues
	if w.spreadUpdates {
		w.publishLocationUpdate(update.TileLocation+qSuffixLow, msg)
		w.publishLocationUpdate(update.TileLocation+qSuffixHigh, msg)
	}
}

func (w *worker) handleInsignificantUpdate(update export.PlaneLocation, msg []byte) {
	updatesInsignificant.Inc()

	if w.spreadUpdates {
		// always publish updates to the high queue.
		w.publishLocationUpdate(update.TileLocation+qSuffixHigh, msg)
	}
}

func (w *worker) publishLocationUpdate(routingKey string, msg []byte) {
	log.Trace().Str("routing-key", routingKey).Bytes("Location", msg).Msg("Publish")

	var sent bool
	if nil != w.router.rmq {
		err := w.router.rmq.Publish(rabbitmq.PlaneWatchExchange, routingKey, amqp.Publishing{
			ContentType:     "application/json",
			ContentEncoding: "utf-8",
			Timestamp:       time.Now(),
			Body:            msg,
		})
		if nil != err {
			log.Warn().Err(err).Msg("Failed to send update to rabbit")
			sent = true
			return
		}
	}

	if nil != w.router.nats {
		err := w.router.nats.Publish(routingKey, msg)
		if nil != err {
			log.Warn().Err(err).Msg("Failed to send update to nats")
			sent = true
			return
		}
	}

	if nil != w.router.redis {
		err := w.router.redis.Publish(routingKey, msg)
		if nil != err {
			log.Warn().Err(err).Msg("Failed to send update to redis")
			sent = true
			return
		}
	}

	if sent {
		log.Trace().Str("routingKey", routingKey).Msg("Sent msg")
		updatesPublished.Inc()
	}
}
