package db

import (
	"context"
)

type RoomUser struct {
	RoomId int
	UserId int
}

func (db *DB) AddUserToRoom(ctx context.Context, roomUser RoomUser) error {
	_, err := db.pool.Exec(context.Background(),
		`
	INSERT INTO room_users (user_id,room_id)
	VALUES ($1,$2)
	`,
		roomUser.UserId, roomUser.RoomId,
	)
	return err
}
