package mode_s

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type featureDescriptionType struct {
	field, meaning string
}

type featureBreakdown struct {
	name       string
	start, end int
	subFields  map[byte][]featureBreakdown
}

var featureDescription = map[string]featureDescriptionType{
	"AA":{field: "Address Announced", meaning: "aircraft identification in All-Call reply - ICAO"},
	"AC":{field: "Altitude Code", meaning: "Aircraft altitude code. All bits are Zeros if altitude information is not available."},
	"AP":{field: "Address/Parity", meaning: "Error detection field. Parity overlaid on the address"},
	"AQ":{field: "Acquisition", meaning: "part of air-to-air protocol"},
	"AB":{field: "Air Speed Bit", meaning:"0=indicated air speed, 1=true air speed"},
	"AS":{field: "True/Indicated Air Speed", meaning:"0=indicated air speed, 1=true air speed"},
	"CA":{field: "Capability", meaning: "aircraft report of system capability"},
	"CC":{field: "Crosslink Capability", meaning:"Indicates XPDR has ability to support crosslink capability"},
	"DF":{field: "Downlink Format", meaning: "downlink descriptor"},
	"DI":{field: "Designator Identification", meaning: "describes content of SD field"},
	"DR":{field: "Downlink Request", meaning: "aircraft requests permission to send data"},
	"FS":{field: "Flight Status", meaning: "aircraft's situation report"},
	"ID":{field: "Identification", meaning: "equivalent to ATCRBS identity number (Squawk)"},
	"II":{field: "Interrogator Identification", meaning: "site number for multisite features"},
	"KE":{field: "Control, ELM", meaning: "part of Extended Length Message protocol"},
	"MA":{field: "Message, Comm-A", meaning: "message to aircraft"},
	"MB":{field: "Message, Comm-B", meaning: "message from aircraft"},
	"MC":{field: "Message, Comm-C", meaning: "long message segment to aircraft"},
	"MD":{field: "Message, Comm-D", meaning: "long message segment from aircraft"},
	"MU":{field: "Message, Comm-U", meaning: "air-to-air message to aircraft"},
	"MV":{field: "Message, Comm-V", meaning: "air-to-air message from aircraft"},
	"NC":{field: "Number, C-segment", meaning: "part of ELM protocol"},
	"ND":{field: "Number, D-segment", meaning: "part of ELM protocol"},
	"PC":{field: "Protocol", meaning: "operating commands for the transponder"},
	"PI":{field: "Parity/Interr.Identity", meaning: "reports source of interrogation. Contains the parity overlaid on the interrogator identity code"},
	"PR":{field: "Probability of Reply", meaning: "used in stochastic acquisition mode"},
	"RC":{field: "Reply Control", meaning: "part of ELM protocol"},
	"RI":{field: "Reply Information", meaning: "aircraft status information for TCAS"},
	"RL":{field: "Reply Length", meaning: "commands air-to-air reply length"},
	"RR":{field: "Reply Request", meaning: "commands details of reply"},
	"SD":{field: "Special Designator", meaning: "control codes to transponder"},
	"SL":{field: "Sensitivity level, ACAS", meaning: "Reports the current operating sensitivity level of TCAS"},
	"SS":{field: "Surveillance Status", meaning: ""},
	"UF":{field: "Uplink Format", meaning: "format descriptor"},
	"UM":{field: "Utility Message", meaning: "protocol message"},
	"VS":{field: "Vertical Status", meaning: "aircraft status, airborne (0) or on the ground (1)"},
	"  ":{field: "Padding", meaning:"Unused"},
	"??":{field: "???", meaning:"Unknown"},
	"CCC":{field: "Capability Class Code", meaning:"Capability Class Code"},
	"OMC":{field: "Operational Mode Code", meaning:"Operational Mode Code"},
	"CRC":{field: "CRC", meaning:"CRC Checksum"},
	"TC":{field:"DF 17 Message Type", meaning:"Message Category"},
	"SUB":{field:"DF 17 Message Sub Type", meaning:"Message Sub Type"},
	"DATA":{field:"ADS-B Data", meaning:"ADS-B Data"},
	"CHAR":{field:"Flight Number", meaning:"1 character of the AIS charset"},
	"TI":{field:"UTC Sync Time Bit", meaning:"Indicates if the Time of Applicability of the message is UTC Sync'd. 0=no"},
	"CPR":{field:"CPR Odd/Even", meaning:"CPR Odd/Even Frame Type"},
	"LAT":{field:"CPR Latitude", meaning:"1 of 4 sets of data required to decode planes lat/lon"},
	"LON":{field:"CPR Longitude", meaning:"1 of 4 sets of data required to decode planes lat/lon"},
	"CAT":{field:"Aircraft Category", meaning:"Category field includes DF field"},
	"MOV":{field:"Movement Field", meaning:"Ground Speed"},
	"HB":{field:"Heading Bit", meaning:"There is a heading available"},
	"HD":{field:"Heading Field", meaning:"The direction the plane is facing"},
	"VR":{field:"Vertical Rate", meaning:"How fast the plane is going up or down"},
	"VRS":{field:"Vertical Rate Sign", meaning:"0=up 1=down"},
	"EWD":{field:"East/West Direction", meaning:"Non-zero == negative velocity. 0=east, 1=west"},
	"EWV":{field:"East/West Velocity", meaning:"How fast the plane is going in the indicated direction"},
	"NSD":{field:"North/South Direction", meaning:"Non-zero == negative velocity. 0=north,1=south"},
	"NSV":{field:"North/South Velocity", meaning:"How fast the plane is going in the indicated direction"},
	"VSRC":{field:"Source Antenna", meaning:"Which antenna this signal was transitted from"},
	"HAED":{field:"Height Above Ellipsoid (HAE) Direction", meaning:"Direction indicator: 1=down, 0=up"},
	"HAEV":{field:"Height Above Ellipsoid (HAE) Delta", meaning:"Barometer offset"},
	"EID":{field:"Emergency ID", meaning:"Emergency Table Lookup ID"},

	"IC":{field:"Intent Change", meaning:"If aircraft wants to change altitude etc"},
	"IFR":{field:"Instrument Flight Rules Capability", meaning:"ADSB v1 Only"},

	"NICp":{field:"Navigation Integrity Category", meaning:""},
	"NICb":{field:"Navigation Integrity Category Supplement B", meaning:""},
	"NACv":{field:"Navigation Accuracy Category", meaning:""},
	"NUC":{field:"Navigation Uncertainty Category", meaning:""},
	"SIL":{field:"Surveillance/Source Integrity Level", meaning:"indicates the probability of exceeding the NIC containment radius"},
	"APLW":{field:"Airplane Width and Length", meaning:""},
	"VER":{field:"ADSB Version", meaning:"This airframes ADSB Compatability"},
	"GVA":{field:"Geometric Vertical Accuracy", meaning:""},

	"NTH":{field:"NIC Altiude|Track/Heading", meaning:"Altudude (sub type 0) or track/heading (sub type 1) have been cross checked with other sources"},
	"HRD":{field:"Heading North Info", meaning:"heading based on 0=true north, 1=magnetic north"},

}

