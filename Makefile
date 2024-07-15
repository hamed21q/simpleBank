postgres:
	docker run --name postgres14 -p 5420:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1qaz -d focker.ir/postgres:12-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres14 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:1qaz@localhost:5420/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:1qaz@localhost:5420/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...


.PHONY: postgres createdb dropdb migrateup migratedown sqlc test