# KFG - Declarative shell compiler
# Build configuration

BINARY_NAME=kfg
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)"

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Directories
CMD_DIR=./src/cmd/kfg
BIN_DIR=./bin

.PHONY: all build clean test install lint fmt vet help

all: build

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Built: $(BIN_DIR)/$(BINARY_NAME)"

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)
	@$(GOCMD) clean

## test: Run tests
test:
	@echo "Running tests..."
	cd src && $(GOTEST) -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	cd src && $(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=src/coverage.out -o coverage.html

## install: Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) $(CMD_DIR)

## lint: Run linter
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found, please install it" && exit 1)
	golangci-lint run ./src/...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./src/...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	cd src && $(GOCMD) vet ./...

## mod: Download and tidy dependencies
mod:
	$(GOMOD) download
	$(GOMOD) tidy

## help: Show this help
help:
	@echo "KFG - Declarative Shell Compiler"
	@echo ""
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'