var featureDF17FlightName = []featureBreakdown{
	{name: "CAT", start:37, end: 40},
	{name: "CHAR", start: 40, end: 46},
	{name: "CHAR", start: 46, end: 52},
	{name: "CHAR", start: 52, end: 58},
	{name: "CHAR", start: 58, end: 64},
	{name: "CHAR", start: 64, end: 70},
	{name: "CHAR", start: 70, end: 76},
	{name: "CHAR", start: 76, end: 82},
	{name: "CHAR", start: 82, end: 88},
}
var featureDF17SurfacePosition = []featureBreakdown{
	{name: "MOV", start:37, end: 44},
	{name: "HB", start: 44, end: 45},
	{name: "HD", start: 45, end: 52},
	{name: "??", start: 52, end: 53},
	{name: "CPR", start: 53, end: 54},
	{name: "LAT", start: 54, end: 71},
	{name: "LON", start: 71, end: 88},
}
var featureDF17AirPosition = []featureBreakdown{
	{name: "SS", start:37, end: 39},
	{name: "NICb", start:39, end: 40},
	{name: "AC", start: 40, end: 52},
	{name: "TI", start: 52, end: 53},
	{name: "CPR", start: 53, end: 54},
	{name: "LAT", start: 54, end: 71},
	{name: "LON", start: 71, end: 88},
}

