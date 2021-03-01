package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"path"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/mode_s"
	"runtime/debug"
	"strings"
	"time"
)

// this is a website where you put in one or more Mode S frames and they are decoded
// in a way that is informational

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port",
			Value: "8080",
			Usage: "Port to run the website on",
		},
	}

	app.Action = runHttpServer

	if err := app.Run(os.Args); nil != err {
		fmt.Println(err)
	}
}

func runHttpServer(c *cli.Context) {
	if len(c.Args()) == 0 || "" == c.Args()[0] {
		fmt.Println("First argument needs to be the htdocs folder")
		return
	}
	var htdocsPath string
	htdocsPath = path.Clean(c.Args()[0])


	http.Handle("/css/", http.FileServer(http.Dir(htdocsPath)))
	http.Handle("/js/", http.FileServer(http.Dir(htdocsPath)))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				_, _ = fmt.Fprintln(w, "Failed big time...")
				_, _ = fmt.Fprintln(w, r)
				_, _ = w.Write(debug.Stack())
			}
		}()
		switch r.URL.Path {
		case "/decode":
			pt := tracker.NewTracker()
			var submittedPackets string
			_ = r.ParseForm()
			submittedPackets = r.FormValue("packet")
			if "" == submittedPackets {
				_, _ = fmt.Fprintln(w, "No Packet Provided")
				return
			}
			packets := strings.Split(submittedPackets, ";")
			icaoList := make(map[uint32]uint32)
			for _, packet := range packets {
				packet = strings.TrimSpace(packet)
				if "" == packet {
					continue
				}
				log.Println("Decoding Frame:", packet)
				frame, err := mode_s.DecodeString(packet, time.Now())
				if err != nil {
					_, _ = fmt.Fprintln(w, "Failed to decode.", err)
					return
				}
				if nil == frame {
					_, _ = fmt.Fprintln(w, "Not an AVR Frame", err)
					return
				}
				pt.HandleModeSFrame(frame)
				icaoList[frame.ICAOAddr()] = frame.ICAOAddr()
				frame.Describe(w)
			}

			for _, icao := range icaoList {
				_, _ = fmt.Fprintln(w, "")
				plane := pt.GetPlane(icao)
				encoded, _ := json.MarshalIndent(plane, "", "  ")
				_, _ = fmt.Fprintf(w, "%s", string(encoded))
			}

		case "/":
			http.ServeFile(w, r, path.Join(htdocsPath, "/index.html"))
		default:
			http.NotFound(w, r)
			_, _ = fmt.Fprintln(w, "<br/>\n"+r.RequestURI)
		}

	})

	port := ":" + c.String("port")
	log.Printf("Listening on %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
