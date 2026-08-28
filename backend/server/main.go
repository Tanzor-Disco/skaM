package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
)

func (s *server) getUserSession(r *http.Request) (models.UserSession, error) {
	cookie, err := r.Cookie("session_string")
	if err != nil {
		return models.UserSession{}, fmt.Errorf("r.Cookie: %w", err)
	}
	sessionString := cookie.Value
	userSession, err := s.db.GetUserSessionBySessionString(r.Context(), sessionString)
	if err != nil {
		return models.UserSession{}, fmt.Errorf("GetUserSessionBySessionString: %w", err)
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
	body := newServerResponseBody(false, apperrors.KindErrNone, []models.UserSessionID{userSession.ID})
	sendJSON(w, http.StatusOK, body)
}
