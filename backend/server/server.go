package server 

import (
	"net/http"
	"encoding/json"
	"github.com/Tanzor-Disco/skaM/db"
	"log"
)

type server struct {
	db *db.DB
}

func newServer(URI string) (*server,error) {
	db,err := db.Connect(URI)
	if err != nil {
		return &server{},err
	}
	return &server {
		db:db,
	},nil
}

type serverResponse struct {
	Success bool
	Error string
	Message string
}

func sendJSON(w http.ResponseWriter,success bool,message string, err error,code int) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)

	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
	}
	response := serverResponse {
		Success:success,
		Error:errorMessage,
		Message:message,
	}
	body,err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}
	_,err = w.Write(body)
	if err != nil {
		log.Println(err)
	}

}

func Run(URI string) error {
	server,err := newServer(URI)
	if err != nil {
		return err
	}

	http.HandleFunc("/api/register", server.handleRegister)
	
	return http.ListenAndServe(":8080",nil)
}
