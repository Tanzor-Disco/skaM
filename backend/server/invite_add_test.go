package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
)

func handleInvite(t *testing.T, user *models.User, room *models.Room) (userID models.UserID, sessionString string, roomID models.RoomID, w *httptest.ResponseRecorder) {
	t.Helper()
	if user != nil {
		userID = createTestUser(t, *user)
		sessionString = createTestUserSession(t, userID)
	}
	if room != nil {
		roomID = createTestRoom(t, *room)
	}
	invite := models.NewRoomInvite(roomID, "123")
	createTestRoomInvite(t, invite)
	r := createEncodeRequest(t, invite, &sessionString)
	w = httptest.NewRecorder()
	srv.handleInviteAdd(w, r)
	return
}

func TestInviteAdd_Valid(t *testing.T) {
	user := models.NewUser("testing@mail.com", "test_user", "123")
	room := models.NewRoom("test_room")
	userID, _, roomID, w := handleInvite(t, &user, &room)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	defer tdb.DeleteRoomByRoomID(t, roomID)
	defer tdb.DeleteSessionsByUserID(t, userID)
	defer tdb.DeleteRoomUsersByUserID(t, userID)
	defer tdb.DeleteRoomInvitesByRoomID(t, roomID)

	wanted := wantedResult{
		Code:        http.StatusOK,
		Success:     true,
		ErrKind:     apperrors.KindErrNone,
		CookieExist: false,
	}

	verifyResponse(t, w, wanted)
}

func TestInviteAdd_No_user(t *testing.T) {
	room := models.NewRoom("test_room")
	userID, _, roomID, w := handleInvite(t, nil, &room)
	defer tdb.DeleteRoomByRoomID(t, roomID)
	defer tdb.DeleteSessionsByUserID(t, userID)
	defer tdb.DeleteRoomUsersByUserID(t, userID)
	defer tdb.DeleteRoomInvitesByRoomID(t, roomID)

	wanted := wantedResult{
		Code:        http.StatusUnauthorized,
		Success:     false,
		ErrKind:     apperrors.KindErrInvalidSessionString,
		CookieExist: false,
	}

	verifyResponse(t, w, wanted)
}
