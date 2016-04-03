package mode_s

/*
  This file contains the main messages
*/

import (
	"fmt"
	"io"
	"time"
)

const MODES_UNIT_FEET = 0
const MODES_UNIT_METRES = 1

type Position struct {
	altitude            int32
	lat, lon            int
	rawLatitude         int     /* Non decoded latitude */
	rawLongitude        int     /* Non decoded longitude */
	eastWestDirection   int     /* 0 = East, 1 = West. */
	eastWestVelocity    int     /* E/W velocity. */
	northSouthDirection int     /* 0 = North, 1 = South. */
	northSouthVelocity  int     /* N/S velocity. */
	verticalRateSource  int     /* Vertical rate source. */
	verticalRateSign    int     /* Vertical rate sign. */
	verticalRate        int     /* Vertical rate. */
	velocity            float64 /* Computed from EW and NS velocity. */
	unit                int
}

type df11 struct {
	capability byte //  DF11 Capability Sub Type

}
type df4_5_20_21 struct {
	flightStatus int
	dr           int /* Request extraction of down link request. */
	um           int /* Request extraction of down link request. */
	identity     int
}

type df17 struct {
	messageType    byte   // DF17 Extended Squitter Message Type
	messageSubType byte   // DF17 Extended Squitter Message Sub Type

	headingIsValid int
	heading        float64
	aircraftType   int
	fFlag          int    /* 1 = Odd, 0 = Even CPR message. */
	timeFlag       int    /* UTC synchronized? */
	flight         []byte /* 8 chars flight number. */
}

type Frame struct {
	df11
	df17
	df4_5_20_21
	Position
	mode           string
	timeStamp      time.Time
	raw            string
	message        []byte
	downLinkFormat byte // Down link Format (DF)
	icao           uint32
	crc, checkSum  uint32
	identity       uint32
	flightId       []byte

	// if we have trouble decoding our frame, the message ends up here
	err error
}

var downlinkFormatTable = map[byte]string{
	0:  "short air-air surveillance (TCAS)",
	4:  "Roll Call Reply - Altitude (~100ft accuracy)",
	5:  "Roll Call Reply - Squawk",
	11: "All-Call reply containing aircraft address", // transponder capabilities
	16: "Long air-air surveillance (TCAS)",
	17: "ADS-B",
	18: "TIS-B - Ground Traffic", // ground traffic
	19: "Military Ext. Squitter",
	20: "Roll Call Reply - Altitude (~25ft accuracy)",
	21: "Roll Call Reply - Identity",
	22: "Military",
}

// DownLink Format Sub Type Capability CA
var capabilityTable = map[byte]string{
	0: "Level 1 (Survillance Only)",
	1: "Level 2 (DF0,4,5,11)",
	2: "Level 3 (DF0,4,5,11,20,21)",
	3: "Level 4 (DF0,4,5,11,20,21,24)",
	4: "Level 2+3+4 (DF0,4,5,11,20,21,24,code7 - is on ground)",
	5: "Level 2+3+4 (DF0,4,5,11,20,21,24,code7 - is airborne)",
	6: "Level 2+3+4 (DF0,4,5,11,20,21,24,code7)",
	7: "Level 7 ???",
}

var flightStatusTable = map[int]string{
	0: "Normal, Airborne",
	1: "Normal, On the ground",
	2: "ALERT,  Airborne",
	3: "ALERT,  On the ground",
	4: "ALERT & Special Position Identification. Airborne or Ground",
	5: "Special Position Identification. Airborne or Ground",
	6: "Value 6 is not assigned",
	7: "Value 7 is not assigned",
}

var aisCharset string = "?ABCDEFGHIJKLMNOPQRSTUVWXYZ????? ???????????????0123456789??????"

// prints out a nice debug message
func (f *Frame) Describe(output io.Writer) {
	fmt.Fprintln(output, "----------------------------------------------------")
	fmt.Fprintf(output, "MODE S Packet: ICAO = %x\n", f.icao)
	fmt.Fprintf(output, "Frame mode  : %s\n", f.mode)
	fmt.Fprintf(output, "Time Stamp  : %s\n", f.timeStamp.Format(time.RFC3339Nano))
	fmt.Fprintf(output, "Frame Type  : %d (%s)\n", f.downLinkFormat, f.GetDownLinkFormat())
	if 17 == f.downLinkFormat {
		fmt.Printf("ADSB Frame  : %d (%s)\n", f.messageType, f.MessageTypeString())
		if f.messageType >= 1 && f.messageType <= 4 {
			fmt.Printf("      Flight: %s", string(f.flight))
		} else if f.messageType >= 9 && f.messageType <= 18 {
			var oddEven string = "Odd"
			if 0 == f.fFlag {
				oddEven = "Even"
			}
			fmt.Printf("   LAT LONG : ?? (%s frame)\n", oddEven)
		}
	} else if 11 == f.downLinkFormat {
		fmt.Fprintf(output, "Capability  : %d (%s)\n", f.df11.capability, f.DownLinkCapability())
	}

	fmt.Fprintln(output, "Position:")
	fmt.Fprintf(output, "    Altitude: %d feet", f.altitude)
	fmt.Fprintln(output, "")
}

// determines what type of mode S Message this frame is
func (f *Frame) GetDownLinkFormat() string {

	if description, ok := downlinkFormatTable[f.downLinkFormat]; ok {
		return description
	}
	return "Unknown Downlink Format"
}

// for when mode is 4,5,20,21
func (f *Frame) DownLinkCapability() string {

	if description, ok := capabilityTable[f.df11.capability]; ok {
		return description
	}
	return "Unknown Downlink Capability"
}

func (f *Frame) DownLinkType() byte {
	return f.downLinkFormat
}

func (f *Frame) ICAOAddr() uint32 {
	return f.icao
}

func (f *Frame) Latitude() int {
	return f.rawLatitude
}
func (f *Frame) Longitude() int {
	return f.rawLongitude
}
func (f *Frame) Altitude() int32 {
	return f.altitude
}
func (f *Frame) AltitudeUnits() string {
	if f.unit == MODES_UNIT_METRES {
		return "metres"
	} else {
		return "feet"
	}
}

func (f *Frame) FlightStatusString() string {
	return flightStatusTable[f.flightStatus]
}

func (f *Frame) FlightStatusInt() int {
	return f.flightStatus
}

func (f *Frame) Velocity() float64 {
	return f.velocity
}
func (f *Frame) Heading() float64 {
	return f.heading
}

func (f *Frame) Flight() string {
	return string(f.flightId)
}