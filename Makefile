.DEFAULT_GOAL := build

test:
	go test -count=1 ./...
.PHONY: test

build: test
	go build ./...
.PHONY: build