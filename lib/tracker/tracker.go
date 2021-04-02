package tracker

import (
	"fmt"
	"plane.watch/lib/tracker/mode_s"
	"plane.watch/lib/tracker/sbs1"
	"sync"
	"sync/atomic"
	"time"
)

const (
	LogLevelQuiet = 0
	LogLevelError = 1
	LogLevelInfo  = 2
	LogLevelDebug = 3
)

type (
	Tracker struct {
		planeList sync.Map

		// configurable options
		logLevel int

		// pruneTick is how long between pruning attempts
		// pruneAfter is how long we wait from the last message before we remove it from the tracker
		pruneTick, pruneAfter time.Duration

		// Input Handling
		producers   []Producer
		middlewares []Middleware
		sinks       []Sink

		producerWaiter sync.WaitGroup

		decodeWorkerCount   int
		decodingQueue       chan Frame
		decodingQueueWaiter sync.WaitGroup

		eventSync    sync.Mutex
		eventsOpen   bool
		events       chan Event
		eventsWaiter sync.WaitGroup

		pruneExitChan chan bool

		refLat, refLon float64

		startTime time.Time
		numFrames uint64
	}
)

var (
	Levels = [4]string{
		"Quiet",
		"Error",
		"Info",
		"Debug",
	}
)

// NewTracker creates a new tracker with which we can populate with plane tracking data
func NewTracker(opts ...Option) *Tracker {
	t := &Tracker{
		logLevel:          LogLevelQuiet,
		producers:         []Producer{},
		middlewares:       []Middleware{},
		decodeWorkerCount: 5,
		pruneTick:         10 * time.Second,
		pruneAfter:        5 * time.Minute,
		decodingQueue:     make(chan Frame, 1000), // a nice deep buffer
		events:            make(chan Event, 10000),
		eventsOpen:        true,
		pruneExitChan:     make(chan bool),

		startTime: time.Now(),
	}

	for _, opt := range opts {
		opt(t)
	}

	// Process our event queue and send them to all the Sinks that are currently listening to us
	go t.processEvents()

	t.decodingQueueWaiter.Add(t.decodeWorkerCount)
	for i := 0; i < t.decodeWorkerCount; i++ {
		go t.decodeQueue()
	}

	go t.prunePlanes()

	return t
}

func (t *Tracker) debugMessage(sfmt string, a ...interface{}) {
	if t.logLevel >= LogLevelDebug {
		t.AddEvent(NewLogEvent(LogLevelDebug, "Tracker", fmt.Sprintf("DEBUG: "+sfmt, a...)))
	}
}

func (t *Tracker) infoMessage(sfmt string, a ...interface{}) {
	if t.logLevel >= LogLevelInfo {
		t.AddEvent(NewLogEvent(LogLevelInfo, "Tracker", fmt.Sprintf("INFO : "+sfmt, a...)))
	}
}

func (t *Tracker) errorMessage(sfmt string, a ...interface{}) {
	t.AddEvent(NewLogEvent(LogLevelError, "Tracker", fmt.Sprintf("ERROR : "+sfmt, a...)))
}

