package main

import (
	"github.com/codegangsta/cli"
	"os"
	"net/http"
	"fmt"
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
	if "" == c.Args()[0] {
		fmt.Println("First argument needs to be the htdocs folder")
	}
	http.Handle("/", http.FileServer(http.Dir(c.Args()[0])))

	http.HandleFunc("/decode", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Woot");
	})

	http.ListenAndServe(":8080", nil)

}