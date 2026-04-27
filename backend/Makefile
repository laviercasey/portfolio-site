.PHONY: build run test lint migrate-up migrate-down

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./... -v -race

lint:
	golangci-lint run ./...

migrate-up:
	goose -dir migrations postgres "$$DATABASE_URL" up

migrate-down:
	goose -dir migrations postgres "$$DATABASE_URL" down
