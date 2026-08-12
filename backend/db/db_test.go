package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Tanzor-Disco/skaM/internal/testutils"

	"github.com/lpernett/godotenv"
)

var db *DB
var tdb *testutils.TestDB

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env_test")
	if err != nil {
		log.Fatal(err)
	}
	URI := os.Getenv("URI")
	db, err = Connect(URI)
	if err != nil {
		log.Fatal(err)
	}
	tdb, err = testutils.Connect(URI)
	if err != nil {
		log.Fatal(err)
	}
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
	user := User{
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
	user := User{
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
	user := User{
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
