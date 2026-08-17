package db

import (
	"context"
)

type RoomUser struct {
	RoomID int
	UserID int
}

func NewRoomUser(roomID, userID int) RoomUser {
	return RoomUser {
		RoomID: roomID,
		UserID: userID,
	}
}

func (db *DB) AddUserToRoom(ctx context.Context, roomUser RoomUser) error {
	_, err := db.pool.Exec(ctx,
		`
	INSERT INTO room_users (user_id,room_id)
	VALUES ($1,$2)
	`,
		roomUser.UserID, roomUser.RoomID,
	)
	return err
}
