SERVICES := gateway auth users vehicles bookings repair inventory parts-marketplace payments payroll lookup notifications search staff
BUILD_DIR := build

.PHONY: all build clean dev dev-down test lint migrate-up migrate-down seed help

SERVICE ?= auth
COMPOSE_PROJECT_NAME ?= autolab

all: build

# ============================================================================
# Development
# ============================================================================
dev: ## Start all services in development mode
	docker compose -p $(COMPOSE_PROJECT_NAME) up --build -d
	@echo "Gateway: http://localhost:4000"
	@echo "Auth:     http://localhost:8081"
	@echo "Lookup:   http://localhost:8090"

dev-down: ## Stop all development services
	docker compose -p $(COMPOSE_PROJECT_NAME) down

dev-logs: ## Tail logs from all services
	docker compose -p $(COMPOSE_PROJECT_NAME) logs -f

dev-logs-%: ## Tail logs from a specific service
	docker compose -p $(COMPOSE_PROJECT_NAME) logs -f $*

# ============================================================================
# Build
# ============================================================================
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

build: $(BUILD_DIR) $(SERVICES:%=build-%)
	@echo "Done"

build-%: ## Build a specific service
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o $(BUILD_DIR)/$* ./services/$*

# ============================================================================
# Database
# ============================================================================
migrate-up: ## Run all pending database migrations
	@echo "Running database migrations..."
	psql -h localhost -U postgres -d autolab -f deployments/migrations/init/000_extensions.sql
	@echo "Migrations complete"

migrate-down: ## Drop all tables (DANGER)
	@echo "WARNING: This will drop all tables!"
	@read -p "Are you sure? [y/N] " confirm; \
	if [ "$$confirm" = "y" ]; then \
		psql -h localhost -U postgres -d autolab -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"; \
	fi

migrate-new: ## Create a new migration file
	@read -p "Migration name: " name; \
	timestamp=$$(date +%Y%m%d%H%M%S); \
	touch "deployments/migrations/$${timestamp}_$${name}.sql"; \
	echo "Created: $${timestamp}_$${name}.sql"

seed: ## Seed the database with test data
	@echo "Seeding database..."
	psql -h localhost -U postgres -d autolab -f deployments/seed.sql 2>/dev/null || true
	@echo "Seed complete"

db-reset: migrate-down migrate-up seed ## Reset and re-seed the database

# ============================================================================
# Code Generation
# ============================================================================
.PHONY: generate generate-% generate-federation generate-federation-%

generate: ## Generate GraphQL code for all services
	powershell -ExecutionPolicy Bypass -File .scripts/generate.ps1

generate-%: ## Generate GraphQL code for a specific service
	powershell -ExecutionPolicy Bypass -File .scripts/generate.ps1 -Service $*

generate-federation: ## Regenerate federation types (all services)
	powershell -ExecutionPolicy Bypass -File .scripts/generate-federation.ps1

generate-federation-%: ## Regenerate federation types for a specific service
	powershell -ExecutionPolicy Bypass -File .scripts/generate-federation.ps1 -Service $*

# ============================================================================
# Testing
# ============================================================================
test: ## Run all tests
	go test ./... -v -count=1

test-%: ## Run tests for a specific service
	go test ./services/$*/... -v -count=1

test-cover: ## Run all tests with coverage
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# ============================================================================
# Linting & Quality
# ============================================================================
lint: ## Run linters
	@echo "Running golangci-lint..."
	golangci-lint run ./... --timeout 5m

tidy: ## Tidy Go module dependencies
	go mod tidy
	go mod verify

fmt: ## Format Go code
	go fmt ./...

vet: ## Run Go vet
	go vet ./...

# ============================================================================
# Service Management (hot reload with air)
# ============================================================================
run-%: ## Run a specific service with hot reload
	cd services/$* && air

# ============================================================================
# Cleanup
# ============================================================================
clean: ## Clean build artifacts
	rm -rf build/ coverage.out coverage.html
	@echo "Cleaned build artifacts"

down: dev-down clean ## Full cleanup

# ============================================================================
# Help
# ============================================================================
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
