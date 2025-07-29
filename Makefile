include .env
export

BIN := "./bin/banners"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/banners

version: build
	$(BIN) version

run:
	docker-compose up -d

down:
	docker-compose down -v

test:
	go test -race -count 100 ./internal/... 

test-int: run
	sleep 30 ;\
	set -e ;\
	test_status_code=0 ;\
	docker-compose run tests go test github.com/tolikproh/banners-rotation/cmd/tests/integration || test_status_code=$$? ;\
	docker-compose down ;\
	docker container prune -f ;\
	exit $$test_status_code ;

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.8

lint: install-lint-deps
	golangci-lint run ./...


migrate:
	goose -dir ./migrations up

migdown:
	goose -dir ./migrations down

.PHONY: build run version test install-lint-deps lint migrate migdown
