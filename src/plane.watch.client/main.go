package main

import (
	"github.com/codegangsta/cli"
	"os"
	"log"
	"mode_s"
	"encoding/json"
	"time"
	"tracker"
)

var (
	pw_host, pw_user, pw_pass string
	pw_port string
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
			Value: "plane.watch",
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
		cli.StringFlag{
			Name: "pw_port",
			Value: "1001",
			Usage: "How we connect to plane.watch",
			Destination: &pw_port,
			EnvVar: "PW_PORT",
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

func test(c *cli.Context) {
	log.Printf("Testing connection to dump1090 @ %s:%s", dump1090_host, dump1090_port)
	d1090 := NewDump1090Reader(dump1090_host, dump1090_port)
	var err error
	if err = d1090.Connect(); err != nil {
		log.Fatalf("Unable to connect to Dump 1090 %s", err)
	} else {
		d1090.Stop()
	}

	log.Printf("Connecting to plane.watch @ %s: %s", pw_host, pw_port)

	log.Printf("Success. You are ready to go");
}

func run(c *cli.Context) {
	d1090 := NewDump1090Reader(dump1090_host, dump1090_port)
	var err error
	if err = d1090.Connect(); err != nil {
		log.Fatalf("Unable to connect to Dump 1090 %s", err)
		return
	}

	d1090.SetHandler(func(msg string) {
		log.Println("Decoding: ", msg)
		frame, err := mode_s.DecodeString(msg, time.Now())
		if nil != err {
			log.Println(err)
		}
		plane := tracker.HandleModeSFrame(frame)

		if nil != plane {
			b, _ := json.Marshal(plane);
			println(string(b))
		}

		tracker.CleanPlanes()
	})

	select{}

	d1090.Stop()
}