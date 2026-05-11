.PHONY: build run dev test clean db-start db-stop sqlc

# Build the server binary
build:
	go build -o bin/mthen-server ./cmd/server

# Run the server
run:
	go run ./cmd/server

# Development with hot reload (requires air)
dev:
	air

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Generate sqlc code
sqlc:
	sqlc generate

# Start PostgreSQL via Docker
db-start:
	docker run -d --name mthen-postgres \
		-e POSTGRES_USER=mthen \
		-e POSTGRES_PASSWORD=mthen \
		-e POSTGRES_DB=mthen \
		-p 5432:5432 \
		postgres:17-alpine

# Stop PostgreSQL
db-stop:
	docker stop mthen-postgres && docker rm mthen-postgres

# Run migrations
db-migrate:
	cd ../mthen-db && go run ./cmd/migrate

# Docker build
docker-build:
	docker build -t mthen-api .

# Docker run
docker-run:
	docker run -p 8080:8080 --env-file .env mthen-api
