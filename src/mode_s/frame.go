package mode_s

/*
  This file contains the main messages
*/

import (
	"time"
	"strings"
	"fmt"
)

const (
	MODES_UNIT_FEET = 0
	MODES_UNIT_METRES = 1
	DF17_FRAME_ID_CAT = "Aircraft Identification and Category"
	DF17_FRAME_SURFACE_POS = "Surface Position"
	DF17_FRAME_AIR_POS_BARO = "Airborne Position (with Barometric Altitude)"
	DF17_FRAME_AIR_VELOCITY = "Airborne Velocity"
	DF17_FRAME_AIR_VELOCITY_UNKNOWN = "Airborne Velocity (unknown sub type)"
	DF17_FRAME_AIR_POS_GNSS = "Airborne Position (with GNSS Height)"
	DF17_FRAME_TEST_MSG = "Test Message"
	DF17_FRAME_TEST_MSG_SQUAWK = "Test Message with Squawk"
	DF17_FRAME_SURFACE_SYS_STATUS = "Surface System Status"
	DF17_FRAME_EMERG_PRIORITY = "Extended Squitter Aircraft Status (Emergency Or Priority)"
	DF17_FRAME_EMERG_PRIORITY_UNKNOWN = "Unknown Emergency or Priority message"
	DF17_FRAME_TCAS_RA = "Extended Squitter Aircraft Status (1090ES TCAS Resolution Advisory)"
	DF17_FRAME_STATE_STATUS = "Target State and Status Message"
	DF17_FRAME_AIRCRAFT_OPER = "Aircraft Operational Status Message"
)

type Position struct {
	validAltitude       bool
	altitude            int32
	unit                int

	rawLatitude         int     /* Non decoded latitude */
	rawLongitude        int     /* Non decoded longitude */

	eastWestDirection   int     /* 0 = East, 1 = West. */
	eastWestVelocity    int     /* E/W velocity. */
	northSouthDirection int     /* 0 = North, 1 = South. */
	northSouthVelocity  int     /* N/S velocity. */
	validVelocity       bool
	velocity            float64 /* Computed from EW and NS velocity. */
	superSonic          bool

	verticalRateSource  int     /* Vertical rate source. */
	verticalRate        int     /* Vertical rate. */
	validVerticalRate   bool

	onGround            bool    /* VS Bit */
	validVerticalStatus bool

	heading             float64
	validHeading        bool

	haeDirection        byte    //up or down increments of 25
	haeDelta            int
	validHae            bool
}

type df17 struct {
	messageType             byte   // DF17 Extended Squitter Message Type
	messageSubType          byte   // DF17 Extended Squitter Message Sub Type

	cprFlagOddEven          int    /* 1 = Odd, 0 = Even CPR message. */
	timeFlag                int    /* UTC synchronized? */
	flight                  []byte /* 8 chars flight number. */

	validCompatibilityClass bool
	compatibilityClass      int
	cccHasOperationalTcas   *bool
	cccHas1090EsIn          bool
	cccHasAirRefVel         *bool  // supports Air Referenced Velocity
	cccHasLowTxPower        *bool
	cccHasTargetStateRpt    *bool  // supports Target State Report
	cccHasTargetChangeRpt   *bool  // supports Target Change Report
	cccHasUATReceiver       bool
	validNacV               bool

	operationalModeCode     int
	adsbVersion             byte
	nacP                    byte   // Navigation accuracy category - position
	geoVertAccuracy         byte   // geometric vertical accuracy
	sil                     byte
	airframe_width_len      byte
	nicCrossCheck           byte   // whether or not the alt or heading is cross checked
	northReference          byte   // 0=true north, 1 = magnetic north

	surveillanceStatus      byte
	nicSupplementA          byte
	nicSupplementB          byte
	nicSupplementC          byte
	containmentRadius       int

	intentChange            byte
	ifrCapability           byte
	nacV                    byte
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
	catType, catSubType                byte
	catValid                           bool
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
		0:  "Short air-air surveillance (TCAS)",
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
		6:  "Downed Aircraft",
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

	adsbCompatibilityVersion = []string{
		0: "Conformant to DO-260/ED-102 and DO-242",
		1: "Conformant to DO-260A and DO-242A",
		2: "Conformant to DO-260B/ED-102A and DO-242B",
		3: "reserved",
		4: "reserved",
		5: "reserved",
		6: "reserved",
		7: "reserved",
	}

	surveillanceStatus = []string{
		0: "No condition information",
		1: "Permanent alert (emergency condition)",
		2: "Temporary alert (change in Mode A identity code other than emergency condition)",
		3: "SPI condition",
	}
)

