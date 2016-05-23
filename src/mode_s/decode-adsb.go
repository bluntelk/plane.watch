package mode_s

// extended squitter decoding

import (
	"math"
)

func (f *Frame) decodeAdsbLatLon() {
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

func (f *Frame) decodeAdsb() {

	// Down Link Format 17 Message Types
	f.messageType = f.message[4] >> 3
	f.messageSubType = f.message[4] & 7

	if f.messageType >= 1 && f.messageType <= 4 {
		/* Aircraft Identification and Category */
		f.decodeFlightNumber()

		f.catType = 4 - f.messageType
		f.catSubType = f.messageSubType
		f.catValid = true

	} else if f.messageType >= 5 && f.messageType <= 8 {
		// surface position
		f.onGround = true
		f.validVerticalStatus = true

		f.decodeAdsbLatLon()
		f.decodeMovementField()

		if f.message[5] & 0x08 != 0 {
			f.heading = float64(((((f.message[5] << 4) | (f.message[6] >> 4)) & 0x007F) * 45) >> 4)
			f.validHeading = true
		}

	} else if (f.messageType >= 9 && f.messageType <= 18) {
		/* Airborne position Message */
		f.timeFlag = int(f.message[6] & (1 << 3))
		f.onGround = false
		f.validVerticalStatus = true
		f.validAltitude = true

		field := ((int32(f.message[5]) << 4) | (int32(f.message[6]) >> 4)) & 0x0FFF
		f.altitude = decodeAC12Field(field)
		f.decodeAdsbLatLon()
	} else if f.messageType == 19 {
		/* Airborne Velocity Message */
		f.onGround = false
		f.validVerticalStatus = true

		var verticalRateSign int = int((f.message[8] & 0x8) >> 3)
		f.verticalRate = (int(f.message[8] & 7) << 6) | (int(f.message[9] & 0xfc) >> 2)
		if f.verticalRate != 0 {
			f.verticalRate--
			if verticalRateSign != 0 {
				f.verticalRate = 0 - f.verticalRate
			}
			f.verticalRate = f.verticalRate * 64
			f.validVerticalRate = true
		}

		if f.messageSubType == 1 || f.messageSubType == 2 {
			f.eastWestDirection = int(f.message[5] & 4) >> 2
			f.eastWestVelocity = int(((f.message[5] & 3) << 8) | f.message[6])
			f.northSouthDirection = int((f.message[7] & 0x80) >> 7)
			f.northSouthVelocity = (int(f.message[7] & 0x7f) << 3) | (int(f.message[8] & 0xe0) >> 5)
			f.verticalRateSource = int((f.message[8] & 0x10) >> 4)
			/* Compute velocity and angle from the two speed components. */

			if f.messageSubType == 2 {
				// supersonic - unit is 4 knots
				f.eastWestVelocity = f.eastWestVelocity << 2
				f.northSouthVelocity = f.northSouthVelocity << 2
				f.superSonic = true
			}

			f.velocity = math.Sqrt(float64((f.northSouthVelocity * f.northSouthVelocity) + (f.eastWestVelocity * f.eastWestVelocity)))
			f.validVelocity = true

			if f.velocity != 0 {
				var heading float64
				f.eastWestVelocity -= 1
				f.northSouthVelocity -= 1
				if f.eastWestDirection != 0 {
					// GO WEST! (0=east, 1=west)
					f.eastWestVelocity *= -1
				}
				if f.northSouthDirection != 0 {
					// Going Down South! (0=north, 1=south)
					f.northSouthVelocity *= -1
				}
				heading = math.Atan2(float64(f.eastWestVelocity), float64(f.northSouthVelocity))
				/* Convert to degrees. */
				f.heading = heading * 360 / (math.Pi * 2)
				/* We don't want negative values but a 0-360 scale. */
				if f.heading < 0 {
					f.heading += 360
				}
				f.validHeading = true
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
				f.validVelocity = true
			}

			if f.message[5] & 4 != 0 {
				f.heading = (360.0 / 128.0) * float64(((f.message[5] & 3) << 5) | (f.message[6] >> 3))
				f.validHeading = true
			}
		}

		if (f.message[10] > 0) {
			f.validHae = true
			f.haeDirection = (f.message[10] & 0x80) >> 7;
			var multiplier int = -25;
			if f.haeDirection == 0 {
				multiplier = 25
			}
			f.haeDelta = multiplier * int((f.message[10] & 0x7f) - 1);
		}

	} else if f.messageType == 23 && f.messageSubType == 7 {
		// TEST MESSAGE with  squawk - decode it!
		f.decodeSquawkIdentity(5, 6)
	} else if f.messageType == 28 {
		if f.messageSubType == 1 {
			// EMERGENCY (or priority), EMERGENCY, THERE'S AN EMERGENCY GOING ON
			f.decodeSquawkIdentity(5, 6)
			var emergencyId int = int((f.message[5] & 0xE0) >> 5)
			f.alert = emergencyId != 0
			f.special = emergencyStateTable[emergencyId]

			// can get the Mode A Address too
			//mode_a_code = (short) (msg[2]|((msg[1]&0x1F)<<8));

		} else if f.messageSubType == 2 {
			// TCAS Resolution Advisory
		}

	} else if f.messageType == 29 {
	} else if f.messageType == 31 {
		// Operational Status Message
		if f.messageSubType == 0 {
			f.validVerticalStatus = true
			f.onGround = false
		} else if f.messageSubType == 1 {
			f.validVerticalStatus = true
			f.onGround = true
		}
	}
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
		f.velocity = float64(gSpeed)
		f.validVelocity = true
	}
}