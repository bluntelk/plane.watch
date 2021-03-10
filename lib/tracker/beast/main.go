package beast

import (
	"plane.watch/lib/tracker/mode_s"
	"time"
)

type (
)

func NewFrame(rawBytes []byte, t time.Time) *mode_s.Frame {

	// decode beast into AVR

	decodeFrame := ""
	return mode_s.NewFrame(decodeFrame, t)
}
