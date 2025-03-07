.PHONY: run lint

# Run examples
run:
	DEBUG=1 go run main.go

# Run linter
lint:
	@which golangci-lint > /dev/null; if [ $$? -eq 0 ]; then \
		echo "Running golangci-lint..."; \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi 