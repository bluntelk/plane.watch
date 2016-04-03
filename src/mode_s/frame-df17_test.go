package mode_s

import (
	"fmt"
	"testing"
	"time"
)

func TestDecodeDF17BaroAlt1(t *testing.T) {
	frame, err := DecodeString("*8d7c4a08581fa28e6038b87a2e88;", time.Now())
	if nil != err {
		t.Error(err)
		return
	}

	if 5250 != frame.Altitude() {
		t.Errorf("DF17 Baro Alt: Failed to decode Altitude field correctly 5250 != %d", frame.Altitude())
	}

	sicao := fmt.Sprintf("%x", frame.ICAOAddr())
	if "7c4a08" != sicao {
		t.Errorf("Did not correctly decode the ICAO address: %s != 7c4a08", sicao)
	}
}

// mlat format test
func TestDecodeDF17BaroAlt2(t *testing.T) {
	frame, err := DecodeString("@000000EF31C08d8960c66055972f34137e0be0a2;", time.Now())
	if nil != err {
		t.Error(err)
		return
	}

	if 16025 != frame.Altitude() {
		t.Errorf("DF17 Baro Alt: Failed to decode Altitude field correctly 16025 != %d", frame.Altitude())
	}

	sicao := fmt.Sprintf("%x", frame.ICAOAddr())
	if "8960c6" != sicao {
		t.Errorf("Did not correctly decode the ICAO address: 8960c6 != %s", sicao)
	}
}

// mlat format test
func TestDecodeDF17BaroAlt3(t *testing.T) {
	frame, err := DecodeString("@000A237DD8708d7c6c9a583fa2c5422ad9e99abb;", time.Now())
	if nil != err {
		t.Error(err)
		return
	}

	if 11650 != frame.Altitude() {
		t.Errorf("DF17 Baro Alt: Failed to decode Altitude field correctly 16025 != %d", frame.Altitude())
	}

	sicao := fmt.Sprintf("%x", frame.ICAOAddr())
	if "7c6c9a" != sicao {
		t.Errorf("Did not correctly decode the ICAO address: 7c6c9a != %s", sicao)
	}
}