func (t *Tracker) numPlanes() int {
	count := 0
	t.planeList.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func (t *Tracker) GetPlane(icao uint32) *Plane {
	plane, ok := t.planeList.Load(icao)
	if ok {
		return plane.(*Plane)
	}
	t.infoMessage("Plane %06X has made an appearance", icao)

	p := newPlane(icao)
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
	plane.setLastSeen(frame.TimeStamp())
	plane.incMsgCount()

	debugMessage := func(sfmt string, a ...interface{}) {
		planeFormat = fmt.Sprintf("DF%02d - \033[0;97mPlane (\033[38;5;118m%s %-8s\033[0;97m)", frame.DownLinkType(), plane.IcaoIdentifierStr(), plane.FlightNumber())
		t.debugMessage(planeFormat+sfmt, a...)
	}

	// determine what to do with our given frame
	switch frame.DownLinkType() {
	case 0:
		// grab the altitude
		if frame.AltitudeValid() {
			alt, _ := frame.Altitude()
			plane.setAltitude(alt, frame.AltitudeUnits())
		}
		if frame.VerticalStatusValid() {
			onGround, _ := frame.OnGround()
			plane.setGroundStatus(onGround)
		}
		plane.setLocationUpdateTime(frame.TimeStamp())
		debugMessage(" is at %d %s \033[0m", plane.Altitude(), plane.AltitudeUnits())

		hasChanged = true

	case 1, 2, 3:
		hasChanged = true
		if frame.VerticalStatusValid() {
			onGround, _ := frame.OnGround()
			plane.setGroundStatus(onGround)
		}
		plane.setLocationUpdateTime(frame.TimeStamp())
		if frame.Alert() {
			plane.setSpecial("Alert")
		}
	case 6, 7, 8, 9, 10, 12, 13, 14, 15, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31:
		debugMessage(" \033[38;5;52mIgnoring Mode S Frame: %d (%s)\033[0m\n", frame.DownLinkType(), frame.DownLinkFormat())
		break
	case 11:
		if frame.VerticalStatusValid() {
			onGround, _ := frame.OnGround()
			plane.setGroundStatus(onGround)
			hasChanged = plane.OnGround() != onGround
		}
	case 4, 5:
		hasChanged = true
		if frame.VerticalStatusValid() {
			onGround, _ := frame.OnGround()
			plane.setGroundStatus(onGround)
		}
		if frame.Alert() {
			plane.setSpecial("Alert")
		}
		if frame.AltitudeValid() {
			alt, _ := frame.Altitude()
			plane.setAltitude(alt, frame.AltitudeUnits())
		}
		plane.setLocationUpdateTime(frame.TimeStamp())
		plane.setFlightStatus(frame.FlightStatus(), frame.FlightStatusString())

		if 5 == frame.DownLinkType() || 21 == frame.DownLinkType() {
			plane.setSquawkIdentity(frame.SquawkIdentity())
		}
		hasChanged = true
		debugMessage(" is at %d %s and flight status is: %s. \033[2mMode S Frame: %d \033[0m",
			plane.Altitude(), plane.AltitudeUnits(), plane.FlightStatus(), frame.DownLinkType())
		break
	case 16:
		hasChanged = true
		if frame.AltitudeValid() {
			alt, _ := frame.Altitude()
			plane.setAltitude(alt, frame.AltitudeUnits())
		}
		if frame.VerticalStatusValid() {
			onGround, _ := frame.OnGround()
			plane.setGroundStatus(onGround)
		}
		plane.setLocationUpdateTime(frame.TimeStamp())

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
				plane.setFlightNumber(frame.FlightNumber())
				if frame.ValidCategory() {
					plane.setAirFrameCategory(frame.Category())
				}
				hasChanged = true
				break
			}
		case mode_s.DF17FrameSurfacePos: // "Surface Position"
			{
				if frame.HeadingValid() {
					heading, _ := frame.Heading()
					plane.setHeading(heading)
				}
				if frame.VelocityValid() {
					velocity, _ := frame.Velocity()
					plane.setVelocity(velocity)
				}
				if frame.VerticalStatusValid() {
					onGround, _ := frame.OnGround()
					plane.setGroundStatus(onGround)
				}

				if frame.IsEven() {
					_ = plane.setCprEvenLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				} else {
					_ = plane.setCprOddLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				}
				if err := plane.decodeCpr(t.refLat, t.refLon, frame.TimeStamp()); nil != err {
					debugMessage("%s", err)
				}
				plane.setLocationUpdateTime(frame.TimeStamp())

				debugMessage(" is on the ground and has heading %s and is travelling at %0.2f knots\033[0m", plane.HeadingStr(), plane.Velocity())
				hasChanged = true
				break
			}
		case mode_s.DF17FrameAirPositionBarometric: // "Airborne Position (with Barometric altitude)"
			{
				if frame.VerticalStatusValid() {
					onGround, _ := frame.OnGround()
					plane.setGroundStatus(onGround)
				}
				plane.setLocationUpdateTime(frame.TimeStamp())
				hasChanged = true

				if frame.IsEven() {
					_ = plane.setCprEvenLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				} else {
					_ = plane.setCprOddLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				}

				altitude, _ := frame.Altitude()
				plane.setAltitude(altitude, frame.AltitudeUnits())
				if err := plane.decodeCpr(t.refLat, t.refLon, frame.TimeStamp()); nil != err {
					debugMessage("%s", err)
				}

				dt := plane.DistanceTravelled()
				if dt.Valid() {
					debugMessage(" travelled %0.2f metres %0.2f seconds", dt.metres, dt.duration)
				}

				break
			}
		case mode_s.DF17FrameAirVelocity: // "Airborne velocity"
			{
				if frame.HeadingValid() {
					heading, _ := frame.Heading()
					plane.setHeading(heading)
				}
				if frame.VelocityValid() {
					velocity, _ := frame.Velocity()
					plane.setVelocity(velocity)
				}
				if frame.VerticalStatusValid() {
					onGround, _ := frame.OnGround()
					plane.setGroundStatus(onGround)
				}
				if frame.VerticalRateValid() {
					vr, _ := frame.VerticalRate()
					plane.setVerticalRate(vr)
				}
				plane.setLocationUpdateTime(frame.TimeStamp())

				headingStr := "unknown heading"
				if plane.HasHeading() {
					headingStr = fmt.Sprintf("heading %0.2f", plane.Heading())
				}
				debugMessage(" has %s and is travelling at %0.2f knots\033[0m", headingStr, plane.Velocity())
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
					plane.setSquawkIdentity(frame.SquawkIdentity())
				}
				break
			}
		case mode_s.DF17FrameSurfaceSystemStatus: //, "Surface System status":
			{
				debugMessage("\033[2m Ignoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		case mode_s.DF17FrameEmergencyPriority: //, "Extended Squitter Aircraft status (Emergency)":
			{
				debugMessage("\033[2m %s\033[0m", messageType)
				if frame.Alert() {
					plane.setSpecial(frame.Special() + " " + frame.Emergency())
				}
				plane.setSquawkIdentity(frame.SquawkIdentity())
				hasChanged = true
				break
			}
		case mode_s.DF17FrameTcasRA: //, "Extended Squitter Aircraft status (1090ES TCAS RA)":
			{
				debugMessage("\033[2m Ignoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		case mode_s.DF17FrameTargetStateStatus: //, "Target State and status Message":
			{
				debugMessage("\033[2m Ignoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		case mode_s.DF17FrameAircraftOperational: //, "Aircraft Operational status Message":
			{
				debugMessage("\033[2m Ignoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
				break
			}
		}

	case 20, 21:
		switch frame.BdsMessageType() {
		case mode_s.BdsElsDataLinkCap: // 1.0
			plane.setSquawkIdentity(frame.SquawkIdentity())
		case mode_s.BdsElsGicbCap: // 1.7
			if frame.AltitudeValid() {
				alt, _ := frame.Altitude()
				plane.setAltitude(alt, frame.AltitudeUnits())
			}
		case mode_s.BdsElsAircraftIdent: // 2.0
			plane.setFlightNumber(frame.FlightNumber())
		default:
			// let's see if we can decode more BDS info
			// TODO: Decode Other BDS frames
		}
	}
	t.AddEvent(newPlaneLocationEvent(plane))
	if hasChanged {
		return plane
	} else {
		return nil
	}
}

func (t *Tracker) HandleSbs1Frame(frame *sbs1.Frame) *Plane {
	var hasChanged bool
	plane := t.GetPlane(frame.IcaoInt)
	plane.setLastSeen(frame.TimeStamp())
	plane.incMsgCount()
	if frame.HasPosition {
		if err := plane.addLatLong(frame.Lat, frame.Lon, frame.Received); nil != err {
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

func (t *Tracker) prunePlanes() {
	ticker := time.NewTicker(t.pruneTick)
	for {
		select {
		case <-ticker.C:
			// prune the planes in the list if they have not been seen > 5 minutes
			oldest := time.Now().Add(-t.pruneAfter)
			t.EachPlane(func(p *Plane) bool {
				if p.LastSeen().Before(oldest) {
					t.planeList.Delete(p.icaoIdentifier)

					// now send an event
					t.AddEvent(newPlaneActionEvent(p, false, true))
				}

				return true
			})

			t.AddEvent(t.newInfoEvent())
		case <-t.pruneExitChan:
			return
		}
	}
}

func (t *Tracker) newInfoEvent() *InfoEvent {
	return &InfoEvent{
		receivedFrames: atomic.LoadUint64(&t.numFrames),
		numReceivers:   len(t.producers),
		uptime:         time.Now().Sub(t.startTime).Seconds(),
	}
}
