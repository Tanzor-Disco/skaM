package db

import (
	"context"
	"errors"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
)

func (db *DB) CreateUser(ctx context.Context, user models.User) error {
	_, err := db.pool.Exec(ctx,
		`
	INSERT INTO users (email,username,password_hash,last_year_active)
	VALUES ($1,$2,$3,$4)
	`,
		user.Email, user.Username, user.PasswordHash, user.LastYearActive,
	)

	var PgErr *pgconn.PgError
	if errors.As(err, &PgErr) && PgErr.Code == pgerrcode.UniqueViolation {
		err = apperrors.ErrEmailTaken
	}
	return err
}

func (db *DB) GetUsers() ([]models.User, error) {
	rows, err := db.pool.Query(context.Background(), `
	SELECT * FROM users
	`)
	if err != nil {
		return make([]models.User, 0), err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash, &user.LastYearActive); err != nil {
			log.Println(err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

// GetUserByEmail sends a DB query that returns an models.User instance
// In case the user with such email doesn't exist it returns apperrors.ErrUserNotFound
func (db *DB) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := db.pool.QueryRow(context.Background(),
		`
		SELECT * FROM users WHERE email = $1
		`,
		email).Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash, &user.LastYearActive)
	if errors.Is(err, pgx.ErrNoRows) {
		err = apperrors.ErrUserNotFound
	}
	return user, err
}

func (db *DB) GetUserByID(ctx context.Context, userID int64) (models.User, error) {
	var user models.User
	err := db.pool.QueryRow(ctx,
		`
		SELECT * FROM users WHERE id = $1
		`,
		userID).Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash, &user.LastYearActive)
	if errors.Is(err, pgx.ErrNoRows) {
		err = apperrors.ErrUserNotFound
	}
	return user, err

}
