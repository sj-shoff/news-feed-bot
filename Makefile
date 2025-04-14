include .env

build:
	docker compose build news-feed-bot

run:
	docker compose up news-feed-bot

down:
	docker compose down

migrate:
	migrate -path ./migrations/schema -database 'postgres://postgres:${POSTGRES_PASSWORD}@0.0.0.0:5432/postgres?sslmode=disable' up

logs:
	docker compose logs -f news-feed-bot