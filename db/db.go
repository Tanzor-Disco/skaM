package db

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"context"
)

type User struct {
	Name string
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

func Connect(URI string) (DB,error) {
	pool,err := pgxpool.Connect(context.Background(),URI)
	if err != nil {
		return DB{},err
	}
	return DB {
		pool:pool,
	}, nil
}

func (db DB) CreateUser(user User) error {
	_,err := db.pool.Exec(context.Background(),
	`
	INSERT INTO users (username,email,password_hash,last_year_active)
	VALUES ($1,$2,$3,$4)
	`,
	user.Name,user.Email,user.PasswordHash,user.LastYearActive,
	)
	return err
}

func (db DB) CreateRoom(room Room) error {
	_,err := db.pool.Exec(context.Background(),
	`
	INSERT INTO rooms (room_name)
	VALUES ($1)
	`,
	room.Name,
	)
	return err
}

func (db DB) CreateMessage(message Message) error {
	_,err := db.pool.Exec(context.Background(),
	`
	INSERT INTO messages (user_id,room_id,message_text)
	VALUES ($1,$2,$3)
	`,
	message.UserId,message.RoomId,message.Text,
	)
	return err
}

func (db DB) AddUserToRoom(roomUser RoomUser) error {
	_,err := db.pool.Exec(context.Background(),
	`
	INSERT INTO messages (user_id,room_id)
	VALUES ($1,$2)
	`,
	roomUser.UserId,roomUser.RoomId,
	)
	return err
}


