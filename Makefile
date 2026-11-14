run:
	go run cmd/server/main.go

docker.up:
	docker-compose up -d

docker.down:
	docker-compose down

docker.build:
	@docker build -f build/server/Dockerfile -t rubenadi/balance-ledger:latest .

docker.run:
	@docker run -d --name balance-ledger \
		--add-host=host.docker.internal:host-gateway \
		-e PORT="8002" \
		-e DB_HOST="host.docker.internal" \
		-e DB_PORT="5433" \
		-e DB_USER="postgres" \
		-e DB_PASSWORD="secret" \
		-e DB_NAME="balance_ledger" \
		-e DB_SSL_MODE="disable" \
		-e TOKEN_SYMMETRIC_KEY="change-me-32-byte-dev-key-000000" \
		-p 8002:8002 \
		rubenadi/balance-ledger:latest

docker.stop:
	@docker rm -f balance-ledger

sqlc:
	docker run --rm -v "$(CURDIR)":/src -w /src sqlc/sqlc:1.30.0 generate

migrate.up:
	migrate -path db/migrations -database "postgresql://postgres:secret@localhost:5433/balance_ledger?sslmode=disable" up

migrate.down:
	migrate -path db/migrations -database "postgresql://postgres:secret@localhost:5433/balance_ledger?sslmode=disable" down