package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/internal/testutils"
	"github.com/Tanzor-Disco/skaM/models"
)

func handleLoginDataDecode(t *testing.T, user userLoginData) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	r := createDecodeRequest(t, user, nil)
	srv.handleLogin(w, r)
	return w
}

func handleLoginDataRaw(t *testing.T, user string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	r := createRawRequest(t, user, nil)
	srv.handleLogin(w, r)
	return w
}

func TestHandleLogin_Valid(t *testing.T) {
	user := models.User{
		Id:             0,
		Email:          "testmail@test.com",
		Username:       "test_valid_user",
		PasswordHash:   testutils.GetPasswordHash(t, "123"),
		LastYearActive: 2026,
	}
	err := database.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	defer tdb.DeleteUsersByEmail(t, user.Email)
	userID := tdb.GetUserIDByEmail(t, user.Email)
	defer tdb.DeleteSessionsByUserID(t, userID)

	userLogin := userLoginData{
		Email:    user.Email,
		Password: "123",
	}
	w := handleLoginDataDecode(t, userLogin)
	wanted := wantedResult{
		Code:        http.StatusOK,
		Success:     true,
		ErrKind:     apperrors.KindErrNone,
		CookieExist: true,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleLogin_NonExistent(t *testing.T) {
	userLogin := userLoginData{
		Email:    "non_existent@test.com",
		Password: "123",
	}
	w := handleLoginDataDecode(t, userLogin)
	wanted := wantedResult{
		Code:        http.StatusUnauthorized,
		Success:     false,
		ErrKind:     apperrors.KindErrWrongLoginData,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleLogin_InvalidJSON(t *testing.T) {
	userLogin := ""
	w := handleLoginDataRaw(t, userLogin)
	wanted := wantedResult{
		Code:        http.StatusBadRequest,
		Success:     false,
		ErrKind:     apperrors.KindErrInternal,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleLogin_Wrong_Password(t *testing.T) {
	user := models.User{
		Id:             0,
		Email:          "testmail@test.com",
		Username:       "test_wrong_password",
		PasswordHash:   testutils.GetPasswordHash(t, "123"),
		LastYearActive: 2026,
	}
	err := database.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	defer tdb.DeleteUsersByEmail(t, user.Email)
	userID := tdb.GetUserIDByEmail(t, user.Email)
	defer tdb.DeleteSessionsByUserID(t, userID)

	userLogin := userLoginData{
		Email:    user.Email,
		Password: "456",
	}
	w := handleLoginDataDecode(t, userLogin)
	wanted := wantedResult{
		Code:        http.StatusUnauthorized,
		Success:     false,
		ErrKind:     apperrors.KindErrWrongLoginData,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}
