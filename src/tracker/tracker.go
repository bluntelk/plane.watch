package tracker

import (
	"log"
	"fmt"
	"mode_s"
	"os"
	"time"
	"io"
)

func SetDebugOutput(out io.Writer) {
	log.SetOutput(out)
}

func init() {
	SetDebugOutput(os.Stdout)
}

func HandleModeSFrame(frame mode_s.Frame, debug bool) *Plane {
	//		frame.Describe(output)
	icao := frame.ICAOAddr()
	if 0 == icao {
		return nil
	}
	var planeFormat string
	var hasChanged bool

	plane := GetPlane(icao)
	planeFormat = fmt.Sprintf("DF%02d - \033[0;97mPlane (\033[38;5;118m%s %-8s\033[0;97m)", frame.DownLinkType(), plane.Icao, plane.Flight.Identifier)
	plane.MarkFrameTime()

	// determine what to do with our given frame
	switch frame.DownLinkType() {
	case 0:
		// grab the altitude
		plane.Location.Altitude = frame.Altitude()
		plane.Location.onGround = frame.OnGround()
		plane.Location.TimeStamp = time.Now()
		if debug {
			log.Printf(planeFormat + " is at %d %s \033[0m", plane.Location.Altitude, plane.Location.AltitudeUnits)
		}

		hasChanged = true

	case 1, 2, 3:
		hasChanged = true
		plane.Location.onGround = frame.OnGround()
		plane.Location.TimeStamp = time.Now()
		if frame.Alert() {
			plane.Special = "Alert"
		}
	case 6, 7, 8, 9, 10, 12, 13, 14, 15, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31:
		if debug {
			log.Printf(planeFormat + " \033[38;5;52mIgnoring Mode S Frame: %d (%s)\033[0m\n", frame.DownLinkType(), frame.DownLinkFormat())
		}
		break
	case 11:
		if frame.ValidVerticalStatus() {
			hasChanged = plane.Location.onGround != frame.OnGround()
			plane.Location.onGround = frame.OnGround()
		}
	case 4, 5, 20, 21:
		hasChanged = true
		plane.Location.onGround = frame.OnGround()
		if frame.Alert() {
			plane.Special = "Alert"
		}
		plane.Location.Altitude = frame.Altitude()
		plane.Location.onGround = false
		plane.Location.AltitudeUnits = frame.AltitudeUnits()
		plane.Location.TimeStamp = time.Now()
		plane.Flight.Status = frame.FlightStatusString()
		plane.Flight.StatusId = frame.FlightStatus()
		if 5 == frame.DownLinkType() || 21 == frame.DownLinkType() {
			plane.SquawkIdentity = frame.SquawkIdentity()
		}
		hasChanged = true
		if debug {
			log.Printf(planeFormat + " is at %d %s and flight status is: %s. \033[2mMode S Frame: %d \033[0m",
				plane.Location.Altitude, plane.Location.AltitudeUnits, plane.Flight.Status, frame.DownLinkType())
		}
		break
	case 16:
		hasChanged = true
		plane.Location.Altitude = frame.Altitude()
		plane.Location.onGround = frame.OnGround()
		plane.Location.TimeStamp = time.Now()

	case 17, 18: // ADS-B
		if debug {
			frame.Describe(os.Stdout)
		}

		// i am using the text version because it is easier to program with.
		// if performance is an issue, change over to byte comparing
		messageType := frame.MessageTypeString()
		switch messageType {
		case mode_s.DF17_FRAME_ID_CAT: // "Aircraft Identification and Category"
			{
				plane.Flight.Identifier = frame.FlightNumber()
				hasChanged = true
				break
			}
		case mode_s.DF17_FRAME_SURFACE_POS: // "Surface Position"
			{
				plane.Location.Heading = frame.Heading()
				plane.Location.Velocity = frame.Velocity()
				plane.Location.onGround = true
				plane.Location.hasHeading = true
				plane.Location.TimeStamp = time.Now()

				if debug {
					log.Printf(planeFormat + " is on the ground and has heading %0.2f and is travelling at %0.2f knots\033[0m", plane.Location.Heading, plane.Location.Velocity)
				}
				hasChanged = true
				break
			}
		case mode_s.DF17_FRAME_AIR_POS_BARO:// "Airborne Position (with Barometric Altitude)"
			{
				plane.Location.onGround = false
				plane.Location.TimeStamp = time.Now()
				hasChanged = true

				if frame.IsEven() {
					plane.SetCprEvenLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				} else {
					plane.SetCprOddLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				}

				plane.DecodeCpr(frame.Altitude(), frame.AltitudeUnits())

				break
			}
		case mode_s.DF17_FRAME_AIR_VELOCITY: // "Airborne Velocity"
			{
				plane.Location.Heading = frame.Heading()
				plane.Location.Velocity = frame.Velocity()
				plane.Location.VerticalRate = frame.VerticalRate()
				plane.Location.onGround = false
				plane.Location.hasHeading = true
				plane.Location.TimeStamp = time.Now()
				if debug {
					log.Printf(planeFormat + " has heading %0.2f and is travelling at %0.2f knots\033[0m", plane.Location.Heading, plane.Location.Velocity)
				}
				hasChanged = true
				break
			}
		case mode_s.DF17_FRAME_AIR_POS_GNSS: // "Airborne Position (GNSS Height)"
			{
				if debug {
					log.Printf(planeFormat + "\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				}
				break
			}
		case mode_s.DF17_FRAME_TEST_MSG: //, "Test Message":
			if debug {
				log.Printf(planeFormat + "\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
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
				if debug {
					log.Printf(planeFormat + "\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				}
				break
			}
		case mode_s.DF17_FRAME_EXT_SQUIT_EMERG: //, "Extended Squitter Aircraft Status (Emergency)":
			{
				if debug {
					log.Printf(planeFormat + "\033[2m %s\033[0m", messageType)
				}
				plane.Special = "Emergency"
				plane.SquawkIdentity = frame.SquawkIdentity()
				hasChanged = true
				break
			}
		case mode_s.DF17_FRAME_EXT_SQUIT_STATUS: //, "Extended Squitter Aircraft Status (1090ES TCAS RA)":
			{
				if debug {
					log.Printf(planeFormat + "\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				}
				break
			}
		case mode_s.DF17_FRAME_STATE_STATUS: //, "Target State and Status Message":
			{
				if debug {
					log.Printf(planeFormat + "\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				}
				break
			}
		case mode_s.DF17_FRAME_AIRCRAFT_OPER: //, "Aircraft Operational Status Message":
			{
				if debug {
					log.Printf(planeFormat + "\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				}
				break
			}
		}

	}
	SetPlane(plane)
	if hasChanged {
		log.Println(plane.String())
		return &plane
	} else {
		return nil
	}
}