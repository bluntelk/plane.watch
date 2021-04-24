package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"plane.watch/lib/sink"
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
		&cli.StringFlag{
			Name:  "port",
			Value: "8080",
			Usage: "Port to run the website on",
		},
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "output debug info on the CLI",
		},
	}

	app.Action = runHttpServer

	if err := app.Run(os.Args); nil != err {
		fmt.Println(err)
	}
}

func runHttpServer(c *cli.Context) error {
	var htdocsPath string
	var err error
	var files fs.FS
	if c.NArg() == 0 {
		println("Using our embedded filesystem")
		files, err = fs.Sub(embeddedHtdocs, "htdocs")
		if nil != err {
			panic(err)
		}
	} else {
		println("Using the files in dir:", htdocsPath)
		htdocsPath = path.Clean(c.Args().First())
		files = os.DirFS(htdocsPath)
	}

	http.Handle("/", http.FileServer(http.FS(files)))

	http.HandleFunc("/decode", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				_, _ = fmt.Fprintln(w, "<pre>Failed big time...")
				_, _ = fmt.Fprintln(w, r)
				_, _ = w.Write(debug.Stack())
			}
		}()
		switch r.URL.Path {
		case "/decode":
			pt := tracker.NewTracker()
			if c.Bool("verbose") {
				pt.AddSink(sink.NewLoggerSink(sink.WithLogOutput(os.Stdout)))
			}
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
				pt.GetPlane(frame.Icao()).HandleModeSFrame(frame, nil, nil)
				icaoList[frame.Icao()] = frame.Icao()
				frame.Describe(w)
			}

			for _, icao := range icaoList {
				_, _ = fmt.Fprintln(w, "")
				plane := pt.GetPlane(icao)
				encoded, _ := json.MarshalIndent(plane, "", "  ")
				_, _ = fmt.Fprintf(w, "%s", string(encoded))
			}

			pt = nil
		default:
			http.NotFound(w, r)
			_, _ = fmt.Fprintln(w, "<br/>\n"+r.RequestURI)
		}

	})

	port := ":" + c.String("port")
	log.Printf("Listening on %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
	return nil
}
