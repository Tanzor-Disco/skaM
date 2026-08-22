package server

import (
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
)

func (s *server) handleRoot(w http.ResponseWriter, r *http.Request) {
	_, err := s.getUserSession(r)
	if err != nil {
		log.Printf("getUserSession: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidSessionString, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}
	body := newServerResponseBody[any](true, apperrors.KindErrNone, nil)
	sendJSON(w, http.StatusOK, body)
}
