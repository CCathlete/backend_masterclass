package main

import (
	"backend-masterclass/api"
	"backend-masterclass/db/sqlc"
	"backend-masterclass/token"
	u "backend-masterclass/util"
	"database/sql"
	"fmt"
)

func main() {
	cfg := u.Must(u.LoadConfig("./")).(u.Config)
	fmt.Println(cfg)
	conn := u.Must(sqlc.ConnectToDB(cfg)).(*sql.DB)
	defer conn.Close()

	store := sqlc.NewStore(conn)
	pasetoTokenMaker := u.Must(token.NewPasetoMaker(
		cfg.TokenKey)).(token.Maker)
	server := api.NewServer(store, cfg, pasetoTokenMaker)

	// --------Just to see that it's quick to replace----------------
	// jwtTokenMaker := u.Must(token.NewJWTMaker(
	// 	cfg.TokenKey)).(token.Maker)
	// server := api.NewServer(store, cfg, jwtTokenMaker)
	// ------------------------------------------------------------
	u.Must(nil, server.Start(cfg.ServerAddress))
}
