package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
)

func (s *server) handleMainRoomInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userSession, err := s.getUserSession(r)
	if err != nil {
		log.Printf("getUserSession: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidSessionString, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}
	roomIDString := r.URL.Query().Get("room_id")
	if roomIDString == "" {
		log.Printf("couldn't get roomIDString")
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidQueryParam, nil)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}
	roomID, err := strconv.ParseInt(roomIDString, 10, 64)
	if err != nil {
		log.Printf("strconv.Atoi: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidQueryParam, nil)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}

	if err = s.db.ValidateRoomUser(ctx, userSession.UserID, roomID); err != nil {
		log.Printf("ValidateRoomUser: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidSessionString, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}

	roomInvite, err := s.db.GetRoomInviteByroomID(ctx, roomID)
	if err != nil {
		log.Printf("GetRoomInviteByroomID: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	inviteURL := s.baseURL + "/main/invite?token=" + roomInvite.Token
	body := newServerResponseBody(true, apperrors.KindErrNone, []string{inviteURL})
	sendJSON(w, http.StatusOK, body)
}
