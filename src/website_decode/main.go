package main

import (
	"github.com/codegangsta/cli"
	"os"
	"net/http"
	"fmt"
	"mode_s"
	"time"
	"path"
	"runtime/debug"
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
			var packet string
			r.ParseForm()
			packet = r.FormValue("packet")
			if "" == packet {
				fmt.Fprintln(w, "No Packet Provided")
				return
			}
			println("Decoding Frame:", packet)
			frame, err := mode_s.DecodeString(packet, time.Now())
			if err != nil {
				fmt.Fprintln(w, "Failed to decode.", err)
				return
			}
			frame.Describe(w)
		case "/":
			http.ServeFile(w, r, path.Join(htdocsPath, "/index.html"))
		default:
			http.NotFound(w, r)
			fmt.Fprintln(w, "<br/>\n" + r.RequestURI)
		}

	})

	http.ListenAndServe(":8080", nil)

}