package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"plane.watch/lib/export"
	"plane.watch/lib/logging"
	"plane.watch/lib/rabbitmq"
	"time"
)

type (
	rabbit struct {
		rmq  *rabbitmq.RabbitMQ
		conf *rabbitmq.Config

		queues map[string]*amqp.Queue

		samples map[string]planeLocations
	}
	planeLocations []planeLocation

	planeLocation struct {
		export.PlaneLocation
		original []byte
	}

)

func main() {
	app := cli.NewApp()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)


	app.Flags = []cli.Flag{
		&cli.BoolFlag{Name: "debug"},
		&cli.BoolFlag{Name: "quiet"},
	}
	app.Action = run

	app.Before = func(c *cli.Context) error {
		logging.ConfigureForCli()
		logging.SetVerboseOrQuiet(c.Bool("debug"), c.Bool("quiet"))
		return nil
	}

	if err := app.Run(os.Args); nil != err {
		log.Error().Err(err).Send()
	}
}

func (r *rabbit) connect(timeout time.Duration) error {
	rabbitConfig := rabbitmq.Config{
		Host:     "localhost",
		Port:     "5672",
		Vhost:    "pw",
		User:     "guest",
		Password: "guest",
		Ssl:      rabbitmq.ConfigSSL{},
	}

	log.Info().Str("host", rabbitConfig.String()).Msg("Connecting to RabbitMQ")
	r.rmq = rabbitmq.New(rabbitConfig)
	connected := make(chan bool)
	go r.rmq.Connect(connected)
	select {
	case <-connected:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("failed to connect to rabbit in a timely manner")
	}
}

func (r *rabbit) makeQueue(name string) error {
	q, err := r.rmq.QueueDeclare(name, 60000)
	if nil != err {
		return err
	}
	r.queues[name] = &q
	return nil
}

func (r *rabbit) subscribe(exchange, queue string) error {
	return r.rmq.QueueBind(queue, "location-updates", exchange)
}

func (r *rabbit) handleMsg(msg []byte) error {
	var err error
	log.Debug().Str("payload", string(msg)).Send()

	update := planeLocation{}
	if err = json.Unmarshal(msg, &update); nil != err {
		return err
	}

	if "" == update.Icao {
		return nil
	}
	update.original = msg

	if _, ok := r.samples[update.Icao]; !ok {
		r.samples[update.Icao] = planeLocations{}
	}

	r.samples[update.Icao] = append(r.samples[update.Icao], update)

	if len(r.samples[update.Icao]) > 100 {
		r.analyze(r.samples[update.Icao])
		r.samples[update.Icao] = planeLocations{}
	}

	return nil
}

func (r *rabbit) showCollectionStatus() {
	maxIcao := ""
	maxUpdates := 0
	sumUpdates := 0
	for icao, updates := range r.samples {
		if len(updates) > maxUpdates {
			maxUpdates = len(updates)
			maxIcao = icao
		}
		sumUpdates = sumUpdates + len(updates)
	}
	log.Info().
		Str("ICAO", maxIcao).
		Int("Updates", maxUpdates).
		Int("Total Msgs", sumUpdates).
		Int("Total Planes", len(r.samples)).
		Msg("stats")
}

func (r *rabbit) analyze(list planeLocations) {
	dedupe := make(map[string]planeLocation)
	noDupeNumbers := make([]int, 0)
	dupCount := 0
	for i, update := range list {
		if _, ok := dedupe[string(update.original)]; ok {
			// duplicate message
			dupCount++
			log.Debug().Int("index", i).Msg("Duplicate update")
		} else {
			noDupeNumbers = append(noDupeNumbers, i)
			dedupe[string(update.original)] = update
		}
	}
	log.Info().
		Str("ICAO", list[0].Icao).
		Int("Duplicate Count", dupCount).
		Msgf("Found %d duplicates", dupCount)


	dedupeList := make([]planeLocation, len(noDupeNumbers))
	for i, idx := range noDupeNumbers {
		dedupeList[i] = list[idx]
	}

	b, _ := json.MarshalIndent(dedupeList, "", "  ")
	_ = ioutil.WriteFile("plane.filter.json", b, 0644)
}

func run(c *cli.Context) error {
	var err error
	// connect to rabbitmq, create ourselves 2 queues
	r := rabbit{
		queues:  map[string]*amqp.Queue{},
		samples: map[string]planeLocations{},
	}
	if err = r.connect(time.Second); nil != err {
		return err
	}

	if err = r.makeQueue("filter-in"); nil != err {
		return err
	}
	//if err = r.makeQueue("filter-out"); nil != err {
	//	return err
	//}
	if err = r.subscribe("plane.watch.data", "filter-in"); nil != err {
		return err
	}
	ch, err := r.rmq.Consume("filter-in", "plane.filter")
	if nil != err {
		return err
	}
	last := time.Now()
	for msg := range ch {
		var gErr error
		if gErr = r.handleMsg(msg.Body); nil != gErr {
			log.Error().Err(gErr).Send()

			if gErr = msg.Nack(false, false); nil != gErr {
				log.Error().Err(gErr).Send()
			}
		} else {
			if gErr = msg.Ack(false); nil != gErr {
				log.Error().Err(gErr).Send()
			}
		}
		if time.Now().After(last.Add(time.Second)) {
			r.showCollectionStatus()
			last = time.Now()
		}
	}

	return nil
}
