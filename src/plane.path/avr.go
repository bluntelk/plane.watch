package main

import (
	"fmt"
	"github.com/kpawlik/geojson"
	"github.com/urfave/cli"
	"io/ioutil"
	"mode_s"
	"time"
	"tracker"
)

func parseAvr(c *cli.Context) error {
	inFileName := c.Args().First()
	outFileName := c.Args().Get(1)
	verbose := c.GlobalBool("v")

	tracker.MaxLocationHistory = -1

	if "" == inFileName {
		println("Usage: <file with avr frames to read> <file to export geojson to or stdout if omitted>")
		return nil
	}

	jobChan := make(chan mode_s.ReceivedFrame, 1000)
	resultChan := make(chan mode_s.Frame, 1000)
	errorChan := make(chan error, 1000)
	exitChan := make(chan bool)

	go handleReceived(resultChan, verbose)
	go handleErrors(errorChan, verbose)
	if !verbose {
		tracker.SetDebugOutput(ioutil.Discard)
	}
	go mode_s.DecodeStringWorker(jobChan, resultChan, errorChan)
	inFile, err := readFile(inFileName)
	if nil != err {
		panic(err)
	}

	go func() {
		var ts time.Time
		for line := range inFile {
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

	fmt.Printf("We have %d points tracked\n", tracker.PointCounter)
	fc := geojson.NewFeatureCollection([]*geojson.Feature{})
	var coordCounter, planeCounter int

	tracker.Each(func(p tracker.Plane) {
		if 0 == len(p.LocationHistory) {
			return
		}
		planeCounter++
		coords := make(geojson.Coordinates, 0, len(p.LocationHistory))
		for _, l := range p.LocationHistory {
			if l.Latitude == 0.0 && l.Longitude == 0.0 {
				continue
			}
			coordCounter++
			coords = append(coords, geojson.Coordinate{geojson.CoordType(l.Longitude), geojson.CoordType(l.Latitude)})
		}
		props := make(map[string]interface{})
		props["icao"] = p.Icao
		if len(coords) > 1 {
			fc.AddFeatures(geojson.NewFeature(geojson.NewLineString(coords), props, p.IcaoIdentifier))
		}
	})

	fmt.Printf("We have %d coords tracked from %d planes\n", coordCounter, planeCounter)

	return writeResult(outFileName, fc)
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
					fmt.Printf("Results: %d (buf %d)\r", resultCounter, len(results))
				}
			}
		}
	}
}
func handleErrors(errors chan error, verbose bool) {
	for {
		select {
		case err := <-errors:
			if verbose {
				fmt.Println(err)
			}
		}
	}
}
