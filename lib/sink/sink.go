package sink

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"plane.watch/lib/dedupe"
	"plane.watch/lib/export"
	"plane.watch/lib/logging"
	"plane.watch/lib/monitoring"
	"plane.watch/lib/rabbitmq"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/beast"
	"plane.watch/lib/tracker/mode_s"
	"plane.watch/lib/tracker/sbs1"
	"regexp"

	//"regexp"
	"strings"
	"time"
)

type (
	Destination interface {
		PublishJson(queue string, msg []byte) error
		PublishText(queue string, msg []byte) error
		Stop()
		monitoring.HealthCheck
	}

	Sink struct {
		fsm    *dedupe.ForgetfulSyncMap
		config *Config
		dest   Destination
		events chan tracker.Event

		sendFrameAll    func(tracker.Frame, *tracker.FrameSource) error
		sendFrameDedupe func(tracker.Frame, *tracker.FrameSource) error
	}
)

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

func stripAnsi(str string) string {
	return re.ReplaceAllString(str, "")
}

func NewSink(conf *Config, dest Destination) tracker.Sink {
	s := Sink{
		fsm:    dedupe.NewForgetfulSyncMap(10*time.Second, 60*time.Second),
		config: conf,
		dest:   dest,
		events: make(chan tracker.Event),
	}
	if _, ok := s.config.queue[QueueTypeLogs]; ok {
		logging.AddLogDestination(&s)
	}

	s.sendFrameAll = s.sendFrameEvent(QueueTypeAvrAll, QueueTypeBeastAll, QueueTypeSbs1All)
	s.sendFrameDedupe = s.sendFrameEvent(QueueTypeAvrReduce, QueueTypeBeastReduce, QueueTypeSbs1Reduce)
	return &s
}

// Write is for the logs sending
func (s *Sink) Write(b []byte) (int, error) {
	return len(b), s.dest.PublishText(QueueTypeLogs, []byte(stripAnsi(string(b))))
}

func (s *Sink) Listen() chan tracker.Event {
	return s.events
}

func (s *Sink) Stop() {
	close(s.events)
	s.config.Finish()
	s.dest.Stop()
}

func (s *Sink) sendLocationEvent(routingKey string, le *tracker.PlaneLocationEvent) error {
	jsonBuf, err := s.trackerMsgJson(le)
	if nil != jsonBuf && nil == err {
		if s.fsm.HasKey(string(jsonBuf)) {
			// sending a message we have already sent!
			return nil
		}
		s.fsm.AddKey(string(jsonBuf))

		err = s.dest.PublishJson(routingKey, jsonBuf)
	}
	return err
}

func (s *Sink) trackerMsgJson(le *tracker.PlaneLocationEvent) ([]byte, error) {
	var err error
	plane := le.Plane()
	if nil == plane {
		return nil, errors.New("no plane")
	}

	callSign := strings.TrimSpace(plane.FlightNumber())
	eventStruct := export.PlaneLocation{
		New:             le.New(),
		Removed:         le.Removed(),
		Icao:            plane.IcaoIdentifierStr(),
		Lat:             plane.Lat(),
		Lon:             plane.Lon(),
		Heading:         plane.Heading(),
		Altitude:        int(plane.Altitude()),
		VerticalRate:    plane.VerticalRate(),
		AltitudeUnits:   plane.AltitudeUnits(),
		Velocity:        plane.Velocity(),
		CallSign:        &callSign,
		FlightStatus:    plane.FlightStatus(),
		OnGround:        plane.OnGround(),
		Airframe:        plane.AirFrame(),
		AirframeType:    plane.AirFrameType(),
		Squawk:          plane.SquawkIdentityStr(),
		Special:         plane.Special(),
		AircraftWidth:   plane.AirFrameWidth(),
		AircraftLength:  plane.AirFrameLength(),
		Registration:    plane.Registration(),
		HasLocation:     plane.HasLocation(),
		HasHeading:      plane.HasHeading(),
		HasVerticalRate: plane.HasVerticalRate(),
		HasVelocity:     plane.HasVelocity(),
		SourceTag:       s.config.sourceTag,
		TileLocation:    plane.GridTileLocation(),
		LastMsg:         plane.LastSeen().UTC(),
		TrackedSince:    plane.TrackedSince().UTC(),
		SignalRssi:      plane.SignalLevel(),
	}

	var jsonBuf []byte
	jsonBuf, err = json.MarshalIndent(&eventStruct, "", "  ")
	if nil != err {
		log.Error().Err(err).Msg("could not create json bytes for sending")
		return nil, err
	} else {
		return jsonBuf, nil
	}
}

