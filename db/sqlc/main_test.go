package sqlc_test

import (
	"backend-masterclass/db/sqlc"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
)

var testQueries *sqlc.Queries
var testDB = must(sqlc.ConnectToDB()).(*sql.DB)

func TestMain(m *testing.M) {
	fmt.Println("Connecting to db...")

	testQueries = sqlc.New(testDB)

	os.Exit(m.Run())
}

func must(value any, err error) any {
	if err != nil {
		log.Fatalln(err)
	}

	return value
}
