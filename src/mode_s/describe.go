package mode_s

import (
	"fmt"
	"io"
	"strconv"
)

func (frame *Frame) Describe(output io.Writer) {
	fmt.Fprintln(output, "----------------------------------------------------")
	fmt.Fprintf(output, "MODE S Packet:\n")
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
	if f.alert {
		fmt.Fprintf(output, "Plane showing Alert!")
		fmt.Fprintln(output, "")
	}
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
	case 9,10,11,12,13,14,15,16,17,18:
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
		}
	default:
		fmt.Fprintln(output, "Packet Type Not Yet Decoded")
	}

	fmt.Fprintln(output, "")
}

func (f *Frame) showFlightNumber(output io.Writer) {
	fmt.Fprintf(output, "Flight Number   : %s", f.FlightNumber())
	fmt.Fprintln(output, "")
}

func (f *Frame) DescribeFlightStatus(output io.Writer) {
	fmt.Fprintf(output, "Flight Status: %s", flightStatusTable[f.flightStatus])
	fmt.Fprintln(output, "")
}
func (f *Frame) DescribeIdentity(output io.Writer) {
	fmt.Fprintf(output, "Flight Identity: %s", flightStatusTable[f.flightStatus])
	fmt.Fprintln(output, "")
}
func (f *Frame) DescribePosition(output io.Writer) {
	fmt.Fprintln(output, "Position:")
	fmt.Fprintf(output, "    Altitude: %d feet", f.altitude)
	fmt.Fprintln(output, "")
}

// determines what type of mode S Message this frame is
func (f *Frame) DownLinkFormat() string {

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