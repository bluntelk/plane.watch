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
}

var featureDescription = map[string]featureDescriptionType{
	"AA":{field: "Address Announced", meaning: "aircraft identification in All-Call reply - ICAO"},
	"AC":{field: "Altitude Code", meaning: "aircraft altitude code"},
	"AP":{field: "Address/Parity", meaning: "error detection field"},
	"AQ":{field: "Acquisition", meaning: "part of air-to-air protocol"},
	"CA":{field: "Capability", meaning: "aircraft report of system capability"},
	"CC":{field: "Crosslink Capability", meaning:""}, // added
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
	"PI":{field: "Parity/Interr.Identity", meaning: "reports source of interrogation"},
	"PR":{field: "Probability of Reply", meaning: "used in stochastic acquisition mode"},
	"RC":{field: "Reply Control", meaning: "part of ELM protocol"},
	"RI":{field: "Reply Information", meaning: "aircraft status information for TCAS"},
	"RL":{field: "Reply Length", meaning: "commands air-to-air reply length"},
	"RR":{field: "Reply Request", meaning: "commands details of reply"},
	"SD":{field: "Special Designator", meaning: "control codes to transponder"},
	"SL":{field: "Sensitivity level, ACAS", meaning: "control codes to transponder"},
	"UF":{field: "Uplink Format", meaning: "format descriptor"},
	"UM":{field: "Utility Message", meaning: "protocol message"},
	"VS":{field: "Vertical Status", meaning: "aircraft status, airborne (0) or on the ground (1)"},
	"??":{field: "???", meaning:"Unknown"},
	"CRC":{field: "CRC", meaning:"CRC Checksum"},
}

var frameFeatures = map[byte][]featureBreakdown{

	0: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "VS", start:5, end: 6},
		{name: "CC", start:6, end: 7},
		{name: "??", start:7, end: 8},
		{name: "SL", start:8, end: 11},
		{name: "??", start:11, end: 20},
		{name: "AC", start:20, end: 32},
		{name: "CRC", start:32, end: 56},
	},
	1: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "FS", start:5, end: 8},
		{name: "??", start:8, end: 32},
		{name: "CRC", start:32, end: 56},
	},
	2: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "FS", start:5, end: 8},
		{name: "??", start:8, end: 32},
		{name: "CRC", start:32, end: 56},
	},
	3: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "FS", start:5, end: 8},
		{name: "??", start:8, end: 32},
		{name: "CRC", start:32, end: 56},
	},
	4: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "FS", start:5, end: 8},
		{name: "DR", start:8, end: 13},
		{name: "UM", start:13, end: 19},
		{name: "AC", start:19, end: 32},
		{name: "CRC", start:32, end: 56},
	},
	5: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "FS", start:5, end: 8},
		{name: "DR", start:8, end: 13},
		{name: "UM", start:13, end: 19},
		{name: "ID", start:19, end: 32},
		{name: "CRC", start:32, end: 56},
	},
	6: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	7: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	8: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	9: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	10: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	11: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "CA", start:5, end: 8},
		{name: "AA", start:8, end: 32},
		{name: "CRC", start:32, end: 56},
	},
	12: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	13: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	14: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	15: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	16: []featureBreakdown{ //RI,MV Fields
		{name: "DF", start:0, end: 5},
		{name: "VS", start:5, end: 6},
		{name: "CC", start:6, end: 7},
		{name: "??", start:7, end: 8},
		{name: "SL", start:8, end: 11},
		{name: "??", start:11, end: 19},
		{name: "AC", start:19, end: 32},
		{name: "??", start:32, end: 88},
		{name: "CRC", start:88, end: 112},
	},
	17: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "??", start:5, end: 88},
		{name: "CRC", start:88, end: 112},
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
	22: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	23: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	24: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	25: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	26: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	27: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	28: []featureBreakdown{
		{name: "DF", start:0, end: 5},
		{name: "??", start:5, end: 32},
		{name: "CRC", start:32, end: 56},
	},
	29: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	30: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
	31: []featureBreakdown{
		{name: "DF", start:0, end: 5},
	},
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
	fmt.Fprintf(output, "  CPR Frame          : %s\n", oddEven)
	fmt.Fprintf(output, "  CPR Latitude       : %d\n", f.rawLatitude)
	fmt.Fprintf(output, "  CPR Longitude      : %d\n", f.rawLongitude)
	fmt.Fprintln(output, "")
}

