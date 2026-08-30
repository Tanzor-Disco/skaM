package room

import (
	"github.com/Tanzor-Disco/skaM/db"
)

type RoomHandler struct {
	db      *db.DB
	baseURL string
}

func NewRoomHandler(db *db.DB, baseURL string) RoomHandler {
	return RoomHandler{
		db:      db,
		baseURL: baseURL,
	}
}
