export PATH := $(PATH):$(shell go env GOPATH)/bin

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1

lint:
	golangci-lint run

build:
	go build -o bin/tl .

.PHONY: lint-install lint build
