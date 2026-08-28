package models

import (
	"time"
)

type UserSessionID int64

// UserSession describes a row in user_sessions table
type UserSession struct {
	ID            UserSessionID
	UserID        UserID
	SessionString string
	ExpiresAt     time.Time
}

// NewUserSession creates a new instance of UserSession that expires in 7 days
func NewUserSession(userID UserID, sessionString string) UserSession {
	return UserSession{
		ID:            0,
		UserID:        userID,
		SessionString: sessionString,
		ExpiresAt:     time.Now().Add(7 * 24 * time.Hour),
	}
}
