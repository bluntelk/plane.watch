export GOBIN=$(shell pwd)/bin

.PHONY: all

all:
	go install ./...

race:
	go install -race ./...
