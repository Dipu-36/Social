include .env
MIGRATIONS_PATH = ./cmd/migrate/migrations
CONTAINER_NAME = postgres-db
VERSION ?= 6
.PHONY: migration
migration:
	@migrate create -seq -ext sql -dir "$(MIGRATIONS_PATH)" $(filter-out $@, $(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path "$(MIGRATIONS_PATH)" -database "$(DB_ADDR)" up  

.PHONY: migrate-down
migrate-down:
	@migrate -path "$(MIGRATIONS_PATH)" -database "$(DB_ADDR)" down $(filter-out $@, $(MAKECMDGOALS))

.PHONY: start-container
start-container:
	docker start "$(CONTAINER_NAME)"

.PHONY: stop-container
stop-container:
	docker stop "$(CONTAINER_NAME)"

.PHONY: force-container
force-container:
	@migrate -path "$(MIGRATIONS_PATH)" -database "$(DB_ADDR)" force $(VERSION)

