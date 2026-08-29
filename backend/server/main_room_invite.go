package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
)

func (s *server) handleMainRoomInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userSession, err := s.getUserSession(r)
	if err != nil {
		log.Printf("handleMainRoomInvite: %v", err)
		body := newServerResponseBody[any](apperrors.KindErrInvalidSessionString, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}
	roomIDString := r.URL.Query().Get("room_id")
	if roomIDString == "" {
		log.Printf("handleMainRoomInvite: couldn't get roomIDString")
		body := newServerResponseBody[any](apperrors.KindErrInvalidQueryParam, nil)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}
	roomID, err := strconv.ParseInt(roomIDString, 10, 64)
	if err != nil {
		log.Printf("handleMainRoomInvite: strconv.Atoi: %v", err)
		body := newServerResponseBody[any](apperrors.KindErrInvalidQueryParam, nil)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}

	if err = s.db.ValidateRoomUser(ctx, userSession.UserID, models.RoomID(roomID)); err != nil {
		log.Printf("handleMainRoomInvite: ValidateRoomUser: %v", err)
		body := newServerResponseBody[any](apperrors.KindErrInvalidSessionString, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}

	roomInvite, err := s.db.GetRoomInviteByroomID(ctx, models.RoomID(roomID))
	if err != nil {
		log.Printf("handleMainRoomInvite: %v", err)
		body := newServerResponseBody[any](apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	inviteURL := s.baseURL + "/invite?token=" + roomInvite.Token
	body := newServerResponseBody(apperrors.KindErrNone, []string{inviteURL})
	sendJSON(w, http.StatusOK, body)
}
