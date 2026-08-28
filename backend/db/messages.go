package db

import (
	"context"
	"github.com/Tanzor-Disco/skaM/models"
)

// messageData is used to transport the message info to frontend
type MessageData struct {
	ID     models.MessageID
	UserID models.UserID
	RoomID models.RoomID
	Text   string
	Author string
}

func NewMessageData(ID models.MessageID, userID models.UserID, roomID models.RoomID, text string, author string) MessageData {
	return MessageData{
		ID:     ID,
		UserID: userID,
		RoomID: roomID,
		Text:   text,
		Author: author,
	}
}

func (db *DB) CreateMessage(ctx context.Context, message models.Message) (models.MessageID, error) {
	var messageID models.MessageID
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

func (db *DB) GetMessageData(ctx context.Context, roomID models.RoomID) ([]MessageData, error) {
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
