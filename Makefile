# Root commands only delegate: every app keeps its own Makefile so it can move to its own repo.
COMPOSE := docker compose --env-file .env -f deploy/docker-compose.yml

-include .env
TEST_DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(or $(POSTGRES_PORT),5432)/$(POSTGRES_DB)?sslmode=disable

.PHONY: dev down gen test

## dev: run postgres and the api (http://localhost:8080)
dev: .env
	$(COMPOSE) up --build

## down: stop everything (data stays in the docker volumes)
down:
	$(COMPOSE) down

## gen: regenerate code from api/openapi.yaml
gen:
	$(MAKE) -C apps/api gen

## test: run every app's tests; the api tests use the dev postgres
test: .env
	$(COMPOSE) up -d --wait postgres
	TEST_DATABASE_URL='$(TEST_DATABASE_URL)' $(MAKE) -C apps/api test
	@if command -v flutter >/dev/null; then $(MAKE) -C apps/mobile test; else echo "flutter not installed: skipping apps/mobile"; fi

.env:
	cp .env.example .env
	@echo "Created .env from .env.example; change the passwords."
