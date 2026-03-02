# Run unit tests
test:
    go test ./...

# Run unit tests with verbose output
test-v:
    go test -v ./...

# Run integration tests (requires .env with TICKTICK_API_TOKEN)
test-integration:
    go test -tags integration -v ./...

# Run all tests
test-all:
    go test -tags integration -v ./...

# Lint
lint:
    golangci-lint run ./...

# Build (verifies compilation, removes binary)
build:
    go build ./cmd/tick && rm -f tick

# Run with arguments
run *ARGS:
    go run ./cmd/tick {{ARGS}}
