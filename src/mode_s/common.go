package mode_s

import ()

// decode an AC12 Altitude field
func decodeAC12Field(AC12Field int32) int32 {
	var q_bit int32 = AC12Field & 0x10
	var n int32

	if q_bit != 0 {
		/// N is the 11 bit integer resulting from the removal of bit Q at bit 4
		n = ((AC12Field & 0x0FE0) >> 1) | (AC12Field & 0x000F)
		// The final altitude is the resulting number multiplied by 25, minus 1000.

		return ((n * 25) - 1000)
	} else {
		// Make N a 13 bit Gillham coded altitude by inserting M=0 at bit 6
		n = ((AC12Field & 0x0FC0) << 1) | (AC12Field & 0x003F)
		n = gillhamToAltitude(n)
		if n < -12 {
			n = 0
		}
		return 100 * n
	}
}

// this code liberally lifted from: http://www.ccsinfo.com/forum/viewtopic.php?p=77544
func gillhamToAltitude(i16GillhamValue int32) int32 {
	var i32Result int32
	var i16TempResult int32

	// Convert Gillham value using gray code to binary conversion algorithm.
	i16TempResult = i16GillhamValue ^ (i16GillhamValue >> 8)
	i16TempResult ^= (i16TempResult >> 4)
	i16TempResult ^= (i16TempResult >> 2)
	i16TempResult ^= (i16TempResult >> 1)

	// Convert gray code converted binary to altitude offset.
	i16TempResult -= (((i16TempResult >> 4) * 6) + (((i16TempResult % 16) / 5) * 2))

	// Convert altitude offset to altitude.
	i32Result = (i16TempResult - 13) * 100

	return i32Result
}

//
//=========================================================================
//
// Decode the 7 bit ground movement field PWL exponential style scale
//
func decodeMovementField(movement uint64) float64 {
	var gspeed uint64

	// Note : movement codes 0,125,126,127 are all invalid, but they are
	//        trapped for before this function is called.

	if movement > 123 {
		gspeed = 199 // > 175kt
	} else if movement > 108 {
		gspeed = ((movement - 108) * 5) + 100
	} else if movement > 93 {
		gspeed = ((movement - 93) * 2) + 70
	} else if movement > 38 {
		gspeed = (movement - 38) + 15
	} else if movement > 12 {
		gspeed = ((movement - 11) >> 1) + 2
	} else if movement > 8 {
		gspeed = ((movement - 6) >> 2) + 1
	} else {
		gspeed = 0
	}

	return float64(gspeed)
}
