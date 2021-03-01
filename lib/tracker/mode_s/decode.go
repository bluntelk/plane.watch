package mode_s

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	modesLongMsgBytes  = 14
	modesShortMsgBytes = 7
	modesLongMsgBits   = modesLongMsgBytes * 8
	modesShortMsgBits  = modesShortMsgBytes * 8
)

//type icaoList struct {
//	icao     int
//	lastSeen time.Time
//}

type (
	ReceivedFrame struct {
		Frame string
		Time  time.Time
	}
	AvrDecoder struct {
		beastAvrBootTime time.Time
	}
)

var beastAvrBootTime time.Time

func init() {
	beastAvrBootTime = time.Now()
}

func NewAvrDecoder() *AvrDecoder {
	return &AvrDecoder{beastAvrBootTime: time.Now()}
}

func (d *AvrDecoder) DecodeAvrString(rawFrame string, t time.Time) (*Frame, error) {
	frame, err := prepareFrame(rawFrame, t)
	if nil != err {
		return nil, err
	}
	if frame.isNoOp() {
		return nil, nil
	}
	frame.beastAvrBoot = &d.beastAvrBootTime
	err = frame.Parse()
	return frame, err
}

func DecodeStringWorker(jobs <-chan ReceivedFrame, results chan<- *Frame, errors chan<- error) {
	for s := range jobs {
		frame, err := DecodeString(s.Frame, s.Time)
		if nil != err {
			errors <- err
		} else {
			results <- frame
		}
	}
}

func DecodeString(rawFrame string, t time.Time) (*Frame, error) {
	frame, err := prepareFrame(rawFrame, t)
	if nil != err {
		return nil, err
	}
	if frame.isNoOp() {
		return nil, nil
	}
	frame.beastAvrBoot = &beastAvrBootTime
	err = frame.Parse()
	return frame, err
}

func prepareFrame(rawFrame string, t time.Time) (*Frame, error) {
	var frame Frame

	encodedFrame := strings.TrimFunc(rawFrame, func(r rune) bool {
		return unicode.IsSpace(r) || ';' == r
	})

	// let's ensure that we have some correct data...
	if "" == encodedFrame {
		return nil, errors.New("cannot decode empty string")
	}

	if len(encodedFrame) < 14 {
		return nil, fmt.Errorf("frame (%s) too short to be a Mode S frame", rawFrame)
	}

	// determine what type of frame we are dealing with
	if encodedFrame[0] == '@' {
		// Beast Timestamp+AVR format
		frame.mode = "MLAT"
	} else {
		frame.mode = "NORMAL"
	}

	frame.timeStamp = t
	// ensure we have a timestamp
	frameStart := 0
	if "MLAT" == frame.mode {
		frameStart = 13
		// try and use the provided timestamp
		frame.beastTimeStamp = encodedFrame[1:12]
	} else if "*" == encodedFrame[0:1] {
		frameStart = 1
	}
	frame.raw = encodedFrame[frameStart:]
	frame.full = encodedFrame

	return &frame, nil
}

func (f *Frame) Parse() error {
	var err error
	f.parseBeastTimeStamp()

	err = f.parseRawToMessage()
	if nil != err {
		return err
	}

	f.decodeDownLinkFormat()

	// now see if the message we got matches up with the DF format we decoded
	if int(f.getMessageLengthBytes()) != len(f.message) {
		return fmt.Errorf("f has incorrect length %d != %d", f.getMessageLengthBytes(), len(f.message))
	}

	err = f.checkCrc()
	if nil != err {
		return err
	}

	// decode the specific DF type
	switch f.downLinkFormat {
	case 0: // Airborne position, baro altitude only
		f.decodeVerticalStatus()
		f.decodeCrossLinkCapability()
		f.decodeSensitivityLevel()
		f.decodeReplyInformation()
		err = f.decode13bitAltitudeCode()
	case 4:
		f.decodeFlightStatus()
		f.decodeDownLinkRequest()
		f.decodeUtilityMessage()
		err = f.decode13bitAltitudeCode()
	case 5: //DF_5
		f.decodeFlightStatus()
		f.decodeDownLinkRequest()
		f.decodeUtilityMessage()
		f.decodeSquawkIdentity(2, 3) // gillham encoded squawk
	case 11: //DF_11
		f.decodeICAO()
		f.decodeCapability()
	case 16: //DF_16
		f.decodeVerticalStatus()
		err = f.decode13bitAltitudeCode()
		f.decodeReplyInformation()
		f.decodeSensitivityLevel()
	case 17: //DF_17
		f.decodeICAO()
		f.decodeCapability()
		f.decodeAdsb()
	case 18: //DF_18
		f.decodeCapability() // control field
		if 0 == f.ca {
			f.decodeICAO()
			f.decodeAdsb()
		}
	case 20: //DF_20
		f.decodeFlightStatus()
		err = f.decode13bitAltitudeCode()
		f.decodeCommB()
	case 21: //DF_21
		f.decodeFlightStatus()
		f.decodeSquawkIdentity(2, 3) // gillham encoded squawk
		f.decodeCommB()
	}

	return err
}

