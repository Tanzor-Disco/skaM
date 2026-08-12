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
	baseURL  string
	SMTPData models.SMTPData
}

func newServer(serverData models.ServerData) (*server, error) {
	db, err := db.Connect(serverData.URI)
	if err != nil {
		return &server{}, err
	}
	return &server{
		db:       db,
		baseURL:  serverData.BaseURL,
		SMTPData: serverData.SMTPData,
	}, nil
}

type serverResponseBody struct {
	Success   bool   `json:"success"`
	ErrorKind string `json:"error_kind"`
}

func newServerResponseBody(success bool, errorKind string) serverResponseBody {
	return serverResponseBody{
		Success:   success,
		ErrorKind: errorKind,
	}
}

func sendJSON(w http.ResponseWriter, code int, body serverResponseBody) {
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

	http.HandleFunc("/api/register", server.handleRegister)
	http.HandleFunc("/api/verify/email/", server.AddUserToMainDB)

	return http.ListenAndServe(":8080", nil)
}
