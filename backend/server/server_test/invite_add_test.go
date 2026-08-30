package servertest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
	"github.com/Tanzor-Disco/skaM/server/invite"
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
	currInvite := models.NewRoomInvite(roomID, "123")
	createTestRoomInvite(t, currInvite)
	r := createEncodeRequest(t, currInvite, &sessionString)
	w = httptest.NewRecorder()

	inviteHandler := invite.NewInviteHandler(srv.DB)
	inviteHandler.HandleInviteAdd(w, r)
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
		ErrKind:     apperrors.KindErrInvalidSessionString,
		CookieExist: false,
	}

	verifyResponse(t, w, wanted)
}

func TestInviteAdd_Duplicate(t *testing.T) {
	user := models.NewUser("testing@mail.com", "test_user", "123")
	room := models.NewRoom("test_room")

	userID := createTestUser(t, user)
	roomID := createTestRoom(t, room)
	roomUser := models.NewRoomUser(roomID, userID)
	sessionString := createTestUserSession(t, userID)

	createTestRoomUser(t, roomUser)

	currInvite := models.NewRoomInvite(roomID, "123")
	createTestRoomInvite(t, currInvite)
	r := createEncodeRequest(t, currInvite, &sessionString)
	w := httptest.NewRecorder()

	inviteHandler := invite.NewInviteHandler(srv.DB)
	inviteHandler.HandleInviteAdd(w, r)

	defer tdb.DeleteUsersByEmail(t, user.Email)
	defer tdb.DeleteRoomByRoomID(t, roomID)
	defer tdb.DeleteSessionsByUserID(t, userID)
	defer tdb.DeleteRoomUsersByUserID(t, userID)
	defer tdb.DeleteRoomInvitesByRoomID(t, roomID)

	wanted := wantedResult{
		Code:        http.StatusInternalServerError,
		ErrKind:     apperrors.KindErrUniqueViolation,
		CookieExist: false,
	}

	verifyResponse(t, w, wanted)
}
