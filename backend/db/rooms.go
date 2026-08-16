package db

import (
	"context"
)

type Room struct {
	Name string
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
