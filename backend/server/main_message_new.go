package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
)

type MessageRegisterRequest struct {
	RoomID models.RoomID
	Text   string
}

func (s *server) handleMainNewMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userSession, err := s.getUserSession(r)
	if err != nil {
		log.Printf("getUserSession: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidSessionString, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}

	var messageRequest MessageRegisterRequest
	err = json.NewDecoder(r.Body).Decode(&messageRequest)
	if err != nil {
		log.Printf("json.NewDecoder: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidJSON, nil)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}

	message := models.NewMessage(userSession.UserID, messageRequest.RoomID, messageRequest.Text)
	messageID, err := s.db.CreateMessage(ctx, message)
	if err != nil {
		log.Printf("CreateMessage: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}
	user, err := s.db.GetUserByID(ctx, userSession.UserID)
	if err != nil {
		log.Printf("GetUserByID: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}
	messageData := db.NewMessageData(messageID, userSession.UserID, messageRequest.RoomID, messageRequest.Text, user.Username)
	body := newServerResponseBody(false, apperrors.KindErrNone, []db.MessageData{messageData})
	sendJSON(w, http.StatusOK, body)
}
