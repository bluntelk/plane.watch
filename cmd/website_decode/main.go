package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"path"
	"plane.watch/pkg/mode_s"
	"plane.watch/pkg/tracker"
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
			Name: "port",
			Value: "8080",
			Usage: "Port to run the website on",
		},
	}

	app.Action = runHttpServer

	app.Run(os.Args)
}

func runHttpServer(c *cli.Context) {
	if len(c.Args()) == 0 || "" == c.Args()[0] {
		fmt.Println("First argument needs to be the htdocs folder")
		return;
	}
	var htdocsPath string
	htdocsPath = path.Clean(c.Args()[0])

	http.Handle("/css/", http.FileServer(http.Dir(htdocsPath)))
	http.Handle("/js/", http.FileServer(http.Dir(htdocsPath)))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintln(w, "Failed big time...");
				fmt.Fprintln(w, r)
				w.Write(debug.Stack())
			}
		}()
		switch r.URL.Path {
		case "/decode":
			tracker.NukePlanes();
			var submittedPackets string
			r.ParseForm()
			submittedPackets = r.FormValue("packet")
			if "" == submittedPackets {
				fmt.Fprintln(w, "No Packet Provided")
				return
			}
			packets := strings.Split(submittedPackets, ";")
			icaoList := make(map[uint32]uint32)
			for _, packet := range packets {
				packet = strings.TrimSpace(packet)
				if "" == packet {
					continue;
				}
				println("Decoding Frame:", packet)
				frame, err := mode_s.DecodeString(packet, time.Now())
				if err != nil {
					fmt.Fprintln(w, "Failed to decode.", err)
					return
				}
				tracker.HandleModeSFrame(frame, false);
				icaoList[frame.ICAOAddr()] = frame.ICAOAddr()
				frame.Describe(w)
			}

			for _, icao := range icaoList {
				fmt.Fprintln(w, "")
				plane := tracker.GetPlane(icao)
				encoded, _ := json.MarshalIndent(plane, "", "  ")
				fmt.Fprintf(w, "%s", string(encoded))
			}

		case "/":
			http.ServeFile(w, r, path.Join(htdocsPath, "/index.html"))
		default:
			http.NotFound(w, r)
			fmt.Fprintln(w, "<br/>\n" + r.RequestURI)
		}

	})

	port := ":" + c.String("port")
	log.Printf("Listening on %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))

}