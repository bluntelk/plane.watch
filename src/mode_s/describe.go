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
	"TI":{field:"Time Bit", meaning:"UTC Time"},
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
	"SRC":{field:"Source Antenna", meaning:"Which antenna this signal was transitted from"},
	"HAED":{field:"Height Above Ellipsoid (HAE) Direction", meaning:"Direction indicator: 1=down, 0=up"},
	"HAEV":{field:"Height Above Ellipsoid (HAE) Delta", meaning:"Barometer offset"},
	"EID":{field:"Emergency ID", meaning:"Emergency Table Lookup ID"},

	"NICp":{field:"Navigation Integrity Category", meaning:""},
	"NACv":{field:"Navigation Accuracy Category", meaning:""},
	"NUC":{field:"Navigation Uncertainty Category", meaning:""},
	"SIL":{field:"Surveillance/Source Integrity Level", meaning:""},
	"APLW":{field:"Airplane Width and Length", meaning:""},
	"VER":{field:"ADSB Version", meaning:"This airframes ADSB Compatability"},
	"GVA":{field:"Geometric Vertical Accuracy", meaning:""},
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
	{name: "SUB", start:37, end: 40},
	{name: "AC", start: 40, end: 52},
	{name: "TI", start: 52, end: 53},
	{name: "CPR", start: 53, end: 54},
	{name: "LAT", start: 54, end: 71},
	{name: "LON", start: 71, end: 88},
}
var featureDF17AirVelocity = []featureBreakdown{
	{name: "SUB", start:37, end: 40},
	{name: "??", start: 40, end: 45},
	{name: "EWD", start: 45, end: 46},
	{name: "EWV", start: 46, end: 56},
	{name: "NSD", start: 56, end: 57},
	{name: "NSV", start: 57, end: 67},
	{name: "SRC", start: 67, end: 68},
	{name: "VRS", start: 68, end: 69},
	{name: "VR", start: 69, end: 78},
	{name: "??", start: 78, end: 80},
	{name: "HAED", start: 80, end: 81},
	{name: "HAEV", start: 81, end: 88},
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
	23: []featureBreakdown{
		{name: "SUB", start:37, end: 40},
		{name: "ID", start: 40, end: 53},
		{name: "  ", start: 53, end: 88},
	},
	28: []featureBreakdown{
		{name: "SUB", start:37, end: 40},
		{name: "??", start: 40, end: 88, subFields:map[byte][]featureBreakdown{
			0:[]featureBreakdown{
				{name: "??", start: 40, end: 88},
			},
			1:[]featureBreakdown{// EMERGENCY (or priority), Status
				{name: "EID", start: 40, end: 43},
				{name: "ID", start: 43, end: 56},
				{name: "  ", start: 56, end: 88},
			},
			2:[]featureBreakdown{// TCAS Resolution Advisory
				{name: "??", start: 40, end: 88},
			},
		},
		},
	},
	29: []featureBreakdown{
		{name: "SUB", start:37, end: 40},
		{name: "??", start: 40, end: 88},
	},
	31: []featureBreakdown{
		{name: "SUB", start:37, end: 40},
		{name: "CCC", start: 40, end: 56, subFields:map[byte][]featureBreakdown{
			0:[]featureBreakdown{ // airborne
				{name: "CCC", start: 40, end: 56},
			},
			1:[]featureBreakdown{ //surface
				{name: "??", start: 40, end: 44},
				{name: "CCC", start: 44, end: 52},
				{name: "APLW", start: 52, end: 56},
			},
		},
		},
		{name: "OMC", start: 56, end: 72},
		{name: "VER", start: 72, end: 75}, //VERSION
		{name: "NICp", start: 75, end: 76}, //nic_suppl - Navigation Integrity Category
		{name: "NACv", start: 76, end: 80}, //nac_pos
		{name: "GVA", start: 80, end: 82}, // geometric_vertical_accuracy
		{name: "SIL", start: 82, end: 84}, // sil
		{name: "??", start: 84, end: 85}, //nic_trk_hdg
		{name: "??", start: 85, end: 86}, // hrd
		{name: "??", start: 86, end: 88},
	},
}

