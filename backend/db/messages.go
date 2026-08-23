package db

import (
	"context"
)

type Message struct {
	Id     int
	UserId int64
	RoomId int64
	Text   string
}

func NewMessage(userID, roomID int64, text string) Message {
	return Message{
		UserId: userID,
		RoomId: roomID,
		Text:   text,
	}
}

type MessageData struct {
	ID     int64
	UserID int64
	RoomID int64
	Text   string
	Author string
}

func NewMessageData(ID, userID, roomID int64, text string, author string) MessageData {
	return MessageData{
		ID:     ID,
		UserID: userID,
		RoomID: roomID,
		Text:   text,
		Author: author,
	}
}

func (db *DB) CreateMessage(ctx context.Context, message Message) (int64, error) {
	var messageID int64
	err := db.pool.QueryRow(ctx,
		`
	INSERT INTO messages (user_id,room_id,message_text)
	VALUES ($1,$2,$3)
	RETURNING id
	`,
		message.UserId, message.RoomId, message.Text,
	).Scan(&messageID)
	return messageID, err
}

func (db *DB) GetMessageData(ctx context.Context, roomID int) ([]MessageData, error) {
	rows, err := db.pool.Query(ctx,
		`
	SELECT messages.*, users.username
	FROM messages JOIN users ON messages.user_id = users.id
	WHERE messages.room_id = $1
	ORDER BY messages.id
	`, roomID)
	if err != nil {
		return make([]MessageData, 0), err
	}
	var messages []MessageData
	for rows.Next() {
		var message MessageData
		err := rows.Scan(&message.ID, &message.UserID, &message.RoomID, &message.Text, &message.Author)
		if err != nil {
			return make([]MessageData, 0), err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
