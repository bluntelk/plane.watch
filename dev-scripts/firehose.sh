bin/pw_ingest \
  --fetch=beast://crawled.mapwithlove.com:3004 \
  --sink=amqp://guest:guest@localhost:5672/pw \
  --debug \
  --rabbitmq-test-queues \
  simple