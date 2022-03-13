# Plane.Watch

Development is now continued over at https://github.com/plane-watch/pw-pipeline/

This repo has a few of the things I have done around plane.watch.

There are some tools (in cmd/) for dealing with various things

using golang 1.12 you should have no problems compiling this go modules
enabled project.



Some Links for More Information

* http://airmetar.main.jp/radio/ADS-B%20Decoding%20Guide.pdf
* https://mode-s.org/decode/book-the_1090mhz_riddle-junzi_sun.pdf
* https://pypi.org/project/pyModeS/
* https://mode-s.org/decode/content/mode-s/6-els.html
* https://www.eurocontrol.int/sites/default/files/content/documents/nm/asterix/archives/asterix-cat062-system-track-data-part9-v1.10-122009.pdf

## Building

### Development

    make

That's it. It runs the tests and builds the binaries and puts them into `bin/`

If you want to build a specific binary

    go build plane.watch/cmd/pw_ingest

or you can run it with

    go run plane.watch/cmd/pw_ingest

### Building Docker Containers

    docker build -t plane.watch/pw_ws_broker:latest -f docker/pw_ws_broker/Dockerfile .
    docker build -t plane.watch/pw_router:latest -f docker/pw_router/Dockerfile .
    docker build -t plane.watch/pw_ingest:latest -f docker/pw_ingest/Dockerfile .

## Commands
This repo contains many commands that can be run

### pw_ingest

This binary is used to accepts many different feeds and put them onto a message bus.

It accepts its inputs in URL form

valid schemes are
* av1
* beast
* sbs1

Fetchers and Listeners do not accept usernames/passwords.

Sinks do accept usernames and passwords, which is needed for rabbitmq

You can specify a `tag`, `refLat` and `refLon` in the fetchers and listeners
* `tag` flows through to the output
* `refLat` and `refLon` are used to calculate surface position

Examples:
* --fetch=beast://crawled.mapwithlove.com:3004?tag=firehose
* --fetch=avr://localhost:30002?tag=local-receiver
* --listen=beast://0.0.0.0:3005?tag=rando
* --sink=amqp://guest:guest@localhost:5672/pw

### pw_router

This binary is used to reduce the incoming feed of location updates down to only updates that indicate a "significant" change. 
A significant change is where:
* Aircraft heading changes by at least 1 degree
* The vertical and horizontal velocity or altitude changes
* The flight metadata (Flight number, status, on ground status, special or squawk codes) changes

Location updates can be read from an AMQP connection and fed back onto the same but to the `reducer-out` topic.

```
NAME:
   Plane Watch Router (pw_router) - Reads location updates from AMQP and publishes only significant updates.

USAGE:
   pw_router [global options] command [command options] [arguments...]

VERSION:
   1.0.0

DESCRIPTION:
   This program takes a stream of plane tracking data (location updates) from an AMQP message bus  and filters messages and only returns significant changes for each aircraft.

   example: ./pw_router --rabbitmq="amqp://guest:guest@localhost:5672" --source-route-key=location-updates --num-workers=8 --prom-metrics-port=9601

COMMANDS:
   daemon
   cli
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --rabbitmq value                   Rabbitmq URL for reaching and publishing updates. (default: "amqp://guest:guest@rabbitmq:5672/pw") [$RABBITMQ]
   --source-route-key value           Name of the routing key to read location updates from. (default: "location-updates-enriched") [$SOURCE_ROUTE_KEY]
   --destination-route-key value      Name of the routing key to publish significant updates to. (default: "location-updates-enriched-reduced") [$DEST_ROUTE_KEY]
   --num-workers value                Number of workers to process updates. (default: 4) [$NUM_WORKERS]
   --spread-updates                   publish location updates to their respective tileXX_high and tileXX_low routing keys as well (default: false) [$DEBUG]
   --update-age value                 seconds to keep an update before aging it out of the cache. (default: 30) [$UPDATE_AGE]
   --update-age-sweep-interval value  Seconds between cache age sweeps.. (default: 5) [$UPDATE_SWEEP]
   --register-test-queues             Subscribes a bunch of queues to our routing keys (default: false)
   --very-verbose                     Enable trace level debugging (default: false)
   --debug                            Show Extra Debug Information (default: false) [$DEBUG]
   --quiet                            Only show important messages (default: false) [$QUIET]
   --monitoring-port value            Port to listen on for prometheus app metrics. (default: 9601) [$MONITORING_PORT]
   --help, -h                         show help (default: false)
   --version, -v                      print the version (default: false)
```

### pw_ws_broker
This is our plane.watch websocket broker. Connect to it and speak its language to get your location information
