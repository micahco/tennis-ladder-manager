.PHONY: build
build:
	GODEBUG=asyncpreemptoff=1 go build -ldflags="-s" -o=./bin/tlm ./cmd/tlm


.PHONY: kill
kill:
	pkill -f python


SERVERS = services/quote_server.py services/validator.py services/artwork.py services/statistics.py
.PHONY: services
services: kill
	@echo -e "\nSTARTING MICROSERVICES:\n"
	$(foreach server, $(SERVERS), .venv/bin/python $(server) &)


.PHONY: run
run:
	GODEBUG=asyncpreemptoff=1 go run ./cmd/tlm	


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