var featureDF17AirVelocityUnknown = []featureBreakdown{
	{name: "SUB", start:37, end: 40},
	{name: "??", start:40, end: 88},
}
var featureDF17AirVelocityGroundSpeed = []featureBreakdown{
	{name: "SUB", start:37, end: 40},
	{name: "IC", start: 40, end: 41},
	{name: "IFR", start: 41, end: 42},
	{name: "NACv", start: 42, end: 45},
	{name: "EWD", start: 45, end: 46},
	{name: "EWV", start: 46, end: 56},
	{name: "NSD", start: 56, end: 57},
	{name: "NSV", start: 57, end: 67},
	{name: "VSRC", start: 67, end: 68},
	{name: "VRS", start: 68, end: 69},
	{name: "VR", start: 69, end: 78},
	{name: "??", start: 78, end: 80},
	{name: "HAED", start: 80, end: 81},
	{name: "HAEV", start: 81, end: 88},
}
var featureDF17AirVelocityAirSpeed = []featureBreakdown{
	{name: "SUB", start:37, end: 40},
	{name: "IC", start: 40, end: 41},
	{name: "IFR", start: 41, end: 42},
	{name: "NACv", start: 42, end: 45},
	{name: "HB", start:45, end: 46},
	{name: "HD", start:46, end: 56},
	{name: "AB", start:56, end: 57},
	{name: "AS", start:57, end: 67},
	{name: "VSRC", start:67, end: 68},
	{name: "VRS", start:68, end: 69},
	{name: "VR", start:69, end: 78},
	{name: "??", start:78, end: 80},
	{name: "HAED", start: 80, end: 81},
	{name: "HAEV", start: 81, end: 88},
}

var featureDF17AirVelocity = []featureBreakdown{
	{name: "??", start:37, end: 88, subFields:map[byte][]featureBreakdown{
		0:featureDF17AirVelocityUnknown,
		1:featureDF17AirVelocityGroundSpeed,
		2:featureDF17AirVelocityGroundSpeed,
		3:featureDF17AirVelocityAirSpeed,
		4:featureDF17AirVelocityAirSpeed,
		5:featureDF17AirVelocityUnknown,
		6:featureDF17AirVelocityUnknown,
		7:featureDF17AirVelocityUnknown,
	},
	},
}

var asdbFeatures = map[byte][]featureBreakdown{
	1: featureDF17FlightName,
	2: featureDF17FlightName,
	3: featureDF17FlightName,
	4: featureDF17FlightName,
	5: featureDF17SurfacePosition,
	6: featureDF17SurfacePosition,
	7: featureDF17SurfacePosition,
	8: featureDF17SurfacePosition,
	9: featureDF17AirPosition,
	10: featureDF17AirPosition,
	11: featureDF17AirPosition,
	12: featureDF17AirPosition,
	13: featureDF17AirPosition,
	14: featureDF17AirPosition,
	15: featureDF17AirPosition,
	16: featureDF17AirPosition,
	17: featureDF17AirPosition,
	18: featureDF17AirPosition,
	19: featureDF17AirVelocity,
	23: {
		{name: "SUB", start: 37, end: 40},
		{name: "ID", start: 40, end: 53},
		{name: "  ", start: 53, end: 88},
	},
	28: {
		{name: "SUB", start: 37, end: 40},
		{name: "??", start: 40, end: 88, subFields: map[byte][]featureBreakdown{
			0: {
				{name: "??", start: 40, end: 88},
			},
			1: { // EMERGENCY (or priority), Status
				{name: "EID", start: 40, end: 43},
				{name: "ID", start: 43, end: 56},
				{name: "  ", start: 56, end: 88},
			},
			2: { // TCAS Resolution Advisory
				{name: "??", start: 40, end: 88},
			},
			3: {
				{name: "??", start: 40, end: 88},
			},
			4: {
				{name: "??", start: 40, end: 88},
			},
			5: {
				{name: "??", start: 40, end: 88},
			},
			6: {
				{name: "??", start: 40, end: 88},
			},
			7: {
				{name: "??", start: 40, end: 88},
			},
		},
		},
	},
	29: {
		{name: "SUB", start: 37, end: 40},
		{name: "??", start: 40, end: 88},
	},
	31: {
		{name: "SUB", start: 37, end: 40},
		{name: "CCC", start: 40, end: 56, subFields: map[byte][]featureBreakdown{
			0: { // airborne
				{name: "CCC", start: 40, end: 56},
			},
			1: { //surface
				{name: "??", start: 40, end: 44},
				{name: "CCC", start: 44, end: 52},
				{name: "APLW", start: 52, end: 56},
			},
		},
		},
		{name: "OMC", start: 56, end: 72},
		{name: "VER", start: 72, end: 75},  //VERSION
		{name: "NICp", start: 75, end: 76}, //Navigation Integrity Category Supplement A
		{name: "NACv", start: 76, end: 80}, //Navigation Accuracy Category Position
		{name: "GVA", start: 80, end: 82},  // geometric_vertical_accuracy
		{name: "SIL", start: 82, end: 84},  // sil
		{name: "NTH", start: 84, end: 85},  //nic_trk_hdg
		{name: "HRD", start: 85, end: 86},  // hrd
		{name: "??", start: 86, end: 88},
	},
}

