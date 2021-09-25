docker run -d \
  --hostname=plane.watch \
  --name=encrichment-shaft \
  -p 15672:15672 \
  -p 5672:5672 \
  rabbitmq:3-management