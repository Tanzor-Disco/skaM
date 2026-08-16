package db

import (
	"context"
)

type Message struct {
	UserId int
	RoomId int
	Text   string
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
