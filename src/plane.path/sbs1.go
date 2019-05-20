package main

import (
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"sbs1"
	"time"
	"tracker"
)

func parseSbs(c *cli.Context) error {
	stdOut := c.GlobalBool("stdout")
	verbose := c.GlobalBool("v")
	if c.NArg() == 0 {
		return fmt.Errorf("you need to specify some files")
	}

	var outFileName string
	var dataFiles []string
	if stdOut {
		dataFiles = c.Args()
	} else {
		outFileName = c.Args().First()
		dataFiles = c.Args()[1:]
	}

	if 0 == len(dataFiles) {
		return fmt.Errorf("you need to specify some files")
	}

	if !verbose {
		tracker.SetDebugOutput(ioutil.Discard)
	}

	inputLines, errChan := readFiles(dataFiles)

	trackingChan := make(chan sbs1.Frame, 50)
	go handleSbs1Tracking(trackingChan, verbose)
	go func() {
		for err := range errChan {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	var lineCounter uint
	for line := range inputLines {
		lineCounter++
		if 0 == lineCounter%10000 {
			fmt.Fprintf(os.Stderr, "Processing line: %d\r", lineCounter)
		}
		frame, err := sbs1.Parse(line)
		if nil != err {
			if verbose {
				fmt.Fprintln(os.Stderr, "Failed to parse SBS1 Frame.", err)
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
