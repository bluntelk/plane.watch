package sink

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"plane.watch/lib/logging"
	"plane.watch/lib/rabbitmq"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/beast"
	"plane.watch/lib/tracker/mode_s"
	"plane.watch/lib/tracker/sbs1"
	"regexp"
	"strings"
	"time"
)

const (
	QueueTypeBeastAll    = "beast-all"
	QueueTypeBeastReduce = "beast-reduce"
	QueueTypeAvrAll      = "avr-all"
	QueueTypeAvrReduce   = "avr-reduce"
	QueueTypeSbs1All     = "sbs1-all"
	QueueTypeSbs1Reduce  = "sbs1-reduce"
	QueueTypeDecodedJson = "decoded-json"
	QueueTypeLogs        = "logs"
	QueueLocationUpdates = "location-updates"
)

type (
	RabbitMqSink struct {
		Config
		mq       *rabbitmq.RabbitMQ
		exchange string

		sendFrameAll    func(tracker.Frame, *tracker.FrameSource) error
		sendFrameDedupe func(tracker.Frame, *tracker.FrameSource) error
	}

	rabbitFrameMsg struct {
		Type, RouteKey string
		Body           []byte
		Source         *tracker.FrameSource
	}

	planeLocation struct {
		New, Removed      bool
		Icao              string
		Lat, Lon, Heading float64
		Velocity          float64
		Altitude          int
		VerticalRate      int
		AltitudeUnits     string
		FlightNumber      string
		FlightStatus      string
		OnGround          bool
		Airframe          string
		AirframeType      string
		HasLocation       bool
		HasHeading        bool
		HasVerticalRate   bool
		HasVelocity       bool
		SourceTag         string
		Squawk            string
		Special           string
		TrackedSince      time.Time
		LastMsg           time.Time
	}
)

var AllQueues = [...]string{
	QueueTypeBeastAll,
	QueueTypeBeastReduce,
	QueueTypeAvrAll,
	QueueTypeAvrReduce,
	QueueTypeSbs1All,
	QueueTypeSbs1Reduce,
	QueueTypeDecodedJson,
	QueueTypeLogs,
	QueueLocationUpdates,
}

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

func stripAnsi(str string) string {
	return re.ReplaceAllString(str, "")
}
func NewRabbitMqSink(opts ...Option) (*RabbitMqSink, error) {
	r := &RabbitMqSink{
		exchange: "plane.watch.data",
	}
	r.queue = map[string]string{}
	r.sendFrameAll = r.sendFrameEvent(QueueTypeAvrAll, QueueTypeBeastAll, QueueTypeSbs1All)
	r.sendFrameDedupe = r.sendFrameEvent(QueueTypeAvrReduce, QueueTypeBeastReduce, QueueTypeSbs1Reduce)

	for _, opt := range opts {
		opt(&r.Config)
	}

	var err error
	if r.mq, err = r.connect(time.Second * 5); nil != err {
		return nil, err
	}

	if err = r.setup(); nil != err {
		return nil, err
	}

	if _, ok := r.queue[QueueTypeLogs]; ok {
		logging.AddLogDestination(r)
	}

	// setup a hook for messages
	return r, nil
}

func (r *RabbitMqSink) Write(b []byte) (int, error) {
	err := r.mq.Publish(r.exchange, QueueTypeLogs, amqp.Publishing{
		ContentType:     "text/plain",
		ContentEncoding: "utf-8",
		Timestamp:       time.Now(),
		Body:            []byte(stripAnsi(string(b))),
	})
	return len(b), err
}

func WithRabbitVhost(vhost string) Option {
	return func(config *Config) {
		config.vhost = vhost
	}
}

func WithRabbitQueues(queues []string) Option {
	return func(conf *Config) {
		if 0 == len(queues) {
			WithAllRabbitQueues()(conf)
			return
		}

		for _, requestedQueue := range queues {
			found := false
			for _, validQueue := range AllQueues {
				if requestedQueue == validQueue {
					conf.queue[validQueue] = validQueue
					found = true
					break
				}
			}
			if !found {
				log.Error().Msgf("Error: Unknown Queue Type: %s", requestedQueue)
			}
		}
	}
}
func WithAllRabbitQueues() Option {
	return func(conf *Config) {
		conf.queue[QueueTypeAvrAll] = QueueTypeAvrAll
		conf.queue[QueueTypeAvrReduce] = QueueTypeAvrReduce
		conf.queue[QueueTypeBeastAll] = QueueTypeBeastAll
		conf.queue[QueueTypeBeastReduce] = QueueTypeBeastReduce
		conf.queue[QueueTypeSbs1All] = QueueTypeSbs1All
		conf.queue[QueueTypeSbs1Reduce] = QueueTypeSbs1Reduce
		conf.queue[QueueTypeDecodedJson] = QueueTypeDecodedJson
		conf.queue[QueueTypeLogs] = QueueTypeLogs
		conf.queue[QueueLocationUpdates] = QueueLocationUpdates
	}
}

func (r *RabbitMqSink) Stop() {
	r.mq.Disconnect()
}

