package server

import (
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"log"
	"net/http"
)

func (s *server) handleMainRooms(w http.ResponseWriter, r *http.Request) {
	_, err := s.getUserSession(r)
	if err != nil {
		log.Printf("getUserSession: %v", err)
		body := newServerResponseBody(false, apperrors.KindErrInternal)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}
	//TODO: add userSession; get the list of rooms from room_users; then get the info from rooms, then send it to frontend

}
