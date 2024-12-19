package sqlc

import (
	u "backend-masterclass/util"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var cfg = must(u.LoadConfig("../..")).(u.Config)
var testDB = must(ConnectToDB(cfg)).(*sql.DB)

func TestMain(m *testing.M) {
	fmt.Println("Connecting to db...")

	testQueries = New(testDB)

	os.Exit(m.Run())
}

func must(value any, err error) any {
	if err != nil {
		log.Fatalln(err)
	}

	return value
}
