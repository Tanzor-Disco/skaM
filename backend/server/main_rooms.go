package server

import (
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"log"
	"net/http"
)

func (s *server) handleMainRooms(w http.ResponseWriter, r *http.Request) {
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

	userSession, err := s.db.GetUserSessionBySessionString(ctx, cookie.Value)
	if err != nil {
		log.Printf("handleMain: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	roomIDS, err := s.db.GetRoomIDSByUserID(ctx, userSession.UserID)
	if err != nil {
		log.Printf("handleMain: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	rooms, err := s.db.GetRoomsByRoomIDS(ctx, roomIDS)
	if err != nil {
		log.Printf("handleMain: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	body := newServerResponseBody(true, apperrors.KindErrNone, rooms)
	sendJSON(w, http.StatusOK, body)
}
