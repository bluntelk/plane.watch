FROM golang:1.17-alpine AS build

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy the code into the container
COPY go.mod /app
COPY go.sum /app
RUN go mod verify && go mod download

COPY . /app

# Build the application
RUN go build plane.watch/cmd/pw_ws_broker

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/pw_ws_broker /app
EXPOSE 80

HEALTHCHECK --interval=15s --timeout=3s CMD wget --spider http://localhost:9603/status

# Command to run
CMD /app/pw_ws_broker daemon
