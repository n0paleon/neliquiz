# Makefile

.PHONY: run test build deploy

include .env

MIGRATION_DIR ?= ./migrations
DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
ENV_PATH := .env

run:
	set ENV_PATH=$(ENV_PATH) && go run ./cmd/api/main.go

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

test-watch:
	gotestsum -f testname --format-icons hivis --watch ./...
