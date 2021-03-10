package tracker

import (
	"fmt"
	"io"
	"plane.watch/lib/tracker/mode_s"
	"plane.watch/lib/tracker/sbs1"
	"sync"
)

const (
	LogLevelQuiet = 0
	LogLevelError = 1
	LogLevelInfo  = 2
	LogLevelDebug = 3
)

type Tracker struct {
	logs      io.Writer
	planeList sync.Map

	logLevel  int

	// Input Handling
	producers   []Producer
	middlewares []Middleware

	producerWaiter sync.WaitGroup

	decodeWorkerCount   uint
	decodingQueue       chan Frame
	decodingQueueWaiter sync.WaitGroup
}

var DefaultTracker = NewTracker()

// NewTracker creates a new tracker with which we can populate with plane tracking data
func NewTracker(opts ...Option) *Tracker {
	t := &Tracker{
		logs:     io.Discard,
		logLevel: LogLevelQuiet,
		producers:     []Producer{},
		middlewares:   []Middleware{},
		decodingQueue: make(chan Frame, 1000), // a nice deep buffer
	}

	for _, opt := range opts {
		opt(t)
	}

	for i := 0; i < 5; i++ {
		go t.decodeQueue()
	}

	return t
}

func (t *Tracker) SetLoggerOutput(out io.Writer) {
	t.logs = out
}

func (t *Tracker) debugMessage(sfmt string, a ...interface{}) {
	if t.logLevel >= LogLevelDebug {
		_, _ = fmt.Fprintf(t.logs, "DEBUG: "+sfmt+"\n", a...)
	}
}

func (t *Tracker) infoMessage(sfmt string, a ...interface{}) {
	if t.logLevel >= LogLevelInfo {
		_, _ = fmt.Fprintf(t.logs, "INFO : "+sfmt+"\n", a...)
	}
}

func (t *Tracker) errorMessage(sfmt string, a ...interface{}) {
	_, _ = fmt.Fprintf(t.logs, "ERROR: "+sfmt+"\n", a...)
}

func SetLoggerOutput(out io.Writer) { DefaultTracker.SetLoggerOutput(out) }

func HandleModeSFrame(frame *mode_s.Frame) *Plane { return DefaultTracker.HandleModeSFrame(frame) }

func (t *Tracker) GetPlane(icao uint32) *Plane {
	plane, ok := t.planeList.Load(icao)
	if ok {
		return plane.(*Plane)
	}

	p := NewPlane(icao)
	t.planeList.Store(icao, p)
	return p
}

func (t *Tracker) EachPlane(pi PlaneIterator) {
	t.planeList.Range(func(key, value interface{}) bool {
		return pi(value.(*Plane))
	})
}

