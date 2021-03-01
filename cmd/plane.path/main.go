package main

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/kpawlik/geojson"
	"github.com/urfave/cli"
	"os"
	"plane.watch/lib/tracker"
	"strings"
	"time"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Failed to run: %s", r)
		}
	}()

	app := cli.NewApp()

	app.Version = "v0.0.2"
	app.Name = "Plane Watch Flight Path Renderer"
	app.Usage = "Reads AVR frames or SBS1 data from a file and generates a GeoJSON file"
	app.Authors = []cli.Author{
		{Name: "Jason Playne", Email: "jason@jasonplayne.com"},
	}
	cli.VersionFlag = cli.BoolFlag{Name: "version, V"}

	app.Commands = []cli.Command{
		{
			Name:      "avr",
			Usage:     "Renders all the plane paths found in an AVR file",
			Action:    parseAvr,
			ArgsUsage: "[outfile if not --stdout] input-file.txt [input-file.gz [input-file.bz2]]",
		},
		{
			Name:      "sbs",
			Aliases:   []string{"sbs1"},
			Usage:     "Renders all the plane paths found in an SBS file",
			Action:    parseSbs,
			ArgsUsage: "[outfile if not --stdout] input-file.txt [input-file.gz [input-file.bz2]]",
		},
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "v",
			Usage: "verbose debugging output",
		},
		cli.BoolFlag{
			Name:  "stdout",
			Usage: "Output to stdout instead of to a file (disables any other output)",
		},
	}

	tracker.MaxLocationHistory = -1

	if err := app.Run(os.Args); nil != err {
		fmt.Println(err)
	}
}

func readFiles(inFileNames []string) (chan string, chan error) {
	outChan := make(chan string, 50)
	errChan := make(chan error, 50)

	go func() {
		for _, inFileName := range inFileNames {
			inFile, err := os.Open(inFileName)
			if err != nil {
				errChan <- fmt.Errorf("failed to open file {%s}: %s", inFileName, err)
				continue
			}

			isGzip := strings.ToLower(inFileName[len(inFileName)-2:]) == "gz"
			isBzip2 := strings.ToLower(inFileName[len(inFileName)-3:]) == "bz2"

			var scanner *bufio.Scanner
			if isGzip {
				gzipFile, err := gzip.NewReader(inFile)
				if nil != err {
					errChan <- err
					continue
				}
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
			_ = inFile.Close()
		}
		// and wait for our outChan to be empty
		for len(outChan) > 0 {
			time.Sleep(100 * time.Millisecond)
		}
		close(outChan)
	}()

	return outChan, nil
}

func writeResult(outFileName string) error {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{})
	var coordCounter, planeCounter int
	var trackCounter int

	addFeature := func(coordinates geojson.Coordinates, p *tracker.Plane) {
		trackCounter++
		props := make(map[string]interface{})
		props["icao"] = p.Icao
		props["trackNo"] = trackCounter
		if len(coordinates) > 1 {
			fc.AddFeatures(geojson.NewFeature(geojson.NewLineString(coordinates), props, fmt.Sprintf("%s_%d", p.Icao, trackCounter)))
		}
	}

	tracker.Each(func(p *tracker.Plane) bool {
		var coords geojson.Coordinates
		if 0 == len(p.LocationHistory) {
			return true
		}
		planeCounter++
		numLocations := len(p.LocationHistory)
		coords = make(geojson.Coordinates, 0, numLocations)
		for index, l := range p.LocationHistory {
			if l.Latitude == 0.0 && l.Longitude == 0.0 {
				continue
			}

			coordCounter++
			coords = append(coords, geojson.Coordinate{geojson.CoordType(l.Longitude), geojson.CoordType(l.Latitude)})

			if l.TrackFinished {
				addFeature(coords, p)
				coords = make(geojson.Coordinates, 0, numLocations-index)
			}
		}
		addFeature(coords, p)
		return true
	})
	_, _ = fmt.Fprintf(os.Stderr, "We have %d coords tracked over %d tracks from %d planes\n", coordCounter, trackCounter, planeCounter)

	jsonContent, err := json.Marshal(fc)
	//jsonContent, err := json.MarshalIndent(fc, "", "  ")
	if nil != err {
		panic(err)
	}
	if outFileName != "" {
		f, err := os.Create(outFileName)
		if nil == err {
			_, _ = f.Write(jsonContent)
			_ = f.Close()
			fmt.Println("Wrote content to file: " + outFileName)
			return nil
		}
	}
	// no writing to a file? output to screen!
	fmt.Println("\n" + string(jsonContent))
	return nil
}
