.PHONY: up down migrate postgres recreate-db build logs install-migrate create-migration migrate-up migrate-down migrate-force migrate-version

DC=docker compose
DB_USER=authentication_user
DB_PASS=authentication_password
DB_NAME=authentication_db
DB_HOST=localhost
DB_PORT=5432
MIGRATIONS_DIR=migrations
PROJECT_NAME=authentication-system
NETWORK=$(PROJECT_NAME)_authentication-network
MIGRATE=migrate
MIGRATIONS_URL=postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

up:
	$(DC) up

down:
	$(DC) down -v

build:
	$(DC) build

rebuild:
	$(DC) up --build

logs:
	$(DC) logs -f

postgres:
	$(DC) up -d postgres
	@until docker exec $$(docker ps -q -f name=postgres) pg_isready -U $(DB_USER) -d $(DB_NAME); do \
		echo "Waiting for postgres..."; \
		sleep 1; \
	done

recreate-db: postgres
	docker exec $$(docker ps -q -f name=postgres) dropdb -U $(DB_USER) --if-exists $(DB_NAME)
	docker exec $$(docker ps -q -f name=postgres) createdb -U $(DB_USER) $(DB_NAME)

migrate: postgres
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(MIGRATIONS_URL)" up

seed:
	@echo "Seeding database..."
	docker exec -i $$(docker ps -q -f name=postgres) psql -U $(DB_USER) -d $(DB_NAME) < scripts/seed.sql
	@echo "Database seeded successfully!"

restart-app:
	$(DC) restart app

app-logs:
	$(DC) logs -f app

db-logs:
	$(DC) logs -f postgres

start: migrate seed build up

reset-db: recreate-db migrate seed

install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

create-migration:
	@read -p "Enter migration name: " name; \
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(MIGRATIONS_URL)" up

migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(MIGRATIONS_URL)" down

migrate-force:
	@read -p "Enter version: " version; \
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(MIGRATIONS_URL)" force $$version

migrate-version:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(MIGRATIONS_URL)" version

.DEFAULT_GOAL := start
