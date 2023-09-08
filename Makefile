DB_URL=postgresql://root:1qaz@localhost:5432/simple_bank?sslmode=disable

postgres:
	sudo docker run --name postgres  -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1qaz -d postgres
createDb:
	sudo docker exec -it postgres createdb --username=root --owner=root simple_bank
dropDb:
	sudo docker exec -it postgres dropdb simple_bank
migrateUp:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migrateUp1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1
migrateDown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
migrateDown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1
sqlc:
	slqc generate
test:
	export TEST_DATABASE_URL="$DB_URL"
	go test ./... -count=1 -v -cover
.PHONY: postgres createDb dropDb migrateUp migrateDown migrateUp1 migrateDown1 sqlc test