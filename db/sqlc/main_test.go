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

func TestMain(m *testing.M) {
	fmt.Println("Connecting to db...")

	conn := must(sqlc.ConnectToDB()).(*sql.DB)

	testQueries = sqlc.New(conn)

	os.Exit(m.Run())
}

func must(value any, err error) any {
	if err != nil {
		log.Fatalln(err)
	}

	return value
}
