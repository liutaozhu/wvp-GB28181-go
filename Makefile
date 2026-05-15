.PHONY: build run test clean docker-build swag

# Build the binary
build:
	@echo "Building WVP-PRO-GO..."
	CGO_ENABLED=0 go build -ldflags="-s -w" -o wvp ./cmd/wvp/
	@echo "Build complete: ./wvp"

# Run the application
run:
	@echo "Starting WVP-PRO-GO..."
	go run ./cmd/wvp/

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f wvp
	rm -f coverage.out
	go clean

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t wvp-pro-go:latest .

# Generate Swagger documentation
swag:
	@echo "Generating Swagger docs..."
	swag init -g cmd/wvp/main.go -o docs/swagger --parseDependency

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Cross-compile for Windows
build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o wvp.exe ./cmd/wvp/
	@echo "Build complete: ./wvp.exe"

# Cross-compile for Linux ARM64
build-linux-arm64:
	@echo "Building for Linux ARM64..."
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o wvp-arm64 ./cmd/wvp/
	@echo "Build complete: ./wvp-arm64"

# Show help
help:
	@echo "WVP-PRO-GO Makefile targets:"
	@echo "  build           - Build the binary (Linux)"
	@echo "  run             - Run the application"
	@echo "  test            - Run tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  docker-build    - Build Docker image"
	@echo "  swag            - Generate Swagger docs"
	@echo "  deps            - Install dependencies"
	@echo "  build-windows   - Cross-compile for Windows"
	@echo "  build-linux-arm64 - Cross-compile for Linux ARM64"
	@echo "  help            - Show this help"
