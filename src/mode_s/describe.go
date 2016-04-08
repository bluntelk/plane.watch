package mode_s

import (
	"fmt"
	"io"
	"strconv"
)



// prints out a nice debug message
func (f *Frame) Describe(output io.Writer) {
	fmt.Fprintln(output, "----------------------------------------------------")
	fmt.Fprintf(output, "MODE S Packet:\n")
	fmt.Fprintf(output, "Downlink Format : %d (%s)\n", f.downLinkFormat, f.DownLinkFormat())
	fmt.Fprintf(output, "Frame           : %s\n", f.raw)
	fmt.Fprintf(output, "ICAO            : %6X\n", f.icao)
	fmt.Fprintf(output, "Frame mode      : %s\n", f.mode)
	//fmt.Fprintf(output, "Time Stamp      : %s\n", f.timeStamp.Format(time.RFC3339Nano))

	switch f.downLinkFormat {
	case 0:
		f.DescribePosition(output)
	case 1, 2, 3, 4:
		f.DescribeIdentity(output)
	case 5:
		f.DescribeFlightStatus(output)
	case 11:
		fmt.Fprintf(output, "Capability  : %d (%s)\n", f.df11.capability, f.DownLinkCapability())
	case 17:
		fmt.Fprintf(output, "ADS-B Frame   : %d (%s)\n", f.messageType, f.MessageTypeString())
		fmt.Fprintf(output, "ADS-B Sub Type: %d\n", f.messageSubType)

		switch f.messageType {
		case 1, 2, 3, 4:
			fmt.Fprintf(output, "      Flight: %s", string(f.flight))
		case 9, 10, 11, 12, 13, 14, 15, 16, 17, 18:
			var oddEven string = "Odd"
			if 0 == f.cprFlagOddEven {
				oddEven = "Even"
			}
			fmt.Fprintf(output, "   CPR FRAME : %s (one half of a location)\n", oddEven)
		}

		fmt.Fprintf(output, f.BitStringDF17())
	default:
		fmt.Fprintln(output, "Not Decoded")
	}
}

func (f *Frame) DescribeFlightStatus(output io.Writer) {
	fmt.Fprintln(output, "Flight Status: %s", flightStatusTable[f.flightStatus])
	fmt.Fprintln(output, "")
}
func (f *Frame) DescribeIdentity(output io.Writer) {
	fmt.Fprintln(output, "Flight Identity: %s", flightStatusTable[f.flightStatus])
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