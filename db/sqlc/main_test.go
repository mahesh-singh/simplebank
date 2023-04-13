package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/mahesh-singh/simplebank/util"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Failed to load config ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
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
