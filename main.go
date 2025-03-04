package main

import (
	"database/sql"
	"log"

	"github.com/badermezzi/KubeGoBank/api"
	db "github.com/badermezzi/KubeGoBank/db/sqlc"
	"github.com/badermezzi/KubeGoBank/util"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err) 
	}

	connection, err := sql.Open(config.DBDriver, config.DBSource) // Open database connection
	if err != nil {
		log.Fatal("cannot connect to db:", err) // Exit if connection fails
	}

	store := db.NewStore(connection)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
