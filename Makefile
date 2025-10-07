.PHONY: test-unit test clean build

# Run unit tests with verbose output
test-unit:
	go test -v ./...

# Run all tests (alias for test-unit for now)
test: test-unit

# Clean test cache and build artifacts
clean:
	go clean -testcache
	go clean

# Build the project
build:
	go build -v ./...
