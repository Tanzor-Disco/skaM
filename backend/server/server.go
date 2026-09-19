package server

import (
	"net/http"

	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/models"
	"github.com/Tanzor-Disco/skaM/server/invite"
	"github.com/Tanzor-Disco/skaM/server/login"
	"github.com/Tanzor-Disco/skaM/server/message"
	"github.com/Tanzor-Disco/skaM/server/page_main"
	"github.com/Tanzor-Disco/skaM/server/register"
	"github.com/Tanzor-Disco/skaM/server/room"
	"github.com/Tanzor-Disco/skaM/server/root"
	"github.com/Tanzor-Disco/skaM/server/static"
	"github.com/Tanzor-Disco/skaM/server/websocket"
)

type Server struct {
	DB       *db.DB
	SMTPData models.SMTPData
	BaseURL  string
	WsHub    websocket.WsHub
}

// newServer creates a new server instance that stores db struct, BaseURL, SMTPData, WsHub
func NewServer(serverData models.ServerData) (*Server, error) {
	db, err := db.Connect(serverData.URI)
	if err != nil {
		return &Server{}, err
	}
	return &Server{
		DB:       db,
		SMTPData: serverData.SMTPData,
		BaseURL:  serverData.BaseURL,
		WsHub:    websocket.NewWsHub(),
	}, nil
}

// Run launches the server
// it creates handlers for different endpoint groups
// sets up handler functions
// initiates a periodic data cleanups as goroutines
func Run(serverData models.ServerData) error {
	server, err := NewServer(serverData)
	if err != nil {
		return err
	}

	roomHandler := room.NewRoomHandler(server.DB, server.BaseURL)
	messageHandler := message.NewMessageHandler(server.DB, &server.WsHub)
	loginHandler := login.NewLoginHandler(server.DB)
	mainHandler := page_main.NewMainHandler(server.DB)
	inviteHandler := invite.NewInviteHandler(server.DB)
	rootHandler := root.NewRootHandler(server.DB)
	registerHandler := register.NewRegisterHandler(server.DB, server.BaseURL, server.SMTPData)
	wsHandler := websocket.NewWsHandler(server.DB, &server.WsHub, server.BaseURL)

	http.HandleFunc("/", static.HandleStatic)

	http.HandleFunc("/api/register", registerHandler.HandleRegister)

	http.HandleFunc("/api/verify/email/", registerHandler.AddUserToMainDB)

	http.HandleFunc("/api/login", loginHandler.HandleLogin)

	http.HandleFunc("/api/main", mainHandler.HandleMain)

	http.HandleFunc("/api/main/new/room", roomHandler.HandleMainNewRoom)
	http.HandleFunc("/api/main/rooms", roomHandler.HandleMainRooms)
	http.HandleFunc("/api/main/room/invite", roomHandler.HandleMainRoomInvite)
	http.HandleFunc("/api/info/room", roomHandler.HandleInfoRoom)

	http.HandleFunc("/api/root", rootHandler.HandleRoot)

	http.HandleFunc("/api/main/messages", messageHandler.HandleMainMessages)
	http.HandleFunc("/api/main/new/message", messageHandler.HandleMainNewMessage)

	http.HandleFunc("/api/invite/add", inviteHandler.HandleInviteAdd)

	http.HandleFunc("/api/ws", wsHandler.HandleWs)

	go server.DB.PeriodicPendingDelete()
	go server.DB.PeriodicDeleteAllExpiredSessions()

	return http.ListenAndServe(":8080", nil)
}
