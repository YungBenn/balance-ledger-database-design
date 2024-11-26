run:
	go run cmd/server/main.go

docker.up:
	docker-compose up -d

docker.down:
	docker-compose down

docker.build:
	@docker build -f build/server/Dockerfile -t rubenadi/balance-ledger:latest .

docker.run:
	@docker run -d \
		-e PORT="8002" \
		-e DB_HOST="localhost" \
		-e DB_PORT="5433" \
		-e DB_USER="postgres" \
		-e DB_PASSWORD="secret" \
		-e DB_NAME="balance_ledger" \
		-e DB_SSL_MODE="disable" \
		-p 8002:8002 \
		rubenadi/balance-ledger:latest

docker.stop:
	@docker stop rubenadi/balance-ledger:latest

sqlc:
	docker run --rm -v /home/rubenadi/tuts/balance-ledger-database-design:/src -w /src kjconroy/sqlc generate

migrate.up:
	migrate -path db/migrations -database "postgresql://postgres:secret@localhost:5433/balance_ledger?sslmode=disable" up

migrate.down:
	migrate -path db/migrations -database "postgresql://postgres:secret@localhost:5433/balance_ledger?sslmode=disable" down