postgres:
	docker run --name postgres14 --network bank-network -p 5432:5432  -e POSTGRES_USER=root  -e POSTGRES_PASSWORD=159159 -d postgres:14.16-alpine3.20

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank

killdb:
	docker kill postgres14

dropdb:
	docker exec -it postgres14 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:159159@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:159159@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:159159@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:159159@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -build_flags=--mod=mod -destination db/mock/store.go -package mockdb github.com/badermezzi/KubeGoBank/db/sqlc Store

.PHONY: createdb dropdb postgres migrateup migratedown sqlc test server mock migrateup1 migratedown1