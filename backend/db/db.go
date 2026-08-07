package db

import (
	"context"
	"log"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgerrcode"
	"errors"
)


type User struct {
	Username string
	Email string
	PasswordHash string
	LastYearActive int
}

type Room struct {
	Name string
}

type Message struct {
	UserId int
	RoomId int
	Text string
}

type RoomUser struct {
	RoomId int
	UserId int
}

type DB struct {
	pool *pgxpool.Pool
}

func Connect(URI string) (*DB,error) {
	pool,err := pgxpool.New(context.Background(),URI)
	if err != nil {
		return &DB{},err
	}
	return &DB {
		pool:pool,
	}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) CreateUser(user User) error {
	_,err := db.pool.Exec(context.Background(),
	`
	INSERT INTO users (email,username,password_hash,last_year_active)
	VALUES ($1,$2,$3,$4)
	`,
	user.Email,user.Username,user.PasswordHash,user.LastYearActive,
	)

	var PgErr *pgconn.PgError
	if errors.As(err, &PgErr) && PgErr.Code == pgerrcode.UniqueViolation {
		err = apperrors.ErrEmailTaken
	}
	return err
}

func (db *DB) CreateRoom(room Room) error {
	_,err := db.pool.Exec(context.Background(),
	`
	INSERT INTO rooms (room_name)
	VALUES ($1)
	`,
	room.Name,
	)
	return err
}

func (db *DB) CreateMessage(message Message) error {
	_,err := db.pool.Exec(context.Background(),
	`
	INSERT INTO messages (user_id,room_id,message_text)
	VALUES ($1,$2,$3)
	`,
	message.UserId,message.RoomId,message.Text,
	)
	return err
}

func (db *DB) AddUserToRoom(roomUser RoomUser) error {
	_,err := db.pool.Exec(context.Background(),
	`
	INSERT INTO messages (user_id,room_id)
	VALUES ($1,$2)
	`,
	roomUser.UserId,roomUser.RoomId,
	)
	return err
}

func (db *DB) GetUsers() ([]User,error) {
	rows,err := db.pool.Query(context.Background(),`
	SELECT * FROM users
	`)
	if err != nil {
		return make([]User,0),err
	}
	defer rows.Close()
	var users []User
	var id int 
	for rows.Next() {
		var user User
		 if err := rows.Scan(&id,&user.Email,&user.Username,&user.PasswordHash,&user.LastYearActive);err != nil {
			 log.Println(err)
			 continue
		 }
		 users = append(users,user)
	}
	return users,nil
}

