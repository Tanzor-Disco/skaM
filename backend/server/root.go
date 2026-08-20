package server

import (
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
)

func (s *server) handleRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, err := r.Cookie("session_string")
	if err == http.ErrNoCookie {
		log.Printf("handleMain: r.Cookie: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrNoSessionString, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}

	if err != nil {
		log.Printf("handleMain: r.Cookie: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	_, err = s.db.GetUserSessionBySessionString(ctx, cookie.Value)
	if err != nil {
		log.Printf("handleMain: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}

	body := newServerResponseBody[any](true, apperrors.KindErrNone, nil)
	sendJSON(w, http.StatusOK, body)
}
