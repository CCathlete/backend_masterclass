package main

import (
	"backend-masterclass/api"
	"backend-masterclass/db/sqlc"
	u "backend-masterclass/util"
	"database/sql"
)

func main() {
	cfg := u.Must(u.LoadConfig("./")).(u.Config)
	conn := u.Must(sqlc.ConnectToDB(cfg)).(*sql.DB)
	defer conn.Close()

	store := sqlc.NewStore(conn)
	server := api.NewServer(store)

	u.Must(nil, server.Start(cfg.ServerAddress))
}
