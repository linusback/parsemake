BINARY_NAME=parsemake
SHELL=/usr/bin/env bash

all: test build

build:
	go build -ldflags='-s -w' -o dist/${BINARY_NAME} ./parsemake.go

test:
	go test ./...

run: build
	./dist/${BINARY_NAME}

dev:
	go run parsemake.go

clean:
	rm -f dist/${BINARY_NAME}
	go mod tidy
	go clean

vet:
	go vet

lint:
	golangci-lint run --enable-all