var frameFeatures = map[byte][]featureBreakdown{

	0: {
		{name: "DF", start: 0, end: 5},
		{name: "VS", start: 5, end: 6},
		{name: "CC", start: 6, end: 7},
		{name: "  ", start: 7, end: 8},
		{name: "SL", start: 8, end: 11},
		{name: "  ", start: 11, end: 13},
		{name: "RI", start: 13, end: 17},
		{name: "  ", start: 17, end: 20},
		{name: "AC", start: 20, end: 32},
		{name: "AP", start: 32, end: 56},
	},
	4: {
		{name: "DF", start: 0, end: 5},
		{name: "FS", start: 5, end: 8},
		{name: "DR", start: 8, end: 13},
		{name: "UM", start: 13, end: 19},
		{name: "AC", start: 19, end: 32},
		{name: "AP", start: 32, end: 56},
	},
	5: {
		{name: "DF", start: 0, end: 5},
		{name: "FS", start: 5, end: 8},
		{name: "DR", start: 8, end: 13},
		{name: "UM", start: 13, end: 19},
		{name: "ID", start: 19, end: 32},
		{name: "AP", start: 32, end: 56},
	},

	11: {
		{name: "DF", start: 0, end: 5},
		{name: "CA", start: 5, end: 8},
		{name: "AA", start: 8, end: 32},
		{name: "PI", start: 32, end: 56},
	},

	16: {
		{name: "DF", start: 0, end: 5},
		{name: "VS", start: 5, end: 6},
		{name: "  ", start: 6, end: 8},
		{name: "SL", start: 8, end: 11},
		{name: "  ", start: 11, end: 13},
		{name: "RI", start: 13, end: 17},
		{name: "  ", start: 17, end: 19},
		{name: "AC", start: 19, end: 32},
		{name: "MV", start: 32, end: 88},
		{name: "AP", start: 88, end: 112},
	},
	17: {
		{name: "DF", start: 0, end: 5},
		{name: "CA", start: 5, end: 8},
		{name: "AA", start: 8, end: 32},
		{name: "TC", start: 32, end: 37},
		{name: "ME", start: 40, end: 88, subFields: asdbFeatures},
		{name: "PI", start: 88, end: 112},
	},
	18: {
		{name: "DF", start: 0, end: 5},
		{name: "??", start: 5, end: 88},
		{name: "CRC", start: 88, end: 112},
	},
	19: {
		{name: "DF", start: 0, end: 5},
	},
	20: {
		{name: "DF", start: 0, end: 5},
		{name: "FS", start: 5, end: 8},
		{name: "DR", start: 8, end: 13},
		{name: "UM", start: 13, end: 19},
		{name: "AC", start: 19, end: 32},
		{name: "MB", start: 32, end: 88},
		{name: "AP", start: 88, end: 112},
	},
	21: {
		{name: "DF", start: 0, end: 5},
		{name: "FS", start: 5, end: 8},
		{name: "DR", start: 8, end: 13}, //
		{name: "UM", start: 13, end: 19},
		{name: "ID", start: 19, end: 32},
		{name: "MB", start: 32, end: 88},
		{name: "AP", start: 88, end: 112},
	},
	24: {
		{name: "DF", start: 0, end: 2},
		{name: "  ", start: 2, end: 3},
		{name: "KE", start: 3, end: 4},
		{name: "ND", start: 4, end: 8},
		{name: "MD", start: 8, end: 88},
		{name: "AP", start: 88, end: 112},
	},
}

