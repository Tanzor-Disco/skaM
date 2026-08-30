package login

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Tanzor-Disco/skaM/auth"
	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/internal/utils"
	"github.com/Tanzor-Disco/skaM/models"
	"golang.org/x/crypto/bcrypt"
)

type LoginHandler struct {
	db *db.DB
}

func NewLoginHandler(db *db.DB) LoginHandler {
	return LoginHandler{
		db: db,
	}
}

type UserLoginData struct {
	Email    string
	Password string
}

// handleLogin handles POST requests to /api/login
// It processes sent JSON, gets the user from users db and checks that the fields match
// After the verification is successful it send the user cookie that contains unique sessionString
// It creates an entry in user_sessions DB that contains the user id and sessionString
func (h *LoginHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user UserLoginData
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("handleLogin: error decoding json: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusBadRequest, body)
		return
	}

	userDB, err := h.db.GetUserByEmail(user.Email)
	if errors.Is(err, apperrors.ErrUserNotFound) {
		body := utils.NewServerResponseBody[any](apperrors.KindErrWrongLoginData, nil)
		utils.SendJSON(w, http.StatusUnauthorized, body)
		return
	}
	if err != nil {
		log.Printf("handleLogin: GetUserByEmail: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDB.PasswordHash), []byte(user.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		body := utils.NewServerResponseBody[any](apperrors.KindErrWrongLoginData, nil)
		utils.SendJSON(w, http.StatusUnauthorized, body)
		return
	}
	if err != nil {
		log.Printf("hanldeLogin: CompareHashAndPassword: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}

	sessionString, err := auth.CreateSessionString()
	if err != nil {
		log.Printf("handleLogin: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}
	userSession := models.NewUserSession(userDB.Id, sessionString)
	err = h.db.CreateSession(r.Context(), userSession)
	if err != nil {
		log.Printf("handleLogin: %v", err)
		body := utils.NewServerResponseBody[any](apperrors.KindErrInternal, nil)
		utils.SendJSON(w, http.StatusInternalServerError, body)
		return
	}
	cookie := auth.CreateSessionCookie(sessionString)
	http.SetCookie(w, &cookie)

	body := utils.NewServerResponseBody[any](apperrors.KindErrNone, nil)
	utils.SendJSON(w, http.StatusOK, body)
}
