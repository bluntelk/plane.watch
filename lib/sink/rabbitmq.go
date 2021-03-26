package sink

import (
	"plane.watch/lib/tracker"
)

type (
	RabbitMqSink struct {
		Config
	}
)

func NewRabbitMqSink(opts ...Option) *RabbitMqSink {
	r := &RabbitMqSink{}
	for _, opt := range opts {
		opt(&r.Config)
	}
	return r
}

func (r *RabbitMqSink) OnEvent(e tracker.Event) {

}



//func getRabbitConnection(timeout int64) (*RabbitMQ, error) {
//	if "" == pwUser {
//		log.Fatalln("No User Specified For Plane.Watch Rabbit Config")
//	}
//	if "" == pwPass {
//		log.Fatalln("No Password Specified For Plane.Watch Rabbit Config")
//	}
//
//	var rabbitConfig RabbitMQConfig
//	rabbitConfig.Host = pwHost
//	rabbitConfig.Port = pwPort
//	rabbitConfig.User = pwUser
//	rabbitConfig.Password = pwPass
//	rabbitConfig.Vhost = pwVhost
//
//	log.Printf("Connecting to plane.watch RabbitMQ @ %s", rabbitConfig)
//	rabbit := NewRabbitMQ(rabbitConfig)
//	connected := make(chan bool)
//	go rabbit.Connect(connected)
//	select {
//	case <-connected:
//		return rabbit, nil
//	case <-time.After(time.Duration(timeout) * time.Second):
//		return nil, fmt.Errorf("failed to connect to rabbit in a timely manner")
//	}
//}
//
//func failOnError(err error, msg string) {
//	if err != nil {
//		log.Fatalf("%s: %s", msg, err)
//		//panic(fmt.Sprintf("%s: %s", msg, err))
//	}
//}
//
//// test makes sure that our setup is working
//func test(c *cli.Context) {
//	log.Printf("Testing connection to dump1090 @ %s:%s", dump1090Host, dump1090Port)
//	d1090 := NewDump1090Reader(dump1090Host, dump1090Port)
//	var err error
//	if err = d1090.Connect(); err != nil {
//		log.Fatalf("Unable to connect to Dump 1090 %s", err)
//	} else {
//		d1090.Stop()
//	}
//
//	rabbit, err := getRabbitConnection(10)
//	failOnError(err, "Unable to connect to rabbit")
//	defer rabbit.Disconnect()
//
//	log.Printf("Success. You are ready to go")
//}