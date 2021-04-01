package beast

import (
	"bytes"
	"testing"
)

var (
	beastModeAc     = []byte{0x1A, 0x31, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	beastModeSShort = []byte{0x1a, 0x32, 0x22, 0x1b, 0x54, 0xf0, 0x81, 0x2b, 0x26, 0x5d, 0x7c, 0x49, 0xf8, 0x28, 0xe9, 0x43}
	beastModeSLong  = []byte{0x1a, 0x33, 0x22, 0x1b, 0x54, 0xac, 0xc2, 0xe9, 0x28, 0x8d, 0x7c, 0x49, 0xf8, 0x58, 0x41, 0xd2, 0x6c, 0xca, 0x39, 0x33, 0xe4, 0x1e, 0xcf}
)

func TestNewBeastMsgModeAC(t *testing.T) {
	f := newBeastMsg(beastModeAc)

	if nil == f {
		t.Error("Did not get a beast message")
		return
	}

	if 0x31 != f.msgType {
		t.Error("Incorrect msg type")
	}
}

func TestNewBeastMsgModeSShort(t *testing.T) {
	f := newBeastMsg(beastModeSShort)

	if nil == f {
		t.Error("Did not get a beast message")
		return
	}

	if !bytes.Equal(beastModeSShort, f.raw) {
		t.Error("Failed to copy the beast message correctly")
	}

	if 0x32 != f.msgType {
		t.Error("Incorrect msg type")
	}

	// check time stamp
	if 6 != len(f.mlatTimestamp) {
		t.Errorf("Incorrect timestamp len. expected 6, got %d", len(f.mlatTimestamp))
	}
	// check signal level - should be 0xBF
	if 38 != f.signalLevel {
		t.Errorf("Did not get the signal level correctly. expected 93: got %d", f.signalLevel)
	}
	// make sure we decode into a mode_s.Frame
	if 7 != len(f.body) {
		t.Errorf("Incorrect body len. expected 7, got %d", len(f.body))
	}
}

func TestNewBeastMsgModeSLong(t *testing.T) {
	f := newBeastMsg(beastModeSLong)

	if nil == f {
		t.Error("Did not get a beast message")
		return
	}

	if !bytes.Equal(beastModeSLong, f.raw) {
		t.Error("Failed to copy the beast message correctly")
	}

	if 0x33 != f.msgType {
		t.Error("Incorrect msg type")
	}

	// check time stamp
	if 6 != len(f.mlatTimestamp) {
		t.Errorf("Incorrect timestamp len. expected 6, got %d", len(f.mlatTimestamp))
	}
	// check signal level - should be 0xBF
	if 40 != f.signalLevel {
		t.Errorf("Did not get the signal level correctly. expected 93: got %d", f.signalLevel)
	}
	// make sure we decode into a mode_s.Frame
	if 14 != len(f.body) {
		t.Errorf("Incorrect body len. expected 7, got %d", len(f.body))
	}
}
