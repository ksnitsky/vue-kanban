.PHONY: dev backend-dev frontend-dev build test clean migrate deps

# =============================================================================
# Container Runtime Detection
# =============================================================================

# Приоритет: docker > podman
DOCKER := $(shell command -v docker 2>/dev/null)
PODMAN := $(shell command -v podman 2>/dev/null)
DOCKER_COMPOSE := $(shell command -v docker-compose 2>/dev/null)
PODMAN_COMPOSE := $(shell command -v podman-compose 2>/dev/null)

# Выбираем доступный runtime
CONTAINER_RUNTIME := $(or $(DOCKER),$(PODMAN))
COMPOSE_CMD := $(or $(DOCKER_COMPOSE),$(PODMAN_COMPOSE))

# Проверка наличия runtime
ifeq ($(CONTAINER_RUNTIME),)
$(error Neither docker nor podman found. Please install one of them.)
endif

ifeq ($(COMPOSE_CMD),)
$(error Neither docker-compose nor podman-compose found. Please install one of them.)
endif

# Вывод информации
$(info Using container runtime: $(CONTAINER_RUNTIME))
$(info Using compose command: $(COMPOSE_CMD))

# =============================================================================
# Development
# =============================================================================

# Запустить всё (DB + backend + frontend)
dev: db-up
	@echo "Starting backend and frontend..."
	@make -j2 backend-dev frontend-dev

# Только backend
backend-dev:
	cd backend && go run cmd/server/main.go

# Только frontend
frontend-dev:
	cd frontend && bun run dev

# =============================================================================
# Build
# =============================================================================

build:
	$(CONTAINER_RUNTIME) build -t kanban .

# =============================================================================
# Testing
# =============================================================================

test:
	cd backend && go test ./... -v

test-unit:
	cd backend && go test ./internal/... -v -short

test-integration:
	cd backend && go test ./tests/integration/... -v

test-coverage:
	cd backend && go test ./... -coverprofile=coverage.out
	cd backend && go tool cover -html=coverage.out

# =============================================================================
# Database
# =============================================================================

db-up:
	$(COMPOSE_CMD) up -d

db-down:
	$(COMPOSE_CMD) down

db-reset:
	$(COMPOSE_CMD) down -v
	$(COMPOSE_CMD) up -d

db-logs:
	$(COMPOSE_CMD) logs -f

# =============================================================================
# Dependencies
# =============================================================================

deps:
	cd backend && go mod download
	cd frontend && bun install

deps-backend:
	cd backend && go mod download

deps-frontend:
	cd frontend && bun install

# =============================================================================
# Cleanup
# =============================================================================

clean:
	$(COMPOSE_CMD) down -v
	rm -rf backend/dist

# =============================================================================
# Migrations (dbmate)
# =============================================================================

migrate-up:
	$(COMPOSE_CMD) up dbmate

migrate-down:
	$(COMPOSE_CMD) run --rm dbmate down

migrate-create:
	$(COMPOSE_CMD) run --rm dbmate new $(name)

migrate-status:
	$(COMPOSE_CMD) run --rm dbmate status
