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
	"strings"
	"time"
	"tracker"
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
		{Name: "Jason Playne", Email: "jason@jasonplayne.com"},
	}
	cli.VersionFlag = cli.BoolFlag{Name: "version, V"}

	app.Commands = []cli.Command{
		{
			Name:   "avr",
			Usage:  "Renders all the plane paths found in an AVR file",
			Action: parseAvr,
		},
		{
			Name:   "sbs",
			Aliases: []string{"sbs1"},
			Usage:  "Renders all the plane paths found in an SBS file",
			Action: parseSbs,
		},
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "v",
			Usage: "verbose debugging output",
		},
	}

	tracker.MaxLocationHistory = -1

	app.Run(os.Args)
}

func readFile(inFileName string) (chan string, error) {
	outChan := make(chan string, 50)

	inFile, err := os.Open(inFileName)
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


	jsonContent, err := json.Marshal(fc)
	//jsonContent, err := json.MarshalIndent(fc, "", "  ")
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
