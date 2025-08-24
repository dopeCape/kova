# Colors for terminal output
YELLOW=\033[1;33m
GREEN=\033[1;32m
BLUE=\033[1;34m
RED=\033[1;31m
NC=\033[0m

# Database URL (can be overridden)
DB_URL ?= postgresql://admin:password123@localhost:5432/mydb?sslmode=disable

.PHONY: help build install dev clean docker-build docker-push tui api dashboard test generate migrate-up migrate-down build-cli run-cli all seed

# Helper function to check root privileges
define check_root
	@if [ "$$(id -u)" != "0" ]; then \
		echo "$(RED)Error: This command requires sudo$(NC)"; \
		echo "Run: $(YELLOW)sudo make $(1)$(NC)"; \
		exit 1; \
	fi
endef


# Helper function to check if kova is installed
define check_installed
	@if [ ! -f /opt/kova/docker-compose.yml ]; then \
		echo "$(RED)Kova not installed. Run 'make install' first.$(NC)"; \
		exit 1; \
	fi
endef


# Development commands
dev:
	air --build.cmd "go build -o bin/api cmd/api/main.go" --build.bin "./bin/api"

api:
	go run cmd/api/main.go

dashboard:
	cd ../../apps/web && npm run dev

dev-env:
	cd ../../ && docker-compose up -d postgres redis
	@echo "$(GREEN)Development services started.$(NC)"

seed:
	go run cmd/seeder/main.go

dev-env-stop:
	cd ../../ && docker-compose down

# Database commands
generate:
	sqlc generate

migrate-up: 
	@echo "$(YELLOW)Running migrations up...$(NC)"
	migrate -path internal/store/migrations -database "$(DB_URL)" up
	@echo "$(GREEN)Migrations completed!$(NC)"

migrate-down: 
	@echo "$(YELLOW)Running migrations down...$(NC)"
	migrate -path internal/store/migrations -database "$(DB_URL)" down 1
	@echo "$(GREEN)Migration rolled back!$(NC)"

db-shell:
	@echo "$(BLUE)Connecting to database...$(NC)"
	@docker exec -it kova-postgres psql -U kova -d kova 2>/dev/null || echo "$(RED)PostgreSQL container not running$(NC)"

db-reset:
	@echo "$(YELLOW)This will DELETE ALL DATA! Type 'yes' to continue:$(NC)" && read -r confirm
	@[ "$$confirm" = "yes" ] && $(MAKE) migrate-down && $(MAKE) migrate-up && echo "$(GREEN)Database reset complete$(NC)" || echo "$(BLUE)Cancelled$(NC)"

# CLI commands
build-cli:
	go build -o bin/cli cmd/cli/main.go

run-cli:
	go run cmd/cli/main.go

tui:
	$(call check_root,tui)
	@echo "$(BLUE)Starting Kova interactive installer...$(NC)"
	go run cmd/cli/main.go

# Building
build:
	@echo "$(BLUE)Building Kova binaries...$(NC)"
	go build -o ../../bin/kova cmd/cli/main.go && go build -o ../../bin/kova-api cmd/api/main.go
	@echo "$(GREEN)Build complete. Binaries in ./bin/$(NC)"

docker-build:
	@echo "$(BLUE)Building Docker images...$(NC)"
	cd ../../ && docker build -t kova/api:latest -f apps/core/Dockerfile . && docker build -t kova/dashboard:latest -f apps/web/Dockerfile .
	@echo "$(GREEN)Docker images built successfully$(NC)"

docker-push:
	docker push kova/api:latest && docker push kova/dashboard:latest

# Installation
install:
	$(call check_root,install)
	@echo "$(BLUE)Running direct installation...$(NC)"
	go run cmd/cli/main.go install

install-dev:
	$(call check_root,install-dev)
	@echo "$(BLUE)Installing with development settings...$(NC)"
	ADMIN_EMAIL=admin@example.com ADMIN_PASSWORD=dev123 ADMIN_USERNAME=admin go run cmd/cli/main.go install

# Testing
test:
	@echo "$(BLUE)Running tests...$(NC)" && go test ./...

test-verbose:
	@echo "$(BLUE)Running tests with verbose output...$(NC)" && go test -v ./...

# Maintenance
clean:
	@echo "$(BLUE)Cleaning build artifacts...$(NC)" && rm -rf ../../bin/ && go clean
	@cd ../../apps/web && rm -rf .next node_modules || true
	@docker system prune -f && echo "$(GREEN)Clean complete$(NC)"

logs:
	$(call check_installed)
	docker-compose -f /opt/kova/docker-compose.yml logs -f

status:
	$(call check_installed)
	docker-compose -f /opt/kova/docker-compose.yml ps

# Service management
stop:
	$(call check_installed)
	docker-compose -f /opt/kova/docker-compose.yml stop && echo "$(GREEN)Kova services stopped$(NC)"

start:
	$(call check_installed)
	docker-compose -f /opt/kova/docker-compose.yml start && echo "$(GREEN)Kova services started$(NC)"

restart:
	$(call check_installed)
	docker-compose -f /opt/kova/docker-compose.yml restart && echo "$(GREEN)Kova services restarted$(NC)"

update:
	$(call check_installed)
	@echo "$(BLUE)Updating Kova installation...$(NC)" && cd /opt/kova && docker-compose pull && docker-compose up -d
	@echo "$(GREEN)Kova updated successfully$(NC)"

backup:
	$(call check_installed)
	mkdir -p ../../backups && docker exec kova-postgres pg_dump -U kova kova > ../../backups/kova-db-$$(date +%Y%m%d-%H%M%S).sql
	tar -czf ../../backups/kova-data-$$(date +%Y%m%d-%H%M%S).tar.gz /data/kova/ && echo "$(GREEN)Backup created in ./backups/$(NC)"

# Code quality
fmt:
	@echo "$(BLUE)Formatting Go code...$(NC)" && go fmt ./...

lint:
	@echo "$(BLUE)Running linter...$(NC)" && golangci-lint run
