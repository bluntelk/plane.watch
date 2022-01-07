package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"plane.watch/lib/export"
)

type (
	worker struct {
		rabbit         *rabbit
		destRoutingKey string
	}
)

const SigHeadingChange = 1.0 // at least 1.0 degrees change.

func (w *worker) isSignificant(history planeLocationLast) bool {
	// check the currentUpdate vs lastUpdate, if any of the following have changed,
	// then emit an event onto the locate-updates-reduced queue.
	// - Heading, VerticalRate, Velocity, Altitude, FlightNumber, FlightStatus, OnGround, Special, Squawk
	candidate := history.candidateUpdate
	last := history.lastSignificantUpdate

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

	if candidate.FlightNumber != last.FlightNumber {
		log.Debug().
			Str("aircraft", candidate.Icao).
			Str("last", last.FlightNumber).
			Str("current", candidate.FlightNumber).
			Int64("diff_time", int64(candidate.LastMsg.Sub(last.LastMsg))).
			Msg("Significant FlightNumber change.")
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
	}

	updatesIgnored.Inc()

	log.Debug().
		Str("aircraft", candidate.Icao).
		Msg("Ignoring insignificant event.")

	return false
}

func (w *worker) logUpdate(update export.PlaneLocation) string {
	s := fmt.Sprint("(", update.Icao, ",", update.Heading, ",", update.Velocity, ",", update.VerticalRate, ",", update.Altitude, ",", update.FlightNumber, ",", update.FlightStatus, ",", update.OnGround, ",", update.Special, ",", update.Squawk, ")")
	return s
}

func (w *worker) run(ctx context.Context, ch <-chan amqp.Delivery) {
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				log.Error().Msg("Worker ending due to error.")
				return
			}

			var gErr error
			if gErr = w.handleMsg(msg.Body); nil != gErr {
				log.Error().Err(gErr).Send()

				if gErr = msg.Nack(false, false); nil != gErr {
					log.Error().Err(gErr).Send()
				}
			} else {
				if gErr = msg.Ack(false); nil != gErr {
					log.Error().Err(gErr).Send()
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (w *worker) handleMsg(msg []byte) error {
	var err error

	update := export.PlaneLocation{}
	if err = json.Unmarshal(msg, &update); nil != err {
		updatesError.Inc()
		return err
	}

	if "" == update.Icao {
		updatesError.Inc()
		return nil
	}

	updatesProcessed.Inc()

	// if this is the first time in a while we've seen this Icao
	if _, ok := w.rabbit.syncSamples.Load(update.Icao); !ok {

		w.rabbit.syncSamples.Store(update.Icao, planeLocationLast{
			lastSignificantUpdate: update,
		})

		log.Debug().
			Str("aircraft", update.Icao).
			Msg("First time seeing aircraft.")

		return nil // can't check this for significance
	} else {
		// we have existing data, check to make sure we
		record, _ := w.rabbit.syncSamples.Load(update.Icao)
		tRecord := record.(planeLocationLast)
		tRecord.candidateUpdate = update
		w.rabbit.syncSamples.Store(update.Icao, tRecord)
	}

	planeRecord, _ := w.rabbit.syncSamples.Load(update.Icao)

	if w.isSignificant(planeRecord.(planeLocationLast)) {
		updatesSignificant.Inc()

		// if it's significant, roll the values through and emit an event.
		sigRecord, _ := w.rabbit.syncSamples.Load(update.Icao)
		tSigRecord := sigRecord.(planeLocationLast)
		tSigRecord.lastSignificantUpdate = tSigRecord.candidateUpdate
		w.rabbit.syncSamples.Store(update.Icao, tSigRecord)

		if err == nil {
			// emit the new lastSignificant
			w.publishLocationUpdate(w.destRoutingKey, msg)               // all low speed messages
			w.publishLocationUpdate(update.TileLocation+qSuffixLow, msg) // tile based low speed messages
			updatesPublished.Inc()
		} else {
			log.Info().Msg("Error Marshalling update to JSON.")
		}
	}

	// send this to the tile HIGH queue
	w.publishLocationUpdate(update.TileLocation+qSuffixHigh, msg) // tile based low speed messages

	return nil
}

func (w *worker) publishLocationUpdate(routingKey string, msg []byte) {
	err := w.rabbit.rmq.Publish("plane.watch.data", routingKey, amqp.Publishing{
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		Timestamp:       time.Now(),
		Body:            msg,
	})
	if nil != err {
		log.Warn().Err(err).Msg("Failed to send update to rabbit")
	} else {
		log.Debug().Msgf("Sent msg to routing key %s", routingKey)
	}

}
