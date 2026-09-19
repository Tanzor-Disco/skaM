package room

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Tanzor-Disco/skaM/auth"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/internal/utils"
	"github.com/Tanzor-Disco/skaM/models"
)

// HandleMainRoomInvite handles the http requests sent to /api/main/room/invite
// it checks the user session, get room id from request query param,
// checks if the user is in the requested room
// gets the room invite, create an invite URL, sends it as the response
func (h *RoomHandler) HandleMainRoomInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userSession, err := auth.GetUserSession(r, h.db)
	if err != nil {
		log.Printf("handleMainRoomInvite: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidSessionString, nil)
		utils.SendJSON(w, http.StatusUnauthorized, body)
		return
	}
	roomIDString := r.URL.Query().Get("room_id")
	if roomIDString == "" {
		log.Printf("handleMainRoomInvite: couldn't get roomIDString")
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidQueryParam, nil)
		utils.SendJSON(w, http.StatusBadRequest, body)
		return
	}
	roomID, err := strconv.ParseInt(roomIDString, 10, 64)
	if err != nil {
		log.Printf("handleMainRoomInvite: strconv.Atoi: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidQueryParam, nil)
		utils.SendJSON(w, http.StatusBadRequest, body)
		return
	}

	if err = h.db.ValidateRoomUser(ctx, userSession.UserID, models.RoomID(roomID)); err != nil {
		log.Printf("handleMainRoomInvite: ValidateRoomUser: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidSessionString, nil)
		utils.SendJSON(w, http.StatusUnauthorized, body)
		return
	}

	roomInvite, err := h.db.GetRoomInviteByroomID(ctx, models.RoomID(roomID))
	if err != nil {
		log.Printf("handleMainRoomInvite: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}

	inviteURL := h.baseURL + "/invite?token=" + roomInvite.Token
	body := utils.NewServerResponseBody(apperrors.KindErrNone, []string{inviteURL})
	utils.SendJSON(w, http.StatusOK, body)
}
