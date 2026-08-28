package server

import (
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
)

func (s *server) handleInfoRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, err := s.getUserSession(r)
	if err != nil {
		log.Printf("handleInfoRoom: getUserSession: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidSessionString, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}
	inviteToken := r.URL.Query().Get("token")
	if inviteToken == "" {
		log.Printf("handleInfoRoom: URL.Query: empty token")
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidInviteToken, nil)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}
	roomInvite, err := s.db.GetRoomInviteByInviteToken(ctx, inviteToken)
	if err != nil {
		log.Printf("handleInfoRoom: GetRoomInviteByInviteToken: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidInviteToken, nil)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}
	room, err := s.db.GetRoomByRoomID(ctx, roomInvite.RoomID)
	if err != nil {
		log.Printf("handleInfoRoom: GetRoomByRoomID: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}
	body := newServerResponseBody(false, apperrors.KindErrNone, []models.Room{room})
	sendJSON(w, http.StatusOK, body)
}