func (s *Sink) sendFrameEvent(queueAvr, queueBeast, queueSbs1 string) func(tracker.Frame, *tracker.FrameSource) error {
	return func(ourFrame tracker.Frame, source *tracker.FrameSource) error {
		var err error
		var body []byte
		if nil == ourFrame {
			return nil
		}

		sendMessage := func(info rabbitFrameMsg) error {
			if _, ok := s.config.queue[info.RouteKey]; !ok {
				return nil
			}
			body, err = json.Marshal(info)
			if nil != err {
				return err
			}
			return s.dest.PublishJson(info.RouteKey, body)
		}

		switch ourFrame.(type) {
		case *mode_s.Frame:
			err = sendMessage(rabbitFrameMsg{Type: "avr", Body: ourFrame.Raw(), RouteKey: queueAvr, Source: source})
		case *beast.Frame:
			err = sendMessage(rabbitFrameMsg{Type: "beast", Body: ourFrame.Raw(), RouteKey: queueBeast, Source: source})
			err = sendMessage(rabbitFrameMsg{Type: "avr", Body: ourFrame.(*beast.Frame).AvrFrame().Raw(), RouteKey: queueAvr, Source: source})
		case *sbs1.Frame:
			err = sendMessage(rabbitFrameMsg{Type: "sbs1", Body: ourFrame.Raw(), RouteKey: queueSbs1, Source: source})
		}
		return err
	}
}

func (s *Sink) OnEvent(e tracker.Event) {
	var err error
	switch e.(type) {
	case *tracker.LogEvent:
		err = s.dest.PublishJson(QueueTypeLogs, []byte(e.String()))

	case *tracker.PlaneLocationEvent:
		le := e.(*tracker.PlaneLocationEvent)
		var jsonBuf []byte
		jsonBuf, err = s.trackerMsgJson(le)
		if nil != jsonBuf && nil == err {
			err = s.dest.PublishJson(QueueLocationUpdates, jsonBuf)
			if nil != s.config.stats.planeLoc {
				s.config.stats.planeLoc.Inc()
			}
		}

	case *tracker.FrameEvent:
		//println("Got a Frame!")
		ourFrame := e.(*tracker.FrameEvent).Frame()
		source := e.(*tracker.FrameEvent).Source()
		err = s.sendFrameAll(ourFrame, source)
		if nil != s.config.stats.frame {
			s.config.stats.frame.Inc()
		}

	case *tracker.DedupedFrameEvent:
		ourFrame := e.(*tracker.DedupedFrameEvent).Frame()
		source := e.(*tracker.DedupedFrameEvent).Source()
		err = s.sendFrameDedupe(ourFrame, source)
		if nil != s.config.stats.dedupeFrame {
			s.config.stats.dedupeFrame.Inc()
		}
	}

	if nil != err {
		log.Error().
			Err(err).
			Str("event-type", e.Type()).
			Str("event", e.String()).
			Msg("Unable to handle event")
	}
	if err == rabbitmq.ErrNilChannel {
		panic(err)
	}
}

func (s *Sink) HealthCheckName() string {
	return s.dest.HealthCheckName()
}

func (s *Sink) HealthCheck() bool {
	return s.dest.HealthCheck()
}