func (f *Frame) Describe(output io.Writer) {
	fprintf(output, "MODE S Packet:\n")
	fprintf(output, "Length              : %d bits\n", f.getMessageLengthBits())
	fprintf(output, "Frame               : %s\n", f.raw)
	fprintf(output, "DF: Downlink Format : (%d) %s\n", f.downLinkFormat, f.DownLinkFormat())
	// decode the specific DF type
	switch f.downLinkFormat {
	case 0:
		f.showVerticalStatus(output)
		f.showCrossLinkCapability(output)
		f.showSensitivityLevel(output)
		f.showReplyInformation(output)
		f.showAltitude(output)
	case 4:
		f.showFlightStatus(output)
		f.showDownLinkRequest(output)
		f.showUtilityMessage(output)
		f.showAltitude(output)
	case 5:
		f.showFlightStatus(output)
		f.showDownLinkRequest(output)
		f.showUtilityMessage(output)
		f.showIdentity(output)
	case 11:
		f.showCapability(output)
		f.showICAO(output)
	case 16:
		f.showVerticalStatus(output)
		f.showSensitivityLevel(output)
		f.showReplyInformation(output)
		f.showAltitude(output)
	case 17:
		f.showCapability(output)
		f.showICAO(output)
		f.showAdsb(output)
	case 18: //DF_18
		//f.showCapability() // control field
		if 0 == f.ca {
			f.showCapability(output)
			f.showICAO(output)
			f.showAdsb(output)
		} else {
			fprintln(output, "Unable to decode DF18 Capability:", f.ca)
		}
	case 20: //DF_20
		f.showFlightStatus(output)
		f.showAltitude(output)
		f.showFlightNumber(output)
		f.showBdsData(output)
	case 21: //DF_21
		f.showFlightStatus(output)
		f.showIdentity(output) // gillham encoded squawk
		f.showFlightNumber(output)
		f.showBdsData(output)
	}

	f.showBitString(output)

}

func (f *Frame) showVerticalStatus(output io.Writer) {
	if !f.VerticalStatusValid() {
		return
	}
	if f.onGround {
		fprintln(output, "VS: Vertical Status : On The Ground")
	} else {
		fprintln(output, "VS: Vertical Status : Airborne")
	}
}

func (f *Frame) showVerticalRate(output io.Writer) {
	if f.validVerticalRate {
		fprintf(output, "  Vertical Rate     : %d\n", f.verticalRate)
	} else {
		fprintln(output, "  Vertical Rate     : Invalid")
		fprintln(output, "")
	}
}

func (f *Frame) showCrossLinkCapability(output io.Writer) {
	fprintf(output, "CC: CrossLink Cap   : %d\n", f.cc)
}

func (f *Frame) showAltitude(output io.Writer) {
	if f.validAltitude {
		fprintf(output, "AC: Altitude        : %d %s (q bit: %t, m bit: %t)\n", f.altitude, f.AltitudeUnits(), f.acQ, f.acM)
	} else {
		fprintln(output, "AC: Altitude        : Invalid")
	}
}

func (f *Frame) showContainmentRadius(output io.Writer) {
	r, err := f.ContainmentRadiusLimit(true)
	if nil != err {
		fprintf(output, "  Containment Radius: %s\n", err)
	} else {
		fprintf(output, "  Containment Radius: %0.2f metres\n", r)
	}
}

