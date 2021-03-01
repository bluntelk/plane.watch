package main

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"plane.watch/lib/tracker/mode_s"
	"strings"
	"time"
)

func getFileReader(filePath string) (io.Reader, error) {
	f, err := os.Open(filePath)
	if nil != err {
		return nil, err
	}
	if strings.HasSuffix(filePath, ".gz") {
		println("Reading Gzip file")
		return gzip.NewReader(f)
	}
	if strings.HasSuffix(filePath, ".bz2") {
		println("Reading Bzip2 file")
		return bzip2.NewReader(f), nil
	}
	println("Reading plain text file")
	return f, nil
}

func gatherSamples(filePath string) {

	f, err := getFileReader(filePath)

	if nil != err {
		println(err)
		return
	}
	println("Processing file...")

	countMap := make(map[byte]uint32)
	df17Map := make(map[byte]uint32)
	bdsMap := make(map[string]uint32)
	samples := make(map[byte][]string)
	existingSamples := make(map[string]bool)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		frame, err := mode_s.DecodeString(line, time.Now())
		if nil != err {
			println("Error! ", line, err.Error())
			continue
		}

		countMap[frame.DownLinkType()]++

		switch frame.DownLinkType() {
		case 17:
			df17Map[frame.MessageType()]++
			key := fmt.Sprintf("DF17/%d", frame.MessageType())
			if _, ok := existingSamples[key]; ok {
				continue
			}
			existingSamples[key] = true
		case 20,21:
			bdsMap[frame.BdsMessageType()]++
			if "0.0" == frame.BdsMessageType() {
				continue
			}
		}

		if len(samples[frame.DownLinkType()]) < 100 {
			if _, exist := existingSamples[line]; !exist {
				samples[frame.DownLinkType()] = append(samples[frame.DownLinkType()], line)
				existingSamples[line] = true
			}
		}
	}

	println("Frame Type Counts")
	for k, c := range countMap {
		println("DF", k, "=\t", c)
	}
	println("DF17 Frame Breakdown")
	for k, c := range df17Map {
		println("DF17 Type", k, "=\t", c)
	}
	println("DF 20/21 BDS Frame Breakdown")
	for k, c := range bdsMap {
		println("BDS Type", k, "=\t", c)
	}

	println("Sample Frames")
	for k, s := range samples {
		println(k, ":", "['" + strings.Join(s, "', '") + "'],")
	}
}

func showTypes(filePath string) {
	f, err := getFileReader(filePath)
	if nil != err {
		println(err)
		return
	}
	println("Processing file...")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		frame, err := mode_s.DecodeString(line, time.Now())
		if nil != err {
			//println("Error! ", line, err.Error())
			continue
		}
		switch frame.DownLinkType() {
		case 17:
			fmt.Printf("DF%02d\tMT%02d\tST%02d\t%s\t%s\n",frame.DownLinkType(), frame.MessageType(), frame.MessageSubType(), frame.ICAOString(), line)
		case 20, 21:
			fmt.Printf("DF%02d\tBDS%s\tST%02d\t%s\t%s\n",frame.DownLinkType(), frame.BdsMessageType(), frame.MessageSubType(), frame.ICAOString(), line)
		default:
			fmt.Printf("DF%02d\tMT%02d\tST%02d\t%s\t%s\n",frame.DownLinkType(), frame.MessageType(), frame.MessageSubType(), frame.ICAOString(), line)

		}

	}

}

func main() {
	if len(os.Args) < 2 {
		println("first arg must be file of stored AVR packets")
		return
	}
	var cmd string
	if 3 <= len(os.Args) {
		cmd = os.Args[2]
	}

	switch cmd {
	case "type":
		showTypes(os.Args[1])
	case "gather":
		gatherSamples(os.Args[1])
	default:
		println("3rd argument must be either type or gather")
	}

}
