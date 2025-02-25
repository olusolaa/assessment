.PHONY: build test lint clean coverage

# Default Go build flags
GOFLAGS := -v

# Build the application
build:
	go build $(GOFLAGS) -o main ./cmd/main.go

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
coverage:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Install golangci-lint if not installed and run it
lint:
	command -v golangci-lint >/dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
	golangci-lint run --timeout=5m

# Run the application
run:
	go run ./cmd/main.go

# Clean build artifacts
clean:
	rm -f main
	rm -f coverage.out
	rm -f coverage.html

# Install dependencies
deps:
	go mod download

# Run security scan using gosec
security:
	command -v gosec >/dev/null 2>&1 || go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec ./...

# Help command
help:
	@echo "Available targets:"
	@echo "  build     - Build the application"
	@echo "  test      - Run tests"
	@echo "  coverage  - Run tests with coverage"
	@echo "  lint      - Run linter"
	@echo "  run       - Run the application"
	@echo "  clean     - Clean build artifacts"
	@echo "  deps      - Install dependencies"
	@echo "  security  - Run security scan"
	@echo "  help      - Show this help message" 