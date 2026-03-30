.PHONY: dev stop restart logs migrate-up migrate-down migrate-status pgadmin run-backend

DB_URL=postgres://snapy:snapy@localhost:5433/ecommerce?sslmode=disable

run-backend:
	docker run -d \
		--name e-commerce-db \
		-e POSTGRES_USER=snapy \
		-e POSTGRES_PASSWORD=snapy \
		-e POSTGRES_DB=ecommerce \
		-p 5433:5432 \
		-v e-commerce_pgdata:/var/lib/postgresql/data \
		postgres:16-alpine || echo "DB container already running"
	@echo "Waiting for postgres to be ready..."
	@ping -n 6 127.0.0.1 >nul
	@echo "Running migrations..."
	goose -dir db/migrations postgres "$(DB_URL)" up
	@echo "Starting backend..."
	go run .

migrate-up:
	goose -dir db/migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir db/migrations postgres "$(DB_URL)" down

migrate-status:
	goose -dir db/migrations postgres "$(DB_URL)" status

stop:
	docker stop e-commerce-db
	docker rm e-commerce-db

restart: stop run-backend

logs:
	docker logs -f e-commerce-db

open-db:
	docker exec -it e-commerce-db psql -U snapy -d ecommerce
