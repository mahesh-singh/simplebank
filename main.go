package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/mahesh-singh/simplebank/api"
	db "github.com/mahesh-singh/simplebank/db/sqlc"
	"github.com/mahesh-singh/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Failed to connect to db", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Failed to start the server", err)
	}
}
