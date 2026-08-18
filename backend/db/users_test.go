package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/Tanzor-Disco/skaM/models"
)

func TestCreateUser_Valid(t *testing.T) {
	user := models.User{
		Username:       "Pidoras",
		Email:          "Pidoras@mail.ru",
		PasswordHash:   "1231231414",
		LastYearActive: 2026,
	}
	err := db.CreateUser(context.Background(), user)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	assertError(t, err, false)

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
	user := models.User{
		Email:          "whatever",
		Username:       "1212141414141423131313212312331",
		PasswordHash:   "213131321",
		LastYearActive: 2026,
	}
	err := db.CreateUser(context.Background(), user)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	assertError(t, err, true)
}

func TestCreateUser_Duplicate(t *testing.T) {
	user := models.User{
		Email:          "whatever",
		Username:       "121214",
		PasswordHash:   "213131321",
		LastYearActive: 2026,
	}
	err := db.CreateUser(context.Background(), user)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	assertError(t, err, false)
	err = db.CreateUser(context.Background(), user)
	assertError(t, err, true)
}

func ValidatePendingUser(t *testing.T, got, wanted models.PendingUser) error {
	if wanted.Email == got.Email &&
		wanted.Username == got.Username &&
		wanted.PasswordHash == got.PasswordHash &&
		wanted.TokenHash == got.TokenHash {
		return nil
	}
	return fmt.Errorf("some of the fields didn't match\n got %+v\n wanted %+v", got, wanted)
}
