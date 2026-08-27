.PHONY: up down test run

up:
	docker compose up --build

down:
	docker compose down

test:
	go test ./...

# Local process (no Docker app container). Requires DATABASE_URL.
# Example: docker compose up -d postgres && set -a && . ./.env && set +a && make run
run:
	@test -n "$$DATABASE_URL" || { echo "DATABASE_URL is required (see .env.example)"; exit 1; }
	go run ./cmd/authlab
