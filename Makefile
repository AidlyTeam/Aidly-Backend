.PHONY: help dev dev-build prod prod-build sqlc swagger

## 🧪 Start project in development mode
dev:
	@echo "[i] Project is starting in development mode...\n"
	@docker compose -f ./deployment/dev.docker-compose.yml up

## 🛠️ Build and start project in development mode
dev-build:
	@echo "[i] Building project in development mode...\n"
	@docker compose -f ./deployment/dev.docker-compose.yml up --build -d

## 🚀 Start project in production mode (detached)
prod:
	@echo "[i] Project is starting in production mode...\n"
	@docker compose -f ./deployment/prod.docker-compose.yml up -d

## 🏗️ Build and start project in production mode (detached)
prod-build:
	@echo "[i] Building project in production mode...\n"
	@docker compose -f ./deployment/prod.docker-compose.yml up --build -d

## 🧬 Generate SQLC code
sqlc:
	@echo "[i] Generating SQLC code...\n"
	@sqlc generate -f ./internal/repos/sqlc.yaml

## 📘 Generate Swagger documentation
swagger:
	@echo "[i] Generating Swagger documentation...\n"
	@swag init --parseVendor -d . -g ./cmd/aidly/main.go

uswagger:
	@echo "Generating swagger..."
	@~/go/bin/swag init  --parseVendor  -d . -g ./cmd/aidly/main.go

## 📋 Show available targets
help:
	@echo "[i] Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  dev          Start project in development mode"
	@echo "  dev-build    Build project in development mode"
	@echo "  prod         Start project in production mode"
	@echo "  prod-build   Build project in production mode"
	@echo "  sqlc         Generate SQLC code"
	@echo "  swagger      Generate Swagger documentation"
	@echo "  help         Show this help"
