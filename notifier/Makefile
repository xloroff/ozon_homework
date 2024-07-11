run:
	go run ./cmd

run-all:
	go run ./cmd

test:
	go test -v -race -count 1 -vet=off ./...

build:
	go build -o ./bin/cart ./bin/cart

lint:
	golangci-lint run --config=./golangci.yaml

PHONY: .lint .test
