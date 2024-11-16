ENTRYPOINT_NAME=main
BINARY_NAME=parsemake

SHELL=/usr/bin/env bash
現   =    		Hello linus

.MSG = "good night"

BENCH=

all: tests build

# perhaps add GOAMD64=v3 to architecture
build:
	go build -ldflags='-s -w' -o dist/${BINARY_NAME} ./cmd/${ENTRYPOINT_NAME}.go

echo-build:
	echo $(現) $(subst /,-,$(BENCH)) ${.MSG}

tests:
	go test ./...

bench:
	go test ./... -bench=. -benchtime 3s -run=^\# -cpu=1,20

bench-prof:
	go test . -bench=${BENCH} -benchtime 3s -run=^\# -cpu=20 -cpuprofile ./tmp/$(subst /,-,$(BENCH))_cpu.prof -memprofile ./tmp/$(subst /,-,$(BENCH))_mem.prof -o ./tmp/$(subst /,-,$(BENCH)).test

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
