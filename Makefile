# Makefile for Lemonade Store DDD Application

# Variables
BINARY_NAME=lemonade-store
GO=go
GOFLAGS=-v
PROTO_DIR=internal/interfaces/grpc/proto
PB_DIR=internal/interfaces/grpc/pb
BUF=buf

# Colors for terminal output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## deps: Download and tidy dependencies
.PHONY: deps
deps:
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(GREEN)Dependencies updated!$(NC)"

## buf-deps: Install buf and required plugins
.PHONY: buf-deps
buf-deps:
	@echo "$(YELLOW)Installing buf...$(NC)"
	@which buf > /dev/null || (echo "Installing buf..." && go install github.com/bufbuild/buf/cmd/buf@latest)
	@echo "$(GREEN)Buf installed!$(NC)"

## proto: Generate Go code from proto files using buf
.PHONY: proto
proto: buf-deps
	@echo "$(YELLOW)Generating proto files...$(NC)"
	@mkdir -p $(PB_DIR)
	$(BUF) generate
	@echo "$(GREEN)Proto generation complete!$(NC)"

## build: Build the binary
.PHONY: build
build: proto
	@echo "$(YELLOW)Building $(BINARY_NAME)...$(NC)"
	$(GO) build $(GOFLAGS) -o bin/$(BINARY_NAME) cmd/grpc/main.go
	@echo "$(GREEN)Build complete! Binary: bin/$(BINARY_NAME)$(NC)"

## run: Run the application
.PHONY: run
run: build
	@echo "$(YELLOW)Starting $(BINARY_NAME)...$(NC)"
	./bin/$(BINARY_NAME)

## dev: Run the application with hot reload (requires air)
.PHONY: dev
dev: proto
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

## test: Run all tests
.PHONY: test
test:
	@echo "$(YELLOW)Running tests...$(NC)"
	$(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@echo "$(GREEN)Tests complete!$(NC)"

## test-unit: Run unit tests only
.PHONY: test-unit
test-unit:
	@echo "$(YELLOW)Running unit tests...$(NC)"
	$(GO) test -v -short ./...
	@echo "$(GREEN)Unit tests complete!$(NC)"

## test-integration: Run integration tests
.PHONY: test-integration
test-integration:
	@echo "$(YELLOW)Running integration tests...$(NC)"
	$(GO) test -v -run Integration ./...
	@echo "$(GREEN)Integration tests complete!$(NC)"

## coverage: Generate test coverage report
.PHONY: coverage
coverage: test
	@echo "$(YELLOW)Generating coverage report...$(NC)"
	$(GO) tool cover -html=coverage.txt -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

## lint: Run linters
.PHONY: lint
lint:
	@echo "$(YELLOW)Running linters...$(NC)"
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...
	@echo "$(YELLOW)Running buf lint...$(NC)"
	cd $(PROTO_DIR) && $(BUF) lint
	@echo "$(GREEN)Linting complete!$(NC)"

## fmt: Format code
.PHONY: fmt
fmt:
	@echo "$(YELLOW)Formatting code...$(NC)"
	$(GO) fmt ./...
	@echo "$(YELLOW)Formatting proto files...$(NC)"
	cd $(PROTO_DIR) && $(BUF) format -w
	@echo "$(GREEN)Formatting complete!$(NC)"

## clean: Clean build artifacts
.PHONY: clean
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	$(GO) clean
	rm -rf bin/
	rm -rf $(PB_DIR)/*.pb.go
	rm -f coverage.txt coverage.html
	@echo "$(GREEN)Clean complete!$(NC)"

## docker-build: Build Docker image
.PHONY: docker-build
docker-build:
	@echo "$(YELLOW)Building Docker image...$(NC)"
	docker build -t $(BINARY_NAME):latest .
	@echo "$(GREEN)Docker build complete!$(NC)"

## docker-run: Run the application in Docker
.PHONY: docker-run
docker-run: docker-build
	@echo "$(YELLOW)Running Docker container...$(NC)"
	docker run -p 50051:50051 $(BINARY_NAME):latest

## grpc-test: Test gRPC endpoints with grpcurl
.PHONY: grpc-test
grpc-test:
	@echo "$(YELLOW)Testing gRPC endpoints...$(NC)"
	@which grpcurl > /dev/null || (echo "Please install grpcurl: brew install grpcurl" && exit 1)
	@echo "$(GREEN)Listing services...$(NC)"
	grpcurl -plaintext localhost:50051 list
	@echo "$(GREEN)Testing RegisterCustomer...$(NC)"
	grpcurl -plaintext -d '{"email":"test@example.com","first_name":"Test","last_name":"User"}' \
		localhost:50051 lemonade.customer.v1.CustomerService/RegisterCustomer

## setup: Initial project setup
.PHONY: setup
setup: deps buf-deps proto
	@echo "$(GREEN)Project setup complete!$(NC)"

## all: Run all steps
.PHONY: all
all: clean deps proto lint test build
	@echo "$(GREEN)All steps complete!$(NC)"

# Include .env file if it exists
-include .env

.PHONY: vendor
vendor:
	$(GO) mod vendor
