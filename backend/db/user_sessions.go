package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Tanzor-Disco/skaM/models"
)

func (db *DB) CreateSession(ctx context.Context, user models.UserSession) error {
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

func (db *DB) GetUserSessionBySessionString(ctx context.Context, sessionString string) (models.UserSession, error) {
	var user models.UserSession
	err := db.pool.QueryRow(ctx,
		`
	SELECT * FROM user_sessions
	WHERE session_string = $1
	`,
		sessionString).Scan(&user.ID, &user.UserID, &user.SessionString, &user.ExpiresAt)
	if err != nil {
		return models.UserSession{}, fmt.Errorf("getUserSessionBySessionString: %w", err)
	}
	return user, nil
}

func (db *DB) deleteAllExpiredSessions() error {
	_, err := db.pool.Exec(context.Background(),
		`
	DELETE FROM user_sessions WHERE expires_at < NOW()
	`)
	if err != nil {
		return fmt.Errorf("deleteAllExpiredSessions: %w", err)
	}
	return nil
}

func (db *DB) PeriodicDeleteAllExpiredSessions() {
	ticker := time.NewTicker(time.Hour * 24 * 7)
	for range ticker.C {
		err := db.deleteAllExpiredSessions()
		if err != nil {
			log.Printf("PeriodicDeleteAllExpiredSessions: %v", err)
		}
	}
}
