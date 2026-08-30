package page_main

import (
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/auth"
	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/internal/utils"
	"github.com/Tanzor-Disco/skaM/models"
)

type MainHandler struct {
	db *db.DB
}

func NewMainHandler(db *db.DB) MainHandler {
	return MainHandler{
		db: db,
	}
}

func (h *MainHandler) HandleMain(w http.ResponseWriter, r *http.Request) {
	userSession, err := auth.GetUserSession(r, h.db)
	if err != nil {
		log.Printf("handleMain: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidSessionString, nil)
		utils.SendJSON(w, http.StatusUnauthorized, body)
		return
	}
	body := utils.NewServerResponseBody(apperrors.KindErrNone, []models.UserSessionID{userSession.ID})
	utils.SendJSON(w, http.StatusOK, body)
}
