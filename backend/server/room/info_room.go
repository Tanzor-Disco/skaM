package room

import (
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/auth"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/internal/utils"
	"github.com/Tanzor-Disco/skaM/models"
)

// handleInfoRoom handles the http requests sent to /api/info/room
// it checks the user session, gets the invite token from the query param,
// gets the invite using the invite token
// get room from rooms using the room id from the invite token
// sends the room as a response
func (h *RoomHandler) HandleInfoRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, err := auth.GetUserSession(r, h.db)
	if err != nil {
		log.Printf("handleInfoRoom: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidSessionString, nil)
		utils.SendJSON(w, http.StatusUnauthorized, body)
		return
	}
	inviteToken := r.URL.Query().Get("token")
	if inviteToken == "" {
		log.Printf("handleInfoRoom: URL.Query: empty token")
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidInviteToken, nil)
		utils.SendJSON(w, http.StatusBadRequest, body)
		return
	}
	roomInvite, err := h.db.GetRoomInviteByInviteToken(ctx, inviteToken)
	if err != nil {
		log.Printf("handleInfoRoom: GetRoomInviteByInviteToken: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidInviteToken, nil)
		utils.SendJSON(w, http.StatusBadRequest, body)
		return
	}
	room, err := h.db.GetRoomByRoomID(ctx, roomInvite.RoomID)
	if err != nil {
		log.Printf("handleInfoRoom: GetRoomByRoomID: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}
	body := utils.NewServerResponseBody(apperrors.KindErrNone, []models.Room{room})
	utils.SendJSON(w, http.StatusOK, body)
}
