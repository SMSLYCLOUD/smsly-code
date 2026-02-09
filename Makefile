# SMSLY Code — Top-level Makefile
# This orchestrates all components in the monorepo

.PHONY: help dev-up dev-down dev-logs dev-reset build test lint clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ========================= Development =========================

dev-up: ## Start all development services
	docker compose -f docker-compose.dev.yml up -d --build
	@echo "SMSLY Code is running at http://localhost:3000"

dev-down: ## Stop all development services
	docker compose -f docker-compose.dev.yml down

dev-logs: ## Follow logs from all services
	docker compose -f docker-compose.dev.yml logs -f

dev-reset: ## Reset everything (destroy data + rebuild)
	docker compose -f docker-compose.dev.yml down -v
	docker compose -f docker-compose.dev.yml up -d --build

# ========================= Build =========================

build: build-git build-api build-web ## Build all components

build-git: ## Build Rust Git engine
	cd smsly-git && cargo build --release

build-api: ## Build Go API server
	cd smsly-code-api && go build -o smsly-code-api ./cmd/server

build-web: ## Build Next.js frontend
	cd smsly-code-web && npm run build

# ========================= Test =========================

test: test-git test-api test-web ## Run all tests

test-git: ## Run Rust tests
	cd smsly-git && cargo test

test-api: ## Run Go tests
	cd smsly-code-api && go test ./...

test-web: ## Run frontend tests
	cd smsly-code-web && npm test

# ========================= Lint =========================

lint: lint-git lint-api lint-web ## Lint all components

lint-git: ## Lint Rust code
	cd smsly-git && cargo clippy -- -D warnings

lint-api: ## Lint Go code
	cd smsly-code-api && golangci-lint run

lint-web: ## Lint frontend code
	cd smsly-code-web && npm run lint

# ========================= Production =========================

prod-up: ## Start production services
	docker compose up -d --build

prod-down: ## Stop production services
	docker compose down

prod-logs: ## Follow production logs
	docker compose logs -f

# ========================= Database =========================

migrate: ## Run database migrations
	cd smsly-code-api && go run ./cmd/migrate up

migrate-down: ## Rollback last migration
	cd smsly-code-api && go run ./cmd/migrate down

# ========================= Maintenance =========================

clean: ## Clean all build artifacts
	cd smsly-git && cargo clean
	cd smsly-code-api && rm -f smsly-code-api
	cd smsly-code-web && rm -rf .next out

backup: ## Create full backup
	./deploy/scripts/backup.sh

restore: ## Restore from backup (use BACKUP=date)
	./deploy/scripts/restore.sh $(BACKUP)
