package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/mahesh-singh/simplebank/api"
	db "github.com/mahesh-singh/simplebank/db/sqlc"
)

var (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8181"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Failed to connect to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Failed to start the server", err)
	}
}
