package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"golang.org/x/crypto/bcrypt"
)

type userData struct {
	Email    string
	Username string
	Password string
}

func getUserData(request *http.Request) (userData, error) {
	var currUser userData
	err := json.NewDecoder(request.Body).Decode(&currUser)
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
		body := newServerResponseBody(false, err, apperrors.KindErrInvalidJSON)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(currUser.Password), bcrypt.DefaultCost)
	if err != nil {
		body := newServerResponseBody(false, err, apperrors.KindErrHashingPassword)
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
			body = newServerResponseBody(false, err, apperrors.KindErrEmailTaken)
			sendJSON(w, http.StatusBadRequest, body)
			return
		}

		body = newServerResponseBody(false, err, apperrors.KindErrDB)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	body := newServerResponseBody(true, nil, apperrors.KindErrNone)
	sendJSON(w, http.StatusOK, body)

}
