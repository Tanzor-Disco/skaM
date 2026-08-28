package server

import (
	"encoding/json"
	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/models"
	"log"
	"net/http"
)

type server struct {
	db       *db.DB
	SMTPData models.SMTPData
	baseURL  string
}

// newServer creates a new server instance that stores db struct and SMTPData
func newServer(serverData models.ServerData) (*server, error) {
	db, err := db.Connect(serverData.URI)
	if err != nil {
		return &server{}, err
	}
	return &server{
		db:       db,
		SMTPData: serverData.SMTPData,
		baseURL:  serverData.BaseURL,
	}, nil
}

// serverResponseBody represents the JSON structure used in server responses
type serverResponseBody[T any] struct {
	Success   bool   `json:"success"`
	ErrorKind string `json:"error_kind"`
	Data      []T    `json:"data"`
}

// newServerResponseBody creates an instance of serverResponseBody
// if no data is sent, nil should be used
func newServerResponseBody[T any](success bool, errorKind string, data []T) serverResponseBody[T] {
	return serverResponseBody[T]{
		Success:   success,
		ErrorKind: errorKind,
		Data:      data,
	}
}

// sendJSON encodes serverResponseBody instance to bytes, sets the header to JSON, sets the selected code
// Writes to ResponseWriter body
func sendJSON[T any](w http.ResponseWriter, code int, body serverResponseBody[T]) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(bodyBytes)
	if err != nil {
		log.Println(err)
	}
}

func Run(serverData models.ServerData) error {
	server, err := newServer(serverData)
	if err != nil {
		return err
	}

	http.HandleFunc("/", handleStatic)

	http.HandleFunc("/api/register", server.handleRegister)
	http.HandleFunc("/api/verify/email/", server.addUserToMainDB)
	http.HandleFunc("/api/login", server.handleLogin)
	http.HandleFunc("/api/main", server.handleMain)
	http.HandleFunc("/api/main/new/room", server.handleMainNewRoom)
	http.HandleFunc("/api/main/rooms", server.handleMainRooms)
	http.HandleFunc("/api/root", server.handleRoot)
	http.HandleFunc("/api/main/messages", server.handleMainMessages)
	http.HandleFunc("/api/main/new/message", server.handleMainNewMessage)
	http.HandleFunc("/api/main/room/invite", server.handleMainRoomInvite)
	http.HandleFunc("/api/info/room", server.handleInfoRoom)
	http.HandleFunc("/api/invite/add", server.handleInviteAdd)

	go server.db.PeriodicPendingDelete()
	go server.db.PeriodicDeleteAllExpiredSessions()

	return http.ListenAndServe(":8080", nil)
}
