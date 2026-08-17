package server

import (
	"net/http"
	"log"
	"encoding/json"
	
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/check/validate"
	"github.com/Tanzor-Disco/skaM/db"
)


type RoomRegisterRequest struct {
	Name string
}

func (s *server) handleMainNewRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userSession,err := s.getUserSession(r)
	if err != nil {
		log.Printf("getUserSession: %v",err)
		body := newServerResponseBody(false,apperrors.KindErrInternal)
		sendJSON(w,http.StatusUnauthorized,body)
		return 
	}
	
	var registerRequest RoomRegisterRequest
	err = json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		log.Printf("couldn't decode json: %v",err)
		body := newServerResponseBody(false,apperrors.KindErrInternal)
		sendJSON(w, http.StatusInternalServerError, body)
		return 
	}

	err = validate.RoomNameLength(registerRequest.Name)
	if err == apperrors.ErrInvalidRoomNameLength {
		log.Printf("validate.RoomNameLength: %v", err)
		body := newServerResponseBody(false,apperrors.KindErrInvalidRoomNameLength)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}
	
	room := db.NewRoom(registerRequest.Name)
	roomID, err := s.db.CreateRoom(ctx, room)
	if err != nil {
		log.Printf("CreateRoom: %v",err)
		body := newServerResponseBody(false,apperrors.KindErrInternal)
		sendJSON(w, http.StatusInternalServerError,body)
		return
	}
	
	roomUser := db.NewRoomUser(roomID, userSession.UserID)
	err = s.db.AddUserToRoom(ctx,roomUser)
	if err != nil {
		log.Printf("AddUserToRoom: %v",err)
		body := newServerResponseBody(false,apperrors.KindErrInternal)
		sendJSON(w,http.StatusInternalServerError,body)
		return
	}

	body := newServerResponseBody(true,apperrors.KindErrNone)
	sendJSON(w,http.StatusOK,body)

}