func (f *Frame) showAdsb(output io.Writer) {
	fmt.Fprintf(output, "ADSB Msg Type   : %d (Sub Type %d): %s\n", f.messageType, f.messageSubType, f.MessageTypeString())

	switch f.messageType {
	case 1, 2, 3, 4:
		f.showFlightNumber(output)
	case 5, 6, 7, 8:
		f.showVelocity(output)
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

func (f *Frame) formatBitString(features []featureBreakdown) string {
	var header, seperator, bits, rawBits, bitFmt, bitDesc, footer, suffix string
	var padLen, realLen, shownBitCount int

	for _, i := range f.message {
		rawBits += fmt.Sprintf("%08s", strconv.FormatUint(uint64(i), 2))
	}

	for _, f := range features {
		padLen = len(f.name)
		realLen = f.end - f.start
		if realLen > padLen {
			padLen = realLen
		}
		shownBitCount += (f.end - f.start)
		bitFmt = fmt.Sprintf(" %%- %ds |", padLen)
		header += fmt.Sprintf(bitFmt, f.name)
		seperator += strings.Repeat("-", padLen + 2) + "+"
		bits += fmt.Sprintf(bitFmt, rawBits[f.start: f.end])
		bitDesc += fmt.Sprintf(bitFmt, strconv.Itoa(f.start))

		if 1 == realLen {
			suffix = ""
		} else {
			suffix = "s"
		}

		feature := featureDescription[f.name]
		footer += fmt.Sprintf(" %s \t%d bit%s\t %s: %s\n", f.name, realLen, suffix, feature.field, feature.meaning)
	}

	return fmt.Sprintf("\n%s\n%s\n%s\n%s\n%s\n\n%s\n%d/%d bits shown\n", header, seperator, bits, seperator, bitDesc, footer, shownBitCount, f.getMessageLengthBits())
	//return "\n" + header + "\n" + seperator + "\n" + bits + "\n" + seperator + "\n" + bitDesc + "\n\n" + footer
}

func (f *Frame) BitStringDF17() string {

	if f.downLinkFormat != 17 {
		return ""
	}
	var header, bits, rawBits string

	switch f.messageType {
	case 9, 10, 11, 12, 13, 14, 15, 16, 17, 18:
		header += " DF   | CA  | ICAO 24bit addr          | DATA                                                                                | CRC                      |\n"
		header += "                                       | TC    | SS | NICsb | ALT     Q      | T | F | LAT-CPR           | LON-CPR           |                          |\n"
		header += "------+-----+--------------------------+-------+----+-------+----------------+---+---+-------------------+-------------------+--------------------------+\n"

		for _, i := range f.message {
			rawBits += fmt.Sprintf("%08s", strconv.FormatUint(uint64(i), 2))
		}

		bits += rawBits[0:5] + " | " // Downlink Format
		bits += rawBits[5:8] + " | " // Capability
		bits += rawBits[8:32] + " | " // ICAO
		// now we are into the packet data

		bits += rawBits[32:37] + " | " // TC - Type code
		bits += rawBits[37:39] + " | " // SS - Surveillance status
		bits += rawBits[39:40] + "     | " // NIC supplement-B
		bits += rawBits[40:47] + " "   // Altitude
		bits += rawBits[47:48] + " "   // Altitude Q Bit
		bits += rawBits[48:52] + " | " // Altitude
		bits += rawBits[52:53] + " | " // Time
		bits += rawBits[53:54] + " | " // F - CPR odd/even frame flag
		bits += rawBits[54:71] + " | " // Latitude in CPR format
		bits += rawBits[71:88] + " | " // Longitude in CPR format
		bits += rawBits[88:] + " | "   // CRC
		bits += "\n"
	default:
		header += " DF   | CA  | ICAO 24bit addr          | DATA                                                     | CRC                      |\n"
		header += "------+-----+--------------------------+-------------------------------------------------------------------------------------+\n"
		for _, i := range f.message {
			rawBits += fmt.Sprintf("%08s", strconv.FormatUint(uint64(i), 2))
		}

		bits += rawBits[0:5] + " | " // Downlink Format
		bits += rawBits[5:8] + " | " // Capability
		bits += rawBits[8:32] + " | " // ICAO
		bits += rawBits[32:88] + " | "   // DATA
		bits += rawBits[88:] + " | "   // CRC
	}

	return header + bits
}