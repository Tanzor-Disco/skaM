package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
)

func handleUser(t *testing.T, user models.RegisterRequest) *httptest.ResponseRecorder {
	t.Helper()
	r := createDecodeRequest(t, user, nil)
	w := httptest.NewRecorder()
	srv.handleRegister(w, r)
	return w
}

func checkStatus(t *testing.T, w *httptest.ResponseRecorder, status int) {
	t.Helper()
	if w.Code != status {
		t.Fatalf("The code of the response didn't match: got %v, wanted %v", w.Code, status)
	}
}

func validateResponseBody(t *testing.T, resp *http.Response, correct serverResponseBody) {
	var body serverResponseBody
	err := json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		t.Fatalf("Failed to decode json: %v", err)
	}
	if body.Success != correct.Success || body.ErrorKind != correct.ErrorKind {
		t.Fatalf("Body of the response is incorrect:\n Got %+v\n Wanted %+v", body, correct)
	}

}

func checkResponse(t *testing.T, user models.RegisterRequest, correct serverResponseBody, wantedStatus int) {
	t.Helper()
	serverRecorder := handleUser(t, user)
	defer tdb.DeletePendingUsersByEmail(t, user.Email)
	serverResp := serverRecorder.Result()
	checkStatus(t, serverRecorder, wantedStatus)
	validateResponseBody(t, serverResp, correct)
}

func TestHandleRequest_Valid(t *testing.T) {
	user := models.RegisterRequest{
		Email:    "pidoras_@mail.ru",
		Username: "pidoras",
		Password: "2313112313",
	}

	correct := newServerResponseBody(true, apperrors.KindErrNone)
	wantedStatus := http.StatusOK
	checkResponse(t, user, correct, wantedStatus)
}

func TestHandleRequest_Empty(t *testing.T) {
	user := models.RegisterRequest{
		Email:    "",
		Username: "",
		Password: "",
	}
	correct := newServerResponseBody(false, apperrors.KindErrInvalidEmail)
	serverRecorder := handleUser(t, user)
	defer tdb.DeletePendingUsersByEmail(t, user.Email)
	serverResp := serverRecorder.Result()
	checkStatus(t, serverRecorder, http.StatusBadRequest)
	validateResponseBody(t, serverResp, correct)
}

func TestHandleRequest_InvalidEmail(t *testing.T) {
	user := models.RegisterRequest{
		Email:    "huylo",
		Username: "testname",
		Password: "123312123",
	}
	correct := newServerResponseBody(false, apperrors.KindErrInvalidEmail)
	wantedStatus := http.StatusBadRequest
	checkResponse(t, user, correct, wantedStatus)
}

func TestHandleRequest_Duplicate(t *testing.T) {
	userPrev := models.RegisterRequest{
		Email:    "pidoras_@mail.ru",
		Username: "pidoras",
		Password: "2313112313",
	}
	user := models.RegisterRequest{
		Email:    "pidoras_@mail.ru",
		Username: "pidoras_new",
		Password: "1321313213131",
	}
	correct := newServerResponseBody(true, apperrors.KindErrNone)
	//creating the first user
	handleUser(t, userPrev)
	//creating a duplicate
	serverRecorder := handleUser(t, user)
	defer tdb.DeletePendingUsersByEmail(t, user.Email)
	serverResp := serverRecorder.Result()
	checkStatus(t, serverRecorder, http.StatusOK)
	validateResponseBody(t, serverResp, correct)
}

func TestHandleRequest_EmailTaken(t *testing.T) {
	userRequest := models.RegisterRequest{
		Email:    "pidoras_@mail.ru",
		Username: "pidoras_new",
		Password: "1321313213131",
	}
	user := models.User{
		Email:          "pidoras_@mail.ru",
		Username:       "pidoras_new",
		PasswordHash:   "1321313213131",
		LastYearActive: 2026,
	}
	srv.db.CreateUser(context.Background(), user)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	correct := newServerResponseBody(false, apperrors.KindErrEmailTaken)
	serverRecorder := handleUser(t, userRequest)
	defer tdb.DeletePendingUsersByEmail(t, userRequest.Email)
	serverResp := serverRecorder.Result()
	checkStatus(t, serverRecorder, http.StatusBadRequest)
	validateResponseBody(t, serverResp, correct)
}

func TestHandleRequest_InvalidUsernameLength(t *testing.T) {
	user := models.RegisterRequest{
		Email:    "example_@mail.ru",
		Username: "125151512515353255151551142144241241",
		Password: "12321331123",
	}
	correct := newServerResponseBody(false, apperrors.KindErrInvalidUsernameLength)
	wantedStatus := http.StatusBadRequest
	checkResponse(t, user, correct, wantedStatus)
}

func TestHandleRequest_InvalidPasswordChars(t *testing.T) {
	user := models.RegisterRequest{
		Email:    "example@mail.ru",
		Username: "huylo",
		Password: "фывфывфйцуйцуол",
	}
	correct := newServerResponseBody(false, apperrors.KindErrForbiddenPasswordChars)
	wantedStatus := http.StatusBadRequest
	checkResponse(t, user, correct, wantedStatus)
}

func TestHandleRequest_InvalidPasswordLength(t *testing.T) {
	user := models.RegisterRequest{
		Email:    "example@mail.ru",
		Username: "huylo",
		Password: "12313131331313113131313132131313131313131313131313131313131331313131313131313131313131313",
	}
	correct := newServerResponseBody(false, apperrors.KindErrInvalidPasswordLength)
	wantedStatus := http.StatusBadRequest
	checkResponse(t, user, correct, wantedStatus)
}
