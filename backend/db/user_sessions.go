package db

import (
	"context"
	"fmt"
	"time"
)

type UserSession struct {
	ID int
	UserID int
	SessionString string
	ExpiresAt time.Time

}

func NewUserSession(userID int,sessionString string) UserSession {
	return UserSession {
		ID:0,
		UserID: userID,
		SessionString: sessionString,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
}

func (db *DB) CreateSession(ctx context.Context, user UserSession) error {
	_, err := db.pool.Exec(ctx,
		`
	INSERT INTO user_sessions (user_id, session_string, expires_at)
	VALUES ($1,$2,$3)
	`,
	user.UserID, user.SessionString, user.ExpiresAt)
	if err != nil {
		return fmt.Errorf("CreateSession: %w", err)
	}
	return nil
}

func (db *DB) GetUserSessionBySessionString(ctx context.Context, sessionString string) (UserSession,error) {
	var user UserSession
	err := db.pool.QueryRow(ctx, 
	`
	SELECT * FROM user_sessions
	WHERE session_string = $1
	`,
	sessionString,).Scan(&user.ID,&user.UserID,&user.SessionString,&user.ExpiresAt)
	if err != nil {
		return UserSession{}, fmt.Errorf("getUserSessionBySessionString: %w",err)
	}
	return user, nil
}
