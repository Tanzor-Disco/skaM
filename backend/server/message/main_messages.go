package message

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/internal/utils"
	"github.com/Tanzor-Disco/skaM/models"
)

// HandleMainMessages handles requests to api/main/messages
// it gets the room id from the request query param,
// does an sql query to get a slice of MessageData
// encodes the slice in json, sends it as a response
func (h *MessageHandler) HandleMainMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomIDString := r.URL.Query().Get("room_id")
	if roomIDString == "" {
		log.Println("r.URL.Query: no room_id parameter found")
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidQueryParam, nil)
		utils.SendJSON(w, http.StatusBadRequest, body)
	}
	roomID, err := strconv.ParseInt(roomIDString, 10, 64)
	if err != nil {
		log.Printf("handleMainMessages: ParseInt: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidQueryParam, nil)
		utils.SendJSON(w, http.StatusBadRequest, body)
		return
	}
	messageData, err := h.db.GetMessageData(ctx, models.RoomID(roomID))
	if err != nil {
		log.Printf("handleMainMessages: GetMessageData: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidQueryParam, nil)
		utils.SendJSON(w, http.StatusBadRequest, body)
		return
	}
	body := utils.NewServerResponseBody(apperrors.KindErrNone, messageData)
	utils.SendJSON(w, http.StatusOK, body)
}
