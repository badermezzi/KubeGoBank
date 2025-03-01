package main

import (
	"database/sql"
	"log"

	"github.com/badermezzi/KubeGoBank/api"
	db "github.com/badermezzi/KubeGoBank/db/sqlc"

	_ "github.com/lib/pq" // PostgreSQL driver
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:159159@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	connection, err := sql.Open(dbDriver, dbSource) // Open database connection
	if err != nil {
		log.Fatal("cannot connect to db:", err) // Exit if connection fails
	}

	store := db.NewStore(connection)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
