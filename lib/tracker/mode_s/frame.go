package mode_s

/*
  This file contains the main messages
*/

import (
	"fmt"
	"regexp"
	"time"
)

const (
	modesUnitFeet                     = 0
	modesUnitMetres                   = 1
	DF17FrameIdCat                    = "Aircraft Identification and Category"
	DF17FrameSurfacePos               = "Surface Position"
	DF17FrameAirPositionBarometric    = "Airborne Position (with Barometric altitude)"
	DF17FrameAirVelocity              = "Airborne velocity"
	DF17FrameAirVelocityUnknown       = "Airborne velocity (unknown sub type)"
	DF17FrameAirPositionGnss          = "Airborne Position (with GNSS Height)"
	DF17FrameTestMessage              = "Test Message"
	DF17FrameTestMessageSquawk        = "Test Message with squawk"
	DF17FrameSurfaceSystemStatus      = "Surface System status"
	DF17FrameEmergencyPriority        = "Extended Squitter Aircraft status (Emergency Or Priority)"
	DF17FrameEmergencyPriorityUnknown = "Unknown Emergency or Priority message"
	DF17FrameTcasRA                   = "Extended Squitter Aircraft status (1090ES TCAS Resolution Advisory)"
	DF17FrameTargetStateStatus        = "Target State and status Message"
	DF17FrameAircraftOperational      = "Aircraft Operational status Message"
)

