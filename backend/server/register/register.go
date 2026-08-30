package register

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Tanzor-Disco/skaM/check/validate"
	"github.com/Tanzor-Disco/skaM/check/verify"
	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/internal/utils"
	"github.com/Tanzor-Disco/skaM/models"
	"golang.org/x/crypto/bcrypt"
)

type RegisterHandler struct {
	db       *db.DB
	baseURL  string
	SMTPData models.SMTPData
}

func NewRegisterHandler(db *db.DB, baseURL string, smtpData models.SMTPData) RegisterHandler {
	return RegisterHandler{
		db:       db,
		baseURL:  baseURL,
		SMTPData: smtpData,
	}
}

func (h *RegisterHandler) checkEmailTaken(email string) error {
	_, err := h.db.GetUserByEmail(email)
	if errors.Is(err, apperrors.ErrUserNotFound) {
		return nil
	}
	if err != nil {
		log.Printf("checkEmailTaken: %v", err)
		return apperrors.ErrDB
	}

	return apperrors.ErrEmailTaken
}

func (s *RegisterHandler) getUserData(request *http.Request) (models.RegisterRequest, error) {
	var currUser models.RegisterRequest
	err := json.NewDecoder(request.Body).Decode(&currUser)
	if err != nil {
		log.Printf("getUserData: Decode: %v", err)
	}
	err = validate.RegisterRequest(currUser)
	if err != nil {
		return currUser, err
	}
	err = s.checkEmailTaken(currUser.Email)
	return currUser, err
}

func getUserDataErrorBody(err error) utils.ServerResponseBody[any] {
	var body utils.ServerResponseBody[any]
	switch {
	case errors.Is(err, apperrors.ErrInvalidEmail):
		body = utils.NewServerResponseBody[any](apperrors.KindErrInvalidEmail, nil)
	case errors.Is(err, apperrors.ErrInvalidEmailLength):
		body = utils.NewServerResponseBody[any](apperrors.KindErrInvalidEmailLength, nil)
	case errors.Is(err, apperrors.ErrInvalidUsernameLength):
		body = utils.NewServerResponseBody[any](apperrors.KindErrInvalidUsernameLength, nil)
	case errors.Is(err, apperrors.ErrInvalidPasswordChars):
		body = utils.NewServerResponseBody[any](apperrors.KindErrForbiddenPasswordChars, nil)
	case errors.Is(err, apperrors.ErrInvalidPasswordLength):
		body = utils.NewServerResponseBody[any](apperrors.KindErrInvalidPasswordLength, nil)
	case errors.Is(err, apperrors.ErrEmailTaken):
		body = utils.NewServerResponseBody[any](apperrors.KindErrEmailTaken, nil)
	case errors.Is(err, apperrors.ErrDB):
		body = utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
	}
	return body
}

func (h *RegisterHandler) registerPendingUser(ctx context.Context, userDB models.PendingUser) error {
	err := h.db.CreatePendingUser(ctx, userDB)
	if err != nil {
		return fmt.Errorf("registerPendingUser: %w", err)
	}
	return nil
}

// handleRegister handles POST requests to /api/register.
// It validates the registration data, hashes the password, creates an email
// verification token and expiration time, stores the pending user in the database,
// and sends the verification email asynchronously.
func (h *RegisterHandler) HandleRegister(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	currUser, err := h.getUserData(req)
	if err != nil {
		body := getUserDataErrorBody(err)
		utils.SendJSON(w, http.StatusBadRequest, body)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(currUser.Password), bcrypt.DefaultCost)
	if err != nil {
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		log.Printf("handleRegister: GenerateFromPassword: %v", err)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}

	expiresAt := time.Now().Add(time.Hour)
	token, tokenHash, err := verify.CreateToken()
	if err != nil {
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		log.Printf("handleRegister: %v", err)
		utils.SendJSON(w, http.StatusInternalServerError, body)
	}
	pendingUserDB := models.PendingUser{
		Email:        currUser.Email,
		Username:     currUser.Username,
		TokenHash:    tokenHash,
		PasswordHash: string(passwordHash),
		ExpiresAt:    expiresAt,
	}
	err = h.registerPendingUser(ctx, pendingUserDB)
	if err != nil {
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		log.Printf("handleRegister: %v", err)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}

	body := utils.NewServerResponseBody[any](apperrors.KindErrNone, nil)
	utils.SendJSON(w, http.StatusOK, body)

	go func() {
		err := verify.SendEmail(token, currUser.Email, h.baseURL, h.SMTPData)
		if err != nil {
			log.Printf("handleRegister: %v", err)
		}
	}()

}

// addUserToMainDB takes user info from pending_users and adds it to users DB
func (h *RegisterHandler) AddUserToMainDB(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	token := req.URL.Query().Get("token")
	sum := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(sum[:])
	user, err := h.db.GetPendingUserByTokenHash(ctx, tokenHash)
	if err != nil {
		log.Printf("error in AddUserToMain: %v", err)
	}
	userDB := models.NewUser(user.Email, user.Username, user.PasswordHash)
	err = h.db.CreateUser(ctx, userDB)
	if err != nil {
		log.Printf("error in AddUserToMain: %v", err)
	}
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}
