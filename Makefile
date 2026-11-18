# Shortcuts for everyday commands: `make dev`, `make check`, ...
# Recipes stay plain (no shell tricks) so they run the same in PowerShell, Git Bash and CI.

.PHONY: db-up db-down dev run fmt vet lint test build check

# Start PostgreSQL in Docker and wait until it is healthy.
db-up:
	docker compose up -d --wait

# Stop PostgreSQL (data stays in the db-data volume).
db-down:
	docker compose down

# Run the API with live reload.
dev: db-up
	go tool air

# Run the API once, without live reload.
run: db-up
	go run ./cmd/api

# Fix formatting.
fmt:
	gofmt -w .

vet:
	go vet ./...

lint:
	golangci-lint run ./...

test:
	go test ./...

# bin/api.exe on Windows, bin/api on Linux / macOS.
build:
	go build -o bin/ ./cmd/api

# Everything CI checks. Run before every commit.
check: vet lint test
