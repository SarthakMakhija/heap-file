.DEFAULT_GOAL := build

test:
	go test ./...
.PHONY: test

build: test
	go build ./...
.PHONY: build