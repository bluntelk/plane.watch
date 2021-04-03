package beast

import (
	"fmt"
	"plane.watch/lib/tracker/mode_s"
	"time"
)

type (
	Frame struct {
		raw           []byte
		msgType       byte
		mlatTimestamp []byte
		signalLevel   byte
		body          []byte

		isRadarCape  bool
		hasDecoded   bool
		decodedModeS *mode_s.Frame
	}
)

func (f *Frame) Icao() uint32 {
	if !f.hasDecoded {
		_, _ = f.Decode()
	}
	return f.decodedModeS.Icao()
}

func (f *Frame) IcaoStr() string {
	if !f.hasDecoded {
		_, _ = f.Decode()
	}
	return f.decodedModeS.IcaoStr()
}

func (f *Frame) Decode() (bool, error) {
	f.hasDecoded = true
	return f.decodedModeS.Decode()
}

func (f *Frame) TimeStamp() time.Time {
	// todo: calculate this off the mlat timestamp
	return time.Now()
}

func (f *Frame) Raw() []byte {
	return f.raw
}

var magicTimestampMLAT = []byte{0xFF, 0x00, 0x4D, 0x4C, 0x41, 0x54}

func newBeastMsg(rawBytes []byte) *Frame {
	if len(rawBytes) <= 8 {
		return nil
	}
	// decode beast into AVR
	if rawBytes[0] != 0x1A {
		// invalid frame
		return nil
	}
	if rawBytes[1] < 0x31 || rawBytes[1] > 0x34 {
		return nil
	}
	return &Frame{
		raw:           rawBytes,
		msgType:       rawBytes[1],
		mlatTimestamp: rawBytes[2:8],
		signalLevel:   rawBytes[8],
		body:          rawBytes[9:],
	}
}

func NewFrame(rawBytes []byte, isRadarCape bool) *Frame {
	if f := newBeastMsg(rawBytes); nil != f {
		f.isRadarCape = isRadarCape
		//if (mm->signalLevel > 0)
		//        printf("RSSI: %.1f dBFS\n", 10 * log10(mm->signalLevel));
		switch f.msgType {
		case 0x31:
			// mode-ac 10 bytes (2+8)
			f.decodeModeAc()
		case 0x32:
			// mode-s short 15 bytes
			f.decodedModeS = f.decodeModeSShort()
		case 0x33:
			// mode-s long 22 bytes
			f.decodedModeS = f.decodeModeSLong()
		case 0x34:
			// signal strength 10 bytes
			f.decodeConfig()
		default:
			return nil
		}
		return f
	}

	return nil
}

func (f *Frame) decodeModeAc() {
	// TODO: Decode ModeAC
}

func (f *Frame) decodeModeSShort() *mode_s.Frame {
	return mode_s.NewFrame(f.avr(), time.Now())
}

func (f *Frame) decodeModeSLong() *mode_s.Frame {
	return mode_s.NewFrame(f.avr(), time.Now())
}

func (f *Frame) decodeConfig() {
	// TODO: Decode RadarCape Config Info
}

func (f *Frame) avr() string {
	return fmt.Sprintf("@%X%X;", f.mlatTimestamp, f.body)
}

// BeastTicksNs returns the number of nanoseconds the beast has been on for (the mlat timestamp is calculated from power on)
func (f *Frame) BeastTicksNs() time.Duration {
	var t uint64
	inc := 40
	for i := 0; i < 6; i++ {
		t = t | uint64(f.mlatTimestamp[i])<<inc
		inc -= 8
	}
	return time.Duration(t * 500)
}

func (f *Frame) String() string {
	msgTypeString := map[byte]string{
		0x31: "MODE_AC",
		0x32: "MODE_S_SHORT",
		0x33: "MODE_S_LONG",
		0x34: "RADARCAPE_STATUS",
	}
	return fmt.Sprintf(
		"Type: %-16s Time: %06X Signal %03d Data: %X",
		msgTypeString[f.msgType],
		f.mlatTimestamp,
		f.signalLevel,
		f.body,
	)
}

func (f *Frame) isMlat() bool {
	for i, b := range magicTimestampMLAT {
		if b != f.raw[i+2] {
			return false
		}
	}
	return true
}

func (f *Frame) AvrFrame() *mode_s.Frame {
	if !f.hasDecoded {
		_, _ = f.Decode()
	}
	return f.decodedModeS
}
