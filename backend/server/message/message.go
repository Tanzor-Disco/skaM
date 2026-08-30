package message

import (
	"github.com/Tanzor-Disco/skaM/db"
)

type MessageHandler struct {
	db *db.DB
}

func NewMessageHandler(db *db.DB) MessageHandler {
	return MessageHandler{
		db: db,
	}
}
