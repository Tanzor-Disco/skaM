package db

import (
	"github.com/Tanzor-Disco/skaM/models"

	"context"
	"log"
	"time"
)

func (db *DB) CreatePendingUser(ctx context.Context, user models.PendingUser) error {
	_, err := db.pool.Exec(ctx,
		`
	INSERT INTO pending_users (email,username,password_hash,token_hash,expires_at)
	VALUES ($1,$2,$3,$4,$5)
	ON CONFLICT(email)
	DO UPDATE SET
	username = EXCLUDED.username,
	password_hash = EXCLUDED.password_hash,
	token_hash = EXCLUDED.token_hash,
	expires_at = EXCLUDED.expires_at
	`, user.Email, user.Username, user.PasswordHash, user.TokenHash, user.ExpiresAt)
	return err
}

func (db *DB) GetPendingUserByTokenHash(ctx context.Context, tokenHash string) (models.PendingUser, error) {
	var user models.PendingUser
	err := db.pool.QueryRow(ctx,
		`
	SELECT * FROM pending_users WHERE token_hash=$1
	`,
		tokenHash).Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash, &user.TokenHash, &user.ExpiresAt)
	return user, err
}

func (db *DB) deleteExpired() error {
	_, err := db.pool.Exec(context.Background(),
		`
	Delete from pending_users WHERE expires_at < CURRENT_TIMESTAMP
	`,
	)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) PeriodicPendingDelete() {
	for {
		err := db.deleteExpired()
		if err != nil {
			log.Printf("failed to delete expired users from pending: %v", err)
		}
		time.Sleep(1 * time.Hour)
	}
}
