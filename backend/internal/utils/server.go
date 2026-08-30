package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// serverResponseBody represents the JSON structure used in server responses
type ServerResponseBody[T any] struct {
	ErrorKind string `json:"error_kind"`
	Data      []T    `json:"data"`
}

// newServerResponseBody creates an instance of serverResponseBody
// if no data is sent, nil should be used
func NewServerResponseBody[T any](errorKind string, data []T) ServerResponseBody[T] {
	return ServerResponseBody[T]{
		ErrorKind: errorKind,
		Data:      data,
	}
}

// sendJSON encodes serverResponseBody instance to bytes, sets the header to JSON, sets the selected code
// Writes to ResponseWriter body
func SendJSON[T any](w http.ResponseWriter, code int, body ServerResponseBody[T]) {
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
