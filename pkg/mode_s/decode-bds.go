package mode_s

import "fmt"

type bds struct {
	major, minor byte
}
var bdsFields = map[string]string{
	"0.5": "Extended Squitter Airborne Position",
	"0.6": "Extended Squitter Surface Position",
	"0.7": "Extended Squitter Status",
	"1.0": "Data Link Capability Report",
	"2.0": "Aircraft Identity",
	"4.0": "Aircraft Intention",
	"4.3": "Aircraft Intention",
	"4.4": "Meteorological Routine Report",
	"4.5": "Meteorological Hazard Report",
	"5.0": "Track and Turn Report",
	"5.1": "Position Report Coarse",
	"5.2": "Position Report Fine",
	"5.3": "Air-referenced State Vector",
	"6.0": "Heading and Speed Report",
}

func (b *bds) DescribeBds() string {
	key := fmt.Sprintf("%x.%x", b.major, b.minor)
	s, ok := bdsFields[key]
	if !ok {
		return key + ": Unknown"
	}
	return key + ": " + s
}


func (f *Frame) decodeDF20DF21() {
	// there is a lot of guess/detective work involved in decoding BDS

	if f.message[4] == 0x10 {
		// BDS 1, 0
		f.major = 1
		f.minor = 0
	}

	if f.message[4] == 0x20 {
		// BDS 2,0 Aircraft Identification
		f.major = 1
		f.minor = 0
		f.decodeFlightNumber()
	}

	// things get a lot murkier from here on in!
	// we should attempt to decode each BDS frame, in turn.
	// if we cannot decode it as a given frame (lots of error checking) fall through to the next type

	// BDS 4,0 - BDS status bits = 1, 14, 26, 37, 48
	// BDS 4,3 - BDS status bits = 1, 13, 26. bits 43-56 are 0's
	// BDS 4,4 - BDS status bits = 5, 24, 35, 47, 50
	// BDS 4,5 - BDS status bits = 1, 4, 7, 10, 13, 16, 27, 39. 52-56 are 0's
	// BDS 5,0 - BDS status bits = 1, 12, 24, 35, 46
	// BDS 5,1 - BDS status bits = 1
	// BDS 5,2 - BDS status bits = 1
	// BDS 5,3 - BDS status bits = 1, 13, 24, 34, 47
	// BDS 6,0 - BDS status bits = 1, 13, 24, 35, 46
}