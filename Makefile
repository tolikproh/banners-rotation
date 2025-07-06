BIN := "./bin/banners"

BANNERS_IMG="banners:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/banners

run: build
	$(BIN) --config ./configs/config.yaml


version: build
	$(BIN) version

 test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.8

lint: install-lint-deps
	golangci-lint run ./...

migrate:
	goose -dir ./migrations up

migdown:
	goose -dir ./migrations down

.PHONY: build run version test install-lint-deps lint migrate migdown
