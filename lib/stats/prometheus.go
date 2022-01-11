package stats

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"net/http"
)

func IncludePrometheusFlags(app *cli.App, defaultPort int) {
	app.Flags = append(app.Flags,
		&cli.IntFlag{
			Name:    "prom-metrics-port",
			Usage:   "Port to listen on for prometheus app metrics.",
			Value:   defaultPort,
			EnvVars: []string{"PROM_METRICS_PORT"},
		},
	)
}

func RunPrometheusWebServer(c *cli.Context) {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(fmt.Sprintf(":%d", c.Int("prom-metrics-port")), nil)
	}()
}
