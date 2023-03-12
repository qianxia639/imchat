package db

import (
	"IMChat/utils/config"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	conf, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatal("load config err: ", err)
	}

	testDB, err = sql.Open(conf.Postgres.Driver, conf.Postgres.Source)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
