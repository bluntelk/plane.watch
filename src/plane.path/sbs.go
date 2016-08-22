package main

import (
	"github.com/urfave/cli"
	"sbs1"
	"fmt"
	"github.com/kpawlik/geojson"
)

func parseSbs(c *cli.Context) error {
	inFileName := c.Args().First()
	outFileName := c.Args().Get(1)
	verbose := c.GlobalBool("v")

	tracking := make(map[string]geojson.Coordinates)

	if "" == inFileName {
		println("Usage: <file with SBS frames to read> <file to export geojson to or stdout if omitted>")
		return nil
	}

	inFile, err := readFile(inFileName)
	if nil != err {
		panic(err)
	}

	var lineCounter uint
	for line := range inFile {
		lineCounter++
		if 0 == lineCounter % 10000 {
			fmt.Printf("Processing line: %d\r", lineCounter)
		}
		frame, err := sbs1.Parse(line)
		if nil != err {
			if verbose {
				fmt.Println("Failed to parse SBS1 Frame.", err)
			}
			continue;
		}

		if frame.HasPosition {
			// add a location for this plane
			if frame.Lat == 0.0 && frame.Lon == 0.0 {
				continue
			}
			tracking[frame.Icao] = append(tracking[frame.Icao], geojson.Coordinate{geojson.CoordType(frame.Lon), geojson.CoordType(frame.Lat)})
		}
	}

	fc := geojson.NewFeatureCollection([]*geojson.Feature{})
	maxIcao := len(tracking)
	var currentIcao int
	for icao, coords := range tracking {
		currentIcao++
		fmt.Printf("Generating GEO Data: %0.1f%%\r", (float64(currentIcao) / float64(maxIcao)) * 100.0)
		props := make(map[string]interface{})
		props["icao"] = icao

		if len(coords) > 1 {
			fc.AddFeatures(geojson.NewFeature(geojson.NewLineString(coords), props, icao))

		}
	}

	fmt.Print("\nWriting result...\r")
	return writeResult(outFileName, fc)
}