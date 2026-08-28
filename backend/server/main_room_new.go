package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/check/invite"
	"github.com/Tanzor-Disco/skaM/check/validate"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
)

type RoomRegisterRequest struct {
	Name string
}

// handleMainNewRoom handles POST requests to /api/main/room/new
// It reads the sent room registration data, gets the user session string
// Validates room name length, then creates entries in rooms with room name
// And in room_users with user id that it got from user_sessions and room id
func (s *server) handleMainNewRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userSession, err := s.getUserSession(r)
	if err != nil {
		log.Printf("getUserSession: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidSessionString, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}

	var registerRequest RoomRegisterRequest
	err = json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		log.Printf("couldn't decode json: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}

	err = validate.RoomNameLength(registerRequest.Name)
	if err == apperrors.ErrInvalidRoomNameLength {
		log.Printf("validate.RoomNameLength: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidRoomNameLength, nil)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}

	room := models.NewRoom(registerRequest.Name)
	roomID, err := s.db.CreateRoom(ctx, room)
	if err != nil {
		log.Printf("CreateRoom: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	inviteToken, err := invite.CreateInviteToken()
	if err != nil {
		log.Printf("handleMainNewRoom: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	roomInvite := models.NewRoomInvite(roomID, inviteToken)
	if err = s.db.CreateRoomInvite(ctx, roomInvite); err != nil {
		log.Printf("handleMainNewRoom: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	roomUser := models.NewRoomUser(roomID, userSession.UserID)
	err = s.db.AddUserToRoom(ctx, roomUser)
	if err != nil {
		log.Printf("AddUserToRoom: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	// returning one member slice with room id for frontend to create a room instance
	body := newServerResponseBody(true, apperrors.KindErrNone, []models.RoomID{roomID})
	sendJSON(w, http.StatusOK, body)
}
