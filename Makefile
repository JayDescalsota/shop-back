
.PHONY: up up-all down logs pull auth dev

# Infra only (postgres + redis) — fast for local dev; run `make auth` separately
up:
	docker compose up -d --build

# Full stack including auth container build
up-all:
	docker compose --profile full up -d --build

setup: up migrate

migrate: migrate-auth

migrate-auth:
	docker compose exec -T postgres psql -U postgres -d postgres < services/auth/migrations/001_users.sql

down:
	docker compose down

logs:
	docker compose logs -f

pull:
	docker compose pull

# Hot-reload auth on the host (requires postgres + redis: `make up`)
auth:
	air -c services/auth/.air.toml

# Infra + auth on host
dev: up auth
