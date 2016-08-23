package main

import (
	"fmt"
	"github.com/urfave/cli"
	"sbs1"
	"time"
	"tracker"
	"io/ioutil"
)

func parseSbs(c *cli.Context) error {
	inFileName := c.Args().First()
	outFileName := c.Args().Get(1)
	verbose := c.GlobalBool("v")

	if "" == inFileName {
		println("Usage: <file with SBS frames to read> <file to export geojson to or stdout if omitted>")
		return nil
	}
	if !verbose {
		tracker.SetDebugOutput(ioutil.Discard)
	}

	inFile, err := readFile(inFileName)
	if nil != err {
		panic(err)
	}

	trackingChan := make(chan sbs1.Frame, 50)
	go handleSbs1Tracking(trackingChan, verbose)

	var lineCounter uint
	for line := range inFile {
		lineCounter++
		if 0 == lineCounter%10000 {
			fmt.Printf("Processing line: %d\r", lineCounter)
		}
		frame, err := sbs1.Parse(line)
		if nil != err {
			if verbose {
				fmt.Println("Failed to parse SBS1 Frame.", err)
			}
			continue
		}
		trackingChan <- frame

	}
	for len(trackingChan) > 0 {
		time.Sleep(100 * time.Millisecond)
	}
	close(trackingChan)

	return writeResult(outFileName)
}

func handleSbs1Tracking(trackingChan chan sbs1.Frame, debug bool) {
	for frame := range trackingChan {
		tracker.HandleSbs1Frame(frame, debug)
	}
}
