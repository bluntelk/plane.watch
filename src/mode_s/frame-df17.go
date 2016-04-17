package mode_s

// extended squitter decoding

import (
	"math"
	"fmt"
)

const (
	DF17_FRAME_ID_CAT = "Aircraft Identification and Category"
	DF17_FRAME_SURFACE_POS = "Surface Position"
	DF17_FRAME_AIR_POS_BARO = "Airborne Position (with Barometric Altitude)"
	DF17_FRAME_AIR_VELOCITY = "Airborne Velocity"
	DF17_FRAME_AIR_POS_GNSS = "Airborne Position (with GNSS Height)"
	DF17_FRAME_TEST_MSG = "Test Message"
	DF17_FRAME_TEST_MSG_SQUAWK = "Test Message with Squawk"
	DF17_FRAME_SURFACE_SYS_STATUS = "Surface System Status"
	DF17_FRAME_EXT_SQUIT_EMERG = "Extended Squitter Aircraft Status (Emergency)"
	DF17_FRAME_EXT_SQUIT_STATUS = "Extended Squitter Aircraft Status (1090ES TCAS RA)"
	DF17_FRAME_STATE_STATUS = "Target State and Status Message"
	DF17_FRAME_AIRCRAFT_OPER = "Aircraft Operational Status Message"
)

func (df *Frame) MessageTypeString() string {
	var name string = "Unknown"
	if df.messageType >= 1 && df.messageType <= 4 {
		name = DF17_FRAME_ID_CAT
	} else if df.messageType >= 5 && df.messageType <= 8 {
		name = DF17_FRAME_SURFACE_POS
	} else if df.messageType >= 9 && df.messageType <= 18 {
		name = DF17_FRAME_AIR_POS_BARO
	} else if df.messageType == 19 && df.messageSubType >= 1 && df.messageSubType <= 4 {
		name = DF17_FRAME_AIR_VELOCITY
	} else if df.messageType >= 20 && df.messageType <= 22 {
		name = DF17_FRAME_AIR_POS_GNSS
	} else if df.messageType == 23 {
		if df.messageSubType == 7 {
			name = DF17_FRAME_TEST_MSG_SQUAWK
		} else {
			name = DF17_FRAME_TEST_MSG
		}
	} else if df.messageType == 24 && df.messageSubType == 1 {
		name = DF17_FRAME_SURFACE_SYS_STATUS
	} else if df.messageType == 28 && df.messageSubType == 1 {
		name = DF17_FRAME_EXT_SQUIT_EMERG
	} else if df.messageType == 28 && df.messageSubType == 2 {
		name = DF17_FRAME_EXT_SQUIT_STATUS
	} else if df.messageType == 29 {
		if (df.messageSubType == 0 || df.messageSubType == 1) {
			name = DF17_FRAME_STATE_STATUS
		} else {
			name = fmt.Sprintf("%s (Unknown Sub Message %d)", DF17_FRAME_STATE_STATUS, df.messageSubType);
		}
	} else if df.messageType == 31 && (df.messageSubType == 0 || df.messageSubType == 1) {
		name = DF17_FRAME_AIRCRAFT_OPER
	}
	return name
}

func (f *Frame) decodeDF17LatLon() {
	var msg6 = int(f.message[6])
	var msg7 = int(f.message[7])
	var msg8 = int(f.message[8])
	var msg9 = int(f.message[9])
	var msg10 = int(f.message[10])

	// CPR LAT/LON
	f.rawLatitude = ((msg6 & 3) << 15) | (msg7 << 7) | (msg8 >> 1)
	f.rawLongitude = ((msg8 & 1) << 16) | (msg9 << 8) | msg10
	f.cprFlagOddEven = int(msg6 & 4) >> 2
}

