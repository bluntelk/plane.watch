FROM golang:1.17-alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ENV RABBITMQ SOURCE_ROUTE_KEY NUM_WORKERS PROM_METRICS_PORT

WORKDIR /app

# Copy the code into the container
COPY . /app

# Build the application
RUN go build /app/cmd/pwreducer

# Command to run
CMD ["/app/docker/pwreducer/entrypoint.sh"]
