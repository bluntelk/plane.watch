package example_finder

import (
	"bytes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/beast"
	"plane.watch/lib/tracker/mode_s"
	"plane.watch/lib/tracker/sbs1"
	"strconv"
)

type (
	Filter struct {
		events chan tracker.Event

		listIcaos    []uint32
		listDfType   []byte
		listDfMeType []byte

		avrOutFile string

		log zerolog.Logger
	}
	Option func(*Filter)
)

// WithDownlinkFormatType adds a type, e.g. for ADSB/DF17 - WithDownlinkFormatType(17)
func WithDownlinkFormatType(dfType byte) Option {
	return func(filter *Filter) {
		filter.listDfType = append(filter.listDfType, dfType)
	}
}

// WithDF17MessageType adds a message type to allow, e.g. for ADSB/DF17 location messages WithDF17MessageType()
func WithDF17MessageType(dfType byte) Option {
	return func(filter *Filter) {
		filter.listDfType = append(filter.listDfType, dfType)
	}
}

// WithDF17MessageTypeLocation adds all location type updates
func WithDF17MessageTypeLocation() Option {
	return func(filter *Filter) {
		filter.listDfType = append(filter.listDfType, 17)
		filter.listDfType = append(filter.listDfType, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 20, 21, 22)
	}
}

// WithPlaneIcao adds a specific plane to allow
func WithPlaneIcao(icao uint32) Option {
	return func(filter *Filter) {
		filter.listIcaos = append(filter.listIcaos, icao)
	}
}

// WithPlaneIcaoStr adds a specific plane to allow
func WithPlaneIcaoStr(icaoStr string) Option {
	return func(f *Filter) {
		icao, err := strconv.ParseUint(icaoStr, 16, 32)
		if nil != err {
			f.log.Error().Err(err).Msg("Unable to understand this ICAO")
		} else {
			f.listIcaos = append(f.listIcaos, uint32(icao))
			f.log.Info().Str("ICAO", icaoStr).Msg("With Plane")
		}
	}
}

func NewFilter(opts ...Option) *Filter {
	f := &Filter{}
	for _, opt := range opts {
		opt(f)
	}
	f.log = log.With().Str("section", "example-finder").Logger()
	println("len(f.listIcaos)==", len(f.listIcaos))
	return f
}

func (f *Filter) Listen() chan tracker.Event {
	return f.events
}

func (f *Filter) Stop() {
	close(f.events)
}

func (f *Filter) String() string {
	return "Example Finder/Filter"
}

func (f *Filter) Handle(frame tracker.Frame, src *tracker.FrameSource) tracker.Frame {
	if nil == frame {
		return nil
	}

	// if we are filtering for one or more planes, then exclude anything that is not
	if len(f.listIcaos) > 0 {
		found := false
		for _, icao := range f.listIcaos {
			if icao == frame.Icao() {
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}

	if len(f.listDfType) > 0 || len(f.listDfMeType) > 0 {

		switch (frame).(type) {
		case *beast.Frame:
			b := frame.(*beast.Frame)
			if f.IsOk(b.AvrFrame()) {
				return frame
			}
		case *mode_s.Frame:
			if f.IsOk(frame.(*mode_s.Frame)) {
				return frame
			}
		case *sbs1.Frame:
			// no SBS1 support
			return nil
		default:
			return nil
		}
		return nil
	}
	return frame
}

func (f *Filter) IsOk(avr *mode_s.Frame) bool {
	if len(f.listDfType) > 0 && !bytes.Contains(f.listDfType, []byte{avr.DownLinkType()}) {
		return false
	}
	if len(f.listDfMeType) > 0 && !bytes.Contains(f.listDfMeType, []byte{avr.MessageType()}) {
		return false
	}
	f.log.Info().Bytes("AVR", avr.Raw()).Str("icao", avr.IcaoStr()).Msgf("Found Frame")
	return true
}
