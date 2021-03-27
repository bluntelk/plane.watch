export GOBIN=$(shell pwd)/bin

.PHONY: all

all:
	go install -race ./...
