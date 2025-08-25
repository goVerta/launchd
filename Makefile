SHELL := /bin/bash

APP := launchd
PKG := ./...

.PHONY: all build test lint clean

all: build

build:
	mkdir -p bin
	go build -o bin/$(APP) ./cmd/launchd


test:
	go test -race -cover $(PKG)

clean:
	rm -rf bin
	rm -f coverage.out
