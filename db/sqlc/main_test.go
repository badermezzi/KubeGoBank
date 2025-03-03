package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" // PostgreSQL driver
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:159159@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries // Package-level variable to hold sqlc Queries
var testDB *sql.DB       // Package-level variable to hold database connection
var testStore Store      // Package-level variable to hold the Store

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource) // Open database connection
	if err != nil {
		log.Fatal("cannot connect to db:", err) // Exit if connection fails
	}

	testQueries = New(testDB)    // Create sqlc Queries instance with the connection
	testStore = NewStore(testDB) // Create a new Store instance

	os.Exit(m.Run()) // Run all tests in the package and exit
}
