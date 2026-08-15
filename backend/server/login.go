package server

import (
	"net/http"
	"encoding/json"
	"log"
	"golang.org/x/crypto/bcrypt"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
)

type userLoginData struct {
	Email string
	Password string
}

func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user userLoginData
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("error decoding json: %v",err)
		body := newServerResponseBody(false,apperrors.KindErrInternal)
		sendJSON(w,http.StatusInternalServerError,body)
		return
	}
	userDB,err := s.db.GetUserByEmail(user.Email)
	if err == apperrors.ErrUserNotFound {
		body := newServerResponseBody(false, apperrors.KindErrWrongLoginData)
		sendJSON(w,http.StatusUnauthorized,body)
		return
	}
	if err != nil {
		log.Printf("error getting the user from db: %v",err)
		body := newServerResponseBody(false,apperrors.KindErrInternal)
		sendJSON(w,http.StatusInternalServerError,body)
		return 
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDB.PasswordHash),[]byte(user.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		body := newServerResponseBody(false,apperrors.KindErrWrongLoginData)
		sendJSON(w,http.StatusUnauthorized,body)
		return
	}
	if err != nil {
		log.Printf("couldn't compare hash and password: %v",err)
		body := newServerResponseBody(false,apperrors.KindErrInternal)
		sendJSON(w,http.StatusInternalServerError,body)
		return
	}

	body := newServerResponseBody(true,apperrors.KindErrNone)
	sendJSON(w,http.StatusOK,body)
}
