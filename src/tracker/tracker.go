package tracker

import (
	"log"
	"fmt"
	"mode_s"
	"time"
	"os"
)

func HandleModeSFrame(frame mode_s.Frame, debug bool) *Plane {
	//		frame.Describe(output)
	icao := frame.ICAOAddr()
	if 0 == icao {
		return nil
	}
	var planeFormat string
	var hasChanged bool

	plane := GetPlane(icao)
	planeFormat = fmt.Sprintf("DF%02d - \033[0;97mPlane (\033[38;5;118m%6x %-8s\033[0;97m)", frame.DownLinkType(), plane.IcaoIdentifier, plane.Flight.Identifier)

	// determine what to do with our given frame
	switch frame.DownLinkType() {
	case 0:
		// grab the altitude
		plane.Location.Altitude = frame.Altitude()
		log.Printf(planeFormat + " is at %d %s \033[0m", plane.Location.Altitude, plane.Location.AltitudeUnits)
		hasChanged = true

	case 1, 2, 3, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 19, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31:
		//log.Printf(planeFormat + " \033[38;5;52mIgnoring Mode S Frame: %d (%s)\033[0m\n", frame.DownLinkType(), frame.DownLinkFormat())
		break
	case 4, 5, 20, 21:
		plane.Location.Altitude = frame.Altitude()
		plane.Location.onGround = false
		plane.Location.AltitudeUnits = frame.AltitudeUnits()
		plane.Flight.Status = frame.FlightStatusString()
		plane.Flight.StatusId = frame.FlightStatusInt()
		if 5 == frame.DownLinkType() || 21 == frame.DownLinkType() {
			plane.SquawkIdentity = frame.SquawkIdentity()
		}
		hasChanged = true
		log.Printf(planeFormat + " is at %d %s and flight status is: %s. \033[2mMode S Frame: %d \033[0m",
			plane.Location.Altitude, plane.Location.AltitudeUnits, plane.Flight.Status, frame.DownLinkType())
		break
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
				plane.Flight.Identifier = frame.GetFlight()
				hasChanged = true
				break
			}
		case mode_s.DF17_FRAME_SURFACE_POS: // "Surface Position"
			{
				plane.Location.Heading = frame.Heading()
				plane.Location.Velocity = frame.Velocity()
				plane.Location.onGround = true
				plane.Location.hasHeading = true

				log.Printf(planeFormat + " is on the ground and has heading %0.2f and is travelling at %0.2f knots\033[0m", plane.Location.Heading, plane.Location.Velocity)
				hasChanged = true
				break
			}
		case mode_s.DF17_FRAME_AIR_POS_BARO:// "Airborne Position (with Barometric Altitude)"
			{
				plane.Location.onGround = false
				hasChanged = true

				if frame.IsEven() {
					plane.SetCprEvenLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				} else {
					plane.SetCprOddLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				}

				err := plane.DecodeCpr(frame.Altitude(), frame.AltitudeUnits())

				if nil == err {
					// we were able to decode!

					log.Printf(planeFormat + " at:  \033[38;5;122m%+03.13f\033[0;97m, \033[38;5;122m%+03.13f\033[0;97m, \033[38;5;226m%d \033[0;97m%s\033[0m",
						plane.Location.Latitude, plane.Location.Longitude, plane.Location.Altitude, plane.Location.AltitudeUnits)
				} else {
					log.Printf(planeFormat + "Failed to Decode Lat/Lon: %s \033[0m", err)
				}

				break
			}
		case mode_s.DF17_FRAME_AIR_VELOCITY: // "Airborne Velocity"
			{
				plane.Location.Heading = frame.Heading()
				plane.Location.Velocity = frame.Velocity()
				plane.Location.onGround = false
				plane.Location.hasHeading = true
				log.Printf(planeFormat + " has heading %0.2f and is travelling at %0.2f knots\033[0m", plane.Location.Heading, plane.Location.Velocity)
				hasChanged = true
				break
			}
		case mode_s.DF17_FRAME_AIR_POS_GNSS: // "Airborne Position (GNSS Height)"
			{
				log.Printf(planeFormat + "\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		case mode_s.DF17_FRAME_TEST_MSG: //, "Test Message":
			{
				log.Printf(planeFormat + "\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		case mode_s.DF17_FRAME_SURFACE_SYS_STATUS: //, "Surface System Status":
			{
				log.Printf(planeFormat + "\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		case mode_s.DF17_FRAME_EXT_SQUIT_EMERG: //, "Extended Squitter Aircraft Status (Emergency)":
			{
				log.Printf(planeFormat + "\033[2m %s\033[0m", messageType)
				plane.Special = "Emergency"
				plane.SquawkIdentity = frame.SquawkIdentity()
				hasChanged = true
				break
			}
		case mode_s.DF17_FRAME_EXT_SQUIT_STATUS: //, "Extended Squitter Aircraft Status (1090ES TCAS RA)":
			{
				log.Printf(planeFormat + "\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		case mode_s.DF17_FRAME_STATE_STATUS: //, "Target State and Status Message":
			{
				log.Printf(planeFormat + "\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		case mode_s.DF17_FRAME_AIRCRAFT_OPER: //, "Aircraft Operational Status Message":
			{
				log.Printf(planeFormat + "\033[2mIgnoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		}

	}
	SetPlane(plane)
	if hasChanged {
		log.Println(plane)
		return plane
	} else {
		return nil
	}
}

func CleanPlanes() {
	//remove planes that have not been seen for a while
	planeAccessMutex.Lock()

	// go through the list and remove planes
	var cutOff, planeCutOff time.Time
	cutOff = time.Now().Add(5 * time.Minute)
	for i, plane := range planeList {
		planeCutOff = plane.lastSeen.Add(5 * time.Minute)
		if planeCutOff.Before(cutOff) {
			delete(planeList, i)
		}
	}

	planeAccessMutex.Unlock()
}