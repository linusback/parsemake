ENTRYPOINT_NAME=main
BINARY_NAME=parsemake

SHELL=/usr/bin/env bash


all: test build

# perhaps add GOAMD64=v3 to architecture
build:
	go build -ldflags='-s -w' -o dist/${BINARYNAME} ./cmd/${ENTRYPOINT_NAME}.go

test:
	go test ./...

bench:
	go test ./... -bench=. -benchtime 3s -run=^\# -cpu=1,20

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

