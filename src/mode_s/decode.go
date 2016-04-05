package mode_s

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	MODES_LONG_MSG_BYTES = 14
	MODES_SHORT_MSG_BYTES = 7
	MODES_LONG_MSG_BITS = (MODES_LONG_MSG_BYTES * 8)
	MODES_SHORT_MSG_BITS = (MODES_SHORT_MSG_BYTES * 8)
)

type icaoList struct {
	icao     int
	lastSeen time.Time
}

var (
	modes_checksum_table = [112]uint32{
		0x3935ea, 0x1c9af5, 0xf1b77e, 0x78dbbf, 0xc397db, 0x9e31e9, 0xb0e2f0, 0x587178,
		0x2c38bc, 0x161c5e, 0x0b0e2f, 0xfa7d13, 0x82c48d, 0xbe9842, 0x5f4c21, 0xd05c14,
		0x682e0a, 0x341705, 0xe5f186, 0x72f8c3, 0xc68665, 0x9cb936, 0x4e5c9b, 0xd8d449,
		0x939020, 0x49c810, 0x24e408, 0x127204, 0x093902, 0x049c81, 0xfdb444, 0x7eda22,
		0x3f6d11, 0xe04c8c, 0x702646, 0x381323, 0xe3f395, 0x8e03ce, 0x4701e7, 0xdc7af7,
		0x91c77f, 0xb719bb, 0xa476d9, 0xadc168, 0x56e0b4, 0x2b705a, 0x15b82d, 0xf52612,
		0x7a9309, 0xc2b380, 0x6159c0, 0x30ace0, 0x185670, 0x0c2b38, 0x06159c, 0x030ace,
		0x018567, 0xff38b7, 0x80665f, 0xbfc92b, 0xa01e91, 0xaff54c, 0x57faa6, 0x2bfd53,
		0xea04ad, 0x8af852, 0x457c29, 0xdd4410, 0x6ea208, 0x375104, 0x1ba882, 0x0dd441,
		0xf91024, 0x7c8812, 0x3e4409, 0xe0d800, 0x706c00, 0x383600, 0x1c1b00, 0x0e0d80,
		0x0706c0, 0x038360, 0x01c1b0, 0x00e0d8, 0x00706c, 0x003836, 0x001c1b, 0xfff409,
		0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000,
		0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000,
		0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000,
	}

//	icaoAddressWhiteList = []icaoList{}
)

type ReceivedFrame struct {
	Frame string
	Time  time.Time
}

func DecodeStringWorker(jobs <-chan ReceivedFrame, results chan <- Frame, errors chan <- error) {
	for s := range jobs {
		frame, err := DecodeString(s.Frame, s.Time)
		if nil != err {
			errors <- err
		} else {
			results <- frame
		}
	}
}

func DecodeString(rawFrame string, t time.Time) (Frame, error) {
	var frame Frame
	var err error

	encodedFrame := strings.TrimRight(rawFrame, "; \n")

	// determine what type of frame we are dealing with
	if encodedFrame[0] == '@' {
		frame.mode = "MLAT"
	} else {
		frame.mode = "NORMAL"
	}

	// ensure we have a timestamp
	frameStart := 0
	if "MLAT" == frame.mode {
		frameStart = 13
		// try and use the provided timestamp
		timeSlice := encodedFrame[1:12]
		frame.SetTimeStamp(timeSlice)
	} else if "*" == encodedFrame[0:1] {
		frameStart = 1
		frame.timeStamp = t
	} else {
		frame.timeStamp = t
	}
	// let's get our frame data in order!

	// let's ensure that we have some correct data...
	if len(rawFrame) < 14 {
		return frame, fmt.Errorf("Frame too short to be a Mode S frame")
	}

	if rawFrame == "*00000000000000;" {
		return frame, fmt.Errorf("Heartbeat Received.")
	}

	frame.raw = encodedFrame[frameStart:len(encodedFrame)]
	err = frame.parseRawToMessage()

	// get the down link format (DF) - first 5 bits
	frame.downLinkFormat = frame.message[0] >> 3

	frame.decodeModeSChecksum()

	if frame.checkSum != frame.crc {
		// todo: make sure we have the right messages that we can check the crc vs checksum
		//return frame, fmt.Errorf("Checksum and CRC do not match %d != %d", frame.checkSum, frame.crc)
	}

	if nil != err {
		return frame, err
	}

	if nil != err {
		return frame, err
	}

	// decode the specific DF type
	switch frame.downLinkFormat {
	case 0: //DF_0
		frame.decodeVSBit()
		frame.decode13BitAltitudeField()
	case 4: //DF_4
		frame.decodeFlightStatus()
		frame.decode13BitAltitudeField()
	case 5: //DF_5
		frame.decodeFlightStatus()
		frame.decodeIdentity(2, 3) // gillham encoded squawk
	case 11: //DF_11
		frame.decodeICAO()
		frame.decodeCapability()
	case 16: //DF_16
		frame.decodeVSBit()
		frame.decode13BitAltitudeField()
	case 17: //DF_17
		frame.decodeICAO()
		frame.decodeCapability()
		frame.decodeDF17()
	case 18: //DF_18
		frame.decodeICAO()
		frame.decodeCapability() // control field
		frame.decodeDF17()
	case 20: //DF_20
		frame.decodeFlightStatus()
		frame.decode13BitAltitudeField()
	case 21: //DF_21
		frame.decodeFlightStatus()
		frame.decodeIdentity(2, 3) // gillham encoded squawk
	}

	return frame, err
}

