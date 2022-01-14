package discord_bot

import (
	"github.com/rs/zerolog"
	"math/rand"
	"plane.watch/lib/ws_client"
	"time"
)

type (
	wsc struct {
		c     *ws_client.Client
		wsLog zerolog.Logger
	}
)

// contains interactions with plane.watch to get the location info
func (b *PwBot) handleWebsocketClient(host string, insecure bool) {
	b.c = ws_client.NewClient(host)
	if insecure {
		b.c.Secure(insecure)
	}
	backoff := 1

	for {
		b.wsLog.Info().Str("host", host).Bool("secure", !insecure).Msg("Connecting...")
		if err := b.c.Connect(); nil != err {
			b.wsLog.Error().Err(err).Int("backoff seconds", backoff).Msg("Cannot connect to websocket")
			time.Sleep(time.Duration(backoff) * time.Second)

			backoff = backoff + backoff + (rand.Intn(10) - 2)
			if backoff > 30 {
				backoff = 30
			}
			continue
		} else {
			backoff = 1
		}

		if err := b.c.Subscribe("all_high"); nil != err {
			b.wsLog.Error().Err(err).Msg("Unable to subscribe to all_high feed")
			// TODO: Handle gracefully
			panic(err)
		}

		for update := range b.c.LocationUpdates() {
			b.wsLog.Debug().Msgf("update: %+v", update)
		}
	}
}