func (df *Frame) MessageTypeString() string {
	var name string = "Unknown"
	if df.messageType >= 1 && df.messageType <= 4 {
		name = DF17_FRAME_ID_CAT
	} else if df.messageType >= 5 && df.messageType <= 8 {
		name = DF17_FRAME_SURFACE_POS
	} else if df.messageType >= 9 && df.messageType <= 18 {
		name = DF17_FRAME_AIR_POS_BARO
	} else if df.messageType == 19 {
		if df.messageSubType >= 1 && df.messageSubType <= 4 {
			name = DF17_FRAME_AIR_VELOCITY
		} else {
			name = DF17_FRAME_AIR_VELOCITY_UNKNOWN
		}
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
	} else if df.messageType == 28 {
		if df.messageSubType == 1 {
			name = DF17_FRAME_EMERG_PRIORITY

		} else if df.messageSubType == 2 {
			name = DF17_FRAME_TCAS_RA
		} else {
			name = DF17_FRAME_EMERG_PRIORITY_UNKNOWN
		}
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

func (f *Frame) DownLinkType() byte {
	return f.downLinkFormat
}

func (f *Frame) ICAOAddr() uint32 {
	return f.icao
}

func (f *Frame) ICAOString() string {
	return fmt.Sprintf("%06X", f.icao)
}

func (f *Frame) Latitude() int {
	return f.rawLatitude
}
func (f *Frame) Longitude() int {
	return f.rawLongitude
}

func (f *Frame) Altitude() (int32, error) {
	if f.validAltitude {
		return f.altitude, nil
	}
	return 0, fmt.Errorf("Altitude is not valid")
}

func (f *Frame) AltitudeUnits() string {
	if f.unit == MODES_UNIT_METRES {
		return "metres"
	} else {
		return "feet"
	}
}

func (f *Frame) AltitudeValid() bool {
	return f.validAltitude
}

func (f *Frame) FlightStatusString() string {
	return flightStatusTable[f.fs]
}

func (f *Frame) FlightStatus() byte {
	return f.fs
}

func (f *Frame) Velocity() (float64, error) {
	if f.validVelocity {
		return f.velocity, nil
	}
	return 0, fmt.Errorf("Velocity is not valid")
}

func (f *Frame) VelocityValid() bool {
	return f.validVelocity
}

func (f *Frame) Heading() (float64, error) {
	if f.validHeading {
		return f.heading, nil
	}
	return 0, fmt.Errorf("Heading is not valid")
}

func (f *Frame) HeadingValid() bool {
	return f.validHeading
}

func (f *Frame) VerticalRate() (int, error) {
	if f.validVerticalRate {
		return f.verticalRate, nil
	}
	return 0, fmt.Errorf("Vertical Rate is not valid")
}

func (f *Frame) VerticalRateValid() bool {
	return f.validVerticalRate
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

func (f *Frame) OnGround() (bool, error) {
	if f.validVerticalStatus {
		return f.onGround, nil
	}
	return false, fmt.Errorf("Vertical Status is not valid")
}
func (f *Frame) VerticalStatusValid() bool {
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

func (df *Frame) MessageType() byte {
	return df.messageType
}

func (df *Frame) MessageSubType() byte {
	return df.messageSubType
}

// Whether or not this frame is even or odd, for CPR Location
func (df *Frame) IsEven() bool {
	return df.cprFlagOddEven == 0
}

func (df *Frame) FlightNumber() string {
	return string(df.flight)
}

/**
 * horizontal containment radius limit in meters.
 * Set NIC supplement A from Operational Status Message for better precision.
 * Otherwise, we'll be pessimistic.
 * Note: For ADS-B versions < 2, this is inaccurate for NIC class 6, since there was
 * no NIC supplement B in earlier versions.
 */
func (f *Frame) ContainmentRadiusLimit(nic_suppl_a bool) (float64, error) {
	var radius float64
	var err error
	if f.downLinkFormat != 17 {
		return radius, fmt.Errorf("ContainmentRadiusLimit Only valid for ADS-B Airborne Position Messages")
	}
	switch (f.messageType) {
	case 0, 18, 22:
		err = fmt.Errorf("Unknown Containment Radius")
	case 9, 20:
		radius = 7.5;
	case 10, 21:
		radius = 25;
	case 11:
		if nic_suppl_a {
			radius = 75
		} else {
			radius = 185.2
		}
	case 12:
		radius = 370.4;
	case 13:
		if 0 == f.nicSupplementB {
			radius = 926
		} else if nic_suppl_a {
			radius = 1111.2
		} else {
			radius = 555.6
		}
	case 14:
		radius = 1852
	case 15:
		radius = 3704
	case 16:
		if nic_suppl_a {
			radius = 7408
		} else {
			radius = 14816
		}
	case 17:
		radius = 37040;
	default:
		radius = 0;
	}

	return radius, err
}

func (f *Frame) NavigationIntegrityCategory(nic_suppl_a bool) (byte, error) {
	var nic byte
	var err error
	if f.downLinkFormat != 17 {
		return nic, fmt.Errorf("ContainmentRadiusLimit Only valid for ADS-B Airborne Position Messages")
	}
	switch (f.messageType) {
	case 0, 18, 22:
		err = fmt.Errorf("Unknown Navigation Integrity Category");
	case 9, 20:
		nic = 11;
	case 10: case 21:
		nic = 10;
	case 11:
		if nic_suppl_a {
			nic = 9;
		}else {
			nic = 8
		}
	case 12:
		nic = 7;
	case 13:
		nic = 6;
	case 14:
		nic = 5;
	case 15:
		nic = 4;
	case 16:
		if nic_suppl_a {
			nic = 3
		} else {
			nic = 2
		}
	case 17:
		nic = 1;
	default:
		nic = 0;
	}

	return nic, err
}

/**
 * Gets the air frames size in metres
 */
func (f *Frame)getAirplaneLengthWidth() (float32, float32, error) {
	if ! (f.messageType == 31 && f.messageSubType == 1) {
		return 0, 0, fmt.Errorf("Can only get aircraft size from ADSB message 31 sub type 1")
	}
	var length, width float32
	var err error

	switch (f.airframe_width_len) {
	case 1:
		length = 15
		width = 23
	case 2:
		length = 25
		width = 28.5
	case 3:
		length = 25
		width = 34
	case 4:
		length = 35
		width = 33
	case 5:
		length = 35
		width = 38
	case 6:
		length = 45
		width = 39.5
	case 7:
		length = 45
		width = 45
	case 8:
		length = 55
		width = 45
	case 9:
		length = 55
		width = 52
	case 10:
		length = 65
		width = 59.5
	case 11:
		length = 65
		width = 67
	case 12:
		length = 75
		width = 72.5
	case 13:
		length = 75
		width = 80
	case 14:
		length = 85
		width = 80
	case 15:
		length = 85
		width = 90
	default:
		err = fmt.Errorf("Unable to determine airframes size")
	}

	return length, width, err
}