var frameFeatures = map[byte][]featureBreakdown{

	0: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "VS", start:5, end: 6},
		{name: "CC", start:6, end: 7},
		{name: "  ", start:7, end: 8},
		{name: "SL", start:8, end: 11},
		{name: "  ", start:11, end: 13},
		{name: "RI", start:13, end: 17},
		{name: "  ", start:17, end: 20},
		{name: "AC", start:20, end: 32},
		{name: "AP", start:32, end: 56},
	},
	4: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "FS", start:5, end: 8},
		{name: "DR", start:8, end: 13},
		{name: "UM", start:13, end: 19},
		{name: "AC", start:19, end: 32},
		{name: "AP", start:32, end: 56},
	},
	5: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "FS", start:5, end: 8},
		{name: "DR", start:8, end: 13},
		{name: "UM", start:13, end: 19},
		{name: "ID", start:19, end: 32},
		{name: "AP", start:32, end: 56},
	},

	11: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "CA", start:5, end: 8},
		{name: "AA", start:8, end: 32},
		{name: "PI", start:32, end: 56},
	},

	16: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "VS", start:5, end: 6},
		{name: "  ", start:6, end: 8},
		{name: "SL", start:8, end: 11},
		{name: "  ", start:11, end: 13},
		{name: "RI", start:13, end: 17},
		{name: "  ", start:17, end: 19},
		{name: "AC", start:19, end: 32},
		{name: "MV", start:32, end: 88},
		{name: "AP", start:88, end: 112},
	},
	17: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "CA", start:5, end: 8},
		{name: "AA", start:8, end: 32},
		{name: "TC", start:32, end: 37},
		{name: "ME", start:40, end: 88, subFields: asdbFeatures},
		{name: "PI", start:88, end: 112},
	},
	18: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "??", start:5, end: 88},
		{name: "CRC", start:88, end: 112},
	},
	19: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	20: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "FS", start:5, end: 8},
		{name: "DR", start:8, end: 13},
		{name: "UM", start:13, end: 19},
		{name: "AC", start:19, end: 32},
		{name: "MB", start:32, end: 88},
		{name: "AP", start:88, end: 112},
	},
	21: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "FS", start:5, end: 8},
		{name: "DR", start:8, end: 13}, //
		{name: "UM", start:13, end: 19},
		{name: "ID", start:19, end: 32},
		{name: "MB", start:32, end: 88},
		{name: "AP", start:88, end: 112},
	},
	24: []featureBreakdown{
		{name: "DF", start:0, end: 2},
		{name: "  ", start:2, end: 3},
		{name: "KE", start:3, end: 4},
		{name: "ND", start:4, end: 8},
		{name: "MD", start:8, end: 88},
		{name: "AP", start:88, end: 112},
	},
}

func (frame *Frame) Describe(output io.Writer) {
	fmt.Fprintf(output, "MODE S Packet:\n")
	fmt.Fprintf(output, "Length              : %d bits\n", frame.getMessageLengthBits())
	fmt.Fprintf(output, "Frame               : %s\n", frame.raw)
	fmt.Fprintf(output, "DF: Downlink Format : %d (%s)\n", frame.downLinkFormat, frame.DownLinkFormat())
	// decode the specific DF type
	switch frame.downLinkFormat {
	case 0:
		frame.showVerticalStatus(output)
		frame.showCrossLinkCapability(output)
		frame.showSensitivityLevel(output)
		frame.showReplyInformation(output)
		frame.showAltitude(output)
	case 4:
		frame.showFlightStatus(output)
		frame.showDownLinkRequest(output)
		frame.showUtilityMessage(output)
		frame.showAltitude(output)
	case 5:
		frame.showFlightStatus(output)
		frame.showDownLinkRequest(output)
		frame.showUtilityMessage(output)
		frame.showIdentity(output)
	case 11:
		frame.showCapability(output)
		frame.showICAO(output)
	case 16:
		frame.showVerticalStatus(output)
		frame.showSensitivityLevel(output)
		frame.showReplyInformation(output)
		frame.showAltitude(output)
	case 17:
		frame.showCapability(output)
		frame.showICAO(output)
		frame.showAdsb(output)
	case 18: //DF_18
		//frame.showCapability() // control field
		if 0 == frame.ca {
			frame.showCapability(output)
			frame.showICAO(output)
			frame.showAdsb(output)
		} else {
			fmt.Fprintln(output, "Unable to decode DF18 Capability:", frame.ca)
		}
	case 20: //DF_20
		frame.showFlightStatus(output)
		frame.showAltitude(output)
		frame.showFlightNumber(output)
		frame.showBdsData(output)
	case 21: //DF_21
		frame.showFlightStatus(output)
		frame.showIdentity(output) // gillham encoded squawk
		frame.showFlightNumber(output)
		frame.showBdsData(output)
	}

	frame.showBitString(output)

}