func (f *Frame) showSurveilanceStatus(output io.Writer) {
	fprintf(output, "  Surveillance      : (Status:%d) %s\n", f.surveillanceStatus, surveillanceStatus[f.surveillanceStatus])
}

func (f *Frame) showNavigationIntegrity(output io.Writer) {
	fprintf(output, "  NIC Supplement B  : %d\n", f.nicSupplementB)
	nic, err := f.NavigationIntegrityCategory(true)
	if nil != err {
		fprintf(output, "  Nav Integrity     : %s\n", err)
	} else {
		fprintf(output, "  Nav Integrity     : %d\n", nic)
	}
}

func (f *Frame) showFlightStatus(output io.Writer) {
	fprintf(output, "FS: Flight Status   : (%d) %s\n", f.fs, flightStatusTable[f.fs])
	if "" != f.special {
		fprintf(output, "FS: Special Status  : %s\n", f.special)
	}
	f.showAlert(output)
	f.showVerticalStatus(output)
}

func (f *Frame) showFlightId(output io.Writer) {
	fprintf(output, "Flight          : %s", f.Flight())
	fprintln(output, "")
}

func (f *Frame) showICAO(output io.Writer) {
	fprintf(output, "AA: ICAO            : %6X", f.icao)
	fprintln(output, "")
}

func (f *Frame) showCapability(output io.Writer) {
	fprintf(output, "CA: Plane Mode S Cap: (%d) %s\n", f.ca, capabilityTable[f.ca])
	f.showVerticalStatus(output)
}

func (f *Frame) showIdentity(output io.Writer) {
	fprintf(output, "ID: Squawk Identity : %04d\n", f.identity)
}

func (f *Frame) showDownLinkRequest(output io.Writer) {
	fprintf(output, "DR: Downlink Request: (%d) %s\n", f.dr, downlinkRequestField[f.dr])
}

func (f *Frame) showUtilityMessage(output io.Writer) {
	fprintf(output, "UM: Utility Request : (%d) %s\n", f.um, utilityMessageField[f.um])
}

func (f *Frame) showHae(output io.Writer) {
	if f.validHae {
		fprintf(output, "  HAE Delta         : %d (Height Above Ellipsoid)\n", f.haeDelta)
	}else {
		fprintln(output, "  HAE Delta         : Unavailable")

	}
}
func (f *Frame) showVelocity(output io.Writer) {
	if f.superSonic {
		fprintln(output, "  Super Sonic?      : Yes!")
	} else {
		fprintln(output, "  Super Sonic?      : No")
	}
	if f.validVelocity {
		fprintf(output, "  Velocity          : %0.2f\n", f.velocity)
		fprintf(output, "  EW/NS VEL         : (East/west: %d) (North/South: %d)\n", f.eastWestVelocity, f.northSouthVelocity)
	} else {
		fprintln(output, "  Velocity          : Invalid")
	}
}

func (f *Frame) showHeading(output io.Writer) {
	if f.validHeading {
		fprintf(output, "  Heading           : %0.2f\n", f.heading)
	} else {
		fprintln(output, "  Heading           : Not Valid")
		fprintln(output, "")

	}
}

func (f *Frame) showIntentChange(output io.Writer) {
	fprintf(output, "  Intent Change     : %t\n", f.intentChange != 0)
}
func (f *Frame) showInstrumentFlightRulesCapability(output io.Writer) {
	fprintf(output, "  IFR Capable       : %t\n", f.ifrCapability != 0)
}

func (f *Frame) showNavAccuracyCat(output io.Writer) {
	if f.validNacV {
		fprintf(output, "  Nav Accuracy Cat  : %d\n", f.nacV)
	}
}

func (f *Frame) showCprLatLon(output io.Writer) {
	fprintln(output, "Before Decoding : Half of vehicle location")
	var oddEven = "Odd"
	if f.IsEven() {
		oddEven = "Even"
	}
	fprintf(output, "  UTC Sync?     : %t\n", f.timeFlag != 0)
	fprintf(output, "  CPR Frame     : %s\n", oddEven)
	fprintf(output, "  CPR Latitude  : %d\n", f.rawLatitude)
	fprintf(output, "  CPR Longitude : %d\n", f.rawLongitude)
	fprintln(output, "")
}

