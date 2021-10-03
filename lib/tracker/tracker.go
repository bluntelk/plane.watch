package tracker

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"plane.watch/lib/tracker/mode_s"
	"plane.watch/lib/tracker/sbs1"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type (
	Tracker struct {
		planeList sync.Map

		// pruneTick is how long between pruning attempts
		// pruneAfter is how long we wait from the last message before we remove it from the tracker
		pruneTick, pruneAfter time.Duration

		// Input Handling
		producers   []Producer
		middlewares []Middleware
		sinks       []Sink

		producerWaiter   sync.WaitGroup
		middlewareWaiter sync.WaitGroup

		decodeWorkerCount   int
		decodingQueue       chan *FrameEvent
		decodingQueueWaiter sync.WaitGroup

		eventSync    sync.RWMutex
		eventsOpen   bool
		events       chan Event
		eventsWaiter sync.WaitGroup

		pruneExitChan chan bool

		startTime time.Time
		numFrames uint64
	}
)

// NewTracker creates a new tracker with which we can populate with plane tracking data
func NewTracker(opts ...Option) *Tracker {
	t := &Tracker{
		producers:         []Producer{},
		middlewares:       []Middleware{},
		decodeWorkerCount: 5,
		pruneTick:         10 * time.Second,
		pruneAfter:        5 * time.Minute,
		decodingQueue:     make(chan *FrameEvent, 1000), // a nice deep buffer
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
	log.Debug().Str("section", "Tracker").Msgf(sfmt, a...)
}

func (t *Tracker) infoMessage(sfmt string, a ...interface{}) {
	log.Info().Str("section", "Tracker").Msgf(sfmt, a...)
}

func (t *Tracker) errorMessage(sfmt string, a ...interface{}) {
	log.Error().Str("section", "Tracker").Msgf(sfmt, a...)
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
	p.tracker = t
	t.planeList.Store(icao, p)
	return p
}

func (t *Tracker) EachPlane(pi PlaneIterator) {
	t.planeList.Range(func(key, value interface{}) bool {
		return pi(value.(*Plane))
	})
}

func (p *Plane) HandleModeSFrame(frame *mode_s.Frame, refLat, refLon *float64) {
	if nil == frame {
		return
	}
	icao := frame.Icao()
	if 0 == icao {
		return
	}
	var planeFormat string
	var hasChanged bool

	p.setLastSeen(frame.TimeStamp())
	p.incMsgCount()

	debugMessage := func(sfmt string, a ...interface{}) {
		planeFormat = fmt.Sprintf("DF%02d - \033[0;97mPlane (\033[38;5;118m%s %-8s\033[0;97m)", frame.DownLinkType(), p.IcaoIdentifierStr(), p.FlightNumber())
		p.tracker.debugMessage(planeFormat+sfmt, a...)
	}

	log.Trace().
		Str("frame", frame.String()).
		Str("icao", frame.IcaoStr()).
		Str("Downlink Type", "DF"+strconv.Itoa(int(frame.DownLinkType()))).
		Int("Downlink Format", int(frame.DownLinkType())).
		Str("DF17 Msg Type", frame.MessageTypeString()).
		Send()

	// determine what to do with our given frame
	switch frame.DownLinkType() {
	case 0:
		// grab the altitude
		if frame.AltitudeValid() {
			alt, _ := frame.Altitude()
			hasChanged = hasChanged || p.setAltitude(alt, frame.AltitudeUnits())
		}
		if frame.VerticalStatusValid() {
			hasChanged = hasChanged || p.setGroundStatus(frame.MustOnGround())
		}
		p.setLocationUpdateTime(frame.TimeStamp())
		debugMessage(" is at %d %s \033[0m", p.Altitude(), p.AltitudeUnits())

	case 1, 2, 3:
		if frame.VerticalStatusValid() {
			hasChanged = hasChanged || p.setGroundStatus(frame.MustOnGround())
		}
		p.setLocationUpdateTime(frame.TimeStamp())
		if frame.Alert() {
			hasChanged = hasChanged || p.setSpecial("alert", "Alert")
		}
	case 6, 7, 8, 9, 10, 12, 13, 14, 15, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31:
		debugMessage(" \033[38;5;52mIgnoring Mode S Frame: %d (%s)\033[0m\n", frame.DownLinkType(), frame.DownLinkFormat())
		break
	case 11:
		if frame.VerticalStatusValid() {
			hasChanged = hasChanged || p.setGroundStatus(frame.MustOnGround())
		}
	case 4, 5:
		if frame.VerticalStatusValid() {
			hasChanged = hasChanged || p.setGroundStatus(frame.MustOnGround())
		}
		if frame.Alert() {
			hasChanged = hasChanged || p.setSpecial("alert", "Alert")
		}
		if frame.AltitudeValid() {
			alt, _ := frame.Altitude()
			hasChanged = hasChanged || p.setAltitude(alt, frame.AltitudeUnits())
		}
		hasChanged = hasChanged || p.setFlightStatus(frame.FlightStatus(), frame.FlightStatusString())

		if 5 == frame.DownLinkType() { // || 21 == frame.DownLinkType()
			hasChanged = hasChanged || p.setSquawkIdentity(frame.SquawkIdentity())
		}

		p.setLocationUpdateTime(frame.TimeStamp())

		debugMessage(" is at %d %s and flight status is: %s. \033[2mMode S Frame: %d \033[0m",
			p.Altitude(), p.AltitudeUnits(), p.FlightStatus(), frame.DownLinkType())
		break
	case 16:
		if frame.AltitudeValid() {
			alt, _ := frame.Altitude()
			hasChanged = hasChanged || p.setAltitude(alt, frame.AltitudeUnits())
		}
		if frame.VerticalStatusValid() {
			hasChanged = hasChanged || p.setGroundStatus(frame.MustOnGround())
		}
		p.setLocationUpdateTime(frame.TimeStamp())

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
				hasChanged = hasChanged || p.setFlightNumber(frame.FlightNumber())
				if frame.ValidCategory() {
					hasChanged = hasChanged || p.setAirFrameCategory(frame.Category())
					hasChanged = hasChanged || p.setAirFrameCategoryType(frame.CategoryType())
				}
				break
			}
		case mode_s.DF17FrameSurfacePos: // "Surface Position"
			{
				if frame.HeadingValid() {
					hasChanged = hasChanged || p.setHeading(frame.MustHeading())
				}
				if frame.VelocityValid() {
					hasChanged = hasChanged || p.setVelocity(frame.MustVelocity())
				}
				if frame.VerticalStatusValid() {
					hasChanged = hasChanged || p.setGroundStatus(frame.MustOnGround())
				}

				if frame.IsEven() {
					_ = p.setCprEvenLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				} else {
					_ = p.setCprOddLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				}
				if err := p.decodeCprFilledRefLatLon(refLat, refLon, frame.TimeStamp()); nil != err {
					debugMessage("%s", err)
				} else {
					hasChanged = true
				}
				p.setLocationUpdateTime(frame.TimeStamp())

				debugMessage(" is on the ground and has heading %s and is travelling at %0.2f knots\033[0m", p.HeadingStr(), p.Velocity())
				break
			}
		case mode_s.DF17FrameAirPositionBarometric, mode_s.DF17FrameAirPositionGnss: // "Airborne Position (with Barometric altitude)"
			{
				if frame.VerticalStatusValid() {
					hasChanged = hasChanged || p.setGroundStatus(frame.MustOnGround())
				}
				p.setLocationUpdateTime(frame.TimeStamp())

				if frame.IsEven() {
					_ = p.setCprEvenLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				} else {
					_ = p.setCprOddLocation(float64(frame.Latitude()), float64(frame.Longitude()), frame.TimeStamp())
				}

				altitude, _ := frame.Altitude()
				p.setAltitude(altitude, frame.AltitudeUnits())
				if err := p.decodeCpr(0, 0, frame.TimeStamp()); nil != err {
					debugMessage("%s", err)
				} else {
					hasChanged = true
				}

				if dt := p.DistanceTravelled(); dt.Valid() {
					debugMessage(" travelled %0.2f metres %0.2f seconds", dt.metres, dt.duration)
				}

				if frame.HasSurveillanceStatus() {
					hasChanged = hasChanged || p.setSpecial("surveillance", frame.SurveillanceStatus())
				} else {
					hasChanged = hasChanged || p.setSpecial("surveillance", "")
				}

				break
			}
		case mode_s.DF17FrameAirVelocity: // "Airborne velocity"
			{
				if frame.HeadingValid() {
					hasChanged = hasChanged || p.setHeading(frame.MustHeading())
				}
				if frame.VelocityValid() {
					hasChanged = hasChanged || p.setVelocity(frame.MustVelocity())
				}
				if frame.VerticalStatusValid() {
					hasChanged = hasChanged || p.setGroundStatus(frame.MustOnGround())
				}
				if frame.VerticalRateValid() {
					hasChanged = hasChanged || p.setVerticalRate(frame.MustVerticalRate())
				}
				p.setLocationUpdateTime(frame.TimeStamp())

				headingStr := "unknown heading"
				if p.HasHeading() {
					headingStr = fmt.Sprintf("heading %0.2f", p.Heading())
				}
				debugMessage(" has %s and is travelling at %0.2f knots\033[0m", headingStr, p.Velocity())
				break
			}
		case mode_s.DF17FrameTestMessage: //, "Test Message":
			debugMessage("\033[2m Ignoring: DF%d %s\033[0m", frame.DownLinkType(), messageType)
			break
		case mode_s.DF17FrameTestMessageSquawk: //, "Test Message":
			{
				if frame.SquawkIdentity() > 0 {
					hasChanged = hasChanged || p.setSquawkIdentity(frame.SquawkIdentity())
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
					hasChanged = hasChanged || p.setSpecial("special", frame.Special())
					hasChanged = hasChanged || p.setSpecial("emergency", frame.Emergency())
				}
				hasChanged = hasChanged || p.setSquawkIdentity(frame.SquawkIdentity())
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
				if frame.VerticalStatusValid() {
					hasChanged = hasChanged || p.setGroundStatus(frame.MustOnGround())
				}

				break
			}
		}

	case 20, 21:
		switch frame.BdsMessageType() {
		case mode_s.BdsElsDataLinkCap: // 1.0
			hasChanged = hasChanged || p.setSquawkIdentity(frame.SquawkIdentity())
		case mode_s.BdsElsGicbCap: // 1.7
			if frame.AltitudeValid() {
				hasChanged = hasChanged || p.setAltitude(frame.MustAltitude(), frame.AltitudeUnits())
			}
		case mode_s.BdsElsAircraftIdent: // 2.0
			hasChanged = hasChanged || p.setFlightNumber(frame.FlightNumber())
		default:
			// let's see if we can decode more BDS info
			// TODO: Decode Other BDS frames
		}
	}

	if hasChanged {
		p.tracker.AddEvent(newPlaneLocationEvent(p))
	}
}

func (p *Plane) HandleSbs1Frame(frame *sbs1.Frame) {
	var hasChanged bool
	p.setLastSeen(frame.TimeStamp())
	p.incMsgCount()
	if frame.HasPosition {
		if err := p.addLatLong(frame.Lat, frame.Lon, frame.Received); nil != err {
			p.tracker.debugMessage("%s", err)
		}

		hasChanged = true
		p.tracker.debugMessage("Plane %s is at %0.4f, %0.4f", frame.IcaoStr(), frame.Lat, frame.Lon)
	}

	if hasChanged {
		p.tracker.AddEvent(newPlaneLocationEvent(p))
	}
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