func (f *Frame) showVerticalStatus(output io.Writer) {
	if !f.VerticalStatusValid() {
		return
	}
	if f.onGround {
		fmt.Fprintln(output, "VS: Vertical Status : On The Ground");
	} else {
		fmt.Fprintln(output, "VS: Vertical Status : Airborne");
	}
}

func (f *Frame) showVerticalRate(output io.Writer) {
	if f.validVerticalRate {
		fmt.Fprintf(output, "  Vertical Rate     : %d\n", f.verticalRate)
	} else {
		fmt.Fprintln(output, "  Vertical Rate     : Invalid\n")
	}
}

func (f *Frame) showCrossLinkCapability(output io.Writer) {
	fmt.Fprintf(output, "CC: CrossLink Cap   : %d\n", f.cc)
}

func (f *Frame) showAltitude(output io.Writer) {
	if f.validAltitude {
		fmt.Fprintf(output, "AC: Altitude        : %d %s (q bit: %t, m bit: %t)\n", f.altitude, f.AltitudeUnits(), f.ac_q, f.ac_m)
	} else {
		fmt.Fprintln(output, "AC: Altitude        : Invalid")
	}
}

func (f *Frame) showFlightStatus(output io.Writer) {
	fmt.Fprintf(output, "FS: Flight Status   : (%d) %s\n", f.fs, flightStatusTable[f.fs])
	if "" != f.special {
		fmt.Fprintf(output, "FS: Special Status  : %s\n", f.special)
	}
	f.showAlert(output)
	f.showVerticalStatus(output)
}

func (f *Frame) showFlightId(output io.Writer) {
	fmt.Fprintf(output, "Flight          : %s", f.Flight())
	fmt.Fprintln(output, "")
}

func (f *Frame) showICAO(output io.Writer) {
	fmt.Fprintf(output, "AA: ICAO            : %6X", f.icao)
	fmt.Fprintln(output, "")
}

func (f *Frame) showCapability(output io.Writer) {
	fmt.Fprintf(output, "CA: Plane Mode S Cap: (%d) %s\n", f.ca, capabilityTable[f.ca])
	f.showVerticalStatus(output)
}

func (f *Frame) showIdentity(output io.Writer) {
	fmt.Fprintf(output, "ID: Squawk Identity : %04d\n", f.identity)
}

func (f *Frame) showDownLinkRequest(output io.Writer) {
	fmt.Fprintf(output, "DR: Downlink Request: (%d) %s\n", f.dr, downlinkRequestField[f.dr])
}

func (f *Frame) showUtilityMessage(output io.Writer) {
	fmt.Fprintf(output, "UM: Utility Request : (%d) %s\n", f.um, utilityMessageField[f.um])
}

func (f *Frame) showHae(output io.Writer) {
	if f.validHae {
		fmt.Fprintf(output, "  HAE Delta         : %d (Height Above Ellipsoid)\n", f.haeDelta)
	}else {
		fmt.Fprintln(output, "  HAE Delta         : Unavailable")

	}
}
func (f *Frame) showVelocity(output io.Writer) {
	if f.superSonic {
		fmt.Fprintln(output, "  Super Sonic?      : Yes!")
	} else {
		fmt.Fprintln(output, "  Super Sonic?      : No")
	}
	if f.validVelocity {
		fmt.Fprintf(output, "  Velocity          : %0.2f\n", f.velocity)
		fmt.Fprintf(output, "  EW/NS VEL         : (East/west: %d) (North/South: %d)\n", f.eastWestVelocity, f.northSouthVelocity)
	} else {
		fmt.Fprintln(output, "  Velocity          : Invalid")
	}
}

func (f *Frame) showHeading(output io.Writer) {
	if f.validHeading {
		fmt.Fprintf(output, "  Heading           : %0.2f\n", f.heading)
	} else {
		fmt.Fprintln(output, "  Heading           : Not Valid\n")
	}
}

func (f *Frame) showCprLatLon(output io.Writer) {
	fmt.Fprintln(output, "Before Decoding : Half of vehicle location")
	var oddEven string = "Odd"
	if f.IsEven() {
		oddEven = "Even"
	}
	fmt.Fprintf(output, "  CPR Frame     : %s\n", oddEven)
	fmt.Fprintf(output, "  CPR Latitude  : %d\n", f.rawLatitude)
	fmt.Fprintf(output, "  CPR Longitude : %d\n", f.rawLongitude)
	fmt.Fprintln(output, "")
}