func (f *Frame)showReplyInformation(output io.Writer) {
	fprintf(output, "RI: TCAS            : (%d) %s\n", f.ri, replyInformationField[f.ri])
}
func (f *Frame)showSensitivityLevel(output io.Writer) {
	fprintf(output, "SL: TCAS            : (%d) %s\n", f.sl, sensitivityLevelInformationField[f.sl])
}

func (f *Frame) showCategory(output io.Writer) {
	if f.ValidCategory() {
		fprintf(output, "CAT: Aircraft Cat   : (%d:%d) %s\n", f.catType, f.catSubType, f.Category())
	}
}

func (f *Frame) showAdsb(output io.Writer) {
	fprintf(output, "ME : ADSB Msg Type  : (%d) %s\n", f.messageType, f.MessageTypeString())

	switch f.messageType {
	case 1, 2, 3, 4:
		f.showCategory(output)
		f.showFlightNumber(output)
	case 5, 6, 7, 8:
		f.showVelocity(output)
		f.showHeading(output)
		f.showCprLatLon(output)
	case 9, 10, 11, 12, 13, 14, 15, 16, 17, 18:
		f.showContainmentRadius(output)
		f.showSurveilanceStatus(output)
		f.showNavigationIntegrity(output)
		f.showAltitude(output)
		f.showCprLatLon(output)
	case 19:
		f.showAdsbMsgSubType(output)
		switch f.messageSubType {
		case 1, 2, 3, 4:
			f.showIntentChange(output)
			f.showInstrumentFlightRulesCapability(output)
			f.showNavAccuracyCat(output)
			f.showHeading(output)
			f.showVelocity(output)
			f.showVerticalRate(output)
		default:
		// unknown sub type
		}
		f.showHae(output)
	case 23:
		f.showAdsbMsgSubType(output)
		if 7 == f.messageSubType {
			f.showIdentity(output)
		}
	case 28:
		f.showAdsbMsgSubType(output)
		if 1 == f.messageSubType {
			f.showIdentity(output)
			f.showAlert(output)
		} else if 2 == f.messageSubType {
			// TCAS RA
		}
	case 29:
	case 31:
		f.showAdsbMsgSubType(output)
		f.showCapabilityClassInfo(output)
		f.showVerticalStatus(output)
		f.showOperationalModeInfo(output)
		f.showAircraftLengthWidth(output)
		f.showAdsbVersion(output)
		f.showNavAccuracyCat(output)
		f.showCrossCheck(output)
		f.showCompassNorth(output)
	default:
		fprintln(output, "Packet Type Not Yet Decoded")
	}

	fprintln(output, "")
}

func (f *Frame) showAdsbMsgSubType(output io.Writer) {
	fprintf(output, "SUB:      Sub Type  : %d \n", f.messageSubType)
}

func (f *Frame) showCapabilityClassInfo(output io.Writer) {
	if f.validCompatibilityClass {
		if nil != f.cccHasLowTxPower {
			fprintf(output, "  Low TX Power      : %t\n", *f.cccHasLowTxPower)
		}
	} else {
		fprintf(output, "Compatibility Class : Unknown\n")
	}
}
func (f *Frame) showOperationalModeInfo(output io.Writer) {

}
func (f *Frame) showAircraftLengthWidth(output io.Writer) {
	length, width, err := f.getAirplaneLengthWidth()
	if nil == err {
		fprintf(output, "    Airframe Size   : width:%0.1f length:%0.1f metres\n", width, length)
	}
}
func (f *Frame) showCrossCheck(output io.Writer) {
	if f.messageSubType == 0 {
		fprintf(output, "NIC Baro CrossCheck : %t\n", f.nicCrossCheck==1)
	} else if f.messageSubType == 1 {
		fprintf(output, "NIC Track CrossCheck: %t\n", f.nicCrossCheck==1)
	}
}
func (f *Frame) showCompassNorth(output io.Writer) {
	if f.northReference != 0 {
		fprintf(output, "  Compass Heading   : Magnetic North\n")
	} else {
		fprintf(output, "  Compass Heading   : True North\n")
	}
}

