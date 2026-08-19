package server

import (
	"fmt"
	"net/http"

	"github.com/Tanzor-Disco/skaM/db"
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
