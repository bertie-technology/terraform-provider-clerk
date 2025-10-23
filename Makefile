.PHONY: build install test clean fmt

# Build the provider
build:
	mkdir -p bin
	go build -o bin/terraform-provider-clerk

# Install dependencies
install:
	go mod download
	go mod tidy

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run with debug mode
debug:
	go run main.go -debug

# Initialize and run example
example: build
	cd examples && \
	terraform init && \
	terraform plan
