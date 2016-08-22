package main

import (
	"github.com/kpawlik/geojson"
	"encoding/json"
	"github.com/urfave/cli"
	"os"
	"bufio"
	"mode_s"
	"time"
	"tracker"
	"io/ioutil"
	"fmt"
	"strings"
	"compress/gzip"
	"compress/bzip2"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Failed to run: %s", r)
		}
	}()

	app := cli.NewApp()

	app.Version = "0.0.1"
	app.Name = "Plane Watch Flight Path Renderer"
	app.Usage = "Reads AVR frames from a file and generates a GeoJSON file"
	app.Authors = []cli.Author{
		{Name:"Jason Playne", Email:"jason@jasonplayne.com"},
	}
	cli.VersionFlag = cli.BoolFlag{Name:"version, V"}

	app.Commands = []cli.Command{
		{
			Name: "avr",
			Usage: "Renders all the plane paths found in an AVR file",
			Action: renderAvr,
		},
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "v",
			Usage: "verbose debugging output",
		},
	}

	app.Run(os.Args)
}

func readFile(inFileName string) (chan string, error) {
	outChan := make(chan string, 50)

	inFile, err := os.Open(inFileName);
	if err != nil {
		return outChan, fmt.Errorf("Failed to open file {%s}: %s", inFileName, err)
	}

	isGzip := strings.ToLower(inFileName[len(inFileName)-2:]) == "gz"
	isBzip2 := strings.ToLower(inFileName[len(inFileName)-3:]) == "bz2"

	go func() {
		defer inFile.Close()

		var scanner *bufio.Scanner
		if isGzip {
			gzipFile, _ := gzip.NewReader(inFile)
			scanner = bufio.NewScanner(gzipFile)
		} else if isBzip2 {
			bzip2File := bzip2.NewReader(inFile)
			scanner = bufio.NewScanner(bzip2File)
		} else {
			scanner = bufio.NewScanner(inFile)
		}
		for scanner.Scan() {
			outChan <- scanner.Text()
		}
		close(outChan)
	}()
	return outChan, nil
}

func renderAvr(c *cli.Context) error {
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
			jobChan <- mode_s.ReceivedFrame{Time:ts, Frame:line}
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
		coords := make(geojson.Coordinates, 0, len(p.LocationHistory));
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

	jsonContent, err := json.MarshalIndent(fc, "", "  ")
	if nil != err {
		panic(err)
	}
	if outFileName != "" {
		f, err := os.Create(outFileName)
		if nil == err {
			f.Write(jsonContent)
			f.Close()
			fmt.Println("Wrote content to file: " + outFileName)
			return nil
		}
	}
	// no writing to a file? output to screen!
	fmt.Println("\n" + string(jsonContent))

	return nil
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
				if resultCounter % 1000 == 0 {
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