func (f *Frame) decodeDownLinkFormat() {
	// DF24 is a little different. if the first two bits of the message are set, it is a DF24 message
	if f.message[0]&0xc0 == 0xc0 {
		f.downLinkFormat = 24
	} else {
		// get the down link format (DF) - first 5 bits
		f.downLinkFormat = f.message[0] >> 3
	}

}

func (f *Frame) parseBeastTimeStamp() {
	if "" == f.beastTimeStamp || "00000000000" == f.beastTimeStamp {
		return
	}
	// MLAT timestamps from Beast AVR are dependent on when the device started ( 500ns intervals / 12mhz)
	// calculated from power on.
	// 48 bits = 2.81474976711e+14
	// max: 2,000,000 seconds
	// Wrinkle: The same 48bites are used in GPS format (from radarcape)
	//   18 bit second of day, 30bit nanosecond
	ticksSinceBoot, err := strconv.ParseInt(f.beastTimeStamp, 16, 64)
	if err != nil {
		log.Printf("Failed to decode beast avr timestamp: %s", err)
		return
	}
	if nil == f.beastAvrBoot {
		panic(errors.New("Should have set f.beastAvrBoot"))
	}

	// handle rollover
	oldestBootTime := time.Now().Add(-2000000 * time.Second)
	if oldestBootTime.After(*f.beastAvrBoot) {
		*f.beastAvrBoot = oldestBootTime
	}

	guess := f.TimeStamp().Add(time.Duration(ticksSinceBoot*500) * time.Nanosecond)
	if guess.Before(*f.beastAvrBoot) {
		*f.beastAvrBoot = guess
	}

	f.beastAvrUptime = time.Duration(ticksSinceBoot * 500)
	//log.Printf("%0.2f seconds, or %s", duration.Seconds(), duration.String())

	f.timeStamp = beastAvrBootTime.Add(time.Duration(ticksSinceBoot * 500))
}

func (f *Frame) TimeStamp() time.Time {
	return f.timeStamp
}

