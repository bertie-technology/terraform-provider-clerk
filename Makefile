.PHONY: build install test testacc clean fmt docs

# Build the provider
build:
	mkdir -p bin
	go build -o bin/terraform-provider-clerk

# Install dependencies
install:
	go mod download
	go mod tidy

# Run unit tests
test:
	go test -v -count=1 -parallel=4 -timeout 5m ./...

# Run acceptance tests
testacc:
	TF_ACC=1 go test -v -count=1 -parallel=4 -timeout 30m ./...

# Format code
fmt:
	go fmt ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Generate documentation
docs:
	@echo "Generating documentation..."
	@rm -rf examples/.terraform examples/.terraform.lock.hcl .terraform-plugin-docs
	@$(shell go env GOPATH)/bin/tfplugindocs generate --provider-name clerk

# Validate documentation
docs-validate:
	@$(shell go env GOPATH)/bin/tfplugindocs validate --provider-name clerk

# Run with debug mode
debug:
	go run main.go -debug

# Initialize and run example
example: build
	cd examples && \
	terraform init && \
	terraform plan