func (r *RabbitMqSink) sendLocationEventToQueue(queue string, le *tracker.PlaneLocationEvent) error {
	var err error
	plane := le.Plane()
	if nil != plane {
		eventStruct := planeLocation{
			New:           le.New(),
			Removed:       le.Removed(),
			Icao:          plane.IcaoIdentifierStr(),
			Lat:           plane.Lat(),
			Lon:           plane.Lon(),
			Heading:       plane.Heading(),
			Altitude:      int(plane.Altitude()),
			VerticalRate:  plane.VerticalRate(),
			AltitudeUnits: plane.AltitudeUnits(),
			Velocity:      plane.Velocity(),
			FlightNumber:  strings.TrimSpace(plane.FlightNumber()),
			FlightStatus:  plane.FlightStatus(),
			OnGround:      plane.OnGround(),
			Airframe:      plane.AirFrame(),
			AirframeType:  plane.AirFrameType(),
			Squawk:        plane.SquawkIdentityStr(),
			Special:       plane.Special(),

			HasLocation:     plane.HasLocation(),
			HasHeading:      plane.HasHeading(),
			HasVerticalRate: plane.HasVerticalRate(),
			HasVelocity:     plane.HasVelocity(),
			SourceTag:       r.Config.sourceTag,
			LastMsg:         plane.LastSeen().UTC(),
			TrackedSince:    plane.TrackedSince().UTC(),
		}

		var jsonBuf []byte
		jsonBuf, err = json.MarshalIndent(&eventStruct, "", "  ")
		if nil == err {
			err = r.mq.Publish(r.exchange, queue, amqp.Publishing{
				ContentType:     "application/json",
				ContentEncoding: "utf-8",
				Timestamp:       time.Now(),
				Body:            jsonBuf,
			})
		}
	}
	return err
}

func (r *RabbitMqSink) sendFrameEvent(queueAvr, queueBeast, queueSbs1 string) func(tracker.Frame, *tracker.FrameSource) error {
	return func(ourFrame tracker.Frame, source *tracker.FrameSource) error {
		var err error
		var body []byte
		if nil == ourFrame {
			return nil
		}

		sendMessage := func(info rabbitFrameMsg) error {
			if _, ok := r.queue[info.RouteKey]; !ok {
				return nil
			}
			body, err = json.Marshal(info)
			if nil != err {
				return err
			}
			return r.mq.Publish(r.exchange, info.RouteKey, amqp.Publishing{
				//ContentType:     "text/plain",
				ContentType:     "application/json",
				ContentEncoding: "utf-8",
				Timestamp:       time.Now(),
				Body:            body,
			})
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

func (r *RabbitMqSink) OnEvent(e tracker.Event) {
	var err error
	switch e.(type) {
	case *tracker.LogEvent:
		err = r.mq.Publish(r.exchange, QueueTypeLogs, amqp.Publishing{
			ContentType:     "text/plain",
			ContentEncoding: "utf-8",
			Timestamp:       time.Now(),
			Body:            []byte(e.String()),
		})
	case *tracker.PlaneLocationEvent:
		le := e.(*tracker.PlaneLocationEvent)
		err = r.sendLocationEventToQueue(QueueLocationUpdates, le)

	case *tracker.FrameEvent:
		//println("Got a Frame!")
		ourFrame := e.(*tracker.FrameEvent).Frame()
		source := e.(*tracker.FrameEvent).Source()
		err = r.sendFrameAll(ourFrame, source)

	case *tracker.DedupedFrameEvent:
		ourFrame := e.(*tracker.DedupedFrameEvent).Frame()
		source := e.(*tracker.DedupedFrameEvent).Source()
		err = r.sendFrameDedupe(ourFrame, source)
	}

	if nil != err {
		fmt.Println(err)
	}
}

func (r *RabbitMqSink) connect(timeout time.Duration) (*rabbitmq.RabbitMQ, error) {
	var rabbitConfig rabbitmq.Config
	rabbitConfig.Host = r.host
	rabbitConfig.Port = r.port
	rabbitConfig.User = r.user
	rabbitConfig.Password = r.pass
	rabbitConfig.Vhost = r.vhost

	log.Info().Str("host", rabbitConfig.String()).Msg("Connecting to RabbitMQ")
	rabbit := rabbitmq.New(rabbitConfig)
	connected := make(chan bool)
	go rabbit.Connect(connected)
	select {
	case <-connected:
		return rabbit, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("failed to connect to rabbit in a timely manner")
	}
}

func (r *RabbitMqSink) setup() error {
	var err error

	// let's make sure all of our queues and exchanges are setup
	if err = r.mq.ExchangeDeclare(r.exchange, amqp.ExchangeDirect); nil != err {
		return err
	}
	for t, q := range r.queue {
		if _, err = r.mq.QueueDeclare(q, r.messageTtlSeconds*1000); nil != err {
			return err
		}

		if err = r.mq.QueueBind(q, t, r.exchange); nil != err {
			return err
		}
	}

	return nil
}
