ENTRYPOINT_NAME=main
BINARY_NAME=parsemake

SHELL=/usr/bin/env bash
MSG="Hello linus"

all: test build

# perhaps add GOAMD64=v3 to architecture
build:
	go build -ldflags='-s -w' -o dist/${BINARY_NAME} ./cmd/${ENTRYPOINT_NAME}.go

tests:
	go test ./...

bench:
	go test ./... -bench=. -benchtime 3s -run=^\# -cpu=1,20

echo:
	echo $(MSG)

echo-multi:
	echo """hello\
	my vr\
	"""

run: build
	./dist/${BINARY_NAME}

dev:
	go run parsemake.go

clean:
	rm -f dist/*
	go mod tidy
	go clean

vet:
	go vet

lint:
	golangci-lint run --enable-all
