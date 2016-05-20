package mode_s

/*
  This file contains the main messages
*/

import (
	"time"
	"strings"
)

const MODES_UNIT_FEET = 0
const MODES_UNIT_METRES = 1

type Position struct {
	altitude            int32
	rawLatitude         int     /* Non decoded latitude */
	rawLongitude        int     /* Non decoded longitude */
	eastWestDirection   int     /* 0 = East, 1 = West. */
	eastWestVelocity    int     /* E/W velocity. */
	northSouthDirection int     /* 0 = North, 1 = South. */
	northSouthVelocity  int     /* N/S velocity. */
	verticalRateSource  int     /* Vertical rate source. */
	verticalRate        int     /* Vertical rate. */
	velocity            float64 /* Computed from EW and NS velocity. */
	unit                int
	onGround            bool    /* VS Bit */
	validVerticalStatus bool
	validAltitude       bool
	superSonic          bool
}

type df17 struct {
	messageType    byte   // DF17 Extended Squitter Message Type
	messageSubType byte   // DF17 Extended Squitter Message Sub Type

						  //headingIsValid int
	heading        float64
	aircraftType   int
	cprFlagOddEven int    /* 1 = Odd, 0 = Even CPR message. */
	timeFlag       int    /* UTC synchronized? */
	flight         []byte /* 8 chars flight number. */
}

type raw_fields struct {
	// fields named what they are. see describe.go for what they mean

	df, vs, ca, cc, sl, ri, dr, um, fs byte
	ac, ap, id, aa, pi                 uint32
	mv, me, mb                         uint64
	md                                 [10]byte

	// altitude decoding
	ac_q, ac_m                         bool

	// adsb decoding
	catType, catSubType byte
	catValid bool
}

type Frame struct {
	raw_fields
	bds
	df17
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
	special        string
	alert          bool
						// if we have trouble decoding our frame, the message ends up here
	err            error
}

