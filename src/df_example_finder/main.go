package main

import (
	"os"
	"bufio"
	"mode_s"
	"time"
	"strings"
	"fmt"
)

func gatherSamples(filePath string) {

	f, err := os.Open(filePath)
	if nil != err {
		println(err)
		return
	}

	countMap := make(map[byte]uint32)
	df17Map := make(map[byte]uint32)
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

		if 17 == frame.DownLinkType() {
			df17Map[frame.MessageType()]++
			key := fmt.Sprintf("DF17/%d", frame.MessageType())
			if _, ok := existingSamples[key]; ok {
				continue
			}
			existingSamples[key] = true
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

	println("Sample Frames")
	for k, s := range samples {
		println(k, ":", "['" + strings.Join(s, "', '") + "'],")
	}
}

func showTypes(filePath string) {
	f, err := os.Open(filePath)
	if nil != err {
		println(err)
		return
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		frame, err := mode_s.DecodeString(line, time.Now())
		if nil != err {
			//println("Error! ", line, err.Error())
			continue
		}

		fmt.Printf("DF%02d\t%s\n",frame.DownLinkType(), line)
	}

}

func main() {
	if len(os.Args) < 2 {
		println("first arg must be file of stored AVS packets")
		return
	}
	var cmd string
	if 3 <= len(os.Args) {
		cmd = os.Args[2]
	}

	switch cmd {
	case "type":
		showTypes(os.Args[1])
	default:
		gatherSamples(os.Args[1])
	}

}
