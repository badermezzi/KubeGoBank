postgres:
	docker run --name postgres14 -p 5432:5432  -e POSTGRES_USER=root  -e POSTGRES_PASSWORD=159159 -d postgres:14.16-alpine3.20

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank

killdb:
	docker kill postgres14

dropdb:
	docker exec -it postgres14 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:159159@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:159159@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: createdb dropdb postgres migrateup migratedown sqlc test