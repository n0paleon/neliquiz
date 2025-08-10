.PHONY: run test build deploy migrate-up migrate-down migrate-force migrate-drop migrate-version test-watch prod-migrate-up

include .env

# Default config (for development)
MIGRATION_DIR ?= ./migrations
ENV_PATH := .env
DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# ======== Run app ========
run:
	ENV_PATH=$(ENV_PATH) go run ./cmd/api/main.go

# ======== Migration (for dev db) ========
migrate-up:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" down 1

migrate-force:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" force $(version)

migrate-drop:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" drop

migrate-version:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" version

# ======== Testing with test DB ========
TEST_DB_NAME = neliquiz_test
TEST_DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(TEST_DB_NAME)?sslmode=$(DB_SSLMODE)

test-watch:
	migrate -path $(MIGRATION_DIR) -database "$(TEST_DB_URL)" drop
	migrate -path $(MIGRATION_DIR) -database "$(TEST_DB_URL)" up
	dotenv -e .env.test -- gotestsum -f testname --format-icons hivis --watch -- -tags=integration ./...

test:
	migrate -path $(MIGRATION_DIR) -database "$(TEST_DB_URL)" drop
	migrate -path $(MIGRATION_DIR) -database "$(TEST_DB_URL)" up
	dotenv -e .env.test -- go test -v -tags=integration ./...
