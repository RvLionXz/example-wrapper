# Makefile for Gemini API Wrapper

# Variables
BINARY=api
MAIN_FILE=cmd/api/main.go

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	go build -o ${BINARY} ${MAIN_FILE}

# Run the application
.PHONY: run
run:
	go run ${MAIN_FILE}

# Install dependencies
.PHONY: deps
deps:
	go mod tidy

# Clean build artifacts
.PHONY: clean
clean:
	rm -f ${BINARY}

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run tests with coverage
.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Vet code
.PHONY: vet
vet:
	go vet ./...

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all       - Build the application (default)"
	@echo "  build     - Build the application"
	@echo "  run       - Run the application"
	@echo "  deps      - Install dependencies"
	@echo "  clean     - Clean build artifacts"
	@echo "  test      - Run tests"
	@echo "  coverage  - Run tests with coverage"
	@echo "  fmt       - Format code"
	@echo "  vet       - Vet code"
	@echo "  help      - Show this help"