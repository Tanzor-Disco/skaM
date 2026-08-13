package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Tanzor-Disco/skaM/internal/testutils"
	"github.com/Tanzor-Disco/skaM/models"
	"github.com/Tanzor-Disco/skaM/db/migrations"

	"github.com/lpernett/godotenv"
)

var db *DB
var tdb *testutils.TestDB

func TestMain(m *testing.M) {
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
	
	//run the tests
	code := m.Run()
	db.pool.Close()
	tdb.Close()
	os.Exit(code)
}

func AssertError(t *testing.T, gotError error, wantError bool) {
	t.Helper()
	if (gotError != nil) != wantError {
		t.Fatalf("wanted the error to not be nil:%v, got %v", wantError, gotError)
	}
}

func TestCreateUser_Valid(t *testing.T) {
	user := models.User {
		Username:       "Pidoras",
		Email:          "Pidoras@mail.ru",
		PasswordHash:   "1231231414",
		LastYearActive: 2026,
	}
	err := db.CreateUser(context.Background(), user)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	AssertError(t, err, false)

	users := tdb.GetUsersByEmail(t, user.Email)
	if len(users) != 1 {
		t.Fatalf("The wrong number of users detected: %d", len(users))
	}

	dUser := users[0]
	if dUser.Username != user.Username || dUser.PasswordHash != user.PasswordHash || dUser.LastYearActive != user.LastYearActive {
		tdb.DeleteUsersByEmail(t, user.Email)
		t.Fatalf("User mismatch\n Tried to register %+v\n Got %+v", user, dUser)
	}
}

func TestCreateUser_Invalid(t *testing.T) {
	user := models.User {
		Email:          "whatever",
		Username:       "1212141414141423131313212312331",
		PasswordHash:   "213131321",
		LastYearActive: 2026,
	}
	err := db.CreateUser(context.Background(), user)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	AssertError(t, err, true)
}

func TestCreateUser_Duplicate(t *testing.T) {
	user := models.User {
		Email:          "whatever",
		Username:       "121214",
		PasswordHash:   "213131321",
		LastYearActive: 2026,
	}
	err := db.CreateUser(context.Background(), user)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	AssertError(t, err, false)
	err = db.CreateUser(context.Background(), user)
	AssertError(t, err, true)
}

func ValidatePendingUser(t *testing.T,got,wanted models.PendingUser) error {
	if wanted.Email == got.Email && 
	wanted.Username == got.Username && 
	wanted.PasswordHash == got.PasswordHash &&
	wanted.TokenHash == got.TokenHash {
		return nil
	}
	return fmt.Errorf("some of the fields didn't match\n got %+v\n wanted %+v",got,wanted)
}

func TestCreatePendingUser_Valid(t *testing.T) {
	user := models.PendingUser {
		Email:          "whatever",
		Username:       "121214",
		PasswordHash:   "213131321",
		TokenHash: "23131",
		ExpiresAt: time.Now(),
	}
	err := db.CreatePendingUser(context.Background(),user)
	defer tdb.DeletePendingUsersByEmail(t,user.Email)
	if err != nil {
		t.Fatal(err)
	}
	gotUser := tdb.GetPendingUserByEmail(t,user.Email)
	err = ValidatePendingUser(t,gotUser,user)
	if err != nil {
		t.Fatalf("err while validating user:\n got %+v\n wanted %+v",gotUser,user)
	}
}

