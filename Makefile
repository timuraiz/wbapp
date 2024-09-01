ifneq (,$(wildcard .env))
    include .env
    export
endif
# Directory where migration files are stored
MIGRATIONS_DIR := db/migrations

# Migration tool commands
MIGRATE_CMD := migrate -path $(MIGRATIONS_DIR) -database $(DB_URL)

# Default target (help command)
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make create-migration name=<name>    Create new migration files"
	@echo "  make migrate-up                     Run all 'up' migrations"
	@echo "  make migrate-down                   Run all 'down' migrations"
	@echo "  make migrate-force                   Force set the migration version"

# Create a new migration file
.PHONY: create-migration
create-migration:
ifndef name
	$(error Migration name is not specified. Use 'make create-migration name=<name>')
endif
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

# Run all "up" migrations
.PHONY: migrate-up
migrate-up:
	$(MIGRATE_CMD) up

# Run all "down" migrations
.PHONY: migrate-down
migrate-down:
	$(MIGRATE_CMD) down

# Force set the migration version (useful for manual overrides)
.PHONY: migrate-force
migrate-force:
ifndef version
	$(error Migration version is not specified. Use 'make migrate-force version=<version>')
endif
	$(MIGRATE_CMD) force $(version)
