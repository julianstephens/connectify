.PHONY: all dev build test fmt vet lint migrate-%

help:
	@echo "Makefile commands:"
	@echo "  dev           - Run server and client in development mode"
	@echo "  build         - Build server and client for production"
	@echo "  test          - Run tests for server and client"
	@echo "  fmt           - Format code for server and client"
	@echo "  lint          - Lint code for server and client"
	@echo "  migrate-up    - Apply database migrations"
	@echo "  migrate-down  - Rollback database migrations"help:

all: build

dev:
	# Run server and client in parallel (use tmux, two terminals, or your IDE)
	@echo "Run server: cd server && go run ./cmd/api"
	@echo "Run client: cd client && npm install && npm run dev"

build:
	cd server && go build -o bin/api ./cmd/api
	cd client && npm run build

test:
	cd server && go test ./...
	cd client && npm test -- --watchAll=false

fmt:
	cd server && go fmt ./...
	cd client && npm run format

lint:
	cd server && golangci-lint run
	cd client && npm run lint

migrate-up:
	@migrate -database $$POSTGRESQL_URL -path server/migrations/ up

migrate-down:
	@migrate -database $$POSTGRESQL_URL -path server/migrations/ down
