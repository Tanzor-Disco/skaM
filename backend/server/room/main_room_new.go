package room

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/auth"
	"github.com/Tanzor-Disco/skaM/check/invite"
	"github.com/Tanzor-Disco/skaM/check/validate"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/internal/utils"
	"github.com/Tanzor-Disco/skaM/models"
)

type RoomRegisterRequest struct {
	Name string
}

// handleMainNewRoom handles POST requests to /api/main/room/new
// It reads the sent room registration data, gets the user session string
// Validates room name length, then creates entries in rooms with room name
// And in room_users with user id that it got from user_sessions and room id
func (h *RoomHandler) HandleMainNewRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userSession, err := auth.GetUserSession(r, h.db)
	if err != nil {
		log.Printf("handleMainNewRoom: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidSessionString, nil)
		utils.SendJSON(w, http.StatusUnauthorized, body)
		return
	}

	var registerRequest RoomRegisterRequest
	err = json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		log.Printf("handleMainNewRoom: Decode: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusBadRequest, body)
		return
	}

	err = validate.RoomNameLength(registerRequest.Name)
	if errors.Is(err, apperrors.ErrInvalidRoomNameLength) {
		log.Printf("handleMainNewRoom: validate.RoomNameLength: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidRoomNameLength, nil)
		utils.SendJSON(w, http.StatusBadRequest, body)
		return
	}

	room := models.NewRoom(registerRequest.Name)
	roomID, err := h.db.CreateRoom(ctx, room)
	if err != nil {
		log.Printf("handleMainNewRoom: CreateRoom: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}

	inviteToken, err := invite.CreateInviteToken()
	if err != nil {
		log.Printf("handleMainNewRoom: CreateRoomInvite: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}

	roomInvite := models.NewRoomInvite(roomID, inviteToken)
	if err = h.db.CreateRoomInvite(ctx, roomInvite); err != nil {
		log.Printf("handleMainNewRoom: CreateRoomInvite: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}

	roomUser := models.NewRoomUser(roomID, userSession.UserID)
	err = h.db.AddUserToRoom(ctx, roomUser)
	if err != nil {
		log.Printf("handleMainNewRoom: AddUserToRoom: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}

	// returning one member slice with room id for frontend to create a room instance
	body := utils.NewServerResponseBody(apperrors.KindErrNone, []models.RoomID{roomID})
	utils.SendJSON(w, http.StatusOK, body)
}
