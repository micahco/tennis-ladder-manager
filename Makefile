.PHONY: build
build:
	go build -ldflags="-s" -o=./bin/tlm ./cmd/tlm

.PHONY: run
run:
	go run ./cmd/tlm

.PHONY: audit
audit:
	@echo "Tidying and verifying module dependencies..."
	go mod tidy
	go mod verify
	@echo "Formatting code..."
	go fmt ./...
	@echo "Vetting code..."
	go vet ./...
	staticcheck ./...
	@echo "Running tests..."
	go test -race -vet=off ./...
