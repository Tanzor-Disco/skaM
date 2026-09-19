package message

import (
	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/server/websocket"
)

// MessageHandler represents a list of dependencies that are passed from the server package to message package
type MessageHandler struct {
	db  *db.DB
	hub *websocket.WsHub
}

func NewMessageHandler(db *db.DB, hub *websocket.WsHub) MessageHandler {
	return MessageHandler{
		db:  db,
		hub: hub,
	}
}