func (f *Frame) SetTimeStamp(timeStamp string) {
	if "" == timeStamp {
		f.timeStamp = time.Now()
	} else if "00000000000" == string(timeStamp) {
		f.timeStamp = time.Now()
	} else {
		// MLAT timestamps from dump 1090 are dependent on when the device started ( 500ns intervals )
		//tmp, err := strconv.ParseInt(timeStamp, 16, 64)
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Printf("To Do: Convert int %d to time stamp\n", tmp)
		f.timeStamp = time.Now()
	}
}

func (f *Frame) TimeStamp() time.Time {
	return f.timeStamp
}

// call after frame.raw is set. does the preparing
func (f *Frame) parseRawToMessage() error {
	frameLen := len(f.raw)

	// cheap bitwise even number check!
	if 0 != (frameLen & 1) {
		return fmt.Errorf("Frame is an odd length (%d), cannot decode unless length is even", frameLen)
	}

	messageLen := frameLen / 2
	f.message = make([]byte, messageLen)
	// the rest of the frame is encoded in 2 char hex values

	index := 0
	for i := 0; i < len(f.raw); i += 2 {
		pair := f.raw[i : i + 2]
		myInt, err := strconv.ParseUint(pair, 16, 8)

		if err != nil {
			return err
		}
		f.message[index] = byte(myInt)
		index++
	}
	return nil
}

func (f *Frame) decodeCapability() {
	f.capability = f.message[0] & 7
}

func (f *Frame) decodeFlightStatus() {
	// first 5 bits are the downlink format
	// bits 5,6,7 are the flight status
	f.df4_5_20_21.flightStatus = int(f.message[0] & 7)
}

// VS == Vertical Status
func (f *Frame) decodeVSBit() {
	// f.message[0] & 4 // VS bit
}

func (f *Frame) decodeFlightId() {
	if (f.message[4] == 32) {
		// Aircraft Identification

		f.flightId[0] = aisCharset[f.message[5] >> 2]
		f.flightId[1] = aisCharset[((f.message[5] & 3) << 4) | (f.message[6] >> 4)]
		f.flightId[2] = aisCharset[((f.message[6] & 15) << 2) | (f.message[7] >> 6)]
		f.flightId[3] = aisCharset[f.message[7] & 63]
		f.flightId[4] = aisCharset[f.message[8] >> 2]
		f.flightId[5] = aisCharset[((f.message[8] & 3) << 4) | (f.message[9] >> 4)]
		f.flightId[6] = aisCharset[((f.message[9] & 15) << 2) | (f.message[10] >> 6)]
		f.flightId[7] = aisCharset[f.message[10] & 63]
	}

}

//func (f *Frame) decodeDF4_5_20_21() error {
//
//	// bits 8,9,10,11,12 (5 bits) are the DR flag
//	f.df4_5_20_21.dr = int(f.message[1]) >> 3 & 31
//	f.df4_5_20_21.um = int(((f.message[1] & 7) << 3) | f.message[2]>>5)
//	return nil
//}


// Determines the ICAO address from bytes 2,3 and 4
func (f *Frame) decodeICAO() {
	switch f.downLinkFormat {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 13, 14, 15, 16, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31:
		f.icao = 0
	case 11, 17, 18:
		a := uint32(f.message[1])
		b := uint32(f.message[2])
		c := uint32(f.message[3])
		f.icao = a << 16 | b << 8 | c
	}
}

func (f *Frame) decodeIdentity(byte1, byte2 int) {
	var a, b, c, d uint32
	var msg2, msg3 uint32

	msg2 = uint32(f.message[byte1])
	msg3 = uint32(f.message[byte2])

	/* In the squawk (identity) field bits are interleaved like that
	* (message bit 20 to bit 32):
	*
	* C1-A1-C2-A2-C4-A4-ZERO-B1-D1-B2-D2-B4-D4
	*
	* So every group of three bits A, B, C, D represent an integer
	* from 0 to 7.
	*
	* The actual meaning is just 4 octal numbers, but we convert it
	* into a base ten number tha happens to represent the four
	* octal numbers.
	*
	* For more info: http://en.wikipedia.org/wiki/Gillham_code */
	a = ((msg3 & 0x80) >> 5) | ((msg2 & 0x02) >> 0) | ((msg2 & 0x08) >> 3)
	b = ((msg3 & 0x02) << 1) | ((msg3 & 0x08) >> 2) | ((msg3 & 0x20) >> 5)
	c = ((msg2 & 0x01) << 2) | ((msg2 & 0x04) >> 1) | ((msg2 & 0x10) >> 4)
	d = ((msg3 & 0x01) << 2) | ((msg3 & 0x04) >> 1) | ((msg3 & 0x10) >> 4)
	f.identity = a * 1000 + b * 100 + c * 10 + d
}

