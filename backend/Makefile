.PHONY: help dev build test clean fmt lint migrate-up migrate-down migrate-create docker-up docker-down

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
dev: ## Start development server with hot reload
	@echo "Starting development server..."
	air

build: ## Build the application
	@echo "Building application..."
	go build -o bin/server cmd/server/main.go

run: ## Run the application without hot reload
	@echo "Running application..."
	go run cmd/server/main.go

# Testing
test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Code quality
fmt: ## Format Go code
	@echo "Formatting code..."
	gofmt -w .
	goimports -w .

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

# Database migrations
migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	migrate -path migrations -database "${DATABASE_URL}" up

migrate-down: ## Run database migrations down (all)
	@echo "Running migrations down..."
	migrate -path migrations -database "${DATABASE_URL}" down

migrate-up-1: ## Run database migrations up by 1 step
	@echo "Running migrations up by 1 step..."
	migrate -path migrations -database "${DATABASE_URL}" up 1

migrate-down-1: ## Run database migrations down by 1 step
	@echo "Running migrations down by 1 step..."
	migrate -path migrations -database "${DATABASE_URL}" down 1

migrate-status: ## Show current migration status
	@echo "Current migration status:"
	migrate -path migrations -database "${DATABASE_URL}" version

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
	migrate create -ext sql -dir migrations $(NAME)

# Database setup
db-reset: migrate-down migrate-up ## Reset database (down then up)

# Docker
docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	docker-compose up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs: ## Show Docker container logs
	docker-compose logs -f

# sqlc
sqlc-generate: ## Generate Go code from SQL queries
	@echo "Generating sqlc code..."
	sqlc generate

# Clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf tmp/
	rm -f coverage.out coverage.html

# Setup development environment
setup: docker-up ## Setup development environment
	@echo "Setting up development environment..."
	@echo "Installing Air for hot reload..."
	go install github.com/air-verse/air@latest
	@echo "Installing sqlc..."
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "Installing golang-migrate..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Installing golangci-lint..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing goimports..."
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "Setup complete!"

# Environment
env: ## Copy .env.example to .env
	@if [ ! -f .env ]; then \
		cp ../.env.example .env; \
		echo ".env file created from .env.example"; \
	else \
		echo ".env file already exists"; \
	fi
