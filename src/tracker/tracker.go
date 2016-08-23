package tracker

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mode_s"
	"os"
	"sbs1"
)

func SetDebugOutput(out io.Writer) {
	log.SetOutput(out)
}

func init() {
	SetDebugOutput(os.Stdout)
}

func HandleModeSFrame(frame mode_s.Frame, debug bool) *Plane {
	icao := frame.ICAOAddr()
	if 0 == icao {
		return nil
	}
	var planeFormat string
	var hasChanged bool

	plane := GetPlane(icao)
	//plane.MarkFrameTime(frame.TimeStamp()) // todo, make this fast!

	debugMessage := func(sfmt string, a ...interface{}) {
		if debug {
			planeFormat = fmt.Sprintf("DF%02d - \033[0;97mPlane (\033[38;5;118m%s %-8s\033[0;97m)", frame.DownLinkType(), plane.Icao, plane.Flight.Identifier)
			fmt.Printf(planeFormat+sfmt+"\n", a...)
		}
	}

	// determine what to do with our given frame
	switch frame.DownLinkType() {
	case 0:
		// grab the altitude
		if frame.AltitudeValid() {
			plane.Location.Altitude, _ = frame.Altitude()
			plane.Location.AltitudeUnits = frame.AltitudeUnits()
		}
		if frame.VerticalStatusValid() {
			plane.Location.onGround, _ = frame.OnGround()
		}
		plane.Location.TimeStamp = frame.TimeStamp()
		debugMessage(" is at %d %s \033[0m", plane.Location.Altitude, plane.Location.AltitudeUnits)

		hasChanged = true

	case 1, 2, 3:
		hasChanged = true
		if frame.VerticalStatusValid() {
			plane.Location.onGround, _ = frame.OnGround()
		}
		plane.Location.TimeStamp = frame.TimeStamp()
		if frame.Alert() {
			plane.Special = "Alert"
		}
	case 6, 7, 8, 9, 10, 12, 13, 14, 15, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31:
		debugMessage(" \033[38;5;52mIgnoring Mode S Frame: %d (%s)\033[0m\n", frame.DownLinkType(), frame.DownLinkFormat())
		break
	case 11:
		if frame.VerticalStatusValid() {
			g, _ := frame.OnGround()
			hasChanged = plane.Location.onGround != g
			plane.Location.onGround = g
		}
	case 4, 5, 20, 21:
		hasChanged = true
		if frame.VerticalStatusValid() {
			plane.Location.onGround, _ = frame.OnGround()
		}
		if frame.Alert() {
			plane.Special = "Alert"
		}
		if frame.AltitudeValid() {
			plane.Location.Altitude, _ = frame.Altitude()
			plane.Location.AltitudeUnits = frame.AltitudeUnits()
		}
		plane.Location.TimeStamp = frame.TimeStamp()
		plane.Flight.Status = frame.FlightStatusString()
		plane.Flight.StatusId = frame.FlightStatus()
		if 5 == frame.DownLinkType() || 21 == frame.DownLinkType() {
			plane.SquawkIdentity = frame.SquawkIdentity()
		}
		hasChanged = true
		debugMessage(" is at %d %s and flight status is: %s. \033[2mMode S Frame: %d \033[0m",
			plane.Location.Altitude, plane.Location.AltitudeUnits, plane.Flight.Status, frame.DownLinkType())
		break
	case 16:
		hasChanged = true
		if frame.AltitudeValid() {
			plane.Location.Altitude, _ = frame.Altitude()
			plane.Location.AltitudeUnits = frame.AltitudeUnits()
		}
		if frame.VerticalStatusValid() {
			plane.Location.onGround, _ = frame.OnGround()
		}
		plane.Location.TimeStamp = frame.TimeStamp()

	case 17, 18: // ADS-B
		//if debug {
		//	frame.Describe(os.Stdout)
		//}

		// i am using the text version because it is easier to program with.
		// if performance is an issue, change over to byte comparing
		messageType := frame.MessageTypeString()
		switch messageType {
		case mode_s.DF17_FRAME_ID_CAT: // "Aircraft Identification and Category"
			{
				plane.Flight.Identifier = frame.FlightNumber()
				if frame.ValidCategory() {
					plane.AirframeCategory = frame.Category()
				}
				hasChanged = true
				break
			}
		case mode_s.DF17_FRAME_SURFACE_POS: // "Surface Position"
			{
				if frame.HeadingValid() {
					plane.Location.Heading, _ = frame.Heading()
					plane.Location.hasHeading = true
				}
				if frame.VelocityValid() {
					plane.Location.Velocity, _ = frame.Velocity()
				}
				if frame.VerticalStatusValid() {
					plane.Location.onGround, _ = frame.OnGround()
				}

				if frame.IsEven() {
					plane.SetCprEvenLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				} else {
					plane.SetCprOddLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				}
				plane.DecodeCpr(frame.TimeStamp())

				plane.Location.TimeStamp = frame.TimeStamp()

				debugMessage(" is on the ground and has heading %0.2f and is travelling at %0.2f knots\033[0m", plane.Location.Heading, plane.Location.Velocity)
				hasChanged = true
				break
			}
		case mode_s.DF17_FRAME_AIR_POS_BARO: // "Airborne Position (with Barometric Altitude)"
			{
				if frame.VerticalStatusValid() {
					plane.Location.onGround, _ = frame.OnGround()
				}
				plane.Location.TimeStamp = frame.TimeStamp()
				hasChanged = true

				if frame.IsEven() {
					plane.SetCprEvenLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				} else {
					plane.SetCprOddLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				}

				altitude, _ := frame.Altitude()
				plane.SetAltitude(altitude, frame.AltitudeUnits())
				plane.DecodeCpr(frame.TimeStamp())

				break
			}
		case mode_s.DF17_FRAME_AIR_VELOCITY: // "Airborne Velocity"
			{
				if frame.HeadingValid() {
					plane.Location.Heading, _ = frame.Heading()
				}
				if frame.VelocityValid() {
					plane.Location.Velocity, _ = frame.Velocity()

				}
				if frame.VerticalRateValid() {
					plane.Location.VerticalRate, _ = frame.VerticalRate()
				}
				if frame.VerticalStatusValid() {
					plane.Location.onGround, _ = frame.OnGround()
				}
				plane.Location.hasHeading = true
				plane.Location.TimeStamp = frame.TimeStamp()
				debugMessage(" has heading %0.2f and is travelling at %0.2f knots\033[0m", plane.Location.Heading, plane.Location.Velocity)
				hasChanged = true
				break
			}
		case mode_s.DF17_FRAME_AIR_POS_GNSS: // "Airborne Position (GNSS Height)"
			{
				debugMessage("\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		case mode_s.DF17_FRAME_TEST_MSG: //, "Test Message":
			if debug {
				debugMessage("\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
			}
			break
		case mode_s.DF17_FRAME_TEST_MSG_SQUAWK: //, "Test Message":
			{
				if frame.SquawkIdentity() > 0 {
					hasChanged = plane.SquawkIdentity != frame.SquawkIdentity()
					plane.SquawkIdentity = frame.SquawkIdentity()
				}
				break
			}
		case mode_s.DF17_FRAME_SURFACE_SYS_STATUS: //, "Surface System Status":
			{
				debugMessage("\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		case mode_s.DF17_FRAME_EMERG_PRIORITY: //, "Extended Squitter Aircraft Status (Emergency)":
			{
				debugMessage("\033[2m %s\033[0m", messageType)
				plane.Special = "Emergency"
				plane.SquawkIdentity = frame.SquawkIdentity()
				hasChanged = true
				break
			}
		case mode_s.DF17_FRAME_TCAS_RA: //, "Extended Squitter Aircraft Status (1090ES TCAS RA)":
			{
				debugMessage("\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		case mode_s.DF17_FRAME_STATE_STATUS: //, "Target State and Status Message":
			{
				debugMessage("\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		case mode_s.DF17_FRAME_AIRCRAFT_OPER: //, "Aircraft Operational Status Message":
			{
				debugMessage("\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		}

	}
	SetPlane(plane, frame.TimeStamp())
	if hasChanged {
		//log.Println(plane.String())
		return &plane
	} else {
		return nil
	}
}

func HandleSbs1Frame(frame sbs1.Frame, debug bool) *Plane {
	var hasChanged bool
	icaoAddr, err := IcaoStringToInt(frame.Icao)
	if err != nil {
		if debug {
			log.Println(err)
		}
		return nil
	}
	plane := GetPlane(icaoAddr)
	if frame.HasPosition {
		plane.AddLatLong(frame.Lat, frame.Lon, frame.Received)
		hasChanged = true
	}
	SetPlane(plane, frame.Received)
	if hasChanged {
		return &plane
	}
	return nil
}

func IcaoStringToInt(icao string) (uint32, error) {
	btoi, err := hex.DecodeString(icao)
	if nil != err {
		return 0, fmt.Errorf("Failed to decode ICAO HEX (%s) into UINT32. %s", icao, err)
	}
	return uint32(btoi[0])<<16 | uint32(btoi[1])<<8 | uint32(btoi[2]), nil
}
