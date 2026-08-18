package testutils

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func GetPasswordHash(t *testing.T, password string) string {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GetPasswordHash: %v", err)
	}
	return string(passwordHash)
}
