package mode_s

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	MODES_LONG_MSG_BYTES = 14
	MODES_SHORT_MSG_BYTES = 7
	MODES_LONG_MSG_BITS = (MODES_LONG_MSG_BYTES * 8)
	MODES_SHORT_MSG_BITS = (MODES_SHORT_MSG_BYTES * 8)
)

const MODES_GENERATOR_POLY uint32 = 0xfff409

type icaoList struct {
	icao     int
	lastSeen time.Time
}

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

	encodedFrame := strings.TrimFunc(rawFrame, func(r rune) bool {
		return unicode.IsSpace(r) || ';' == r
	})

	// let's ensure that we have some correct data...
	if "" == encodedFrame {
		return frame, fmt.Errorf("Cannot Decode Empty String")
	}

	if len(encodedFrame) < 14 {
		return frame, fmt.Errorf("Frame too short to be a Mode S frame")
	}


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

	if rawFrame == "*00000000000000;" {
		return frame, fmt.Errorf("Heartbeat Received.")
	}

	frame.raw = encodedFrame[frameStart:len(encodedFrame)]
	err = frame.parseRawToMessage()
	if nil != err {
		return frame, err
	}

	frame.decodeDownlink()

	// now see if the message we got matches up with the DF format we decoded
	if int(frame.getMessageLengthBytes()) != len(frame.message) {
		return frame, fmt.Errorf("Frame has Incorrect length %d != %d", frame.getMessageLengthBytes(), len(frame.message))
	}

	err = frame.checkCrc()
	if nil != err {
		return frame, err
	}

	// decode the specific DF type
	switch frame.downLinkFormat {
	case 0: // Airborne position, baro altitude only
		frame.decodeVerticalStatus()
		frame.decode13BitAltitudeField()
		frame.decodeRInformationField()
		frame.decodeSLField()
	case 1, 2, 3: // Aircraft Identification and Category
		frame.decodeFlightStatus()
		frame.decodeFlightId()
	case 4:
		frame.decodeFlightStatus()
		frame.decode13BitAltitudeField()
	case 5: //DF_5
		frame.decodeFlightStatus()
		frame.decodeSquawkIdentity(2, 3) // gillham encoded squawk
	case 11: //DF_11
		frame.decodeICAO()
		frame.decodeCapability()
	case 16: //DF_16
		frame.decodeVerticalStatus()
		frame.decode13BitAltitudeField()
		frame.decodeRInformationField()
		frame.decodeSLField()
	case 17: //DF_17
		frame.decodeICAO()
		frame.decodeCapability()
		frame.decodeDF17()
	case 18: //DF_18
		frame.decodeCapability() // control field
		if 0 == frame.capability {
			frame.decodeICAO()
			frame.decodeDF17()
		}
	case 20: //DF_20
		frame.decodeFlightStatus()
		frame.decode13BitAltitudeField()
		frame.decodeCommB()
	case 21: //DF_21
		frame.decodeFlightStatus()
		frame.decodeSquawkIdentity(2, 3) // gillham encoded squawk
		frame.decodeCommB()
	}

	return frame, err
}

func (f *Frame) decodeDownlink() {
	// DF24 is a little different. if the first two bits of the message are set, it is a DF24 message
	if f.message[0] & 0xc0 == 0xc0 {
		f.downLinkFormat = 24
	} else {
		// get the down link format (DF) - first 5 bits
		f.downLinkFormat = f.message[0] >> 3
	}

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

	if ! (messageLen == MODES_SHORT_MSG_BYTES || messageLen == MODES_LONG_MSG_BYTES) {
		return fmt.Errorf("Frame is incorrect length. %d != 7 or 14", messageLen)
	}

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

	switch f.capability {
	case 4:
		f.validVerticalStatus = true
		f.onGround = true
	case 5:
		f.validVerticalStatus = true
		f.onGround = false
	default:
	}
}

func (f *Frame) decodeFlightStatus() {
	// first 5 bits are the downlink format
	// bits 5,6,7 are the flight status
	f.flightStatus = int(f.message[0] & 7)
	if f.flightStatus == 0 || f.flightStatus == 2 {
		f.validVerticalStatus = true
		f.onGround = false
	}
	if f.flightStatus == 1 || f.flightStatus == 3 {
		f.validVerticalStatus = true
		f.onGround = true
	}
	if f.flightStatus == 4 || f.flightStatus == 5 {
		// special pos
		f.validVerticalStatus = true
		f.onGround = false // assume in the air
		f.special = flightStatusTable[f.flightStatus]
	}
	if f.flightStatus == 2 || f.flightStatus == 3 || f.flightStatus == 4 {
		// ALERT!
		f.alert = true
	}
}

// VS == Vertical Status
func (f *Frame) decodeVerticalStatus() {
	if f.message[0] & 4 != 0 {
		f.onGround = true
	}
	f.validVerticalStatus = true
}

// bits 13,14,15 and 16 make up the RI field
func (f *Frame) decodeRInformationField() {
	f.ri = (f.message[1] & 7) << 1 | (f.message[2] & 0x80) >> 7
}
func (f *Frame) decodeSLField() {
	f.sl = (f.message[1] & 0xe0) >> 5
}

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

func (f *Frame) decodeSquawkIdentity(byte1, byte2 int) {
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
	* into a base ten number that happens to represent the four
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
	//if f.downLinkFormat & 0x10 != 0 {
	if f.downLinkFormat & 0x10 != 0 {
		if len(f.raw) == 14 {
			return MODES_SHORT_MSG_BITS
		}
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

func (f *Frame) decodeFlightNumber() {
	f.flight = make([]byte, 8)
	f.aircraftType = int(((0x0E - f.messageType) << 4) | f.messageSubType);
	f.flight[0] = aisCharset[f.message[5] >> 2]
	f.flight[1] = aisCharset[((f.message[5] & 3) << 4) | (f.message[6] >> 4)]
	f.flight[2] = aisCharset[((f.message[6] & 15) << 2) | (f.message[7] >> 6)]
	f.flight[3] = aisCharset[f.message[7] & 63]
	f.flight[4] = aisCharset[f.message[8] >> 2]
	f.flight[5] = aisCharset[((f.message[8] & 3) << 4) | (f.message[9] >> 4)]
	f.flight[6] = aisCharset[((f.message[9] & 15) << 2) | (f.message[10] >> 6)]
	f.flight[7] = aisCharset[f.message[10] & 63]
}

func (f *Frame) decodeCommB() {
	if f.message[4] == 0x20 {
		// BDS 2,0 Aircraft Identification
		f.decodeFlightNumber()
	}
}

func (f *Frame) decodeFlightId() {
	if f.message[4] == 32 && len(f.message) >= 10 {
		// Aircraft Identification
		f.decodeFlightNumber()
	}
}


//func (list *icaoAddressWhiteList)updateLastSeen(addr uint32) {
// look for our addr
//	for i, icao := range list.ica
//}
