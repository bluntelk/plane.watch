FROM golang:1.17-alpine AS build

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy the code into the container
COPY . /app

# Build the application
RUN go build plane.watch/cmd/pw_ws_broker

FROM alpine:latest
ENV SOURCE ROUTE_KEY_LOW ROUTE_KEY_HIGH TEST_WEB
WORKDIR /app
COPY --from=build /app/pw_ws_broker /app
COPY docker/pw_ws_broker/entrypoint.sh /app
EXPOSE 80
EXPOSE 9603

# Command to run
CMD ["/app/entrypoint.sh"]
