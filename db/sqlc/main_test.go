package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn := must(ConnectToDB()).(*sql.DB)

	testQueries = New(conn)

	os.Exit(m.Run())
}

func must(value any, err error) any {
	if err != nil {
		log.Fatalln(err)
	}

	return value
}
