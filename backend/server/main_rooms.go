package server

import (
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"log"
	"net/http"
)

// handleMainRooms handles /api/main/rooms requests.
// It gets the session string from the request, retrieves the user ID
// from user_sessions, gets the user's room IDs,
// and retrieves the room information from rooms.
// If successful, it sends []Room to frontend.
func (s *server) handleMainRooms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userSession, err := s.getUserSession(r)
	if err != nil {
		log.Printf("handleMainRooms: %v", err)
		body := newServerResponseBody[any](apperrors.KindErrInvalidSessionString, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}
	roomIDS, err := s.db.GetRoomIDSByUserID(ctx, userSession.UserID)
	if err != nil {
		log.Printf("handleMain: %v", err)
		body := newServerResponseBody[any](apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	rooms, err := s.db.GetRoomsByRoomIDS(ctx, roomIDS)
	if err != nil {
		log.Printf("handleMain: %v", err)
		body := newServerResponseBody[any](apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	body := newServerResponseBody(apperrors.KindErrNone, rooms)
	sendJSON(w, http.StatusOK, body)
}
