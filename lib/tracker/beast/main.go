package beast

import (
	"fmt"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/mode_s"
	"time"
)

type (
	beastMsg struct {
		raw           []byte
		msgType       byte
		mlatTimestamp []byte
		signalLevel   byte
		body          []byte
	}
)

var magicTimestampMLAT = []byte{0xFF, 0x00, 0x4D, 0x4C, 0x41, 0x54}

func newBeastMsg(rawBytes []byte) *beastMsg {
	// decode beast into AVR
	if rawBytes[0] != 0x1A {
		// invalid frame
		return nil
	}
	if rawBytes[1] < 0x31 || rawBytes[1] > 0x34 {
		return nil
	}
	if len(rawBytes) <= 8 {
		return nil
	}
	return &beastMsg{
		raw:           rawBytes,
		msgType:       rawBytes[1],
		mlatTimestamp: rawBytes[2:8],
		signalLevel:   rawBytes[8],
		body:          rawBytes[9:],
	}
}
func NewFrame(rawBytes []byte) tracker.Frame {
	if msg := newBeastMsg(rawBytes); nil != msg {
		fmt.Println(msg)
		return msg.decode()
	}

	return nil
}

func (bm *beastMsg) decode() tracker.Frame {
	switch bm.msgType {
	case 0x31:
		// mode-ac 10 bytes (2+8)
		return bm.decodeModeAc()
	case 0x32:
		// mode-s short 15 bytes
		return bm.decodeModeSShort()
	case 0x33:
		// mode-s long 22 bytes
		return bm.decodeModeSLong()
	case 0x34:
		// signal strength 10 bytes
		return bm.decodeConfig()
	default:
		return nil
	}
}

func (bm *beastMsg) decodeModeAc() tracker.Frame {
	// TODO: Decode ModeAC
	return nil
}

func (bm *beastMsg) decodeModeSShort() tracker.Frame {
	avr := fmt.Sprintf("%X", bm.body)
	return mode_s.NewFrame(avr, time.Now())
}

func (bm *beastMsg) decodeModeSLong() tracker.Frame {
	avr := fmt.Sprintf("%X", bm.body)
	return mode_s.NewFrame(avr, time.Now())
}

func (bm *beastMsg) decodeConfig() tracker.Frame {
	// TODO: handle 0x34 style messages
	return nil
}

func (bm *beastMsg) String() string {
	msgTypeString := map[byte]string{
		0x31: "MODE_AC",
		0x32: "MODE_S_SHORT",
		0x33: "MODE_S_LONG",
		0x34: "RADARCAPE_STATUS",
	}
	return fmt.Sprintf(
		"Type: %-16s Time: %06X Signal %03d Data: %X",
		msgTypeString[bm.msgType],
		bm.mlatTimestamp,
		bm.signalLevel,
		bm.body,
	)
}

func (bm *beastMsg) isMlat() bool {
	for i, b := range magicTimestampMLAT {
		if b != bm.raw[i+2] {
			return false
		}
	}
	return true
}
