package server

import (
	"context"
	"encoding/json"
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
	r := createDecodeRequest(t, user)
	srv.handleLogin(w, r)
	return w
}

func handleLoginDataRaw(t *testing.T, user string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	r := createRawRequest(t, user)
	srv.handleLogin(w, r)
	return w
}

type wantedLoginResult struct {
	Code        int
	Success     bool
	ErrKind     string
	CookieExist bool
}

func verifyLoginResponse(t *testing.T, w *httptest.ResponseRecorder, wanted wantedLoginResult) {
	t.Helper()
	if w.Code != wanted.Code {
		t.Fatalf("verifyLoginResponse: code of the response doesn't match:\n got %v\n wanted %v", w.Code, wanted.Code)
	}
	var gotBody serverResponseBody
	err := json.NewDecoder(w.Body).Decode(&gotBody)
	if err != nil {
		t.Fatalf("verifyLoginResponse: couldn't decode the server response")
	}
	if gotBody.Success != wanted.Success {
		t.Fatalf("verifyLoginResponse: success of the response doesn't match:\n got %v\n wanted %v", gotBody.Success, wanted.Success)
	}
	if gotBody.ErrorKind != wanted.ErrKind {
		t.Fatalf("verifyLoginResponse: error kind of the response doesn't match:\n got %v\n wanted %v", gotBody.ErrorKind, wanted.ErrKind)
	}
	cookies := w.Result().Cookies()
	if (len(cookies) == 1) != wanted.CookieExist {
		t.Fatalf("got %v cookies, wanted cookies: %v", len(cookies), wanted.CookieExist)
	}
	if len(cookies) == 0 {
		return
	}
	cookie := cookies[0]
	if cookie.Name != "session_string" {
		t.Fatalf("verifyLoginResponse: cookie name mismatch:\n got %v\n wanted %v", cookie.Name, "session_string")
	}
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
	wanted := wantedLoginResult{
		Code:        http.StatusOK,
		Success:     true,
		ErrKind:     apperrors.KindErrNone,
		CookieExist: true,
	}
	verifyLoginResponse(t, w, wanted)
}

func TestHandleLogin_NonExistent(t *testing.T) {
	userLogin := userLoginData{
		Email:    "non_existent@test.com",
		Password: "123",
	}
	w := handleLoginDataDecode(t, userLogin)
	wanted := wantedLoginResult{
		Code:        http.StatusUnauthorized,
		Success:     false,
		ErrKind:     apperrors.KindErrWrongLoginData,
		CookieExist: false,
	}
	verifyLoginResponse(t, w, wanted)
}

func TestHandleLogin_InvalidJSON(t *testing.T) {
	userLogin := ""
	w := handleLoginDataRaw(t, userLogin)
	wanted := wantedLoginResult{
		Code:        http.StatusBadRequest,
		Success:     false,
		ErrKind:     apperrors.KindErrInternal,
		CookieExist: false,
	}
	verifyLoginResponse(t, w, wanted)
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
	wanted := wantedLoginResult{
		Code:        http.StatusUnauthorized,
		Success:     false,
		ErrKind:     apperrors.KindErrWrongLoginData,
		CookieExist: false,
	}
	verifyLoginResponse(t, w, wanted)
}