// returns the AC12 Altitude Field
func (f *Frame) getAC12Field() int32 {
	return ((int32(f.message[5]) << 4) | (int32(f.message[6]) >> 4)) & 0x0FFF
}

func (f *Frame) decodeAC12AltitudeField() {
	field := f.getAC12Field()
	f.altitude = decodeAC12Field(field)
}

// bits 20-32 are the altitude
// the 1 bits are AC13 field
// 00000000 00000000 00011111 1X1Y1111 00000000
func (f *Frame) decode13BitAltitudeField() error {
	var m_bit int = int(f.message[3] & (1 << 6)) // bit 26. 0 == feet, 1 = metres
	var q_bit int = int(f.message[3] & (1 << 4)) // bit 28

	// make sure all the bits are good

	if m_bit == 0 {
		f.unit = MODES_UNIT_FEET
		if q_bit == 1 {
			// 25 ft increments
			/* N is the 11 bit integer resulting from the removal of bit Q and M */
			var msg2 int32 = int32(f.message[2])
			var msg3 int32 = int32(f.message[3])
			var n int32 = int32((msg2 & 31) << 6) | ((msg3 & 0x80) >> 2) | ((msg3 & 0x20) >> 1) | (msg3 & 15)
			/* The final altitude is due to the resulting number multiplied
			 * by 25, minus 1000. */
			f.altitude = (n * 25) - 1000
		} else {
			// altitude reported in feet, 100ft increments
			/* Annex 10 — Aeronautical Telecommunications:
			   SSR automatic pressure-altitude transmission code (pulse position assignment)
			*/
			/* If the M bit (bit 26) and the Q bit (bit 28) equal 0, the altitude shall be coded according to the
			   pattern for Mode C replies of 3.1.1.7.12.2.3. Starting with bit 20 the sequence shall be
			   C1, A1, C2, A2, C4, A4, ZERO, B1, ZERO, B2, D2, B4, D4.
			*/
			var msg2 int32 = int32(f.message[2])
			var msg3 int32 = int32(f.message[3])
			var n int32 = int32((msg2 & 31) << 6) | ((msg3 & 0x80) >> 2) | ((msg3 & 0x20) >> 1) | (msg3 & 15)

			f.altitude = gillhamToAltitude(n)
		}
	} else {
		f.unit = MODES_UNIT_METRES
		/* TODO: Implement altitude when meter unit is selected. */
		var msg2 int32 = int32(f.message[2])
		var msg3 int32 = int32(f.message[3])
		var n int32 = int32((msg2 & 31) << 6) | int32((msg3 & 0x80) >> 2) | int32(msg3 & 15)

		// bits 20,21,22,23,24,25, 27,28,29,30,31,32 are used for altitude
		f.altitude = n
	}
	return nil
}

func (f *Frame) getMessageLengthBits() uint32 {
	if f.downLinkFormat & 0x10 != 0 {
		return MODES_LONG_MSG_BITS
	} else {
		return MODES_SHORT_MSG_BITS
	}
}

func (f *Frame) getMessageLengthBytes() uint32 {
	if f.downLinkFormat & 0x10 != 0 {
		return MODES_LONG_MSG_BYTES
	} else {
		return MODES_SHORT_MSG_BYTES
	}
}

// TODO: Make checksum decoding work correctly!
func (f *Frame) decodeModeSChecksum() {
	//unsigned char *msg, int bits
	var bitmask, counter, bit, offset uint32

	bits := f.getMessageLengthBits()
	if bits == MODES_SHORT_MSG_BITS {
		offset = 56 // default 0
	}

	for j := uint32(0); j < bits; j++ {
		counter = j / uint32(8)
		bit = j % uint32(8)
		bitmask = uint32(1) << (uint32(7) - bit)

		/* If bit is set, xor with corresponding table entry. */
		if uint32(f.message[counter]) & bitmask != 0 {
			f.checkSum ^= modes_checksum_table[j + offset]
		}
	}

	i := f.getMessageLengthBytes() - 3
	a := uint32(f.message[i])
	b := uint32(f.message[i + 1])
	c := uint32(f.message[i + 2])

	f.crc = (a << 16) | (b << 8) | c

}

//func (list *icaoAddressWhiteList)updateLastSeen(addr uint32) {
// look for our addr
//	for i, icao := range list.ica
//}