func (f *Frame) decodeDF17() {

	// Down Link Format 17 Message Types
	f.messageType = f.message[4] >> 3
	f.messageSubType = f.message[4] & 7

	if f.messageType >= 1 && f.messageType <= 4 {
		/* Aircraft Identification and Category */
		f.decodeFlightNumber()
	} else if f.messageType >= 5 && f.messageType <= 8 {
		// surface position
		f.onGround = true
		f.validVerticalStatus = true

		f.decodeDF17LatLon()
		f.decodeMovementField()

		if f.message[5] & 0x08 != 0 {
			f.heading = float64(((((f.message[5] << 4) | (f.message[6] >> 4)) & 0x007F) * 45) >> 4)
		}

	} else if (f.messageType >= 9 && f.messageType <= 18) {
		/* Airborne position Message */
		f.timeFlag = int(f.message[6] & (1 << 3))
		f.onGround = false
		f.validVerticalStatus = true
		f.decodeAC12AltitudeField() // decode altitude and unit
		f.decodeDF17LatLon()
	} else if f.messageType == 19 && f.messageSubType >= 1 && f.messageSubType <= 4 {
		/* Airborne Velocity Message */
		f.onGround = false
		f.validVerticalStatus = true
		if f.messageSubType >= 1 && f.messageSubType <= 4 {
			var verticalRateSign int = int((f.message[8] & 0x8) >> 3)
			f.verticalRate = int(((f.message[8] & 7) << 6) | ((f.message[9] & 0xfc) >> 2))
			if f.verticalRate != 0 {
				f.verticalRate--
				if verticalRateSign != 0 {
					f.verticalRate = 0 - f.verticalRate
				}
				f.verticalRate = f.verticalRate * 64
			}

		}
		if f.messageSubType == 1 || f.messageSubType == 2 {
			f.eastWestDirection = int((f.message[5] & 4) >> 2)
			f.eastWestVelocity = int(((f.message[5] & 3) << 8) | f.message[6])
			f.northSouthDirection = int((f.message[7] & 0x80) >> 7)
			f.northSouthVelocity = int(((f.message[7] & 0x7f) << 3) | ((f.message[8] & 0xe0) >> 5))
			f.verticalRateSource = int((f.message[8] & 0x10) >> 4)
			/* Compute velocity and angle from the two speed components. */

			if f.messageSubType == 2 {
				// supersonic - unit is 4 knots
				f.eastWestVelocity = f.eastWestVelocity << 2
				f.northSouthVelocity = f.northSouthVelocity << 2
				f.superSonic = true
			}

			f.velocity = math.Sqrt(float64((f.northSouthVelocity * f.northSouthVelocity) + (f.eastWestVelocity * f.eastWestVelocity)))
			if f.velocity != 0 {
				var ewv float64 = float64(f.eastWestVelocity)
				var nsv float64 = float64(f.northSouthVelocity)
				var heading float64
				if f.eastWestDirection != 0 {
					ewv *= -1
				}
				if f.northSouthDirection != 0 {
					nsv *= -1
				}
				heading = math.Atan2(ewv, nsv)
				/* Convert to degrees. */
				f.heading = heading * 360 / (math.Pi * 2)
				/* We don't want negative values but a 0-360 scale. */
				if f.heading < 0 {
					f.heading += 360
				}
			} else {
				f.heading = 0
			}
		} else if f.messageSubType == 3 || f.messageSubType == 4 {
			var airspeed int = int(((f.message[7] & 0x7f) << 3) | (f.message[8] >> 5));
			if airspeed != 0 {
				airspeed--;
				if f.messageSubType == 4 {
					// If (supersonic) unit is 4 kts
					f.superSonic = true
					airspeed = airspeed << 2;
				}
				f.velocity = float64(airspeed);
			}

			if f.message[5] & 4 != 0 {
				f.heading = (360.0 / 128.0) * float64(((f.message[5] & 3) << 5) | (f.message[6] >> 3))
			}
		}
	} else if f.messageType == 23 && f.messageSubType == 7 {
		// TEST MESSAGE with  squawk - decode it!
		f.decodeSquawkIdentity(5, 6)
	} else if f.messageType == 28 && f.messageSubType == 1 {
		// EMERGENCY, EMERGENCY, THERE'S AN EMERGENCY GOING ON
		f.decodeSquawkIdentity(5, 6)
		var emergencyId int = int((f.message[5] & 0xE0) >> 5)
		f.alert = emergencyId != 0
		f.special = emergencyStateTable[emergencyId]
	}
}

func (df *Frame) MessageType() byte {
	return df.messageType
}

// Whether or not this frame is even or odd, for CPR Location
func (df *Frame) IsEven() bool {
	return df.cprFlagOddEven == 0
}

func (df *Frame) FlightNumber() string {
	return string(df.flight)
}

//
//=========================================================================
//
// Decode the 7 bit ground movement field PWL exponential style scale
//
func (f *Frame) decodeMovementField() {
	var gSpeed uint64
	movement := uint64(((f.message[4] << 4) | (f.message[5] >> 4)) & 0x007F)
	if (movement > 0) && (movement < 125) {

		if movement > 123 {
			gSpeed = 199 // > 175kt
		} else if movement > 108 {
			gSpeed = ((movement - 108) * 5) + 100
		} else if movement > 93 {
			gSpeed = ((movement - 93) * 2) + 70
		} else if movement > 38 {
			gSpeed = (movement - 38) + 15
		} else if movement > 12 {
			gSpeed = ((movement - 11) >> 1) + 2
		} else if movement > 8 {
			gSpeed = ((movement - 6) >> 2) + 1
		} else {
			gSpeed = 0
		}
		f.velocity =float64(gSpeed)
	}
}
