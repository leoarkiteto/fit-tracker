BINARY=api-server
MAIN=./cmd/api
BUILD_DIR=./bin

.DEFAULT_GOAL := build

.PHONY: build run test clean fmt vet tidy lint help

build:  ## Build the binary
	go build -o $(BUILD_DIR)/$(BINARY) $(MAIN)

run:  ## Run the server (development)
	go run $(MAIN)

test:  ## Run all tests
	go test ./...

clean:  ## Remove build artifacts
	rm -rf $(BUILD_DIR)

fmt:  ## Format Go source code
	go fmt ./...

vet:  ## Run Go vet
	go vet ./...

tidy:  ## Tidy Go modules
	go mod tidy

lint:  ## Run golangci-lint (if installed)
	golangci-lint run ./... || echo "golangci-lint not installed, skipping"

help:  ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'
