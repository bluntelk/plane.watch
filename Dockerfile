FROM golang:1.17-alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ENV SOURCE SINK

WORKDIR /app

# Copy the code into the container
COPY . /app

# Build the application
RUN go build /app/cmd/plane.watch.client

# Command to run
CMD ["/app/entrypoint.sh"]
