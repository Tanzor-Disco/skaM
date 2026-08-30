package invite

import (
	"github.com/Tanzor-Disco/skaM/db"
)

type InviteHandler struct {
	db *db.DB
}

func NewInviteHandler(db *db.DB) InviteHandler {
	return InviteHandler{
		db: db,
	}
}
