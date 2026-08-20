package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/auth"
	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
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
		log.Printf("error decoding json: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}

	userDB, err := s.db.GetUserByEmail(user.Email)
	if err == apperrors.ErrUserNotFound {
		body := newServerResponseBody[any](false, apperrors.KindErrWrongLoginData, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}
	if err != nil {
		log.Printf("error getting the user from db: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDB.PasswordHash), []byte(user.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		body := newServerResponseBody[any](false, apperrors.KindErrWrongLoginData, nil)
		sendJSON(w, http.StatusUnauthorized, body)
		return
	}
	if err != nil {
		log.Printf("couldn't compare hash and password: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	sessionString, err := auth.CreateSessionString()
	if err != nil {
		log.Printf("auth error: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}
	userSession := db.NewUserSession(userDB.Id, sessionString)
	err = s.db.CreateSession(r.Context(), userSession)
	if err != nil {
		log.Printf("db: %v", err)
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}
	cookie := auth.CreateSessionCookie(sessionString)
	http.SetCookie(w, &cookie)

	body := newServerResponseBody[any](true, apperrors.KindErrNone, nil)
	sendJSON(w, http.StatusOK, body)
}
