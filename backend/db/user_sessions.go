package db

import (
	"context"
	"fmt"
	"time"
)

func (db *DB) CreateSession(ctx context.Context, userID int, sessionID string) error {
	expires := time.Now().Add(7 * 24 * time.Hour)
	_, err := db.pool.Exec(ctx,
		`
	INSERT INTO user_sessions (user_id, session_id, expires_at)
	VALUES ($1,$2,$3)
	`,
		userID, sessionID, expires)
	if err != nil {
		return fmt.Errorf("CreateSession: %w", err)
	}
	return nil
}
