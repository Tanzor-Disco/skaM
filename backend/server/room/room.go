package room

import (
	"github.com/Tanzor-Disco/skaM/db"
)

// RegisterHandler represents a list of dependencies that are passed from the server package to message package
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
