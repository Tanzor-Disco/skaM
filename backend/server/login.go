package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/auth"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
	"golang.org/x/crypto/bcrypt"
)

type userLoginData struct {
	Email    string
	Password string
}

// handleLogin handles POST requests to /api/login
// It processes sent JSON, gets the user from users db and checks that the fields match
// After the verification is successful it send the user cookie that contains unique sessionString
// It creates an entry in user_sessions DB that contains the user id and sessionString
func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user userLoginData
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("handleLogin: error decoding json: %v", err)
		body := newServerResponseBody[any](apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}

	userDB, err := s.db.GetUserByEmail(user.Email)
	if errors.Is(err, apperrors.ErrUserNotFound) {
		body := newServerResponseBody[any](apperrors.KindErrWrongLoginData, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}
	if err != nil {
		log.Printf("handleLogin: GetUserByEmail: %v", err)
		body := newServerResponseBody[any](apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDB.PasswordHash), []byte(user.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		body := newServerResponseBody[any](apperrors.KindErrWrongLoginData, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}
	if err != nil {
		log.Printf("hanldeLogin: CompareHashAndPassword: %v", err)
		body := newServerResponseBody[any](apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	sessionString, err := auth.CreateSessionString()
	if err != nil {
		log.Printf("handleLogin: %v", err)
		body := newServerResponseBody[any](apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}
	userSession := models.NewUserSession(userDB.Id, sessionString)
	err = s.db.CreateSession(r.Context(), userSession)
	if err != nil {
		log.Printf("handleLogin: %v", err)
		body := newServerResponseBody[any](apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}
	cookie := auth.CreateSessionCookie(sessionString)
	http.SetCookie(w, &cookie)

	body := newServerResponseBody[any](apperrors.KindErrNone, nil)
	sendJSON(w, http.StatusOK, body)
}