// call after frame.raw is set. does the preparing
func (f *Frame) parseRawToMessage() error {
	frameLen := len(f.raw)

	// cheap bitwise even number check!
	if 0 != (frameLen & 1) {
		return fmt.Errorf("frame is an odd length (%d), cannot decode unless length is even", frameLen)
	}

	messageLen := frameLen / 2

	if !(messageLen == modesShortMsgBytes || messageLen == modesLongMsgBytes) {
		return fmt.Errorf("frame is incorrect length. %d != 7 or 14", messageLen)
	}

	f.message = make([]byte, messageLen)
	// the rest of the frame is encoded in 2 char hex values

	index := 0
	for i := 0; i < len(f.raw); i += 2 {
		pair := f.raw[i : i+2]
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
	f.ca = f.message[0] & 7

	switch f.ca {
	case 4:
		f.validVerticalStatus = true
		f.onGround = true
	case 5:
		f.validVerticalStatus = true
		f.onGround = false
	default:
	}
}
func (f *Frame) decodeCrossLinkCapability() {
	f.cc = f.message[0] & 0x2 >> 1
}

func (f *Frame) decodeFlightStatus() {
	// first 5 bits are the downlink format
	// bits 5,6,7 are the flight status
	f.fs = f.message[0] & 0x7
	if f.fs == 0 || f.fs == 2 {
		f.validVerticalStatus = true
		f.onGround = false
	}
	if f.fs == 1 || f.fs == 3 {
		f.validVerticalStatus = true
		f.onGround = true
	}
	if f.fs == 4 || f.fs == 5 {
		// special pos
		f.validVerticalStatus = true
		f.onGround = false // assume in the air
		f.special = flightStatusTable[f.fs]
	}
	if f.fs == 2 || f.fs == 3 || f.fs == 4 {
		// ALERT!
		f.alert = true
	}
}

// VS == Vertical Status
func (f *Frame) decodeVerticalStatus() {
	f.vs = f.message[0] & 4 >> 2
	f.onGround = f.vs != 0
	f.validVerticalStatus = true
}

// bits 13,14,15 and 16 make up the RI field
func (f *Frame) decodeReplyInformation() {
	f.ri = (f.message[1]&7)<<1 | (f.message[2]&0x80)>>7
}
func (f *Frame) decodeSensitivityLevel() {
	f.sl = (f.message[1] & 0xe0) >> 5
}

func (f *Frame) decodeDownLinkRequest() {
	f.dr = (f.message[1] & 0xf8) >> 3
}

func (f *Frame) decodeUtilityMessage() {
	f.um = (f.message[1]&0x7)<<3 | (f.message[2]&0xe0)>>5
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
		f.icao = a<<16 | b<<8 | c
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
	f.identity = a*1000 + b*100 + c*10 + d
}

// bits 20-32 are the altitude
// the 1 bits are AC13 field
// 00000000 00000000 00011111 1M1Q1111 00000000
func (f *Frame) decode13bitAltitudeCode() error {

	f.ac = uint32(f.message[2]&0xf)<<8 | uint32(f.message[3])

	// altitude
	f.acM = f.ac&0x40 == 0x40 // bit 26 of message. 0 == feet, 1 = metres
	// resolution
	f.acQ = f.ac&0x10 == 0x10 // bit 28 of message. 1 = 25 ft encoding, 0 = Gillham Mode C encoding

	// make sure all the bits are good

	if !f.acM {
		f.unit = modesUnitFeet

		/* N is the 11 bit integer resulting from the removal of bit Q and M */
		var msg2 = int32(f.message[2])
		var msg3 = int32(f.message[3])
		var n = int32((msg2&31)<<6) | ((msg3 & 0x80) >> 2) | ((msg3 & 0x20) >> 1) | (msg3 & 15)

		if f.acQ {
			// 25 ft increments
			/* The final altitude is due to the resulting number multiplied
			 * by 25, minus 1000. */
			f.altitude = (n * 25) - 1000
			f.validAltitude = true
		} else {
			// altitude reported in feet, 100ft increments
			/* Annex 10 — Aeronautical Telecommunications:
			   SSR automatic pressure-altitude transmission code (pulse position assignment)
			*/
			/* If the M bit (bit 26) and the Q bit (bit 28) equal 0, the altitude shall be coded according to the
			   pattern for Mode C replies of 3.1.1.7.12.2.3. Starting with bit 20 the sequence shall be
			   C1, A1, C2, A2, C4, A4, ZERO, B1, ZERO, B2, D2, B4, D4.
			*/

			f.altitude = gillhamToAltitude(n)
			f.validAltitude = true
		}
	} else {
		f.unit = modesUnitMetres
	}
	return nil
}

func (f *Frame) getMessageLengthBits() uint32 {
	//if f.downLinkFormat & 0x10 != 0 {
	if f.downLinkFormat&0x10 != 0 {
		if len(f.raw) == 14 {
			return modesShortMsgBits
		}
		return modesLongMsgBits
	} else {
		return modesShortMsgBits
	}
}

func (f *Frame) getMessageLengthBytes() uint32 {
	if f.downLinkFormat&0x10 != 0 {
		return modesLongMsgBytes
	} else {
		return modesShortMsgBytes
	}
}

func (f *Frame) decodeFlightNumber() {
	f.flight = decodeFlightNumber(f.message[5:11])
}

func decodeFlightNumber(b []byte) []byte {
	if 6 != len(b) {
		panic(fmt.Sprintf("attempting to decode a flight number/callsign with too many bytes (%d)", len(b)))
	}
	callsign := make([]byte, 8)
	callsign[0] = aisCharset[b[0]>>2]
	callsign[1] = aisCharset[((b[0]&3)<<4)|(b[1]>>4)]
	callsign[2] = aisCharset[((b[1]&15)<<2)|(b[2]>>6)]
	callsign[3] = aisCharset[b[2]&63]
	callsign[4] = aisCharset[b[3]>>2]
	callsign[5] = aisCharset[((b[3]&3)<<4)|(b[4]>>4)]
	callsign[6] = aisCharset[((b[4]&15)<<2)|(b[5]>>6)]
	callsign[7] = aisCharset[b[5]&63]
	return callsign
}

//func (f *Frame) decodeFlightId() {
//	if f.message[4] == 32 && len(f.message) >= 10 {
//		// Aircraft Identification
//		f.decodeFlightNumber()
//	}
//}
