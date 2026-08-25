package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
)

type AddRequest struct {
	Token string
}

func (s *server) handleInviteAdd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userSession, err := s.getUserSession(r)
	if err != nil {
		log.Printf("handleInviteAdd: getUserSession: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidSessionString, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}

	var addRequest AddRequest
	if err = json.NewDecoder(r.Body).Decode(&addRequest); err != nil {
		log.Printf("handleInviteAdd: Decode: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidJSON, nil)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}

	roomInvite, err := s.db.GetRoomInviteByInviteToken(ctx, addRequest.Token)
	if err != nil {
		log.Printf("handleInviteAdd: GetRoomInviteByInviteToken: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidInviteToken, nil)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}

	roomUser := db.NewRoomUser(roomInvite.RoomID, userSession.UserID)
	err = s.db.AddUserToRoom(ctx, roomUser)
	if err == apperrors.ErrUniqueViolation {
		log.Printf("handleInviteAdd: AddUserToRoom: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrUniqueViolation, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}
	if err != nil {
		log.Printf("handleInviteAdd: AddUserToRoom: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	body := newServerResponseBody[any](false, apperrors.KindErrNone, nil)
	sendJSON(w, http.StatusOK, body)
}
