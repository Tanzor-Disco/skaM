package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tanzor-Disco/skaM/auth"
	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/internal/testutils"
	"github.com/Tanzor-Disco/skaM/models"
)

func createUserAndSession(t *testing.T, user models.User) (userEmail, sessionString string, userID int) {
	t.Helper()
	err := database.CreateUser(context.Background(), user)
	userID = tdb.GetUserIDByEmail(t, user.Email)
	userEmail = user.Email
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	sessionString, err = auth.CreateSessionString()
	if err != nil {
		t.Fatalf("CreateSessionString: %v", err)
	}
	userSession := db.NewUserSession(userID, sessionString)
	err = database.CreateSession(context.Background(), userSession)
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	return
}

func handleNewRoomDecode(t *testing.T, body RoomRegisterRequest, sessionString *string) *httptest.ResponseRecorder {
	t.Helper()
	r := createDecodeRequest(t, body, sessionString)
	w := httptest.NewRecorder()
	srv.handleMainNewRoom(w, r)
	return w
}

func handleNewRoomRaw(t *testing.T, body string, sessionString *string) *httptest.ResponseRecorder {
	t.Helper()
	r := createRawRequest(t, body, sessionString)
	w := httptest.NewRecorder()
	srv.handleMainNewRoom(w, r)
	return w
}

func TestHandleMainNewRoom_Valid(t *testing.T) {
	user := models.NewUser("email@test.com", "test_username", testutils.GetPasswordHash(t, "123"))

	userEmail, sessionString, userID := createUserAndSession(t, user)
	defer tdb.DeleteUsersByEmail(t, userEmail)
	defer tdb.DeleteSessionsByUserID(t, userID)

	reqBody := RoomRegisterRequest{
		Name: "room_test",
	}
	w := handleNewRoomDecode(t, reqBody, &sessionString)
	defer tdb.DeleteRoomUsersByUserID(t, userID)
	wanted := wantedResult{
		Code:        http.StatusOK,
		Success:     true,
		ErrKind:     apperrors.KindErrNone,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleMainNewRoom_NoSession(t *testing.T) {
	user := models.NewUser("email@test.com", "test_username", testutils.GetPasswordHash(t, "123"))
	userEmail, _, userID := createUserAndSession(t, user)
	defer tdb.DeleteUsersByEmail(t, userEmail)
	defer tdb.DeleteSessionsByUserID(t, userID)

	reqBody := RoomRegisterRequest{
		Name: "room_test",
	}
	w := handleNewRoomDecode(t, reqBody, nil)
	defer tdb.DeleteRoomUsersByUserID(t, userID)
	wanted := wantedResult{
		Code:        http.StatusUnauthorized,
		Success:     false,
		ErrKind:     apperrors.KindErrInternal,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleMainNewRoom_NoJSON(t *testing.T) {
	user := models.NewUser("email@test.com", "test_username", testutils.GetPasswordHash(t, "123"))

	userEmail, sessionString, userID := createUserAndSession(t, user)
	defer tdb.DeleteUsersByEmail(t, userEmail)
	defer tdb.DeleteSessionsByUserID(t, userID)

	reqBody := ""
	w := handleNewRoomRaw(t, reqBody, &sessionString)
	defer tdb.DeleteRoomUsersByUserID(t, userID)
	wanted := wantedResult{
		Code:        http.StatusBadRequest,
		Success:     false,
		ErrKind:     apperrors.KindErrInternal,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)

}

func TestHandleMainNewRoom_LongName(t *testing.T) {
	user := models.NewUser("email@test.com", "test_username", testutils.GetPasswordHash(t, "123"))

	userEmail, sessionString, userID := createUserAndSession(t, user)
	defer tdb.DeleteUsersByEmail(t, userEmail)
	defer tdb.DeleteSessionsByUserID(t, userID)

	reqBody := RoomRegisterRequest{
		Name: "very_long_room_name_for_real_12331313131331331123213231313131313131331",
	}
	w := handleNewRoomDecode(t, reqBody, &sessionString)
	defer tdb.DeleteRoomUsersByUserID(t, userID)
	wanted := wantedResult{
		Code:        http.StatusBadRequest,
		Success:     false,
		ErrKind:     apperrors.KindErrInvalidRoomNameLength,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleMainNewRoom_ExtraFieldsJSON(t *testing.T) {
	user := models.NewUser("email@test.com", "test_username", testutils.GetPasswordHash(t, "123"))

	userEmail, sessionString, userID := createUserAndSession(t, user)
	defer tdb.DeleteUsersByEmail(t, userEmail)
	defer tdb.DeleteSessionsByUserID(t, userID)

	reqBody := `{"Name":"test room","id":123,"admin":true,"random_field":"hello"}`
	w := handleNewRoomRaw(t, reqBody, &sessionString)
	defer tdb.DeleteRoomUsersByUserID(t, userID)
	wanted := wantedResult{
		Code:        http.StatusOK,
		Success:     true,
		ErrKind:     apperrors.KindErrNone,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)

}
