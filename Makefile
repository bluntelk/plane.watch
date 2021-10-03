export GOBIN=$(shell pwd)/bin

.PHONY: all

all: test build

build:
	go install ./...

test:
	go test ./...

race:
	go install -race ./...
