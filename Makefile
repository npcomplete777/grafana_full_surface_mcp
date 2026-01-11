.PHONY: build clean test fmt vet

BINARY_NAME=grafana-mcp-server
BUILD_DIR=bin

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

build-linux:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/server

build-darwin:
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/server

build-all: build-linux build-darwin

clean:
	rm -rf $(BUILD_DIR)

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test -v ./...

run:
	go run ./cmd/server

install:
	go install ./cmd/server
