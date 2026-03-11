.PHONY: help build run test clean docker-build docker-up docker-down swagger fmt lint

PROJECT_NAME ?= $(shell basename "$(CURDIR)" | tr '[:upper:]' '[:lower:]' | sed -E 's/[^a-z0-9]+/-/g; s/^-+|-+$$//g')
COMPOSE_CMD = PROJECT_NAME=$(PROJECT_NAME) docker compose -p $(PROJECT_NAME)

help:
	@echo "BookStore Management API - Makefile"
	@echo ""
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application locally"
	@echo "  make test           - Run tests"
	@echo "  make test-verbose   - Run tests with verbose output"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-up      - Start containers with docker-compose"
	@echo "  make docker-down    - Stop containers"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make deps           - Download dependencies"

build:
	@echo "Building application..."
	go build -o cmd/main/main ./cmd/main

run: build
	@echo "Running application..."
	./cmd/main/main

test:
	@echo "Running tests..."
	go test -v ./...

test-verbose:
	@echo "Running tests with coverage..."
	go test -v -cover ./...

swagger:
	@echo "Generating Swagger documentation..."
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/main/main.go

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Running linter..."
	go vet ./...

docker-build:
	@echo "Building Docker image..."
	$(COMPOSE_CMD) build

docker-up:
	@echo "Starting containers..."
	$(COMPOSE_CMD) up -d

docker-down:
	@echo "Stopping containers..."
	$(COMPOSE_CMD) down

docker-logs:
	@echo "Showing container logs..."
	$(COMPOSE_CMD) logs -f

clean:
	@echo "Cleaning build artifacts..."
	rm -f cmd/main/main
	go clean

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

all: clean deps build test

.DEFAULT_GOAL := help
