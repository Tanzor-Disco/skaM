package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
)

type AddRequest struct {
	Token string
}

// handleInviteAdd handles adding a user to room users table after the invitations has been accepted
// it checks the session of the user, gets the room id through room invite
// it creates a new room user, adds them to the table
// it handles the unique violation in the table separately
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

	roomUser := models.NewRoomUser(roomInvite.RoomID, userSession.UserID)
	err = s.db.AddUserToRoom(ctx, roomUser)

	// a unique violation occurs when the user who is already in the room accepts an invite
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

	body := newServerResponseBody[any](true, apperrors.KindErrNone, nil)
	sendJSON(w, http.StatusOK, body)
}
