package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

const dbDriver = "postgres"

var dbSource = os.Getenv("TEST_DATABASE_URL")

var testQueries *Queries
var dbConnection *sql.DB

func TestMain(m *testing.M) {
	var err error
	dbConnection, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cant connect to the database ", err)
	}

	testQueries = New(dbConnection)
	os.Exit(m.Run())
}
