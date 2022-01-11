FROM golang:1.17-alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ENV SOURCE SINK LISTEN TAG

WORKDIR /app

# Copy the code into the container
COPY . /app

# Build the application
RUN go build /app/cmd/pw_ingest

# Command to run
CMD ["/app/docker/pw_ingest/entrypoint.sh"]
