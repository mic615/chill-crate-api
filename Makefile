# Variables
BINARY_NAME=chill-crate-api
BUILD_DIR=bin

.PHONY: up down restart logs build run test clean tidy fmt linter help

all: build

## build: Build the Go binary
build:
	@echo "Building binary..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/server/main.go 

## run: Build and run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

## test: Run all unit tests
test:
	@echo "Running tests..."
	@go test -v -race ./...

## clean: Remove build artifacts
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)

## tidy: Add missing and remove unused modules
tidy:
	@go mod tidy

## fmt: Run go fmt against all packages
fmt:
	@golangci-lint fmt ./...

## linter: Run golangci-lint
lint:
	@golangci-lint run ./...

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .

up:
	@echo "Starting Docker Compose..."
	@docker-compose up -d

down:
	@echo "Stopping Docker Compose..."
	@docker-compose down

restart: down up ## Restart all containers

logs: ## Tail the API logs
	docker compose logs -f api