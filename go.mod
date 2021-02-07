module plane.watch

go 1.13

require (
	github.com/kpawlik/geojson v0.0.0-20171201195549-1a4f120c6b41
	github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94
	github.com/urfave/cli v1.20.0
)

replace plane.watch/mode_s => ./mode_s
