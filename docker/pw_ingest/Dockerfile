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
COPY docker/pw_ingest/entrypoint.sh /app
EXPOSE 9602

# Command to run
CMD ["/app/entrypoint.sh"]
