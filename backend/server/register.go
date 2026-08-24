package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Tanzor-Disco/skaM/check/validate"
	"github.com/Tanzor-Disco/skaM/check/verify"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *server) checkEmailTaken(email string) error {
	_, err := s.db.GetUserByEmail(email)
	if err == apperrors.ErrUserNotFound {
		return nil
	}
	if err != nil {
		log.Printf("couldn't get user by email: %v", err)
		return apperrors.ErrDB
	}

	return apperrors.ErrEmailTaken
}

func (s *server) getUserData(request *http.Request) (models.RegisterRequest, error) {
	var currUser models.RegisterRequest
	err := json.NewDecoder(request.Body).Decode(&currUser)
	if err != nil {
		log.Printf("couldn't decode request body: %v", err)
	}
	err = validate.RegisterRequest(currUser)
	if err != nil {
		return currUser, err
	}
	err = s.checkEmailTaken(currUser.Email)
	return currUser, err
}

func getUserDataErrorBody(err error) serverResponseBody[any] {
	var body serverResponseBody[any]
	switch err {
	case apperrors.ErrInvalidEmail:
		body = newServerResponseBody[any](false, apperrors.KindErrInvalidEmail, nil)
	case apperrors.ErrInvalidEmailLength:
		body = newServerResponseBody[any](false, apperrors.KindErrInvalidEmailLength, nil)
	case apperrors.ErrInvalidUsernameLength:
		body = newServerResponseBody[any](false, apperrors.KindErrInvalidUsernameLength, nil)
	case apperrors.ErrInvalidPasswordChars:
		body = newServerResponseBody[any](false, apperrors.KindErrForbiddenPasswordChars, nil)
	case apperrors.ErrInvalidPasswordLength:
		body = newServerResponseBody[any](false, apperrors.KindErrInvalidPasswordLength, nil)
	case apperrors.ErrEmailTaken:
		body = newServerResponseBody[any](false, apperrors.KindErrEmailTaken, nil)
	case apperrors.ErrDB:
		body = newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
	}
	return body
}

func (s *server) registerPendingUser(ctx context.Context, userDB models.PendingUser) error {
	err := s.db.CreatePendingUser(ctx, userDB)
	return err
}

// handleRegister handles POST requests to /api/register.
// It validates the registration data, hashes the password, creates an email
// verification token and expiration time, stores the pending user in the database,
// and sends the verification email asynchronously.
func (s *server) handleRegister(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	currUser, err := s.getUserData(req)
	if err != nil {
		body := getUserDataErrorBody(err)
		sendJSON(w, http.StatusBadRequest, body)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(currUser.Password), bcrypt.DefaultCost)
	if err != nil {
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		log.Printf("handleRegister error: %v", err)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	expiresAt := time.Now().Add(time.Hour)
	token, tokenHash, err := verify.CreateToken()
	if err != nil {
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
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
		body := newServerResponseBody[any](false, apperrors.KindErrInternal, nil)
		log.Printf("handleRegister error: %v", err)
		sendJSON(w, http.StatusInternalServerError, body)
		return
	}

	body := newServerResponseBody[any](true, apperrors.KindErrNone, nil)
	sendJSON(w, http.StatusOK, body)

	go func() {
		err := verify.SendEmail(token, currUser.Email, s.baseURL, s.SMTPData)
		if err != nil {
			log.Printf("handleRegister error: %v", err)
		}
	}()

}

// addUserToMainDB takes user info from pending_users and adds it to users DB
func (s *server) addUserToMainDB(w http.ResponseWriter, req *http.Request) {
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
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}
