package room

import (
	"github.com/Tanzor-Disco/skaM/auth"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/internal/utils"

	"log"
	"net/http"
)

// handleMainRooms handles /api/main/rooms requests.
// It gets the session string from the request, retrieves the user ID
// from user_sessions, gets the user's room IDs,
// and retrieves the room information from rooms.
// If successful, it sends []Room to frontend.
func (h *RoomHandler) HandleMainRooms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userSession, err := auth.GetUserSession(r, h.db)
	if err != nil {
		log.Printf("handleMainRooms: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidSessionString, nil)
		utils.SendJSON(w, http.StatusUnauthorized, body)
		return
	}
	roomIDS, err := h.db.GetRoomIDSByUserID(ctx, userSession.UserID)
	if err != nil {
		log.Printf("handleMain: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}

	rooms, err := h.db.GetRoomsByRoomIDS(ctx, roomIDS)
	if err != nil {
		log.Printf("handleMain: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}

	body := utils.NewServerResponseBody(apperrors.KindErrNone, rooms)
	utils.SendJSON(w, http.StatusOK, body)
}
