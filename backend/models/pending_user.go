package models

import (
	"time"
)

type PendingUserID int64

// PendingUser represents a record in pending_users table
// It is used to create a new record in pending_users table or retrieve data from it
type PendingUser struct {
	Id           PendingUserID
	Email        string
	Username     string
	PasswordHash string
	TokenHash    string
	ExpiresAt    time.Time
}
