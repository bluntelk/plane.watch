package sink

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"plane.watch/lib/sink/rabbitmq"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/beast"
	"plane.watch/lib/tracker/mode_s"
	"plane.watch/lib/tracker/sbs1"
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
	}

	frame struct {
		Type, RouteKey string
		Body           []byte
		Source         *tracker.FrameSource
	}

	planeLocation struct {
		New, Removed      bool
		Icao              string
		Lat, Lon, Heading float64
		Altitude          int
		VerticalRate      int
		AltitudeUnits     string
		FlightNumber      string
		FlightStatus      string
		Airframe          string
		HasLocation       bool
		HasHeading        bool
		HasVerticalRate   bool
	}
)

func NewRabbitMqSink(opts ...Option) (*RabbitMqSink, error) {
	r := &RabbitMqSink{
		exchange: "plane.watch.data",
	}
	r.queue = map[string]string{}

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
	return r, nil
}

func WithRabbitVhost(vhost string) Option {
	return func(config *Config) {
		config.vhost = vhost
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
func WithRabbitQueue(msgType, queue string) Option {
	return func(conf *Config) {
		conf.queue[msgType] = queue
	}
}

func (r *RabbitMqSink) OnEvent(e tracker.Event) {
	var err error
	var body []byte
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
				FlightNumber:  plane.FlightNumber(),
				FlightStatus:  plane.FlightStatus(),

				HasLocation:     plane.HasLocation(),
				HasHeading:      plane.HasHeading(),
				HasVerticalRate: plane.HasVerticalRate(),
			}

			var jsonBuf []byte
			jsonBuf, err = json.MarshalIndent(&eventStruct, "", "  ")
			if nil == err {
				err = r.mq.Publish(r.exchange, QueueLocationUpdates, amqp.Publishing{
					ContentType:     "text/plain",
					ContentEncoding: "utf-8",
					Timestamp:       time.Now(),
					Body:            jsonBuf,
				})
			}

		}

	case *tracker.FrameEvent:
		ourFrame := e.(*tracker.FrameEvent).Frame()
		if nil == ourFrame {
			return
		}
		source := e.(*tracker.FrameEvent).Source()

		sendMessage := func(info frame) error {
			body, err = json.Marshal(info)
			if nil != err {
				return err
			}
			return r.mq.Publish(r.exchange, info.RouteKey, amqp.Publishing{
				ContentType:     "text/plain",
				ContentEncoding: "utf-8",
				Timestamp:       time.Now(),
				Body:            body,
			})
		}

		switch ourFrame.(type) {
		case *mode_s.Frame:
			err = sendMessage(frame{Type: "avr", Body: ourFrame.Raw(), RouteKey: QueueTypeAvrAll, Source: source})
		case *beast.Frame:
			err = sendMessage(frame{Type: "beast", Body: ourFrame.Raw(), RouteKey: QueueTypeBeastAll, Source: source})
			err = sendMessage(frame{Type: "avr", Body: ourFrame.(*beast.Frame).AvrFrame().Raw(), RouteKey: QueueTypeAvrAll, Source: source})
		case *sbs1.Frame:
			err = sendMessage(frame{Type: "sbs1", Body: ourFrame.Raw(), RouteKey: QueueTypeSbs1All, Source: source})
		}
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

	//log.Printf("Connecting to RabbitMQ @ %s", rabbitConfig)
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
		if _, err = r.mq.QueueDeclare(q); nil != err {
			return err
		}

		if err = r.mq.QueueBind(q, t, r.exchange); nil != err {
			return err
		}
	}

	return nil
}
