package models

type MessageID int64

// Message is used to describe a row in messages table
// it can be used both for retrieving and registering messages
type Message struct {
	Id     MessageID
	UserId UserID
	RoomId RoomID
	Text   string
}

func NewMessage(userID UserID, roomID RoomID, text string) Message {
	return Message{
		UserId: userID,
		RoomId: roomID,
		Text:   text,
	}
}
