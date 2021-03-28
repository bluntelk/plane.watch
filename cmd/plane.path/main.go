package main

import (
	"encoding/json"
	"fmt"
	"github.com/kpawlik/geojson"
	"github.com/urfave/cli"
	"io"
	"os"
	"plane.watch/lib/tracker"
	"runtime/pprof"
)


func main() {
	app := cli.NewApp()

	app.Version = "v0.1.0"
	app.Name = "Plane Watch flight Path Renderer"
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
			Before:    validateParams,
		},
		{
			Name:      "sbs",
			Aliases:   []string{"sbs1"},
			Usage:     "Renders all the plane paths found in an SBS file",
			Action:    parseSbs1,
			ArgsUsage: "[outfile if not --stdout] input-file.txt [input-file.gz [input-file.bz2]]",
			Before:    validateParams,
		},
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "verbose debugging output",
		},
		cli.BoolFlag{
			Name:  "stdout",
			Usage: "Output to stdout instead of to a file (disables any other output)",
		},
		cli.BoolFlag{
			Name:  "profile",
			Usage: "creates a CPU Profile of the code",
		},
	}

	tracker.MaxLocationHistory = -1
	app.Before = func(context *cli.Context) error {
		if context.Bool("profile") {

			f, err := os.Create("cpuprofile.pprof")
			if err != nil {
				return err
			}
			if err = pprof.StartCPUProfile(f); err != nil {
				return err
			}
		}
		return nil
	}

	app.After = func(context *cli.Context) error {
		if context.Bool("profile") {
			pprof.StopCPUProfile()
			println("To analyze the profile, use this cmd")
			println("go tool pprof -http=:7777 cpuprofile.pprof")
		}
		return nil
	}

	if err := app.Run(os.Args); nil != err {
		fmt.Println(err)
	}
}

func validateParams(c *cli.Context) error {
	stdOut := c.GlobalBool("stdout")
	if c.NArg() == 0 {
		return fmt.Errorf("you need to specify some files")
	}
	if !stdOut {
		if c.NArg() <= 1 {
			return fmt.Errorf("the first file is the output file, the other files are the input files")
		}
	}
	return nil
}

func getFilePaths(c *cli.Context) []string {
	stdOut := c.GlobalBool("stdout")
	if c.NArg() == 0 {
		return []string{}
	}

	if stdOut {
		return c.Args()
	} else {
		return c.Args()[1:]
	}
}

func getOutput(c *cli.Context) (io.WriteCloser, error) {
	stdOut := c.GlobalBool("stdout")
	if c.NArg() == 0 || stdOut {
		return os.Stdout, nil
	} else {
		f, err := os.OpenFile(c.Args().First(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if nil != err {
			return os.Stderr, err
		}
		return f, nil
	}

}

func writeResult(trk *tracker.Tracker, out io.WriteCloser) error {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{})
	var coordCounter, planeCounter int
	var trackCounter int

	addFeature := func(coordinates geojson.Coordinates, p *tracker.Plane) {
		trackCounter++
		props := make(map[string]interface{})
		props["icao"] = p.IcaoIdentifierStr()
		props["trackNo"] = trackCounter
		if len(coordinates) > 1 {
			fc.AddFeatures(geojson.NewFeature(geojson.NewLineString(coordinates), props, fmt.Sprintf("%s_%d", p.IcaoIdentifierStr(), trackCounter)))
		}
	}

	trk.EachPlane(func(p *tracker.Plane) bool {
		var coords geojson.Coordinates
		if 0 == len(p.LocationHistory()) {
			return true
		}
		planeCounter++
		numLocations := len(p.LocationHistory())
		coords = make(geojson.Coordinates, 0, numLocations)
		for index, l := range p.LocationHistory() {
			if l.Lat() == 0.0 && l.Lon() == 0.0 {
				continue
			}

			coordCounter++
			coords = append(coords, geojson.Coordinate{geojson.CoordType(l.Lon()), geojson.CoordType(l.Lat())})

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
	_, _ = out.Write(jsonContent)
	return out.Close()
}
