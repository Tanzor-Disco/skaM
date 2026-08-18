package db

import (
	"context"
	"github.com/Tanzor-Disco/skaM/models"
	"testing"
	"time"
)

func TestCreatePendingUser_Valid(t *testing.T) {
	user := models.PendingUser{
		Email:        "whatever",
		Username:     "121214",
		PasswordHash: "213131321",
		TokenHash:    "23131",
		ExpiresAt:    time.Now(),
	}
	err := db.CreatePendingUser(context.Background(), user)
	defer tdb.DeletePendingUsersByEmail(t, user.Email)
	if err != nil {
		t.Fatal(err)
	}
	gotUser := tdb.GetPendingUserByEmail(t, user.Email)
	err = ValidatePendingUser(t, gotUser, user)
	if err != nil {
		t.Fatalf("err while validating user:\n got %+v\n wanted %+v", gotUser, user)
	}
}
