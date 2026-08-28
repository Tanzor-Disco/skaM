package models

import (
	"time"
)

type UserID int64

// User represents a record in the users table
// It is used to create a new record in users table or retrieve data from it
type User struct {
	Id             UserID
	Email          string
	Username       string
	PasswordHash   string
	LastYearActive int
}

// NewUser is used to create a new User instance
// Id and LastYearActive aren't required
func NewUser(email, username, passwordHash string) User {
	return User{
		Email:          email,
		Username:       username,
		PasswordHash:   passwordHash,
		LastYearActive: time.Now().Year(),
	}
}
