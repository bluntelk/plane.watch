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
	"??":{field: "???", meaning:"Unknown"},
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
	{name: "??", start: 40, end: 69},
	{name: "VRS", start: 68, end: 69},
	{name: "VR", start: 69, end: 78},
	{name: "??", start: 78, end: 88},
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
		{name: "??", start: 40, end: 88},
	},
	28: []featureBreakdown{
		{name: "??", start: 40, end: 88},
	},

}

var frameFeatures = map[byte][]featureBreakdown{

	0: []featureBreakdown{ // DF, VS, CC, SL, RI, AC, AP
		{name: "DF", start:0, end: 5},
		{name: "VS", start:5, end: 6},
		{name: "CC", start:6, end: 7},
		{name: "??", start:7, end: 8},
		{name: "SL", start:8, end: 11},
		{name: "??", start:11, end: 13},
		{name: "RI", start:13, end: 17},
		{name: "??", start:17, end: 20},
		{name: "AC", start:20, end: 32},
		{name: "AP", start:32, end: 56},
	},
	//1: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//	{name: "FS", start:5, end: 8},
	//	{name: "??", start:8, end: 32},
	//	{name: "AP", start:32, end: 56},
	//},
	//2: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//	{name: "FS", start:5, end: 8},
	//	{name: "??", start:8, end: 32},
	//	{name: "AP", start:32, end: 56},
	//},
	//3: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//	{name: "FS", start:5, end: 8},
	//	{name: "??", start:8, end: 32},
	//	{name: "AP", start:32, end: 56},
	//},
	4: []featureBreakdown{ // DF, FS, DR, UM, AC, AP
		{name: "DF", start:0, end: 5},
		{name: "FS", start:5, end: 8},
		{name: "DR", start:8, end: 13},
		{name: "UM", start:13, end: 19},
		{name: "AC", start:19, end: 32},
		{name: "AP", start:32, end: 56},
	},
	5: []featureBreakdown{ // DF, FS, DR, UM, ID, AP
		{name: "DF", start:0, end: 5},
		{name: "FS", start:5, end: 8},
		{name: "DR", start:8, end: 13},
		{name: "UM", start:13, end: 19},
		{name: "ID", start:19, end: 32},
		{name: "AP", start:32, end: 56},
	},
	//6: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//	{name: "AP", start:32, end: 56},
	//},
	//7: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//	{name: "AP", start:32, end: 56},
	//},
	//8: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//	{name: "AP", start:32, end: 56},
	//},
	//9: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//	{name: "AP", start:32, end: 56},
	//},
	//10: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//	{name: "AP", start:32, end: 56},
	//},
	11: []featureBreakdown{ // DF, CA, AA, PI
		{name: "DF", start:0, end: 5},
		{name: "CA", start:5, end: 8},
		{name: "AA", start:8, end: 32},
		{name: "PI", start:32, end: 56},
	},
	//12: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//	{name: "AP", start:32, end: 56},
	//},
	//13: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//	{name: "AP", start:32, end: 56},
	//},
	//14: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//	{name: "AP", start:32, end: 56},
	//},
	//15: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//	{name: "AP", start:32, end: 56},
	//},
	16: []featureBreakdown{//DF(5), VS(1), SL(3), RI(4), AC(13), MV (56), AP(24) ?? 106
		{name: "DF", start:0, end: 5},
		{name: "VS", start:5, end: 6},
		{name: "??", start:6, end: 8},
		{name: "SL", start:8, end: 11},
		{name: "??", start:11, end: 13},
		{name: "RI", start:13, end: 17},
		{name: "??", start:17, end: 19},
		{name: "AC", start:19, end: 32},
		{name: "MV", start:32, end: 88},
		{name: "AP", start:88, end: 112},
	},
	17: []featureBreakdown{ // DF, CA, AA, ME, PI
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
		{name: "??", start:5, end: 88},
		{name: "CRC", start:88, end: 112},
	},
	21: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "??", start:5, end: 88},
		{name: "CRC", start:88, end: 112},
	},
	//22: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//},
	//23: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//},
	24: []featureBreakdown{// DF, KE, ND, MD, AP
		{name: "DF", start:0, end: 5},
	},
	//25: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//},
	//26: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//},
	//27: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//},
	//28: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//	{name: "??", start:5, end: 32},
	//	{name: "CRC", start:32, end: 56},
	//},
	//29: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//},
	//30: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//},
	//31: []featureBreakdown{
	//	{name: "DF", start:0, end: 5},
	//},
}

