package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/internal/testutils"
	"github.com/Tanzor-Disco/skaM/models"
)

func createUserAndSession(t *testing.T, user models.User) (sessionString string, userID models.UserID) {
	t.Helper()
	userID = createTestUser(t, user)
	sessionString = createTestUserSession(t, userID)
	return
}

func handleNewRoomEncode(t *testing.T, body RoomRegisterRequest, sessionString *string) *httptest.ResponseRecorder {
	t.Helper()
	r := createEncodeRequest(t, body, sessionString)
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

	sessionString, userID := createUserAndSession(t, user)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	defer tdb.DeleteSessionsByUserID(t, userID)

	reqBody := RoomRegisterRequest{
		Name: "room_test",
	}
	w := handleNewRoomEncode(t, reqBody, &sessionString)
	defer tdb.DeleteRoomByRoomName(t, reqBody.Name)
	defer tdb.DeleteRoomUsersByUserID(t, userID)
	wanted := wantedResult{
		Code:        http.StatusOK,
		ErrKind:     apperrors.KindErrNone,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleMainNewRoom_NoSession(t *testing.T) {
	user := models.NewUser("email@test.com", "test_username", testutils.GetPasswordHash(t, "123"))
	_, userID := createUserAndSession(t, user)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	defer tdb.DeleteSessionsByUserID(t, userID)

	reqBody := RoomRegisterRequest{
		Name: "room_test",
	}
	w := handleNewRoomEncode(t, reqBody, nil)
	defer tdb.DeleteRoomByRoomName(t, reqBody.Name)
	defer tdb.DeleteRoomUsersByUserID(t, userID)
	wanted := wantedResult{
		Code:        http.StatusUnauthorized,
		ErrKind:     apperrors.KindErrInvalidSessionString,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleMainNewRoom_NoJSON(t *testing.T) {
	user := models.NewUser("email@test.com", "test_username", testutils.GetPasswordHash(t, "123"))

	sessionString, userID := createUserAndSession(t, user)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	defer tdb.DeleteSessionsByUserID(t, userID)

	reqBody := ""
	w := handleNewRoomRaw(t, reqBody, &sessionString)
	defer tdb.DeleteRoomUsersByUserID(t, userID)
	wanted := wantedResult{
		Code:        http.StatusBadRequest,
		ErrKind:     apperrors.KindErrInternal,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)

}

func TestHandleMainNewRoom_LongName(t *testing.T) {
	user := models.NewUser("email@test.com", "test_username", testutils.GetPasswordHash(t, "123"))

	sessionString, userID := createUserAndSession(t, user)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	defer tdb.DeleteSessionsByUserID(t, userID)

	reqBody := RoomRegisterRequest{
		Name: "very_long_room_name_for_real_12331313131331331123213231313131313131331",
	}
	w := handleNewRoomEncode(t, reqBody, &sessionString)
	defer tdb.DeleteRoomUsersByUserID(t, userID)
	wanted := wantedResult{
		Code:        http.StatusBadRequest,
		ErrKind:     apperrors.KindErrInvalidRoomNameLength,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleMainNewRoom_ExtraFieldsJSON(t *testing.T) {
	user := models.NewUser("email@test.com", "test_username", testutils.GetPasswordHash(t, "123"))

	sessionString, userID := createUserAndSession(t, user)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	defer tdb.DeleteSessionsByUserID(t, userID)

	reqBody := `{"Name":"test room","id":123,"admin":true,"random_field":"hello"}`
	w := handleNewRoomRaw(t, reqBody, &sessionString)
	defer tdb.DeleteRoomByRoomName(t, "test room")
	defer tdb.DeleteRoomUsersByUserID(t, userID)
	wanted := wantedResult{
		Code:        http.StatusOK,
		ErrKind:     apperrors.KindErrNone,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)

}
