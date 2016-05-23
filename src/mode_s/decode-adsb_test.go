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

	if a, _ := frame.Altitude(); a != 5250 {
		t.Errorf("DF17 Baro Alt: Failed to decode Altitude field correctly 5250 != %d", a)
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

	if a, _ := frame.Altitude(); a != 16025 {
		t.Errorf("DF17 Baro Alt: Failed to decode Altitude field correctly 16025 != %d", a)
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

	if a, _ := frame.Altitude(); a != 11650 {
		t.Errorf("DF17 Baro Alt: Failed to decode Altitude field correctly 16025 != %d", a)
	}

	sicao := fmt.Sprintf("%x", frame.ICAOAddr())
	if "7c6c9a" != sicao {
		t.Errorf("Did not correctly decode the ICAO address: 7c6c9a != %s", sicao)
	}
}

// airborne velocity
func TestDecodeDF17MT19ST1(t *testing.T) {
	frame, err := DecodeString("8D7C451C99C4182CA0A4164A8C70", time.Now())
	if nil != err {
		t.Error(err.Error())
	}

	if 17 != frame.DownLinkType() {
		t.Error("Strange, I swore that this was an ADS-B frame (type 17)")
	}

	if 5 != frame.ca {
		t.Errorf("Capability should be 5, got %d", frame.ca)
	}

	if "7C451C" != frame.ICAOString() {
		t.Errorf("Failed to decode ICAO. expected 7C451C, got %s", frame.ICAOString())
	}

	if 19 != frame.messageType || 1 != frame.messageSubType {
		t.Errorf("Expected ADS-B Frame 19, subtype 1. Got: %d:%d",frame.messageType, frame.messageSubType)
	}

	if 1 != frame.eastWestDirection {
		t.Errorf("Expected plane to be going west (1), but instead got: %d", frame.eastWestDirection)
	}
	if -23 != frame.eastWestVelocity {
		t.Errorf("Expected plane to be going west @ 23 (-23), got %d", frame.eastWestVelocity)
	}

	if 0 != frame.northSouthDirection {
		t.Errorf("Expected plane to be going north (0), but instead got: %d", frame.northSouthDirection)
	}
	if 356 != frame.northSouthVelocity {
		t.Errorf("Expected plane to be going north @ 356 (356), got %d", frame.northSouthVelocity)
	}
	if frame.superSonic {
		t.Errorf("Wow, this plane is going a lot faster than it should be! why is it thinking it is supersonic?")
	}

}