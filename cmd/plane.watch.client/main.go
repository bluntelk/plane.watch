package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/urfave/cli"
	"log"
	"os"
	"plane.watch/lib/producer"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/mode_s"
	"time"
)

var (
	pwHost, pwUser, pwPass, pwVhost string
	pwPort                          int
	showDebug                       bool
	dump1090Host                    string
	dump1090Port                    string
)

func main() {

	app := cli.NewApp()

	app.Version = "1.0.0"
	app.Name = "Plane Watch Client"
	app.Usage = "Reads from dump1090 and sends it to http://plane.watch/"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "pw_host",
			Value:       "mq.plane.watch",
			Usage:       "How we connect to plane.watch",
			Destination: &pwHost,
			EnvVar:      "PW_HOST",
		},
		cli.StringFlag{
			Name:        "pw_user",
			Value:       "",
			Usage:       "user for plane.watch",
			Destination: &pwUser,
			EnvVar:      "PW_USER",
		},
		cli.StringFlag{
			Name:        "pw_pass",
			Value:       "",
			Usage:       "plane.watch password",
			Destination: &pwPass,
			EnvVar:      "PW_PASS",
		},
		cli.IntFlag{
			Name:        "pw_port",
			Value:       5672,
			Usage:       "How we connect to plane.watch",
			Destination: &pwPort,
			EnvVar:      "PW_PORT",
		},
		cli.StringFlag{
			Name:        "pw_vhost",
			Value:       "/pw_feedin",
			Usage:       "the virtual host on the plane watch rabbit server",
			Destination: &pwVhost,
			EnvVar:      "PW_VHOST",
		},
		cli.StringFlag{
			Name:        "dump1090_host",
			Value:       "localhost",
			Usage:       "The host to read dump1090 from",
			Destination: &dump1090Host,
			EnvVar:      "DUMP1090_HOST",
		},
		cli.StringFlag{
			Name:        "dump1090_port",
			Value:       "30002",
			Usage:       "The port to read dump1090 from",
			Destination: &dump1090Port,
			EnvVar:      "DUMP1090_PORT",
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Show Extra Debug Information",
			Destination: &showDebug,
			EnvVar:      "DEBUG",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "test",
			Usage:  "Tests the configuration",
			Action: test,
		},
		{
			Name:   "run",
			Usage:  "Gather ADSB data and sends it to plane.watch",
			Action: run,
		},
	}

	if err := app.Run(os.Args); nil != err {
		fmt.Println(err)
	}
}

func getRabbitConnection(timeout int64) (*RabbitMQ, error) {
	if "" == pwUser {
		log.Fatalln("No User Specified For Plane.Watch Rabbit Connection")
	}
	if "" == pwPass {
		log.Fatalln("No Password Specified For Plane.Watch Rabbit Connection")
	}

	var rabbitConfig RabbitMQConfig
	rabbitConfig.Host = pwHost
	rabbitConfig.Port = pwPort
	rabbitConfig.User = pwUser
	rabbitConfig.Password = pwPass
	rabbitConfig.Vhost = pwVhost

	log.Printf("Connecting to plane.watch RabbitMQ @ %s", rabbitConfig)
	rabbit := NewRabbitMQ(rabbitConfig)
	connected := make(chan bool)
	go rabbit.Connect(connected)
	select {
	case <-connected:
		return rabbit, nil
	case <-time.After(time.Duration(timeout) * time.Second):
		return nil, fmt.Errorf("failed to connect to rabbit in a timely manner")
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		//panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// test makes sure that our setup is working
func test(c *cli.Context) {
	log.Printf("Testing connection to dump1090 @ %s:%s", dump1090Host, dump1090Port)
	d1090 := NewDump1090Reader(dump1090Host, dump1090Port)
	var err error
	if err = d1090.Connect(); err != nil {
		log.Fatalf("Unable to connect to Dump 1090 %s", err)
	} else {
		d1090.Stop()
	}

	rabbit, err := getRabbitConnection(10)
	failOnError(err, "Unable to connect to rabbit")
	defer rabbit.Disconnect()

	log.Printf("Success. You are ready to go")
}

// run is our method for running things
func run(c *cli.Context) {
	trk := tracker.NewTracker(tracker.WithVerboseOutput())
	trk.AddProducer(producer.NewAvrFetcher(dump1090Host, dump1090Port))
	trk.Wait()

	// TODO: turn the below into a sink
	// trk.AddSink(sink.Redis(host,port))
	// trk.AddSink(sink.RabbitMq(host, port, queue)

	// TODO: Add an API for programatically handling differing inputs
	// /api/producer[CRUD]
	// /api/sink[CRUD]

	return
	d1090 := NewDump1090Reader(dump1090Host, dump1090Port)
	var err error
	if err = d1090.Connect(); err != nil {
		log.Fatalf("Unable to connect to Dump 1090 %s", err)
		return
	}

	rabbit, err := getRabbitConnection(60)
	failOnError(err, "Unable to connect to rabbit")
	defer rabbit.Disconnect()
	err = rabbit.ExchangeDeclare("planes", "topic")
	failOnError(err, "Failed to declare a topic exchange")

	d1090.SetHandler(func(msg string) {
		var publishError error
		if showDebug {
			log.Println("Decoding: ", msg)
		}
		frame, err := mode_s.DecodeString(msg, time.Now())
		if nil != err {
			log.Println(err)
		}
		plane := tracker.HandleModeSFrame(frame)

		if nil != plane {
			planeJson, _ := json.Marshal(plane)
			msg := amqp.Publishing{
				ContentType: "application/json",
				Body:        planeJson,
			}
			if showDebug {
				log.Println("Sending message to plane.watch for plane:", plane.Icao)
			}
			publishError = rabbit.Publish("planes", plane.Icao, msg)
			if nil != publishError {
				log.Println("Failed to publish message to plane.watch for plane", plane.Icao)
			}

		}

		//tracker.CleanPlanes()
	})

	select {}

	d1090.Stop()
}
