.PHONY: all build run test clean

all: build

build:
	go build -o bin/api ./cmd/api

run:
	go run ./cmd/api

test:
	go test -v ./...

clean:
	rm -rf bin/
