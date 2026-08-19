package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/check/validate"
	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
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
		body := newServerResponseBody(false, apperrors.KindErrInternal)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}

	var registerRequest RoomRegisterRequest
	err = json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		log.Printf("couldn't decode json: %v", err)
		body := newServerResponseBody(false, apperrors.KindErrInternal)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}

	err = validate.RoomNameLength(registerRequest.Name)
	if err == apperrors.ErrInvalidRoomNameLength {
		log.Printf("validate.RoomNameLength: %v", err)
		body := newServerResponseBody(false, apperrors.KindErrInvalidRoomNameLength)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}

	room := db.NewRoom(registerRequest.Name)
	roomID, err := s.db.CreateRoom(ctx, room)
	if err != nil {
		log.Printf("CreateRoom: %v", err)
		body := newServerResponseBody(false, apperrors.KindErrInternal)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	roomUser := db.NewRoomUser(roomID, userSession.UserID)
	err = s.db.AddUserToRoom(ctx, roomUser)
	if err != nil {
		log.Printf("AddUserToRoom: %v", err)
		body := newServerResponseBody(false, apperrors.KindErrInternal)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	body := newServerResponseBody(true, apperrors.KindErrNone)
	sendJSON(w, http.StatusOK, body)
}
