package models

import (
	"time"
)

type RegisterRequest struct {
	Email    string
	Username string
	Password string
}

type User struct {
	Id             int
	Email          string
	Username       string
	PasswordHash   string
	LastYearActive int
}

func NewUser(email, username, passwordHash string) User {
	return User{
		Id:             0,
		Email:          email,
		Username:       username,
		PasswordHash:   passwordHash,
		LastYearActive: time.Now().Year(),
	}
}

type PendingUser struct {
	Id           int
	Email        string
	Username     string
	PasswordHash string
	TokenHash    string
	ExpiresAt    time.Time
}
