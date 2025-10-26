# EasyGo Panel Makefile

PROJECT_NAME := easygo
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.CommitHash=$(COMMIT_HASH)

.PHONY: all build clean test install deps run dev help

all: build

## build: Build the application
build:
	@echo "Building $(PROJECT_NAME) $(VERSION)..."
	@mkdir -p build
	CGO_ENABLED=1 go build -ldflags "$(LDFLAGS)" -o build/$(PROJECT_NAME) cmd/easygo/main.go

## build-linux: Build for Linux (AMD64)
build-linux:
	@echo "Building for Linux AMD64..."
	@mkdir -p build
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "$(LDFLAGS)" -o build/$(PROJECT_NAME)-linux-amd64 cmd/easygo/main.go

## build-all: Build for all platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p build
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "$(LDFLAGS)" -o build/$(PROJECT_NAME)-linux-amd64 cmd/easygo/main.go
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o build/$(PROJECT_NAME)-linux-arm64 cmd/easygo/main.go

## test: Run tests
test:
	@echo "Running tests..."
	go test -v ./...

## clean: Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf build/
	go clean

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

## run: Run the application in web mode
run: build
	@echo "Starting EasyGo Panel..."
	./build/$(PROJECT_NAME) web

## dev: Run in development mode with auto-reload
dev:
	@echo "Starting development server..."
	go run cmd/easygo/main.go web --port 8080

## install: Install the application system-wide
install: build-linux
	@echo "Installing $(PROJECT_NAME)..."
	@if [ "$$(id -u)" != "0" ]; then echo "Please run as root"; exit 1; fi
	mkdir -p /opt/easygo
	cp build/$(PROJECT_NAME)-linux-amd64 /opt/easygo/$(PROJECT_NAME)
	chmod +x /opt/easygo/$(PROJECT_NAME)
	ln -sf /opt/easygo/$(PROJECT_NAME) /usr/local/bin/$(PROJECT_NAME)
	@echo "$(PROJECT_NAME) installed to /opt/easygo/"
	@echo "CLI available as: $(PROJECT_NAME)"

## package: Create installation package
package: build-all
	@echo "Creating installation package..."
	@./build.sh

## format: Format Go code
format:
	@echo "Formatting code..."
	go fmt ./...

## lint: Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

## security: Run security checks
security:
	@echo "Running security checks..."
	gosec ./...

## docker: Build Docker image
docker:
	@echo "Building Docker image..."
	docker build -t $(PROJECT_NAME):$(VERSION) .

## help: Show this help
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

# Development helpers
.PHONY: watch
## watch: Auto-rebuild on file changes (requires entr)
watch:
	@echo "Watching for changes..."
	find . -name "*.go" | entr -r make run

# Git hooks
.PHONY: hooks
## hooks: Install git hooks
hooks:
	@echo "Installing git hooks..."
	@cp scripts/pre-commit .git/hooks/
	@chmod +x .git/hooks/pre-commit

# Default target
.DEFAULT_GOAL := help