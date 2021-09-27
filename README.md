# Plane.Watch

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


## plane.watch.client

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