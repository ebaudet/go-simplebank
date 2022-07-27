postgres:
	docker run --name postgres14 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres14 dropdb simple_bank

new_migration:
	migrate create -ext sql -dir db/migration -seq $(filter-out $@,$(MAKECMDGOALS))

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

test_nocache:
	go clean -testcache
	make test

server:
	docker start postgres14
	go run main.go

mock:
	mockgen -destination db/mock/store.go -build_flags=--mod=mod -package mockdb github.com/ebaudet/simplebank/db/sqlc Store

docker_build:
	docker build -t simplebank:latest .

docker_run:
	docker run --name simplebank --network bank-network -p 8080:8080 -e DB_SOURCE=postgresql://root:secret@postgres14:5432/simple_bank?sslmode=disable simplebank:latest

.PHONY: postgres createdb dropdb new_migration migrateup migrateup1 migratedown migratedown1 sqlc test test_nocache server mock docker_build docker_run
