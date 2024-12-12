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
	conn := must(sqlc.ConnectToDB()).(*sql.DB)

	testQueries = sqlc.New(conn)

	fmt.Println("This is the main test.")

	os.Exit(m.Run())
}

func must(value any, err error) any {
	if err != nil {
		log.Fatalln(err)
	}

	return value
}
