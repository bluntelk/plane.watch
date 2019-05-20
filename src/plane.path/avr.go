package main

import (
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"mode_s"
	"os"
	"time"
	"tracker"
)

func parseAvr(c *cli.Context) error {
	stdOut := c.GlobalBool("stdout")
	verbose := c.GlobalBool("v")

	var outFileName string
	var dataFiles []string
	if stdOut {
		dataFiles = c.Args()
	} else {
		outFileName = c.Args().First()
		dataFiles = c.Args()[1:]
	}

	if 0 == len(dataFiles) {
		fmt.Fprintf(os.Stderr, "Usage: %s %s [output-file.json] input [input...]\n", os.Args[0], c.Command.Name)
	}

	if !verbose {
		tracker.SetDebugOutput(ioutil.Discard)
	}

	inputLines, errChan := readFiles(dataFiles)

	jobChan := make(chan mode_s.ReceivedFrame, 1000)
	resultChan := make(chan mode_s.Frame, 1000)
	errorChan := make(chan error, 1000)
	exitChan := make(chan bool)

	go handleReceived(resultChan, verbose)
	go handleErrors(errChan, verbose)
	go handleErrors(errorChan, verbose)
	if !verbose {
		tracker.SetDebugOutput(ioutil.Discard)
	}
	go mode_s.DecodeStringWorker(jobChan, resultChan, errorChan)

	go func() {
		var ts time.Time
		for line := range inputLines {
			ts = ts.Add(500 * time.Millisecond)
			jobChan <- mode_s.ReceivedFrame{Time: ts, Frame: line}
		}

		for len(jobChan) > 0 {
			time.Sleep(time.Second)
		}

		exitChan <- true
	}()

	select {
	case <-exitChan:
		close(jobChan)
		close(resultChan)
		close(errorChan)
	}

	fmt.Fprintf(os.Stderr,"We have %d points tracked\n", tracker.PointCounter)

	return writeResult(outFileName)
}

func handleReceived(results chan mode_s.Frame, verbose bool) {
	var resultCounter int
	for {
		select {
		case frame := <-results:
			resultCounter++
			plane := tracker.HandleModeSFrame(frame, verbose)
			if nil != plane {
				// whee plane changed - now has it moved from its last position?
				if resultCounter%1000 == 0 {
					fmt.Fprintf(os.Stderr,"Results: %d (buf %d)\r", resultCounter, len(results))
				}
			}
		}
	}
}
func handleErrors(errors chan error, verbose bool) {
	for {
		select {
		case err := <-errors:
			if verbose && nil != err{
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}
