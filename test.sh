#!/bin/bash
sudo docker exec -it postgresTest createdb --username=root --owner=root test_simple_bank
export TEST_DATABASE_URL="postgresql://root:test@localhost:5433/test_simple_bank?sslmode=disable"
migrate -path db/migration -database "$TEST_DATABASE_URL" -verbose up 
go test ./... -count=1 -v -cover
sleep 3
sudo docker exec -it postgresTest dropdb test_simple_bank