func (f *Frame)showReplyInformation(output io.Writer) {
	fmt.Fprintf(output, "RI: TCAS            : (%d) %s\n", f.ri, replyInformationField[f.ri])
}
func (f *Frame)showSensitivityLevel(output io.Writer) {
	fmt.Fprintf(output, "SL: TCAS            : (%d) %s\n", f.sl, sensitivityLevelInformationField[f.sl])
}

func (f *Frame) showCategory(output io.Writer) {
	if f.ValidCategory() {
		fmt.Fprintf(output, "Aircraft Cat    : (%d:%d) %s\n", f.catType, f.catSubType, f.Category())
		fmt.Fprintf(output, "Aircraft Cat    : (%d) (second calc)\n", f.aircraftType)
	}
}

func (f *Frame) showAdsb(output io.Writer) {
	fmt.Fprintf(output, "ME: ADSB Msg Type   : %d (Sub Type %d): %s\n", f.messageType, f.messageSubType, f.MessageTypeString())

	switch f.messageType {
	case 1, 2, 3, 4:
		f.showFlightNumber(output)
		f.showCategory(output)
	case 5, 6, 7, 8:
		f.showHeading(output)
		f.showVelocity(output)
		f.showCprLatLon(output)
	case 9, 10, 11, 12, 13, 14, 15, 16, 17, 18:
		f.showAltitude(output)
		f.showCprLatLon(output)
	case 19:
		switch f.messageSubType {
		case 1, 2, 3, 4:
			f.showHeading(output)
			f.showVerticalRate(output)
			f.showVelocity(output)
		default:
			// unknown sub type
		}
		f.showHae(output)
	case 23:
		if 7 == f.messageSubType {
			f.showIdentity(output);
		}
	case 28:
		if 1 == f.messageSubType {
			f.showIdentity(output);
			f.showAlert(output);
		}
	case 29:
	case 31:
		f.showVerticalStatus(output)
		f.showAdsbVersion(output)
	default:
		fmt.Fprintln(output, "Packet Type Not Yet Decoded")
	}

	fmt.Fprintln(output, "")
}

func (f *Frame) showAlert(output io.Writer) {
	if f.alert {
		fmt.Fprintf(output, "Plane showing Alert!\n")
	}
	f.showSpecial(output)
}
func (f *Frame) showSpecial(output io.Writer) {
	if "" != f.special {
		fmt.Fprintf(output, "  Special           : %s\n", f.special)
	}
}

func (f *Frame) showFlightNumber(output io.Writer) {
	fmt.Fprintf(output, "Flight Number   : %s\n", f.FlightNumber())
}

// determines what type of mode S Message this frame is
func (f *Frame) DownLinkFormat() string {

	if description, ok := downlinkFormatTable[f.downLinkFormat]; ok {
		return description
	}
	return "Unknown Downlink Format"
}

func (f *Frame)showAdsbVersion(output io.Writer) {
	fmt.Fprintf(output, "    ADS-B Version   : (%d) %s\n", f.adsbVersion, adsbCompatibilityVersion[f.adsbVersion])
}

func (f *Frame) showBdsData(output io.Writer) {
	fmt.Fprintln(output, "BDS Info")
	fmt.Fprintf(output, "  BDS Msg       : %s\n", f.DescribeBds())
}

func (f *Frame) showBitString(output io.Writer) {
	if features, ok := frameFeatures[f.downLinkFormat]; ok {
		fmt.Fprintln(output, f.formatBitString(features))
	}
}

func (frame *Frame) formatBitString(features []featureBreakdown) string {
	var header, separator, bits, rawBits, bitFmt, bitDesc, footer, suffix string
	var padLen, realLen, shownBitCount, i int

	for _, i := range frame.message {
		rawBits += fmt.Sprintf("%08s", strconv.FormatUint(uint64(i), 2))
	}

	doMakeBitString := func(f featureBreakdown) {
		padLen = len(f.name)
		realLen = f.end - f.start
		if realLen > padLen {
			padLen = realLen
		}
		shownBitCount += (f.end - f.start)
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

	for _, f := range features {
		if 0 == len(f.subFields[frame.messageType]) {
			doMakeBitString(f)
		} else {
			for _, sf := range f.subFields[frame.messageType] {
				if 0 == len(sf.subFields[frame.messageSubType]) {
					doMakeBitString(sf)
				} else {
					for _, ssf := range sf.subFields[frame.messageSubType] {
						doMakeBitString(ssf)
					}
				}
			}
		}
	}

	return fmt.Sprintf("\n%s\n%s\n%s\n%s\n%s\n\n%s\n%d/%d bits shown\n", header, separator, bits, separator, bitDesc, footer, shownBitCount, frame.getMessageLengthBits())
}