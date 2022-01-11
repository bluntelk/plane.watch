FROM golang:1.17-alpine AS build

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy the code into the container
COPY . /app/

# Build the application
RUN go build plane.watch/cmd/pw_router

FROM alpine:latest
ENV RABBITMQ SOURCE_ROUTE_KEY DEST_ROUTE_KEY NUM_WORKERS PROM_METRICS_PORT DEBUG QUIET
WORKDIR /app
COPY --from=build /app/pw_router /app
COPY docker/pw_router/entrypoint.sh /app
EXPOSE 9601

# Command to run
CMD ["/app/entrypoint.sh"]
