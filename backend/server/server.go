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

type serverResponseBody struct {
	Success bool `json:"success"`
	Error string `json:"error"`
	ErrorKind string `json:"error_kind"`
}

func newServerResponseBody(success bool,err error, errorKind string) serverResponseBody {
	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
	}
	return serverResponseBody {
		Success: success, 
		Error: errorMessage,
		ErrorKind: errorKind,
	}
}

func sendJSON(w http.ResponseWriter,code int, body serverResponseBody) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)

	bodyBytes,err := json.Marshal(body)
	if err != nil {
		log.Println(err)
	}
	_,err = w.Write(bodyBytes)
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
