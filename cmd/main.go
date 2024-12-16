package main

import (
	"backend-masterclass/api"
	"backend-masterclass/db/sqlc"
	"backend-masterclass/util"
	"database/sql"
)

const (
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn := util.Must(sqlc.ConnectToDB()).(*sql.DB)
	defer conn.Close()

	store := sqlc.NewStore(conn)
	server := api.NewServer(store)

	util.Must(nil, server.Start(serverAddress))
}
