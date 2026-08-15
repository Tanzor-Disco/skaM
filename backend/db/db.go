package db

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)


type Room struct {
	Name string
}

type Message struct {
	UserId int
	RoomId int
	Text   string
}

type RoomUser struct {
	RoomId int
	UserId int
}

type DB struct {
	pool *pgxpool.Pool
}

func Connect(URI string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), URI)
	if err != nil {
		return &DB{}, err
	}
	return &DB{
		pool: pool,
	}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

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
	`, user.Email, user.Username, user.PasswordHash, user.TokenHash, user.ExpiresAt,)
	return err
}

func (db *DB) CreateRoom(ctx context.Context, room Room) error {
	_, err := db.pool.Exec(context.Background(),
		`
	INSERT INTO rooms (room_name)
	VALUES ($1)
	`,
		room.Name,
	)
	return err
}

func (db *DB) CreateMessage(ctx context.Context, message Message) error {
	_, err := db.pool.Exec(context.Background(),
		`
	INSERT INTO messages (user_id,room_id,message_text)
	VALUES ($1,$2,$3)
	`,
		message.UserId, message.RoomId, message.Text,
	)
	return err
}

func (db *DB) AddUserToRoom(ctx context.Context, roomUser RoomUser) error {
	_, err := db.pool.Exec(context.Background(),
		`
	INSERT INTO messages (user_id,room_id)
	VALUES ($1,$2)
	`,
		roomUser.UserId, roomUser.RoomId,
	)
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
	var id int
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&id, &user.Email, &user.Username, &user.PasswordHash, &user.LastYearActive); err != nil {
			log.Println(err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func (db *DB) GetUserByEmail(email string) (models.User,error){
	var id int
	var user models.User
	err := db.pool.QueryRow(context.Background(),
		`
		SELECT * FROM users WHERE email = $1
		`,
		email).Scan(&id,&user.Email,&user.Username,&user.PasswordHash,&user.LastYearActive)
	if errors.Is(err,pgx.ErrNoRows) {
		err = apperrors.ErrUserNotFound
	}
	return user,err
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
	_,err := db.pool.Exec(context.Background(),
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
			log.Printf("failed to delete expired users from pending: %v",err)
		}
		time.Sleep(1 * time.Hour)
	}
}
