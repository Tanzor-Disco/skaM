package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/validate"
	"github.com/Tanzor-Disco/skaM/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)



func getUserData(request *http.Request) (models.RegisterRequest, error) {
	var currUser models.RegisterRequest
	err := json.NewDecoder(request.Body).Decode(&currUser)
	err = validate.RegisterRequest(currUser)
	return currUser, err
}

func (s *server) registerUser(ctx context.Context, userDB db.User) error {
	err := s.db.CreateUser(ctx,userDB)
	return err
}

func (s *server) handleRegister(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	currUser, err := getUserData(req)
	if err != nil {
		body := newServerResponseBody(false, apperrors.KindErrInvalidJSON)
		log.Printf("handleRegister error: %v",err)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}
	
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(currUser.Password), bcrypt.DefaultCost)
	if err != nil {
		body := newServerResponseBody(false, apperrors.KindErrHashingPassword)
		log.Printf("handleRegister error: %v",err)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	date := time.Now().Year()
	userDB := db.User{
		Email:          currUser.Email,
		Username:       currUser.Username,
		PasswordHash:   string(passwordHash),
		LastYearActive: date,
	}
	err = s.registerUser(ctx,userDB)
	if err != nil {
		var body serverResponseBody
		if errors.Is(err, apperrors.ErrEmailTaken) {
			body = newServerResponseBody(false, apperrors.KindErrEmailTaken)
			log.Printf("handleRegister error: %v",err)
			sendJSON(w, http.StatusBadRequest, body)
			return
		}

		body = newServerResponseBody(false, apperrors.KindErrDB)
		log.Printf("handleRegister error: %v",err)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	body := newServerResponseBody(true, apperrors.KindErrNone)
	sendJSON(w, http.StatusOK, body)

}
