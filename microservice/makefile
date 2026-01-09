.PHONY: run build test coverage coverage-html clean

# Variables
BINARY_NAME=main
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Run the API
run:
	go run main.go

# Build the binary
build:
	go build -o $(BINARY_NAME) .

# Run all tests
test:
	go test -v ./...

coverage:
	go test ./... -coverprofile=$(COVERAGE_FILE)
	go tool cover -func=$(COVERAGE_FILE)

coverage-html: coverage
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report: $(COVERAGE_HTML)"

