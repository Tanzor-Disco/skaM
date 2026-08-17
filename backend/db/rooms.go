package db

import (
	"context"
)

type Room struct {
	Name string
}

func NewRoom(name string) Room {
	return Room {
		Name:name,
	}
}

func (db *DB) CreateRoom(ctx context.Context, room Room) (int, error) {
	var roomID int
	err := db.pool.QueryRow(context.Background(),
		`
	INSERT INTO rooms (room_name)
	VALUES ($1)
	RETURNING id
	`,
	room.Name,
	).Scan(&roomID)
	return roomID, err
}
