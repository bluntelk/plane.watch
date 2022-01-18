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
RUN go build plane.watch/cmd/pw_ingest

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/pw_ingest /app

# limit ourselves to just location updates
ENV PUBLISH_TYPES=location-updates

HEALTHCHECK --interval=1m --timeout=3s CMD wget http://localhost:9602/health-check

# Command to run
CMD /app/pw_ingest daemon
