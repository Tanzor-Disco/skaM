package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
)

func (s *server) handleMainMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomIDString := r.URL.Query().Get("room_id")
	if roomIDString == "" {
		log.Println("r.URL.Query: no room_id parameter found")
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidQueryParam, nil)
		sendJSON(w, http.StatusBadRequest, body)
	}
	roomID, err := strconv.Atoi(roomIDString)
	if err != nil {
		log.Printf("strconv.Atoi: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInvalidQueryParam, nil)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}
	messageData, err := s.db.GetMessageData(ctx, roomID)
	body := newServerResponseBody(false, apperrors.KindErrNone, messageData)
	sendJSON(w, http.StatusOK, body)
}
