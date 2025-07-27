.PHONY: build clean test install

build:
	go build -o gh-cp ./cmd/gh-cp

clean:
	rm -f gh-cp

test:
	go test ./...

install: build
	cp gh-cp $(GOPATH)/bin/

.DEFAULT_GOAL := build