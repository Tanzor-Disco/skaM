package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Tanzor-Disco/skaM/check/validate"
	"github.com/Tanzor-Disco/skaM/check/verify"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
	"golang.org/x/crypto/bcrypt"
)

func getUserData(request *http.Request) (models.RegisterRequest, error) {
	var currUser models.RegisterRequest
	err := json.NewDecoder(request.Body).Decode(&currUser)
	err = validate.RegisterRequest(currUser)
	return currUser, err
}

func getUserDataErrorBody(err error) serverResponseBody {
	var body serverResponseBody
	switch err {
	case apperrors.ErrInvalidEmail:
		body = newServerResponseBody(false, apperrors.KindErrInvalidEmail)
	case apperrors.ErrInvalidEmailLength:
		body = newServerResponseBody(false, apperrors.KindErrInvalidEmailLength)
	case apperrors.ErrInvalidUsernameLength:
		body = newServerResponseBody(false, apperrors.KindErrInvalidUsernameLength)
	case apperrors.ErrInvalidPasswordChars:
		body = newServerResponseBody(false, apperrors.KindErrForbiddenPasswordChars)
	case apperrors.ErrInvalidPasswordLength:
		body = newServerResponseBody(false, apperrors.KindErrInvalidPasswordLength)
	}
	return body
}

func (s *server) registerPendingUser(ctx context.Context, userDB models.PendingUser) error {
	err := s.db.CreatePendingUser(ctx, userDB)
	return err
}

func (s *server) handleRegister(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	currUser, err := getUserData(req)
	if err != nil {
		body := getUserDataErrorBody(err)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(currUser.Password), bcrypt.DefaultCost)
	if err != nil {
		body := newServerResponseBody(false, apperrors.KindErrHashingPassword)
		log.Printf("handleRegister error: %v", err)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	expiresAt := time.Now().Add(time.Hour)
	token, tokenHash, err := verify.CreateToken()
	if err != nil {
		body := newServerResponseBody(false, apperrors.KindErrTokenCreation)
		log.Printf("handleRegister error: %v", err)
		sendJSON(w, http.StatusInternalServerError, body)
	}
	pendingUserDB := models.PendingUser{
		Email:        currUser.Email,
		Username:     currUser.Username,
		TokenHash:    tokenHash,
		PasswordHash: string(passwordHash),
		ExpiresAt:    expiresAt,
	}
	err = s.registerPendingUser(ctx, pendingUserDB)
	if err != nil {
		var body serverResponseBody
		if errors.Is(err, apperrors.ErrEmailTaken) {
			body = newServerResponseBody(false, apperrors.KindErrEmailTaken)
			log.Printf("handleRegister error: %v", err)
			sendJSON(w, http.StatusBadRequest, body)
			return
		}

		body = newServerResponseBody(false, apperrors.KindErrDB)
		log.Printf("handleRegister error: %v", err)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	err = verify.SendEmail(s.baseURL, token, currUser.Email, s.SMTPData)
	if err != nil {
		body := newServerResponseBody(false, apperrors.KindErrSendingEmail)
		log.Printf("handleRegister error: %v", err)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	body := newServerResponseBody(true, apperrors.KindErrNone)
	sendJSON(w, http.StatusOK, body)

}

func (s *server) AddUserToMainDB(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	token := req.URL.Query().Get("token")
	sum := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(sum[:])
	user, err := s.db.GetPendingUserByTokenHash(ctx, tokenHash)
	if err != nil {
		log.Printf("error in AddUserToMain: %v", err)
	}
	userDB := models.NewUser(user.Email, user.Username, user.PasswordHash)
	err = s.db.CreateUser(ctx, userDB)
	if err != nil {
		log.Printf("error in AddUserToMain: %v", err)
	}
}