func (t *Tracker) HandleModeSFrame(frame *mode_s.Frame) *Plane {
	if nil == frame {
		return nil
	}
	icao := frame.Icao()
	if 0 == icao {
		return nil
	}
	var planeFormat string
	var hasChanged bool

	plane := t.GetPlane(icao)
	plane.SetLastSeen(frame.TimeStamp())
	//plane.MarkFrameTime(frame.TimeStamp()) // todo, make this fast!

	debugMessage := func(sfmt string, a ...interface{}) {
		planeFormat = fmt.Sprintf("DF%02d - \033[0;97mPlane (\033[38;5;118m%s %-8s\033[0;97m)", frame.DownLinkType(), plane.Icao, plane.Flight.Identifier)
		t.debugMessage(planeFormat+sfmt, a...)
	}

	// determine what to do with our given frame
	switch frame.DownLinkType() {
	case 0:
		// grab the altitude
		if frame.AltitudeValid() {
			alt, _ := frame.Altitude()
			plane.SetAltitude(alt, frame.AltitudeUnits())
		}
		if frame.VerticalStatusValid() {
			onGround, _ := frame.OnGround()
			plane.SetGroundStatus(onGround)
		}
		plane.SetLocationUpdateTime(frame.TimeStamp())
		debugMessage(" is at %d %s \033[0m", plane.Location.Altitude, plane.Location.AltitudeUnits)

		hasChanged = true

	case 1, 2, 3:
		hasChanged = true
		if frame.VerticalStatusValid() {
			onGround, _ := frame.OnGround()
			plane.SetGroundStatus(onGround)
		}
		plane.SetLocationUpdateTime(frame.TimeStamp())
		if frame.Alert() {
			plane.SetSpecial("Alert")
		}
	case 6, 7, 8, 9, 10, 12, 13, 14, 15, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31:
		debugMessage(" \033[38;5;52mIgnoring Mode S Frame: %d (%s)\033[0m\n", frame.DownLinkType(), frame.DownLinkFormat())
		break
	case 11:
		if frame.VerticalStatusValid() {
			onGround, _ := frame.OnGround()
			plane.SetGroundStatus(onGround)
			hasChanged = plane.GroundStatus() != onGround
		}
	case 4, 5:
		hasChanged = true
		if frame.VerticalStatusValid() {
			onGround, _ := frame.OnGround()
			plane.SetGroundStatus(onGround)
		}
		if frame.Alert() {
			plane.SetSpecial("Alert")
		}
		if frame.AltitudeValid() {
			alt, _ := frame.Altitude()
			plane.SetAltitude(alt, frame.AltitudeUnits())
		}
		plane.SetLocationUpdateTime(frame.TimeStamp())
		plane.SetFlightStatus(frame.FlightStatus(), frame.FlightStatusString())

		if 5 == frame.DownLinkType() || 21 == frame.DownLinkType() {
			plane.SetSquawkIdentity(frame.SquawkIdentity())
		}
		hasChanged = true
		debugMessage(" is at %d %s and flight status is: %s. \033[2mMode S Frame: %d \033[0m",
			plane.Altitude(), plane.AltitudeUnits(), plane.FlightStatus(), frame.DownLinkType())
		break
	case 16:
		hasChanged = true
		if frame.AltitudeValid() {
			alt, _ := frame.Altitude()
			plane.SetAltitude(alt, frame.AltitudeUnits())
		}
		if frame.VerticalStatusValid() {
			onGround, _ := frame.OnGround()
			plane.SetGroundStatus(onGround)
		}
		plane.SetLocationUpdateTime(frame.TimeStamp())

	case 17, 18: // ADS-B
		//if debug {
		//	frame.Describe(os.Stdout)
		//}

		// i am using the text version because it is easier to program with.
		// if performance is an issue, change over to byte comparing
		messageType := frame.MessageTypeString()
		switch messageType {
		case mode_s.DF17FrameIdCat: // "Aircraft Identification and Category"
			{
				plane.SetFlightIdentifier(frame.FlightNumber())
				if frame.ValidCategory() {
					plane.SetAirFrameCategory(frame.Category())
				}
				hasChanged = true
				break
			}
		case mode_s.DF17FrameSurfacePos: // "Surface Position"
			{
				if frame.HeadingValid() {
					heading, _ := frame.Heading()
					plane.SetHeading(heading)
				}
				if frame.VelocityValid() {
					velocity, _ := frame.Velocity()
					plane.SetVelocity(velocity)
				}
				if frame.VerticalStatusValid() {
					onGround, _ := frame.OnGround()
					plane.SetGroundStatus(onGround)
				}

				if frame.IsEven() {
					_ = plane.SetCprEvenLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				} else {
					_ = plane.SetCprOddLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				}
				if err := plane.DecodeCpr(frame.TimeStamp()); nil != err {
					debugMessage("%s", err)
				}
				plane.SetLocationUpdateTime(frame.TimeStamp())

				heading := "unknown heading"
				if plane.HasHeading() {heading = fmt.Sprintf("heading %0.2f", plane.Heading())}
				debugMessage(" is on the ground and has %s and is travelling at %0.2f knots\033[0m", heading, plane.Velocity())
				hasChanged = true
				break
			}
		case mode_s.DF17FrameAirPositionBarometric: // "Airborne Position (with Barometric Altitude)"
			{
				if frame.VerticalStatusValid() {
					onGround, _ := frame.OnGround()
					plane.SetGroundStatus(onGround)
				}
				plane.SetLocationUpdateTime(frame.TimeStamp())
				hasChanged = true

				if frame.IsEven() {
					_ = plane.SetCprEvenLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				} else {
					_ = plane.SetCprOddLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				}

				altitude, _ := frame.Altitude()
				plane.SetAltitude(altitude, frame.AltitudeUnits())
				if err := plane.DecodeCpr(frame.TimeStamp()); nil != err {
					debugMessage("%s", err)
				}

				dt := plane.DistanceTravelled()
				if dt.Valid() {
					debugMessage(" travelled %0.2f metres %0.2f seconds", dt.metres, dt.duration)
				}

				break
			}
		case mode_s.DF17FrameAirVelocity: // "Airborne Velocity"
			{
				if frame.HeadingValid() {
					heading, _ := frame.Heading()
					plane.SetHeading(heading)
				}
				if frame.VelocityValid() {
					velocity, _ := frame.Velocity()
					plane.SetVelocity(velocity)
				}
				if frame.VerticalStatusValid() {
					onGround, _ := frame.OnGround()
					plane.SetGroundStatus(onGround)
				}
				if frame.VerticalRateValid() {
					vr, _ := frame.VerticalRate()
					plane.SetVerticalRate(vr)
				}
				plane.SetLocationUpdateTime(frame.TimeStamp())

				heading := "unknown heading"
				if plane.HasHeading() {heading = fmt.Sprintf("heading %0.2f", plane.Heading())}
				debugMessage(" has %s and is travelling at %0.2f knots\033[0m", heading, plane.Velocity())
				hasChanged = true
				break
			}
		case mode_s.DF17FrameAirPositionGnss: // "Airborne Position (GNSS Height)"
			{
				debugMessage("\033[2m Ignoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		case mode_s.DF17FrameTestMessage: //, "Test Message":
			debugMessage("\033[2m Ignoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
			break
		case mode_s.DF17FrameTestMessageSquawk: //, "Test Message":
			{
				if frame.SquawkIdentity() > 0 {
					hasChanged = plane.SquawkIdentity() != frame.SquawkIdentity()
					plane.SetSquawkIdentity(frame.SquawkIdentity())
				}
				break
			}
		case mode_s.DF17FrameSurfaceSystemStatus: //, "Surface System Status":
			{
				debugMessage("\033[2m Ignoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		case mode_s.DF17FrameEmergencyPriority: //, "Extended Squitter Aircraft Status (Emergency)":
			{
				debugMessage("\033[2m %s\033[0m", messageType)
				plane.SetSpecial("Emergency")
				plane.SetSquawkIdentity(frame.SquawkIdentity())
				hasChanged = true
				break
			}
		case mode_s.DF17FrameTcasRA: //, "Extended Squitter Aircraft Status (1090ES TCAS RA)":
			{
				debugMessage("\033[2m Ignoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		case mode_s.DF17FrameTargetStateStatus: //, "Target State and Status Message":
			{
				debugMessage("\033[2m Ignoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		case mode_s.DF17FrameAircraftOperational: //, "Aircraft Operational Status Message":
			{
				debugMessage("\033[2m Ignoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		}

	case 20, 21:
		switch frame.BdsMessageType() {
		case mode_s.BdsElsDataLinkCap: // 1.0
			plane.SetSquawkIdentity(frame.SquawkIdentity())
		case mode_s.BdsElsGicbCap: // 1.7
			if frame.AltitudeValid() {
				alt, _ := frame.Altitude()
				plane.SetAltitude(alt, frame.AltitudeUnits())
			}
		case mode_s.BdsElsAircraftIdent: // 2.0
			plane.SetFlightIdentifier(frame.FlightNumber())
		default:
			// let's see if we can decode more BDS info
			// TODO: Decode Other BDS frames
		}
	}
	if hasChanged {
		return plane
	} else {
		return nil
	}
}

func HandleSbs1Frame(frame *sbs1.Frame) *Plane {
	return DefaultTracker.HandleSbs1Frame(frame)
}

func (t *Tracker) HandleSbs1Frame(frame *sbs1.Frame) *Plane {
	var hasChanged bool
	plane := t.GetPlane(frame.IcaoInt)
	plane.SetLastSeen(frame.TimeStamp())
	if frame.HasPosition {
		if err := plane.AddLatLong(frame.Lat, frame.Lon, frame.Received); nil != err {
			t.debugMessage("%s", err)
		}
		hasChanged = true
		t.debugMessage("Plane %s is at %0.4f, %0.4f", frame.IcaoStr(), frame.Lat, frame.Lon)
	}
	if hasChanged {
		return plane
	}
	return nil
}
