package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
)

func (s *server) handleMainMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomIDString := r.URL.Query().Get("room_id")
	if roomIDString == "" {
		log.Println("r.URL.Query: no room_id parameter found")
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidQueryParam, nil)
		sendJSON(w, http.StatusBadRequest, body)
	}
	roomID, err := strconv.ParseInt(roomIDString, 10, 64)
	if err != nil {
		log.Printf("handleMainMessages: ParseInt: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidQueryParam, nil)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}
	messageData, err := s.db.GetMessageData(ctx, models.RoomID(roomID))
	if err != nil {
		log.Printf("handleMainMessages: GetMessageData: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidQueryParam, nil)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}
	body := newServerResponseBody(false, apperrors.KindErrNone, messageData)
	sendJSON(w, http.StatusOK, body)
}
