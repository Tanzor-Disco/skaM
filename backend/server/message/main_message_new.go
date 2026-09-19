package message

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/auth"
	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/internal/utils"
	"github.com/Tanzor-Disco/skaM/models"
)

type MessageRegisterRequest struct {
	RoomID models.RoomID
	Text   string
}

// HandleMainNewMessage handles http request sent to /api/new/message
// it checks the user session, decodes MessageRegisterRequest
// creates a new message, inserts it into messages table,
// forms a new instance of MessageData struct
// sends the messageData to all the room users through websocket
func (h *MessageHandler) HandleMainNewMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userSession, err := auth.GetUserSession(r, h.db)
	if err != nil {
		log.Printf("handleMainNewMessage: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidSessionString, nil)
		utils.SendJSON(w, http.StatusUnauthorized, body)
		return
	}

	var messageRequest MessageRegisterRequest
	err = json.NewDecoder(r.Body).Decode(&messageRequest)
	if err != nil {
		log.Printf("handleMainNewMessage: json.NewDecoder: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidJSON, nil)
		utils.SendJSON(w, http.StatusBadRequest, body)
		return
	}

	message := models.NewMessage(userSession.UserID, messageRequest.RoomID, messageRequest.Text)
	messageID, err := h.db.CreateMessage(ctx, message)
	if err != nil {
		log.Printf("handleMainNewMessage: CreateMessage: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}

	user, err := h.db.GetUserByID(ctx, userSession.UserID)
	if err != nil {
		log.Printf("handleMainNewMessage: GetUserByID: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}
	messageData := db.NewMessageData(messageID, userSession.UserID, messageRequest.RoomID, messageRequest.Text, user.Username)
	roomUsers, err := h.db.GetRoomUsersByRoomID(ctx, messageRequest.RoomID)
	if err != nil {
		log.Printf("handleMainNewMessage: GetRoomUsersByRoomID: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}
	h.hub.SendMessageData(roomUsers, messageData)

	body := utils.NewServerResponseBody(apperrors.KindErrNone, []db.MessageData{messageData})
	utils.SendJSON(w, http.StatusOK, body)
}