var (
	downlinkFormatTable = map[byte]string{
		0:  "short air-air surveillance (TCAS)",
		4:  "Roll Call Reply - Altitude (~100ft accuracy)",
		5:  "Roll Call Reply - Squawk",
		11: "All-Call reply containing aircraft address", // transponder capabilities
		16: "Long air-air surveillance (TCAS)",
		17: "ADS-B",
		18: "TIS-B - Ground Traffic", // ground traffic
		19: "Military Ext. Squitter",
		20: "Airborne position, GNSS HAE",
		21: "Roll Call Reply - Identity",
		22: "Military",
		24: "Comm. D Extended Length Message (ELM)",
	}

// DownLink Format Sub Type Capability CA
	capabilityTable = map[byte]string{
		0: "Level 1 no communication capability (Survillance Only)", // 0,4,5,11
		1: "Level 2 Comm-A and Comm-B capability", // DF 0,4,5,11,20,21
		2: "Level 3 Comm-A, Comm-B and uplink ELM capability", // (DF0,4,5,11,20,21)
		3: "Level 4 Comm-A, Comm-B uplink and downlink ELM capability", // (DF0,4,5,11,20,21,24)
		4: "Level 2,3 or 4. can set code 7. is on ground", // DF0,4,5,11,20,21,24,code7
		5: "Level 2,3 or 4. can set code 7. is airborne", // DF0,4,5,11,20,21,24,
		6: "Level 2,3 or 4. can set code 7.",
		7: "Level 7 DR≠0 or FS=3, 4 or 5",
	}

	flightStatusTable = map[byte]string{
		0: "Normal, Airborne",
		1: "Normal, On the ground",
		2: "ALERT, Airborne",
		3: "ALERT, On the ground",
		4: "ALERT, Special Position Identification. Airborne or Ground",
		5: "Normal, Special Position Identification. Airborne or Ground",
		6: "Value 6 is not assigned",
		7: "Value 7 is not assigned",
	}

	emergencyStateTable = map[int]string{
		0:  "No emergency",
		1:  "General emergency (squawk 7700)",
		2:  "Lifeguard/Medical",
		3:  "Minimum fuel",
		4:  "No communications (squawk 7600)",
		5:  "Unlawful interference (squawk 7500)",
		6:  "Reserved",
		7:  "Reserved",
	};

	replyInformationField = map[byte]string{
		0: "No on-board TCAS.",
		1: "Not assigned.",
		2: "On-board TCAS with resolution capability inhibited.",
		3: "On-board TCAS with vertical-only resolution capability.",
		4: "On-board TCAS with vertical and horizontal resolution capability.",
		5: "Not assigned.",
		6: "Not assigned.",
		7: "Not assigned.",
		8: "No maximum airspeed data available.",
		9: "Airspeed is ≤75kts.",
		10: "Airspeed is >75kts and ≤150kts.",
		11: "Airspeed is >150kts and ≤300kts.",
		12: "Airspeed is >300kts and ≤600kts.",
		13: "Airspeed is >600kts and ≤1200kts.",
		14: "Airspeed is >1200kts.",
		15: "Not assigned.",
	}

	sensitivityLevelInformationField = []string{
		"No TCAS sensitivity level reported",
		"TCAS sensitivity level 1. Likely on Ground (or TCAS Broken)",
		"TCAS sensitivity level 2. TA-Only. Pilot Selected",
		"TCAS sensitivity level 3.",
		"TCAS sensitivity level 4.",
		"TCAS sensitivity level 5.",
		"TCAS sensitivity level 6.",
		"TCAS sensitivity level 7.",
	}

	aisCharset string = "?ABCDEFGHIJKLMNOPQRSTUVWXYZ????? ???????????????0123456789??????"

	downlinkRequestField = []string{
		0: "No downlink request.",
		1: "Request to send Comm-B message (B-Bit set).",
		2: "TCAS information available.",
		3: "TCAS information available and request to send Comm-B message.",
		4: "Comm-B broadcast #1 available.",
		5: "Comm-B broadcast #2 available.",
		6: "TCAS information and Comm-B broadcast #1 available.",
		7: "TCAS information and Comm-B broadcast #2 available.",
		8: "Not assigned.",
		9: "Not assigned.",
		10: "Not assigned.",
		11: "Not assigned.",
		12: "Not assigned.",
		13: "Not assigned.",
		14: "Not assigned.",
		15: "Request to send 30 segments signified by 15+n.",
		16: "Request to send 31 segments signified by 15+n.",
		17: "Request to send 32 segments signified by 15+n.",
		18: "Request to send 33 segments signified by 15+n.",
		19: "Request to send 34 segments signified by 15+n.",
		21: "Request to send 35 segments signified by 15+n.",
		22: "Request to send 36 segments signified by 15+n.",
		23: "Request to send 37 segments signified by 15+n.",
		24: "Request to send 38 segments signified by 15+n.",
		25: "Request to send 39 segments signified by 15+n.",
		26: "Request to send 40 segments signified by 15+n.",
		27: "Request to send 41 segments signified by 15+n.",
		28: "Request to send 42 segments signified by 15+n.",
		29: "Request to send 43 segments signified by 15+n.",
		30: "Request to send 44 segments signified by 15+n.",
		31: "Request to send 45 segments signified by 15+n.",

	}

	utilityMessageField = []string{
		0: "No operating ACAS",
		1: "Not assigned",
		2: "ACAS with resolution capability inhibited",
		3: "ACAS with vertical-only resolution capability",
		4: "ACAS with vertical and horizontal resolution capability",
	}

	aircraftCategory = [][]string{
		0:{
			0:"No ADS-B Emitter Category Information",
			1:"Light (< 15500 lbs)",
			2:"Small (15500 to 75000 lbs)",
			3:"Large (75000 to 300000 lbs)",
			4:"High Vortex Large (aircraft such as B-757)",
			5:"Heavy (> 300000 lbs)",
			6:"High Performance (> 5g acceleration and 400 kts)",
			7:"Rotorcraft",
		},
		1:{
			0:"No ADS-B Emitter Category Information",
			1:"Glider / sailplane",
			2:"Lighter-than-air",
			3:"Parachutist / Skydiver",
			4:"Ultralight / hang-glider / paraglider",
			5:"Reserved",
			6:"Unmanned Aerial Vehicle",
			7:"Space / Trans-atmospheric vehicle",
		},
		2:{
			0:"No ADS-B Emitter Category Information",
			1:"Surface Vehicle – Emergency Vehicle",
			2:"Surface Vehicle – Service Vehicle",
			3:"Point Obstacle (includes tethered balloons)",
			4:"Cluster Obstacle",
			5:"Line Obstacle",
			6:"Reserved",
			7:"Reserved",
		},
		3:{
			0:"Reserved",
			1:"Reserved",
			2:"Reserved",
			3:"Reserved",
			4:"Reserved",
			5:"Reserved",
			6:"Reserved",
			7:"Reserved",
		},
	}
)

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
	return flightStatusTable[f.fs]
}

func (f *Frame) FlightStatus() byte {
	return f.fs
}

func (f *Frame) Velocity() float64 {
	return f.velocity
}
func (f *Frame) Heading() float64 {
	return f.heading
}
func (f *Frame) VerticalRate() int {
	return f.verticalRate
}

func (f *Frame) Flight() string {
	flight := string(f.flightId)
	if "" == flight {
		flight = "??????"
	}
	return strings.Trim(flight, " ")
}

func (f *Frame) SquawkIdentity() uint32 {
	return f.identity
}

func (f *Frame) OnGround() bool {
	return f.onGround
}
func (f *Frame) ValidVerticalStatus() bool {
	return f.validVerticalStatus
}
func (f *Frame) Alert() bool {
	return f.alert
}

func (f *Frame) ValidCategory() bool {
	return f.catValid
}

func (f *Frame) Category() string {
	return aircraftCategory[f.catType][f.catSubType]
}
