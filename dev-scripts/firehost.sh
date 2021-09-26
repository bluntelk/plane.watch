bin/plane.watch.client \
  --source=beast://crawled.mapwithlove.com:3004 \
  --sink=amqp://guest:guest@localhost:5672/pw \
  --debug \
  simple