type (
	Position struct {
		validAltitude bool
		altitude      int32
		isGnssAlt     bool
		unit          int

		rawLatitude  int /* Non decoded latitude */
		rawLongitude int /* Non decoded longitude */

		eastWestDirection   int /* 0 = East, 1 = West. */
		eastWestVelocity    int /* E/W velocity. */
		northSouthDirection int /* 0 = North, 1 = South. */
		northSouthVelocity  int /* N/S velocity. */
		validVelocity       bool
		velocity            float64 /* Computed from EW and NS velocity. */
		superSonic          bool

		verticalRateSource int /* Vertical rate source. */
		verticalRate       int /* Vertical rate. */
		validVerticalRate  bool

		onGround            bool /* VS Bit */
		validVerticalStatus bool

		heading      float64
		validHeading bool

		haeDirection byte //up or down increments of 25
		haeDelta     int
		validHae     bool
	}

	df17 struct {
		messageType    byte // DF17 Extended Squitter Message Type
		messageSubType byte // DF17 Extended Squitter Message Sub Type

		cprFlagOddEven int    /* 1 = Odd, 0 = Even CPR message. */
		timeFlag       int    /* UTC synchronized? */
		flight         []byte /* 8 chars flight number. */

		validCompatibilityClass bool
		compatibilityClass      int
		cccHasOperationalTcas   *bool
		cccHas1090EsIn          bool
		cccHasAirRefVel         *bool // supports Air Referenced velocity
		cccHasLowTxPower        *bool
		cccHasTargetStateRpt    *bool // supports Target State Report
		cccHasTargetChangeRpt   *bool // supports Target Change Report
		cccHasUATReceiver       bool
		validNacV               bool

		operationalModeCode int
		adsbVersion         byte
		nacP                byte // Navigation accuracy category - position
		geoVertAccuracy     byte // geometric vertical accuracy
		sil                 byte
		airframeWidthLen    byte
		nicCrossCheck       byte // whether or not the alt or heading is cross checked
		northReference      byte // 0=true north, 1 = magnetic north

		surveillanceStatus byte
		nicSupplementA     byte
		nicSupplementB     byte
		nicSupplementC     byte
		containmentRadius  int

		intentChange  byte
		ifrCapability byte
		nacV          byte
	}

	extendedSquitter struct {
		Df byte   `bits:"0-5" name:"DF" desc:"Downlink Format"`
		Ca byte   `bits:"5-8" name:"CA" desc:"Aircraft System Capability"`
		Aa uint32 `bits:"8-32" name:"AA" desc:"Address Announce, ICAO Identification"`
		Me uint64 `bits:"32-88" name:"ME" desc:"ADSB Message"`
		Pi uint32 `bits:"88-111" name:"PI" desc:"Parity/Interr.Identity: reports source of interrogation. Contains the parity overlaid on the interrogator identity code"`

		df17me1

		df17me19t1
		df17me19t3

		df17me28t1
		df17me28t2

		df17me29t0
		df17me29t1
		df17me29t2

		df17me31t0
		df17me31t1
	}

	// df17me1 is Aircraft Identification and Category
	df17me1 struct {
		CategoryType    byte `bits:"33-38" name:"CAT" desc:"Aircraft Category"`
		CategorySubType byte `bits:"38-41" name:"CAT" desc:"Aircraft Category Sub Type"`
		IdentChar1      byte `bits:"41-47" name:"ident char"`
		IdentChar2      byte `bits:"47-53" name:"ident char"`
		IdentChar3      byte `bits:"53-59" name:"ident char"`
		IdentChar4      byte `bits:"59-65" name:"ident char"`
		IdentChar5      byte `bits:"65-71" name:"ident char"`
		IdentChar6      byte `bits:"71-77" name:"ident char"`
		IdentChar7      byte `bits:"77-83" name:"ident char"`
		IdentChar8      byte `bits:"83-89" name:"ident char"`
	}

	// df17me19t1 Airborne Velocity Message Subtype=1 and 2
	df17me19t1 struct {
		CategoryType    byte   `bits:"33-38"`
		CategorySubType byte   `bits:"38-41"`
		IntentChange    byte   `bits:"41-42" name:"Intent" desc:"Intent Change"`
		ReservedA       byte   `bits:"42-43" name:"Res" desc:"Reserved-A"`
		NacV            byte   `bits:"43-46"`
		EWDir           byte   `bits:"46-47" name:"EW" desc:"(0)East/(1)West Direction"`
		EWVelocity      uint16 `bits:"47-57" name:"EWV" desc:"East/West Velocity"`
		NSDir           byte   `bits:"57-58" name:"NS" desc:"(0)North/(1)South Direction"`
		NSVelocity      uint16 `bits:"58-68" name:"NSV" desc:"North/South Velocity"`
		VertRateSource  byte   `bits:"68-69" name:"VRS" desc:"Vertical Rate Source"`
		VertRateSign    byte   `bits:"69-70" name:"VRS+" desc:"Vertical Rate Sign"`
		VertRate        uint16 `bits:"70-79" name:"VR" desc:"Vertical Rate"`
		ReservedB       byte   `bits:"79-81" name:"Res" desc:"Reserved-B"`
		DiffBaroSign    byte   `bits:"81-82" name:"DB+" desc:"Diff from Baro Altitude Sign. 0=up, 1=down"`
		DiffBaroAlt     byte   `bits:"82-89" name:"DB+" desc:"Diff from Baro Altitude"`
	}

	// df17me19t3 Airborne Velocity Message Subtype=3 and 4
	df17me19t3 struct {
		CategoryType    byte   `bits:"33-38"`
		CategorySubType byte   `bits:"38-41"`
		IntentChange    byte   `bits:"41-42" name:"Intent" desc:"Intent Change"`
		ReservedA       byte   `bits:"42-43" name:"Res" desc:"Reserved-A"`
		NacV            byte   `bits:"43-46"`
		HeadingStatus   byte   `bits:"46-47" name:"HDS" desc:"Heading Status"`
		Heading         uint16 `bits:"47-57" name:"HD" desc:"Heading Status"`
		AirSpeedType    byte   `bits:"57-58" name:"AST" desc:"Air Speed Type"`
		AirSpeed        uint16 `bits:"58-68" name:"AS" desc:"Air Speed"`
		VertRateSource  byte   `bits:"68-69" name:"VRS" desc:"Vertical Rate Source"`
		VertRateSign    byte   `bits:"69-70" name:"VRS+" desc:"Vertical Rate Sign"`
		VertRate        uint16 `bits:"70-79" name:"VR" desc:"Vertical Rate"`
		ReservedB       byte   `bits:"79-81" name:"Res" desc:"Reserved-B"`
		DiffBaroSign    byte   `bits:"81-82" name:"DB+" desc:"Diff from Baro Altitude Sign. 0=up, 1=down"`
		DiffBaroAlt     byte   `bits:"82-89" name:"DB+" desc:"Diff from Baro Altitude"`
	}
	// df17me28t1 is a Emergency / Priority Status and Mode A Code(Subtype=1)
	df17me28t1 struct {
		Type     byte `bits:"33-38" name:"Type" desc:"ADS-B Message Type"`
		SubType  byte `bits:"38-41" name:"Sub" desc:"ADS-B Message Sub Type"`
		Status   byte `bits:"41-44" name:"Status" desc:"Emergency/Priority Status"`
		ModeA    byte `bits:"44-57" name:"ModeA" desc:"Mode A Code"`
		Reserved byte `bits:"57-89" name:"Res" desc:"Reserved"`
	}
	// df17me28t2 is a 1090ES TCAS Resolution Advisory (RA) Broadcast Message (Subtype=2
	df17me28t2 struct {
		Type    byte   `bits:"33-38" name:"Type" desc:"ADS-B Message Type"`
		SubType byte   `bits:"38-41" name:"Sub" desc:"ADS-B Message Sub Type"`
		Ara     uint16 `bits:"41-55" name:"ARA" desc:"Active Resolution Advisories"`
		Racs    byte   `bits:"55-59" name:"RACs" desc:"RACs Record"`
		Rat     byte   `bits:"59-60" name:"RAT" desc:"RA Terminated"`
		Mte     byte   `bits:"60-61" name:"MTE" desc:"Multiple Threat Encounter"`
		Tti     byte   `bits:"61-63" name:"TTI" desc:"Threat Type Indicator"`
		Tid     byte   `bits:"63-89" name:"TID" desc:"Threat Identity Data"`
	}
	df17me29t0 struct {
		// DO-260A
		Type    byte   `bits:"33-38" name:"Type" desc:"ADS-B Message Type"`
		SubType byte   `bits:"38-41" name:"Sub" desc:"ADS-B Message Sub Type"`
		unknown uint64 `bit:"40-89" name:"Unknown"`
	}
	df17me29t1 struct {
		// DO-260B - http://www.anteni.net/adsb/Doc/1090-WP30-18-DRAFT_DO-260B-V42.pdf
		Type            byte `bits:"33-38" name:"Type" desc:"ADS-B Message Type"`
		SubType         byte `bits:"38-41" name:"Sub" desc:"ADS-B Message Sub Type"`
		SilSupplement   byte `bits:"40-41" name:"SIL sup" desc:"SIL Supplement: SIL Per Hour (0) or Per Sample (1)"`
		AltType         byte `bits:"41-42" name:"Sel Alt Type" desc:"Selected Altitude Type 0 = Mode Control Panel/Flight Control Unit, 1=Flight Management System"`
		SelectedAlt     uint `bits:"42-53" name:"Sel Alt" desc:"MCP/FCU Selected Altitude OR FMS Selected Altitude"`
		BaroSetting     uint `bits:"53-62" name:"Baro Set" desc:"Barometric Pressure Setting (minus 800 millibars)"`
		SelHeadingStat  byte `bits:"62-63" name:"Sel Hd Stat" desc:"Selected Heading Status"`
		SelHeadingSign  byte `bits:"63-64" name:"Sign" desc:"Selected Heading Sign"`
		SelHeading      byte `bits:"64-72" name:"Sel Hd" desc:"Selected Heading"`
		NacP            byte `bits:"72-76" name:"NACp" desc:"Navigation Accuracy Category_Position"`
		NicBaro         byte `bits:"76-77" name:"NICbaro" desc:"Navigation Integrity Category_Baro"`
		Sil             byte `bits:"77-79" name:"SIL" desc:"Source Integrity Level"`
		McuFpuStatus    byte `bits:"79-80" name:"MCU/FPU" desc:"MCU/FPU Status"`
		AutoPilot       byte `bits:"80-81" name:"Autopilot" desc:"Autopilot Engaged"`
		VNav            byte `bits:"81-82" name:"VNAV" desc:"VNav Mode Engaged"`
		AltHoldMode     byte `bits:"82-83" name:"Alt Hold" desc:"Altitude Hold Mode"`
		ResAdsR         byte `bits:"83-84" name:"ads-r" desc:"Reserved for ADS-R"`
		ApproachMode    byte `bits:"84-85" name:"Approach Mode" desc:"Approach Mode"`
		TcasOperational byte `bits:"85-86" name:"TCAS" desc:"TCAS Operational"`
		Reserved        byte `bits:"86-89" name:"Res" desc:"Reserved"`
	}
	df17me29t2 struct {
		// DO-260A
		Type     byte   `bits:"33-38" name:"Type" desc:"ADS-B Message Type"`
		SubType  byte   `bits:"38-41" name:"Sub" desc:"ADS-B Message Sub Type"`
		reserved uint64 `bits:"40-89" name:"Unknown"`
	}
	df17me31t0 struct {
		Type            byte   `bits:"33-38" name:"Type" desc:"ADS-B Message Type"`
		SubType         byte   `bits:"38-41" name:"Sub" desc:"ADS-B Message Sub Type"`
		CapClass        uint16 `bits:"41-57" name:"CC" desc:"Capability Class"`
		OperationalMode uint16 `bits:"57-73" name:"OM" desc:"Operational Mode Codes"`
		MopsVer         byte   `bits:"73-76" name:"MOPS" desc:"MOPS Version"`
		NicSuppA        byte   `bits:"76-77" name:"NicA" desc:"NIC Supp-A"`
		NACp            byte   `bits:"77-81" name:"NicA" desc:"NIC Supp-A"`
		Gva             byte   `bits:"81-83" name:"GVA" desc:"Geometric  Vertical  Accuracy"`
		Sil             byte   `bits:"83-85" name:"SIL" desc:"Source Integrity Level"`
		NicBaro         byte   `bits:"85-86" name:"NICbaro" desc:"NICbaro"`
		Hrd             byte   `bits:"86-87" name:"HRD" desc:"Horizontal Reference Direction. 0=True North, 1=Magnetic"`
		SilSup          byte   `bits:"87-88" name:"Sil Sup" desc:"SIL Supp"`
		Reserved        byte   `bits:"88-89" name:"Res" desc:"Reserved"`
	}
	df17me31t1 struct {
		Type            byte   `bits:"33-38" name:"Type" desc:"ADS-B Message Type"`
		SubType         byte   `bits:"38-41" name:"Sub" desc:"ADS-B Message Sub Type"`
		CapClassCodes   uint16 `bits:"41-53" name:"CCC" desc:"Capability Class Codes"`
		LWCodes         byte   `bits:"53-57" name:"L/W C" desc:"L/W Codes"`
		OperationalMode uint16 `bits:"57-73" name:"OM" desc:"Operational Mode Codes"`
		MopsVer         byte   `bits:"73-76" name:"VN" desc:"ADS-B MOPS Compliant Version"`
		NicSuppA        byte   `bits:"76-77" name:"NicA" desc:"NIC Supp-A"`
		NACp            byte   `bits:"77-81" name:"NicA" desc:"NIC Supp-A"`
		ReservedA       byte   `bits:"81-83" name:"Res" desc:"ReservedA"`
		Sil             byte   `bits:"83-85" name:"SIL" desc:"Source Integrity Level"`
		NicBaro         byte   `bits:"85-86" name:"NICbaro" desc:"NICbaro"`
		Hrd             byte   `bits:"86-87" name:"HRD" desc:"Horizontal Reference Direction.  0=True North, 1=Magnetic"`
		SilSup          byte   `bits:"87-88" name:"Sil Sup" desc:"SIL Supp"`
		ReservedB       byte   `bits:"88-89" name:"Res" desc:"ReservedB"`
	}

	rawFields struct {
		// fields named what they are. see describe.go for what they mean

		df, vs, ca, cc, sl, ri, dr, um, fs byte
		ac, ap, id, aa, pi                 uint32
		mv, me, mb                         uint64
		md                                 [10]byte

		// altitude decoding
		acQ, acM bool

		// adsb decoding
		catType, catSubType byte
		catValid            bool
	}

	Frame struct {
		rawFields
		bds
		df17
		Position
		mode string
		// the timestamp we are processing this message at
		timeStamp      time.Time
		beastTimeStamp string
		// beastTicks is the number of ticks since the beast was turned on
		beastTicks uint64
		// beastTicksNs is the number of nanoseconds since the beast was turned on
		beastTicksNs   uint64
		beastAvrUptime time.Duration
		// raw is our semi processed string, full is the original string
		raw, full      string
		message        []byte
		downLinkFormat byte // Down link Format (DF)
		icao           uint32
		crc, checkSum  uint32
		identity       uint32
		special        string
		emergency      string
		alert          bool
		// if we have trouble decoding our frame, the message ends up here
		err error
	}
)

