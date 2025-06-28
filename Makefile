ifneq (,$(wildcard ./.env))
	include .env
	export
endif

install:
	go mod download && go mod verify

.PHONY: build
build:
	go build -trimpath -ldflags="-s -w" -o ./build/baitomebot ./cmd/baitomebot/main.go

.PHONY: run
run:
	go run ./cmd/baitomebot/main.go

.PHONY: test
test:
	go test ./...

