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

	switch f.messageType {
	case 1, 2, 3, 4:
		/* Aircraft Identification and Category */
		f.decodeFlightNumber()

		f.catType = 4 - f.messageType
		f.catSubType = f.message[4] & 7
		f.catValid = true
		f.messageSubType = 0
	case 5, 6, 7, 8:
		// surface position
		f.onGround = true
		f.validVerticalStatus = true
		f.messageSubType = 0

		f.decodeAdsbLatLon()
		f.decodeMovementField()

		if f.message[5] & 0x08 != 0 {
			f.heading = float64(((((f.message[5] << 4) | (f.message[6] >> 4)) & 0x007F) * 45) >> 4)
			f.validHeading = true
		}

	case 9, 10, 11, 12, 13, 14, 15, 16, 17, 18:
		/* Airborne position Message */
		f.messageSubType = 0
		f.timeFlag = int(f.message[6] & (1 << 3))
		f.onGround = false
		f.validVerticalStatus = true
		f.surveillanceStatus = (f.message[4] & 0x06) >> 1
		f.nicSupplementB = f.message[4] & 0x01

		field := ((int32(f.message[5]) << 4) | (int32(f.message[6]) >> 4)) & 0x0FFF
		f.altitude = decodeAC12Field(field)
		f.validAltitude = f.altitude != 0
		f.decodeAdsbLatLon()

	case 19:
		/* Airborne Velocity Message */
		f.onGround = false
		f.validVerticalStatus = true

		f.intentChange = (f.message[5] & 0x80) >> 7
		f.ifrCapability = (f.message[5] & 0x40) >> 6
		f.validNacV = true
		f.nacV = (f.message[5] & 0x38) >> 3

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
			// Ground Speed Message
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
			// Air Speed -- ground speed not available
			var airspeed int = int(((f.message[7] & 0x7f) << 3) | (f.message[8] >> 5));
			if airspeed != 0 {
				airspeed -= 1;
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
	case 20, 21, 22:
	//NoOp -- Airborne Position with GNSS instead of Baro
	case 23:
		if f.messageSubType == 7 {
			// TEST MESSAGE with  squawk - decode it!
			f.decodeSquawkIdentity(5, 6)
		} else {
			//??
		}

	case 24, 25, 26, 27:
	//NoOp
	case 28:
		if f.messageSubType == 1 {
			// EMERGENCY (or priority), EMERGENCY, THERE'S AN EMERGENCY GOING ON
			f.decodeSquawkIdentity(5, 6)
			var emergencyId int = int((f.message[5] & 0xe0) >> 5)
			f.alert = emergencyId != 0
			f.special = emergencyStateTable[emergencyId]

			// can get the Mode A Address too
			//mode_a_code = (short) (msg[2]|((msg[1]&0x1F)<<8));

		} else if f.messageSubType == 2 {
			// TCAS Resolution Advisory
		}
	case 29:
	case 30:
	// NoOp
	case 31:
		// Operational Status Message
		// TODO: Finish this off - it is not in a good working state

		// bool pointer helper
		bp := func(b bool) *bool {return &b}

		if f.messageSubType == 0 {
			// on the ground!
			f.validVerticalStatus = true
			f.onGround = false

			f.compatibilityClass = int(f.message[5]) << 8 | int(f.message[6])
			if f.compatibilityClass & 0xC000 == 0 {
				f.cccHasOperationalTcas = bp((f.compatibilityClass & 0x2000) != 0)
				f.cccHasAirRefVel = bp((f.compatibilityClass & 0x200) != 0)
				f.cccHasTargetStateRpt = bp((f.compatibilityClass & 0x100) != 0)

				changeRpt := (f.compatibilityClass & 0xC0)
				f.cccHasTargetChangeRpt = bp(changeRpt == 1 || changeRpt == 2)
				f.cccHasUATReceiver = (f.compatibilityClass & 0x20) != 0
			}

		} else if f.messageSubType == 1 {
			f.validVerticalStatus = true
			f.onGround = true
			f.compatibilityClass = int(f.message[5]) << 4 | int(f.message[6] & 0xF0) >> 4
			f.airframe_width_len = f.message[6] & 0x0F

			if f.compatibilityClass & 0xC000 == 0 {
				f.cccHasLowTxPower = bp((f.compatibilityClass & 0x200) != 0)
				f.cccHasUATReceiver = (f.compatibilityClass & 0x100) != 0
				f.validNacV = true
				f.nacV = byte((f.compatibilityClass & 0xE0) >> 5)
				f.nicSupplementC = byte((f.compatibilityClass & 0x10) >> 4)
			}
		}
		if f.compatibilityClass & 0xC000 == 0 {
			f.validCompatibilityClass = true
			f.cccHas1090EsIn = (f.compatibilityClass & 0x1000) != 0
		}

		f.operationalModeCode = int(f.message[7]) << 8 | int(f.message[8])
		f.adsbVersion = (f.message[9] & 0xe0) >> 5
		f.nicSupplementA = f.message[9] & 0x10 >> 4

		f.nacP = f.message[9] & 0x0F
		f.geoVertAccuracy = f.message[10] & 0xC0 >> 6
		f.sil = f.message[10] & 0x30 >> 4
		f.nicCrossCheck = f.message[10] & 0x08 >> 3
		f.northReference = f.message[10] & 0x04 >> 2
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