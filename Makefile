# ============================================
# NimbusU Backend - Makefile
# ============================================

.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down migrate-create seed

# Default target
.DEFAULT_GOAL := help

# Variables
SERVICE_NAME := user-service
SERVICE_DIR := services/$(SERVICE_NAME)
BINARY_DIR := $(SERVICE_DIR)/bin
BINARY_NAME := $(SERVICE_NAME)
MIGRATIONS_DIR := $(SERVICE_DIR)/migrations
DATABASE_URL := postgres://nimbusu:password@localhost:5433/user_service_db?sslmode=disable

# ----------------
# Help
# ----------------
help: ## Show this help message
	@echo "NimbusU Backend - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

# ----------------
# Development
# ----------------
install: ## Install dependencies
	@echo "Installing shared module dependencies..."
	cd shared && go mod download
	@echo "Installing user-service dependencies..."
	cd $(SERVICE_DIR) && go mod download
	@echo "Dependencies installed successfully!"

tidy: ## Tidy go modules
	@echo "Tidying shared module..."
	cd shared && go mod tidy
	@echo "Tidying user-service..."
	cd $(SERVICE_DIR) && go mod tidy
	@echo "Modules tidied successfully!"

build: ## Build the user service binary
	@echo "Building $(SERVICE_NAME)..."
	cd $(SERVICE_DIR) && go build -o $(BINARY_DIR)/$(BINARY_NAME) cmd/main.go
	@echo "Build complete: $(SERVICE_DIR)/$(BINARY_DIR)/$(BINARY_NAME)"

run: ## Run the user service
	@echo "Starting $(SERVICE_NAME)..."
	cd $(SERVICE_DIR) && go run cmd/main.go

dev: ## Run with live reload (requires air: go install github.com/cosmtrek/air@latest)
	@echo "Starting $(SERVICE_NAME) in development mode..."
	cd $(SERVICE_DIR) && air

# ----------------
# Testing
# ----------------
test: ## Run all tests
	@echo "Running tests for shared module..."
	cd shared && go test -v ./...
	@echo "Running tests for $(SERVICE_NAME)..."
	cd $(SERVICE_DIR) && go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	cd $(SERVICE_DIR) && go test -v -coverprofile=coverage.out ./...
	cd $(SERVICE_DIR) && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: $(SERVICE_DIR)/coverage.html"

# ----------------
# Database
# ----------------
migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up
	@echo "Migrations applied successfully!"

migrate-down: ## Rollback last migration
	@echo "Rolling back last migration..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

migrate-reset: ## Reset all migrations (down then up)
	@echo "Resetting database..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down -all
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up
	@echo "Database reset complete!"

migrate-create: ## Create new migration (usage: make migrate-create name=add_users_table)
	@if [ -z "$(name)" ]; then \
		echo "Error: Migration name required. Usage: make migrate-create name=add_users_table"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

seed: ## Seed database with default data
	@echo "Seeding database..."
	psql "$(DATABASE_URL)" -f $(SERVICE_DIR)/migrations/seed.sql
	@echo "Database seeded successfully!"

db-shell: ## Open PostgreSQL shell
	psql "$(DATABASE_URL)"

# ----------------
# Docker & Infrastructure
# ----------------
docker-up: ## Start infrastructure (Kafka, PostgreSQL, Redis)
	@echo "Starting infrastructure..."
	docker-compose -f kafka/docker-compose.yaml up -d
	@echo "Infrastructure started!"
	@echo "Kafka: localhost:9092"
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"

docker-down: ## Stop infrastructure
	@echo "Stopping infrastructure..."
	docker-compose -f kafka/docker-compose.yaml down
	@echo "Infrastructure stopped!"

docker-logs: ## View infrastructure logs
	docker-compose -f kafka/docker-compose.yaml logs -f

docker-build: ## Build service Docker image
	@echo "Building Docker image for $(SERVICE_NAME)..."
	cd $(SERVICE_DIR) && docker build -t nimbusu/$(SERVICE_NAME):latest .

# ----------------
# Cleanup
# ----------------
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(SERVICE_DIR)/$(BINARY_DIR)
	rm -f $(SERVICE_DIR)/coverage.out
	rm -f $(SERVICE_DIR)/coverage.html
	@echo "Clean complete!"

clean-all: clean ## Clean all artifacts including dependencies
	@echo "Cleaning all artifacts..."
	cd shared && go clean -cache -modcache
	cd $(SERVICE_DIR) && go clean -cache -modcache
	@echo "All artifacts cleaned!"

# ----------------
# Code Quality
# ----------------
fmt: ## Format code
	@echo "Formatting code..."
	cd shared && go fmt ./...
	cd $(SERVICE_DIR) && go fmt ./...
	@echo "Code formatted!"

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	cd shared && golangci-lint run
	cd $(SERVICE_DIR) && golangci-lint run

vet: ## Run go vet
	@echo "Running go vet..."
	cd shared && go vet ./...
	cd $(SERVICE_DIR) && go vet ./...

# ----------------
# Quick Commands
# ----------------
setup: docker-up migrate-up seed ## Complete setup (start infrastructure, migrate, seed)
	@echo ""
	@echo "Setup complete! You can now run 'make run' to start the service."

start: ## Quick start (assumes infrastructure is running)
	@$(MAKE) build
	@$(MAKE) run

restart: ## Restart the service
	@pkill -f "$(BINARY_NAME)" || true
	@$(MAKE) start

# ----------------
# Environment
# ----------------
env: ## Create .env file from .env.example
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env file created! Please update with your configuration."; \
	else \
		echo ".env file already exists."; \
	fi
