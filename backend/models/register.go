package models

import (
	"time"
)

// RegisterRequest represents data passed from user to server after registration form submit
type RegisterRequest struct {
	Email    string
	Username string
	Password string
}

// User represents a record in the users table
// It is used to create a new record in users table or retrieve data from it
type User struct {
	Id             int64
	Email          string
	Username       string
	PasswordHash   string
	LastYearActive int
}

// NewUser is used to create a new User instance
// Id and LastYearActive aren't required
func NewUser(email, username, passwordHash string) User {
	return User{
		Id:             0,
		Email:          email,
		Username:       username,
		PasswordHash:   passwordHash,
		LastYearActive: time.Now().Year(),
	}
}

// PendingUser represents a record in pending_users table
// It is used to create a new record in pending_users table or retrieve data from it
type PendingUser struct {
	Id           int64
	Email        string
	Username     string
	PasswordHash string
	TokenHash    string
	ExpiresAt    time.Time
}
