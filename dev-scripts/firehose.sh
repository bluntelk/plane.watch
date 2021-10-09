bin/plane.watch.client \
  --fetch=beast://crawled.mapwithlove.com:3004 \
  --sink=amqp://guest:guest@localhost:5672/pw \
  --rabbit-queue=location-updates \
  --debug \
  daemon