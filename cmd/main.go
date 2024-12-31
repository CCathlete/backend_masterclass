package main

import (
	"backend-masterclass/api"
	"backend-masterclass/db/sqlc"
	"backend-masterclass/gapi"
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

	// --------Just to see that it's quick to replace----------------
	// jwtTokenMaker := u.Must(token.NewJWTMaker(
	// 	cfg.TokenKey)).(token.Maker)
	// ------------------------------------------------------------

	// ------------------HTTP (Gin) server--------------------------------
	// RunGinServer(store, cfg, pasetoTokenMaker)
	// -------------------------------------------------------------------

	// ------------------gRPC server--------------------------------------
	RunGRPCServer(store, cfg, pasetoTokenMaker)
	// -------------------------------------------------------------------

}

func RunGinServer(store sqlc.Store, cfg u.Config, tokenMaker token.Maker) {
	server := api.NewServer(store, cfg, tokenMaker)
	u.Must(nil, server.Start(cfg.HTTPServerAddress))
}

func RunGRPCServer(store sqlc.Store, cfg u.Config, tokenMaker token.Maker) {
	server := gapi.NewServer(store, cfg, tokenMaker)
	u.Must(nil, server.Start(cfg.GRPCServerAddress))
}
