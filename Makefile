.PHONY: dev stop restart logs

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
	go run .

stop:
	docker stop e-commerce-db
	docker rm e-commerce-db
	docker volume rm e-commerce_pgdata

restart: stop dev

logs:
	docker logs -f e-commerce-db

open-db:
	docker exec -it e-commerce-db psql -U snapy -d ecommerce