func (frame *Frame) Describe(output io.Writer) {
	fmt.Fprintf(output, "MODE S Packet:\n")
	fmt.Fprintf(output, "Length:         : %d bits\n", frame.getMessageLengthBits())
	fmt.Fprintf(output, "Downlink Format : %d (%s)\n", frame.downLinkFormat, frame.DownLinkFormat())
	fmt.Fprintf(output, "Frame           : %s\n", frame.raw)
	// decode the specific DF type
	switch frame.downLinkFormat {
	case 0: // Airborne position, baro altitude only
		frame.showVerticalStatus(output)
		frame.showAltitude(output)
		frame.showRInformation(output)
		frame.showSLField(output)
	case 1, 2, 3: // Aircraft Identification and Category
		frame.showFlightStatus(output)
		frame.showFlightId(output)
	case 4:
		frame.showFlightStatus(output)
		frame.showAltitude(output)
	case 5: //DF_5
		frame.showFlightStatus(output)
		frame.showIdentity(output)
	case 11: //DF_11
		frame.showICAO(output)
		frame.showCapability(output)
	case 16: //DF_16
		frame.showVerticalStatus(output)
		frame.showAltitude(output)
		frame.showRInformation(output)
		frame.showSLField(output)
	case 17: //DF_17
		frame.showICAO(output)
		frame.showCapability(output)
		frame.showAdsb(output)
	case 18: //DF_18
		//frame.showCapability() // control field
		if 0 == frame.capability {
			frame.showICAO(output)
			frame.showCapability(output)
			frame.showAdsb(output)
		} else {
			fmt.Fprintln(output, "Unable to decode DF18 Capability:", frame.capability)
		}
	case 20: //DF_20
		frame.showFlightStatus(output)
		frame.showAltitude(output)
		frame.showFlightNumber(output)
	case 21: //DF_21
		frame.showFlightStatus(output)
		frame.showIdentity(output) // gillham encoded squawk
		frame.showFlightNumber(output)
	}

	frame.showBitString(output)

}

func (f *Frame) showVerticalStatus(output io.Writer) {
	if !f.validVerticalStatus {
		return
	}
	if f.onGround {
		fmt.Fprintf(output, "Vertical Status : On The Ground");
	} else {
		fmt.Fprintf(output, "Vertical Status : Airborne");
	}
	fmt.Fprintln(output, "")
}
func (f *Frame) showVerticalRate(output io.Writer) {
	fmt.Fprintf(output, "  Vertical Rate : %d", f.verticalRate)
	fmt.Fprintln(output, "")
}

func (f *Frame) showAltitude(output io.Writer) {
	fmt.Fprintf(output, "Altitude        : %d %s", f.altitude, f.AltitudeUnits())
	fmt.Fprintln(output, "")
}

func (f *Frame) showFlightStatus(output io.Writer) {
	fmt.Fprintf(output, "Flight Status   : (%d) %s\n", f.flightStatus, flightStatusTable[f.flightStatus])
	f.showVerticalStatus(output)
	if "" != f.special {
		fmt.Fprintf(output, "Special Status  : %s", f.special)
		fmt.Fprintln(output, "")
	}
	f.showAlert(output)
}

func (f *Frame) showFlightId(output io.Writer) {
	fmt.Fprintf(output, "Flight          : %s", f.Flight())
	fmt.Fprintln(output, "")
}

func (f *Frame) showICAO(output io.Writer) {
	fmt.Fprintf(output, "ICAO            : %6X", f.icao)
	fmt.Fprintln(output, "")
}

func (f *Frame) showCapability(output io.Writer) {
	f.showVerticalStatus(output)
	fmt.Fprintf(output, "Plane Mode S Cap: %s", capabilityTable[f.capability])
	fmt.Fprintln(output, "")
}

func (f *Frame) showIdentity(output io.Writer) {
	fmt.Fprintf(output, "Squawk Identity : %04d", f.identity)
	fmt.Fprintln(output, "")
}
func (f *Frame) showVelocity(output io.Writer) {
	if f.superSonic {
		fmt.Fprintln(output, "  Super Sonic?  : Yes!")
	} else {
		fmt.Fprintln(output, "  Super Sonic?  : No")
	}
	fmt.Fprintf(output, "  Velocity      : %0.2f", f.velocity)
	fmt.Fprintln(output, "")
}
func (f *Frame) showHeading(output io.Writer) {
	fmt.Fprintf(output, "  Heading       : %0.2f", f.heading)
	fmt.Fprintln(output, "")
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

func (f *Frame)showRInformation(output io.Writer) {
	fmt.Fprintf(output, "TCAS            : (%d) %s\n", f.ri, rInformationField[f.ri])
}
func (f *Frame)showSLField(output io.Writer) {
	fmt.Fprintf(output, "TCAS            : (%d) %s\n", f.sl, slInformationField[f.sl])
}

func (f *Frame) showAdsb(output io.Writer) {
	fmt.Fprintf(output, "ADSB Msg Type   : %d (Sub Type %d): %s\n", f.messageType, f.messageSubType, f.MessageTypeString())

	switch f.messageType {
	case 1, 2, 3, 4:
		f.showFlightNumber(output)
	case 5, 6, 7, 8:
		f.showHeading(output)
		f.showVelocity(output)
		f.showCprLatLon(output)
	case 9, 10, 11, 12, 13, 14, 15, 16, 17, 18:
		f.showAltitude(output)
		f.showCprLatLon(output)
	case 19:
		switch f.messageSubType {
		case 1, 2:
			f.showHeading(output)
			f.showVerticalRate(output)
			f.showVelocity(output)
		case 3, 4:
			f.showHeading(output)
			f.showVerticalRate(output)
			f.showVelocity(output)
		default:
			fmt.Println("Unknown Sub Type")
		}
	case 23:
		if 7 == f.messageSubType {
			f.showIdentity(output);
		}
	case 28:
		if 1 == f.messageSubType {
			f.showIdentity(output);
			f.showAlert(output);
		}
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
		fmt.Fprintf(output, "  Special       : %s\n", f.special)
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
				doMakeBitString(sf)
			}
		}
	}

	return fmt.Sprintf("\n%s\n%s\n%s\n%s\n%s\n\n%s\n%d/%d bits shown\n", header, separator, bits, separator, bitDesc, footer, shownBitCount, frame.getMessageLengthBits())
}