func (f *Frame) showAlert(output io.Writer) {
	if f.alert {
		fprintf(output, "Plane showing Alert!\n")
	}
	f.showSpecial(output)
}
func (f *Frame) showSpecial(output io.Writer) {
	if "" != f.special {
		fprintf(output, "  Special           : %s\n", f.special)
	}
}

func (f *Frame) showFlightNumber(output io.Writer) {
	fprintf(output, "    Flight Number   : %s\n", f.FlightNumber())
}

// determines what type of mode S Message this frame is
func (f *Frame) DownLinkFormat() string {

	if description, ok := downlinkFormatTable[f.downLinkFormat]; ok {
		return description
	}
	return "Unknown Downlink Format"
}

func (f *Frame)showAdsbVersion(output io.Writer) {
	fprintf(output, "    ADS-B Version   : (%d) %s\n", f.adsbVersion, adsbCompatibilityVersion[f.adsbVersion])
}

func (f *Frame) showBdsData(output io.Writer) {
	fprintln(output, "BDS Info")
	fprintf(output, "  BDS Msg       : %s\n", f.DescribeBds())
}

func (f *Frame) showBitString(output io.Writer) {
	if features, ok := frameFeatures[f.downLinkFormat]; ok {
		fprintln(output, f.formatBitString(features))
	}
}

func (f *Frame) formatBitString(features []featureBreakdown) string {
	var header, separator, bits, rawBits, bitFmt, bitDesc, footer, suffix string
	var padLen, realLen, shownBitCount, i int

	for _, i := range f.message {
		rawBits += fmt.Sprintf("%08s", strconv.FormatUint(uint64(i), 2))
	}

	doMakeBitString := func(f featureBreakdown) {
		padLen = len(f.name)
		realLen = f.end - f.start
		if realLen > padLen {
			padLen = realLen
		}
		shownBitCount += f.end - f.start
		bitFmt = fmt.Sprintf(" %%- %ds |", padLen)
		header += fmt.Sprintf(bitFmt, f.name)
		separator += strings.Repeat("-", padLen + 2) + "+"
		//bits += fmt.Sprintf(bitFmt, rawBits[f.start: f.end])
		bits += " "
		for i = f.start; i < f.end; i++ {
			if i % 8 == 0 {
				bits += "<span class='byte-start'>" + string(rawBits[i]) + "</span>"
			} else {
				bits += string(rawBits[i])
			}
		}
		bits += strings.Repeat(" ", padLen - (f.end - f.start) + 1) + "|"
		bitDesc += fmt.Sprintf(bitFmt, strconv.Itoa(f.start))

		if 1 == realLen {
			suffix = ""
		} else {
			suffix = "s"
		}

		feature := featureDescription[f.name]
		footer += fmt.Sprintf(" %s \t%d bit%s\t %s: %s\n", f.name, realLen, suffix, feature.field, feature.meaning)
	}

	for _, feat := range features {
		if 0 == len(feat.subFields[f.messageType]) {
			doMakeBitString(feat)
		} else {
			for _, sf := range feat.subFields[f.messageType] {
				if 0 == len(sf.subFields[f.messageSubType]) {
					doMakeBitString(sf)
				} else {
					for _, ssf := range sf.subFields[f.messageSubType] {
						doMakeBitString(ssf)
					}
				}
			}
		}
	}

	return fmt.Sprintf("\n%s\n%s\n%s\n%s\n%s\n\n%s\n%d/%d bits shown\n", header, separator, bits, separator, bitDesc, footer, shownBitCount, f.getMessageLengthBits())
}

func fprintf(output io.Writer, line string, params ...interface{}) {
	if "\n" != line[len(line)-1:1] {
		line += "\n"
	}
	_, _ = fmt.Fprintf(output, line, params...)
}

func fprintln(output io.Writer, params ...interface{}) {
	_, _ = fmt.Fprintln(output, params...)
}
