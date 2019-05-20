module plane.watch

go 1.12

replace mode_s => ../mode_s

replace tracker => ../tracker

replace sbs1 => ../sbs1

require (
	github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94
	github.com/urfave/cli v1.20.0
	mode_s v0.0.0-00010101000000-000000000000
	sbs1 v0.0.0-00010101000000-000000000000 // indirect
	tracker v0.0.0-00010101000000-000000000000
)