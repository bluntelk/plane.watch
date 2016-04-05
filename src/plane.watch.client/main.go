package main

import (
	"github.com/codegangsta/cli"
	"os"
	"log"
	"mode_s"
	"encoding/json"
	"time"
	"tracker"
	"fmt"
	"github.com/streadway/amqp"
)

var (
	pw_host, pw_user, pw_pass, pw_vHost string
	pw_port int
	dump1090_host string
	dump1090_port string
)

func main() {

	app := cli.NewApp()

	app.Version = "1.0.0"
	app.Name = "Plane Watch Client"
	app.Usage = "Reads from dump1090 and sends it to http://plane.watch/"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "pw_host",
			Value: "mq.plane.watch",
			Usage: "How we connect to plane.watch",
			Destination: &pw_host,
			EnvVar: "PW_HOST",
		},
		cli.StringFlag{
			Name: "pw_user",
			Value: "",
			Usage: "user for plane.watch",
			Destination: &pw_user,
			EnvVar: "PW_USER",
		},
		cli.StringFlag{
			Name: "pw_pass",
			Value: "",
			Usage: "plane.watch password",
			Destination: &pw_pass,
			EnvVar: "PW_PASS",
		},
		cli.IntFlag{
			Name: "pw_port",
			Value: 5672,
			Usage: "How we connect to plane.watch",
			Destination: &pw_port,
			EnvVar: "PW_PORT",
		},
		cli.StringFlag{
			Name: "pw_vhost",
			Value: "/pw_feedin",
			Usage: "the virtual host on the plane watch rabbit server",
			Destination: &pw_vHost,
			EnvVar: "PW_VHOST",
		},
		cli.StringFlag{
			Name: "dump1090_host",
			Value: "localhost",
			Usage: "The host to read dump1090 from",
			Destination: &dump1090_host,
			EnvVar: "DUMP1090_HOST",
		},
		cli.StringFlag{
			Name: "dump1090_port",
			Value: "30002",
			Usage: "The port to read dump1090 from",
			Destination: &dump1090_port,
			EnvVar: "DUMP1090_PORT",
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "test",
			Usage: "Tests the configuration",
			Action: test,
		},
		{
			Name: "run",
			Usage: "Gather ADSB data and sends it to plane.watch",
			Action: run,
		},
	}

	app.Run(os.Args)
}

func getRabbitConnection(timeout int64) (*RabbitMQ, error) {
	if "" == pw_user {
		log.Fatalln("No User Specified For Plane.Watch Rabbit Connection")
	}
	if "" == pw_pass {
		log.Fatalln("No Password Specified For Plane.Watch Rabbit Connection")
	}

	var rabbitConfig RabbitMQConfig
	rabbitConfig.Host = pw_host
	rabbitConfig.Port = pw_port
	rabbitConfig.User = pw_user
	rabbitConfig.Password = pw_pass
	rabbitConfig.Vhost = pw_vHost

	log.Printf("Connecting to plane.watch RabbitMQ @ %s", rabbitConfig)
	rabbit := NewRabbitMQ(rabbitConfig)
	connected := make(chan bool)
	go rabbit.Connect(connected)
	select {
	case <-connected:
		return rabbit, nil
	case <-time.After(time.Duration(timeout) * time.Second):
		return nil, fmt.Errorf("Failed to connect to rabbit in a timely manner")
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func test(c *cli.Context) {
	log.Printf("Testing connection to dump1090 @ %s:%s", dump1090_host, dump1090_port)
	d1090 := NewDump1090Reader(dump1090_host, dump1090_port)
	var err error
	if err = d1090.Connect(); err != nil {
		log.Fatalf("Unable to connect to Dump 1090 %s", err)
	} else {
		d1090.Stop()
	}

	rabbit, err := getRabbitConnection(10)
	defer rabbit.Disconnect()
	failOnError(err, "Unable to connect to rabbit")

	log.Printf("Success. You are ready to go");
}

func run(c *cli.Context) {
	d1090 := NewDump1090Reader(dump1090_host, dump1090_port)
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
		log.Println("Decoding: ", msg)
		frame, err := mode_s.DecodeString(msg, time.Now())
		if nil != err {
			log.Println(err)
		}
		plane := tracker.HandleModeSFrame(frame, false)

		if nil != plane {
			planeJson, _ := json.Marshal(plane);

			icao := fmt.Sprintf("%6x", plane.IcaoIdentifier);
			msg := amqp.Publishing{
				ContentType: "application/json",
				Body: planeJson,
			}
			log.Println("Sending message to plane.watch for plane:", icao)
			publishError = rabbit.Publish("planes", icao, msg)
			if nil != publishError {
				log.Println("Failed to publish message to plane.watch for plane", icao)
			}

		}

		tracker.CleanPlanes()
	})

	select {}

	d1090.Stop()
}