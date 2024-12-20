package api

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"backend-masterclass/db/sqlc"
	u "backend-masterclass/util"

	"github.com/gin-gonic/gin"
)

var testQueries *sqlc.Queries
var cfg = must(u.LoadConfig("../..")).(u.Config)
var testDB = must(sqlc.ConnectToDB(cfg)).(*sql.DB)

func TestMain(m *testing.M) {
	log.Println("Connecting to db...")

	testQueries = sqlc.New(testDB)

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())

}

func must(value any, err error) any {
	if err != nil {
		log.Fatalln(err)
	}

	return value
}
