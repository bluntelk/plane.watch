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
RUN go build plane.watch/cmd/pw_discord_bot

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/pw_discord_bot /app
COPY docker/pw_discord_bot/entrypoint.sh /app

# Command to run
CMD ["/app/entrypoint.sh"]
