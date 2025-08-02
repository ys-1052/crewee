.PHONY: help dev build test clean fmt lint migrate-up migrate-down migrate-create docker-up docker-down region-data-convert region-data-migrate python-fmt python-lint sql-fmt

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
	@if [ -z "${DATABASE_URL}" ]; then \
		echo "Error: DATABASE_URL environment variable is not set."; \
		echo "Please set DATABASE_URL (e.g., export DATABASE_URL=\"postgres://user:password@localhost:5432/crewee_dev\")"; \
		exit 1; \
	fi
	migrate -path migrations -database "${DATABASE_URL}" up

migrate-down: ## Run database migrations down (all)
	@echo "Running migrations down..."
	@if [ -z "${DATABASE_URL}" ]; then \
		echo "Error: DATABASE_URL environment variable is not set."; \
		echo "Please set DATABASE_URL (e.g., export DATABASE_URL=\"postgres://user:password@localhost:5432/crewee_dev\")"; \
		exit 1; \
	fi
	migrate -path migrations -database "${DATABASE_URL}" down

migrate-up-1: ## Run database migrations up by 1 step
	@echo "Running migrations up by 1 step..."
	@if [ -z "${DATABASE_URL}" ]; then \
		echo "Error: DATABASE_URL environment variable is not set."; \
		echo "Please set DATABASE_URL (e.g., export DATABASE_URL=\"postgres://user:password@localhost:5432/crewee_dev\")"; \
		exit 1; \
	fi
	migrate -path migrations -database "${DATABASE_URL}" up 1

migrate-down-1: ## Run database migrations down by 1 step
	@echo "Running migrations down by 1 step..."
	@if [ -z "${DATABASE_URL}" ]; then \
		echo "Error: DATABASE_URL environment variable is not set."; \
		echo "Please set DATABASE_URL (e.g., export DATABASE_URL=\"postgres://user:password@localhost:5432/crewee_dev\")"; \
		exit 1; \
	fi
	migrate -path migrations -database "${DATABASE_URL}" down 1

migrate-status: ## Show current migration status
	@echo "Current migration status:"
	@if [ -z "${DATABASE_URL}" ]; then \
		echo "Error: DATABASE_URL environment variable is not set."; \
		echo "Please set DATABASE_URL (e.g., export DATABASE_URL=\"postgres://user:password@localhost:5432/crewee_dev\")"; \
		exit 1; \
	fi
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

# Region master data
region-data-convert: ## Convert Excel file to CSV only
	@echo "Converting local_government_codes.xlsx to CSV..."
	@if [ ! -f "../local_government_codes.xlsx" ]; then \
		echo "Error: local_government_codes.xlsx not found in project root"; \
		exit 1; \
	fi
	@echo "Setting up Python environment..."
	@if [ ! -d "../venv" ]; then \
		python3 -m venv ../venv; \
		source ../venv/bin/activate && pip install pandas openpyxl xlrd; \
	fi
	@echo "Converting Excel to CSV..."
	@cd .. && source venv/bin/activate && python3 scripts/convert_lgcode_to_csv.py
	@echo "CSV conversion complete!"

region-data-migrate: ## Convert region data and generate migration SQL (handles both initial and incremental updates)
	@echo "Converting local_government_codes.xlsx to CSV..."
	@if [ ! -f "../local_government_codes.xlsx" ]; then \
		echo "Error: local_government_codes.xlsx not found in project root"; \
		exit 1; \
	fi
	@echo "Setting up Python environment..."
	@if [ ! -d "../venv" ]; then \
		python3 -m venv ../venv; \
		source ../venv/bin/activate && pip install pandas openpyxl xlrd; \
	fi
	@echo "Converting Excel to CSV..."
	@cd .. && source venv/bin/activate && python3 scripts/convert_lgcode_to_csv.py
	@echo "Generating migration SQL (initial setup or incremental updates)..."
	@cd .. && source venv/bin/activate && python3 scripts/generate_region_diff_migration.py
	@echo "Cleaning up temporary CSV file..."
	@rm -f sql/local_government_codes.csv
	@if ls migrations/*update_regions_data.up.sql >/dev/null 2>&1; then \
		echo "✓ Migration SQL files generated. Run 'make migrate-up' to apply changes."; \
	else \
		echo "No changes detected. Migration files not generated."; \
	fi

# Python code quality
python-fmt: ## Format Python scripts
	@echo "Formatting Python scripts..."
	@if [ ! -d "../venv" ]; then \
		echo "Python environment not found. Run region-data-convert first."; \
		exit 1; \
	fi
	@cd .. && source venv/bin/activate && black scripts/ && isort scripts/
	@echo "Python formatting complete!"

python-lint: ## Lint Python scripts
	@echo "Linting Python scripts..."
	@if [ ! -d "../venv" ]; then \
		echo "Python environment not found. Run region-data-convert first."; \
		exit 1; \
	fi
	@cd .. && source venv/bin/activate && flake8 scripts/
	@echo "Python linting complete!"

# SQL formatting
sql-fmt: ## Format SQL files using sql-formatter
	@echo "Formatting SQL files..."
	@if ! command -v sql-formatter >/dev/null 2>&1; then \
		echo "Installing sql-formatter..."; \
		npm install -g sql-formatter; \
	fi
	@echo "Formatting migration files..."
	@find migrations -name "*.sql" -exec sql-formatter --config ../.sql-formatter.json --fix {} \;
	@echo "Formatting query files..."
	@find sql/queries -name "*.sql" -exec sql-formatter --config ../.sql-formatter.json --fix {} \; 2>/dev/null || true
	@echo "SQL formatting complete!"
