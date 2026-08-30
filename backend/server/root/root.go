package root

import (
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/auth"
	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/internal/utils"
)

type RootHandler struct {
	db *db.DB
}

func NewRootHandler(db *db.DB) RootHandler {
	return RootHandler{
		db: db,
	}
}

func (h *RootHandler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	_, err := auth.GetUserSession(r, h.db)
	if err != nil {
		log.Printf("handleRoot: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInvalidSessionString, nil)
		utils.SendJSON(w, http.StatusUnauthorized, body)
		return
	}
	body := utils.NewServerResponseBody[any](apperrors.KindErrNone, nil)
	utils.SendJSON(w, http.StatusOK, body)
}
