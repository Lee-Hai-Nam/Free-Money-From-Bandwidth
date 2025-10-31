.PHONY: help install build dev clean test lint

help:
	@echo "Available commands:"
	@echo "  make install    - Install dependencies"
	@echo "  make build       - Build the application"
	@echo "  make dev         - Run in development mode"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make lint        - Run linters"

install:
	@echo "Installing dependencies..."
	go mod tidy
	cd frontend && npm install

build:
	@echo "Building application..."
	wails build

dev:
	@echo "Starting development mode..."
	wails dev

clean:
	@echo "Cleaning artifacts..."
	rm -rf build/
	rm -rf frontend/dist/
	rm -rf frontend/node_modules/

test:
	@echo "Running tests..."
	go test ./...

lint:
	@echo "Running linters..."
	golangci-lint run

