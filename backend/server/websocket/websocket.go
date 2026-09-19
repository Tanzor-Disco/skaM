package websocket

import (
	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/internal/utils"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// RootHandler represents a list of dependencies that are passed from the server package to message package
type WsHandler struct {
	db      *db.DB
	hub     *WsHub
	baseURL string
}

func NewWsHandler(db *db.DB, hub *WsHub, baseURL string) WsHandler {
	return WsHandler{
		db:      db,
		hub:     hub,
		baseURL: baseURL,
	}
}

// HandleWs handles http requests sent to api/ws
// it checks the user session, sets up the upgrader
// upgrades connection to websocket, adds user connection to hub.Users
// waits for the connection to close, then deletes the user connection from hub.Users
func (h *WsHandler) HandleWs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionCookie, err := r.Cookie("session_string")
	if err != nil {
		log.Printf("HandleWs: r.Cookie: err")
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidSessionString, nil)
		utils.SendJSON(w, http.StatusUnauthorized, body)
		return
	}
	userSession, err := h.db.GetUserSessionBySessionString(ctx, sessionCookie.Value)
	if err != nil {
		log.Printf("HandleWs: GetUserSessionBySessionString: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidSessionString, nil)
		utils.SendJSON(w, http.StatusUnauthorized, body)
		return
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			allowed := map[string]bool{
				"http://localhost:5173": true,
				"http://localhost:8080":true,
				h.baseURL:               true,
			}
			baseURL := r.Header.Get("Origin")
			_, ok := allowed[baseURL]
			return ok
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("HandleWs: upgrader.Upgrade: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrWebsocketUpgrade, nil)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}

	h.hub.AddUser(userSession.UserID, conn)

	_, _, err = conn.ReadMessage()
	if err != nil {
		log.Printf("HandleWs: conn.ReadMessage: %v", err)
		h.hub.DeleteUser(userSession.UserID)
	}
}
