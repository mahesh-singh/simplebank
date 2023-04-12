package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Failed to connect to db", err)
	}
	if testDB == nil {
		log.Print("testDB is nil")
	} else {
		log.Print("test db is not nil")
	}
	err = testDB.Ping()
	if err != nil {
		log.Fatal("Failed to connect to db", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
