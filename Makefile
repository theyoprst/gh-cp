.PHONY: build clean test install lint lint-fix

build:
	go build -o gh-cp ./cmd/gh-cp

clean:
	rm -f gh-cp

test:
	go test ./...

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

install: build
	cp gh-cp $(shell go env GOPATH)/bin/

.DEFAULT_GOAL := build
