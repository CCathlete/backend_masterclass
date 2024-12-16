package main

import (
	"backend-masterclass/db/sqlc"
	"backend-masterclass/util"
	"database/sql"
)

func main() {
	conn := util.Must(sqlc.ConnectToDB()).(*sql.DB)
	defer conn.Close()
}
