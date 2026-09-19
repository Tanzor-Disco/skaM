package websocket

import (
	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/models"
	"github.com/gorilla/websocket"
	"sync"
)

// WsHub represents a structure that handles users currently connected to websocket
// it also stores mutex to enable the concurrent mutation of the Users
type WsHub struct {
	Users map[models.UserID]*websocket.Conn
	Mutex sync.Mutex
}

func NewWsHub() WsHub {
	return WsHub{
		Users: make(map[models.UserID]*websocket.Conn, 0),
	}
}

func (h *WsHub) AddUser(userID models.UserID, conn *websocket.Conn) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	h.Users[userID] = conn
}

func (h *WsHub) DeleteUser(userID models.UserID) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	delete(h.Users, userID)
}

// SendMessageData sends a certain MessageData struct to all the users from a certain room
// that have an active websocket connection
func (h *WsHub) SendMessageData(roomUsers []models.RoomUser, message db.MessageData) {
	for _, roomUser := range roomUsers {
		conn, ok := h.Users[roomUser.UserID]
		if !ok {
			continue
		}
		conn.WriteJSON(message)
	}
}