var (
	downlinkFormatTable = map[byte]string{
		0:  "Short air-air surveillance (TCAS)",
		4:  "Roll Call Reply - altitude (~100ft accuracy)",
		5:  "Roll Call Reply - squawk",
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
		0: "Level 1 no communication capability (Survillance Only)",    // 0,4,5,11
		1: "Level 2 Comm-A and Comm-B capability",                      // DF 0,4,5,11,20,21
		2: "Level 3 Comm-A, Comm-B and uplink ELM capability",          // (DF0,4,5,11,20,21)
		3: "Level 4 Comm-A, Comm-B uplink and downlink ELM capability", // (DF0,4,5,11,20,21,24)
		4: "Level 2,3 or 4. can set code 7. is on ground",              // DF0,4,5,11,20,21,24,code7
		5: "Level 2,3 or 4. can set code 7. is airborne",               // DF0,4,5,11,20,21,24,
		6: "Level 2,3 or 4. can set code 7.",
		7: "Level 7 DR≠0 or FS=3, 4 or 5",
	}

	flightStatusTable = map[byte]string{
		0: "Normal, Airborne",
		1: "Normal, On the ground",
		2: "ALERT, Airborne",
		3: "ALERT, On the ground",
		4: "ALERT, special Position Identification. Airborne or Ground",
		5: "Normal, special Position Identification. Airborne or Ground",
		6: "Value 6 is not assigned",
		7: "Value 7 is not assigned",
	}

	emergencyStateTable = map[int]string{
		0: "No emergency",
		1: "General emergency (squawk 7700)",
		2: "Lifeguard/Medical",
		3: "Minimum fuel",
		4: "No communications (squawk 7600)",
		5: "Unlawful interference (squawk 7500)",
		6: "Downed Aircraft",
		7: "Reserved",
	}

	replyInformationField = map[byte]string{
		0:  "No on-board TCAS.",
		1:  "Not assigned.",
		2:  "On-board TCAS with resolution capability inhibited.",
		3:  "On-board TCAS with vertical-only resolution capability.",
		4:  "On-board TCAS with vertical and horizontal resolution capability.",
		5:  "Not assigned.",
		6:  "Not assigned.",
		7:  "Not assigned.",
		8:  "No maximum airspeed data available.",
		9:  "Airspeed is ≤75kts.",
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

	aisCharset = "?ABCDEFGHIJKLMNOPQRSTUVWXYZ????? ???????????????0123456789??????"

	downlinkRequestField = []string{
		0:  "No downlink request.",
		1:  "Request to send Comm-B message (B-Bit set).",
		2:  "TCAS information available.",
		3:  "TCAS information available and request to send Comm-B message.",
		4:  "Comm-B broadcast #1 available.",
		5:  "Comm-B broadcast #2 available.",
		6:  "TCAS information and Comm-B broadcast #1 available.",
		7:  "TCAS information and Comm-B broadcast #2 available.",
		8:  "Not assigned.",
		9:  "Not assigned.",
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
		0: {
			0: "No ADS-B Emitter Category Information",
			1: "Light (< 15500 lbs)",
			2: "Small (15500 to 75000 lbs)",
			3: "Large (75000 to 300000 lbs)",
			4: "High Vortex Large (aircraft such as B-757)",
			5: "Heavy (> 300000 lbs)",
			6: "High Performance (> 5g acceleration and 400 kts)",
			7: "Rotorcraft",
		},
		1: {
			0: "No ADS-B Emitter Category Information",
			1: "Glider / sailplane",
			2: "Lighter-than-air",
			3: "Parachutist / Skydiver",
			4: "Ultralight / hang-glider / paraglider",
			5: "Reserved",
			6: "Unmanned Aerial Vehicle",
			7: "Space / Trans-atmospheric vehicle",
		},
		2: {
			0: "No ADS-B Emitter Category Information",
			1: "Surface Vehicle – Emergency Vehicle",
			2: "Surface Vehicle – Service Vehicle",
			3: "Point Obstacle (includes tethered balloons)",
			4: "Cluster Obstacle",
			5: "Line Obstacle",
			6: "Reserved",
			7: "Reserved",
		},
		3: {
			0: "Reserved",
			1: "Reserved",
			2: "Reserved",
			3: "Reserved",
			4: "Reserved",
			5: "Reserved",
			6: "Reserved",
			7: "Reserved",
		},
	}

	// TC -> CA
	wakeVortex = [5][8]string{
		0: {},
		1: {},
		2: {
			1: "Surface emergency vehicle",
			2: "?",
			3: "Surface service vehicle",
			4: "Ground obstruction",
			5: "Ground obstruction",
			6: "Ground obstruction",
			7: "Ground obstruction",
		},
		3: {
			1: "Glider, sailplane",
			2: "Lighter-than-air",
			3: "Parachutist, skydiver",
			4: "Ultralight, hang-glider, paraglider",
			5: "Reserved",
			6: "Unmanned aerial vehicle",
			7: "Space or transatmospheric vehicle",
		},
		4: {
			1: "Light (less than 7000 kg)",
			2: "Medium 1 (between 7000 kg and 34000 kg)",
			3: "Medium 2 (between 34000 kg to 136000 kg)",
			4: "High vortex aircraft",
			5: "Heavy (larger than 136000 kg)",
			6: "High performance (>5 g acceleration) and high speed (>400 kt)",
			7: "Rotocraft",
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

func (f *Frame) MessageTypeString() string {
	name := "Unknown"
	if f.messageType >= 1 && f.messageType <= 4 {
		name = DF17FrameIdCat
	} else if f.messageType >= 5 && f.messageType <= 8 {
		name = DF17FrameSurfacePos
	} else if f.messageType >= 9 && f.messageType <= 18 {
		name = DF17FrameAirPositionBarometric
	} else if f.messageType == 19 {
		if f.messageSubType >= 1 && f.messageSubType <= 4 {
			name = DF17FrameAirVelocity
		} else {
			name = DF17FrameAirVelocityUnknown
		}
	} else if f.messageType >= 20 && f.messageType <= 22 {
		name = DF17FrameAirPositionGnss
	} else if f.messageType == 23 {
		if f.messageSubType == 7 {
			name = DF17FrameTestMessageSquawk
		} else {
			name = DF17FrameTestMessage
		}
	} else if f.messageType == 24 && f.messageSubType == 1 {
		name = DF17FrameSurfaceSystemStatus
	} else if f.messageType == 28 {
		if f.messageSubType == 1 {
			name = DF17FrameEmergencyPriority
			f.decodeFlightNumber()

		} else if f.messageSubType == 2 {
			name = DF17FrameTcasRA
		} else {
			name = DF17FrameEmergencyPriorityUnknown
		}
	} else if f.messageType == 29 {
		if f.messageSubType == 0 || f.messageSubType == 1 {
			name = DF17FrameTargetStateStatus
		} else {
			name = fmt.Sprintf("%s (Unknown Sub Message %d)", DF17FrameTargetStateStatus, f.messageSubType)
		}
	} else if f.messageType == 31 && (f.messageSubType == 0 || f.messageSubType == 1) {
		name = DF17FrameAircraftOperational
	}
	return name
}

func (f *Frame) DownLinkType() byte {
	return f.downLinkFormat
}

func (f *Frame) Icao() uint32 {
	return f.icao
}

func (f *Frame) Raw() []byte {
	if nil == f {
		return []byte{}
	}
	return []byte(f.raw)
}

func (f *Frame) IcaoStr() string {
	if nil == f {
		return ""
	}
	return fmt.Sprintf("%06X", f.icao)
}

func (f *Frame) Latitude() int {
	if nil == f {
		return -1
	}
	return f.rawLatitude
}
func (f *Frame) Longitude() int {
	if nil == f {
		return -1
	}
	return f.rawLongitude
}

func (f *Frame) Altitude() (int32, error) {
	if f.validAltitude {
		return f.altitude, nil
	}
	return 0, fmt.Errorf("altitude is not valid")
}
func (f *Frame) MustAltitude() int32 {
	if f.validAltitude {
		return f.altitude
	}
	panic("altitude is not valid")
}

func (f *Frame) AltitudeUnits() string {
	if nil == f {
		return "metres"
	}
	if f.unit == modesUnitMetres {
		return "metres"
	} else {
		return "feet"
	}
}

func (f *Frame) AltitudeValid() bool {
	if nil == f {
		return false
	}
	return f.validAltitude
}

func (f *Frame) FlightStatusString() string {
	if nil == f {
		return ""
	}
	return flightStatusTable[f.fs]
}

func (f *Frame) FlightStatus() byte {
	if nil == f {
		return 255
	}
	return f.fs
}

func (f *Frame) Velocity() (float64, error) {
	if f.validVelocity {
		return f.velocity, nil
	}
	return 0, fmt.Errorf("velocity is not valid")
}

func (f *Frame) MustVelocity() float64 {
	if f.validVelocity {
		return f.velocity
	}
	panic("velocity is not valid")
}

func (f *Frame) VelocityValid() bool {
	if nil == f {
		return false
	}
	return f.validVelocity
}

func (f *Frame) Heading() (float64, error) {
	if f.validHeading {
		return f.heading, nil
	}
	return 0, fmt.Errorf("heading is not valid")
}
func (f *Frame) MustHeading() float64 {
	if f.validHeading {
		return f.heading
	}
	panic("heading is not valid")
}

func (f *Frame) HeadingValid() bool {
	if nil == f {
		return false
	}
	return f.validHeading
}

func (f *Frame) VerticalRate() (int, error) {
	if f.VerticalRateValid() {
		return f.verticalRate, nil
	}
	return 0, fmt.Errorf("vertical rate (VR) is not valid")
}
func (f *Frame) MustVerticalRate() int {
	if f.VerticalRateValid() {
		return f.verticalRate
	}
	panic("vertical rate (VR) is not valid")
}

func (f *Frame) VerticalRateValid() bool {
	if nil == f {
		return false
	}
	return f.validVerticalRate
}

//func (f *Frame) flight() string {
//	flight := string(f.flightId)
//	if "" == flight {
//		flight = "??????"
//	}
//	return strings.Trim(flight, " ")
//}

func (f *Frame) SquawkIdentity() uint32 {
	if nil == f {
		return 0
	}
	return f.identity
}

func (f *Frame) OnGround() (bool, error) {
	if f.VerticalStatusValid() {
		return f.onGround, nil
	}
	return false, fmt.Errorf("vertical status (VS) is not valid")
}
func (f *Frame) MustOnGround() bool {
	if f.VerticalStatusValid() {
		return f.onGround
	}
	panic("vertical status (VS) is not valid")
}
func (f *Frame) VerticalStatusValid() bool {
	if nil == f {
		return false
	}
	return f.validVerticalStatus
}
func (f *Frame) Alert() bool {
	if nil == f {
		return false
	}
	return f.alert
}

func (f *Frame) ValidCategory() bool {
	if nil == f {
		return false
	}
	return f.catValid
}

func (f *Frame) Category() string {
	if !f.ValidCategory() {
		return ""
	}
	return aircraftCategory[f.catType][f.catSubType]
}
func (f *Frame) CategoryType() string {
	return fmt.Sprintf("%d/%d", f.catType, f.catSubType)
}

func (f *Frame) MessageType() byte {
	return f.messageType
}

func (f *Frame) MessageSubType() byte {
	return f.messageSubType
}

// Whether or not this frame is even or odd, for CPR location
func (f *Frame) IsEven() bool {
	return f.cprFlagOddEven == 0
}

func (f *Frame) FlightNumber() string {
	return string(f.flight)
}
func (f *Frame) Special() string {
	return f.special
}
func (f *Frame) HasSurveillanceStatus() bool {
	return f.surveillanceStatus > 0
}
func (f *Frame) SurveillanceStatus() string {
	if int(f.surveillanceStatus) < len(surveillanceStatus) {
		return surveillanceStatus[f.surveillanceStatus]
	}
	return ""
}
func (f *Frame) Emergency() string {
	return f.emergency
}

// the first character can be * or @ (or left out)
// if the entire string is then 0's, it's a noop
var noopRw = regexp.MustCompile("^[*@]?0+$")

func (f *Frame) isNoOp() bool {
	if "" == f.raw {
		return true
	}
	if len(f.raw) > 16 { // if we have a frame that is at least the right size, it is not a heart beat
		return false
	}
	return noopRw.MatchString(f.raw)
}

/**
 * horizontal containment radius limit in meters.
 * Set NIC supplement A from Operational status Message for better precision.
 * Otherwise, we'll be pessimistic.
 * Note: For ADS-B versions < 2, this is inaccurate for NIC class 6, since there was
 * no NIC supplement B in earlier versions.
 */
func (f *Frame) ContainmentRadiusLimit(nicSupplA bool) (float64, error) {
	var radius float64
	var err error
	if f.downLinkFormat != 17 {
		return radius, fmt.Errorf("ContainmentRadiusLimit Only valid for ADS-B Airborne Position Messages")
	}
	switch f.messageType {
	case 0, 18, 22:
		err = fmt.Errorf("unknown containment radius")
	case 9, 20:
		radius = 7.5
	case 10, 21:
		radius = 25
	case 11:
		if nicSupplA {
			radius = 75
		} else {
			radius = 185.2
		}
	case 12:
		radius = 370.4
	case 13:
		if 0 == f.nicSupplementB {
			radius = 926
		} else if nicSupplA {
			radius = 1111.2
		} else {
			radius = 555.6
		}
	case 14:
		radius = 1852
	case 15:
		radius = 3704
	case 16:
		if nicSupplA {
			radius = 7408
		} else {
			radius = 14816
		}
	case 17:
		radius = 37040
	default:
		radius = 0
	}

	return radius, err
}

func (f *Frame) NavigationIntegrityCategory(nicSupplA bool) (byte, error) {
	var nic byte
	var err error
	if f.downLinkFormat != 17 {
		return nic, fmt.Errorf("ContainmentRadiusLimit Only valid for ADS-B Airborne Position Messages")
	}
	switch f.messageType {
	case 0, 18, 22:
		err = fmt.Errorf("unknown navigation integrity category")
	case 9, 20:
		nic = 11
	case 10:
	case 21:
		nic = 10
	case 11:
		if nicSupplA {
			nic = 9
		} else {
			nic = 8
		}
	case 12:
		nic = 7
	case 13:
		nic = 6
	case 14:
		nic = 5
	case 15:
		nic = 4
	case 16:
		if nicSupplA {
			nic = 3
		} else {
			nic = 2
		}
	case 17:
		nic = 1
	default:
		nic = 0
	}

	return nic, err
}

/**
 * Gets the air frames size in metres
 */
func (f *Frame) getAirplaneLengthWidth() (float32, float32, error) {
	if !(f.messageType == 31 && f.messageSubType == 1) {
		return 0, 0, fmt.Errorf("can only get aircraft size from ADSB message 31 sub type 1")
	}
	var length, width float32
	var err error

	switch f.airframeWidthLen {
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
		err = fmt.Errorf("unable to determine airframes size")
	}

	return length, width, err
}
