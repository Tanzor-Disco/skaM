package db

import (
	"log"
	"os"
	"testing"

	"github.com/Tanzor-Disco/skaM/db/migrations"
	"github.com/Tanzor-Disco/skaM/internal/testutils"

	"github.com/lpernett/godotenv"
)

var db *DB
var tdb *testutils.TestDB

func setUpDB() (*DB, *testutils.TestDB) {
	//load the environment
	err := godotenv.Load("../.env_test")
	if err != nil {
		log.Fatal(err)
	}
	testURI := os.Getenv("TEST_URI")

	//reset the test db and update to the latest migration
	migrations.Reset(testURI)

	//connect the main db
	db, err = Connect(testURI)
	if err != nil {
		log.Fatal(err)
	}

	//connect the testdb
	tdb, err = testutils.Connect(testURI)
	if err != nil {
		log.Fatal(err)
	}

	return db, tdb
}

func TestMain(m *testing.M) {
	db, tdb := setUpDB()

	//run the tests
	code := m.Run()
	db.pool.Close()
	tdb.Close()
	os.Exit(code)
}

func assertError(t *testing.T, gotError error, wantError bool) {
	t.Helper()
	if (gotError != nil) != wantError {
		t.Fatalf("wanted the error to not be nil:%v, got %v", wantError, gotError)
	}
}
