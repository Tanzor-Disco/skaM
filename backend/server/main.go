package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
)

func (s *server) getUserSession(r *http.Request) (db.UserSession, error) {
	cookie, err := r.Cookie("session_string")
	if err != nil {
		return db.UserSession{}, fmt.Errorf("r.Cookie: %w", err)
	}
	sessionString := cookie.Value
	userSession, err := s.db.GetUserSessionBySessionString(r.Context(), sessionString)
	if err != nil {
		return db.UserSession{}, fmt.Errorf("GetUserSessionBySessionString: %w", err)
	}
	return userSession, nil
}

func (s *server) handleMain(w http.ResponseWriter, r *http.Request) {
	userSession, err := s.getUserSession(r)
	if err != nil {
		log.Printf("getUserSession: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidSessionString, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}
	body := newServerResponseBody(false, apperrors.KindErrNone, []int64{userSession.ID})
	sendJSON(w, http.StatusOK, body)
}
