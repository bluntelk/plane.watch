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
		close(outChan)
	}()
	return outChan, nil
}

func writeResult(outFileName string, fc *geojson.FeatureCollection) error {
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
