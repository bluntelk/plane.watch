export GOBIN=$(shell pwd)/bin

.PHONY: all

race:
	go install -race ./...

all:
	go install ./...
