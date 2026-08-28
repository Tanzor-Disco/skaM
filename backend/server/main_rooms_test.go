package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tanzor-Disco/skaM/auth"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
)

func createUserSessionRooms(t *testing.T, user models.User, rooms []models.Room) (userID models.UserID, sessionString string, roomIDs []models.RoomID) {
	t.Helper()
	ctx := context.Background()
	err := database.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("createUserSessionRooms: CreateUser: %v", err)
	}
	userID = tdb.GetUserIDByEmail(t, user.Email)
	sessionString, err = auth.CreateSessionString()
	if err != nil {
		t.Fatalf("createUserSessionRooms: %v", err)
	}
	userSession := models.NewUserSession(userID, sessionString)
	err = database.CreateSession(ctx, userSession)
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	for _, room := range rooms {
		roomID, err := database.CreateRoom(ctx, room)
		if err != nil {
			t.Fatalf("createUserSessionRooms: CreateRoom: %v", err)
		}
		roomIDs = append(roomIDs, roomID)
	}
	for _, roomID := range roomIDs {
		roomUser := models.NewRoomUser(roomID, userID)
		database.AddUserToRoom(ctx, roomUser)
	}
	return
}

func handleRooms(t *testing.T, sessionString *string) *httptest.ResponseRecorder {
	r := createRawRequest(t, "", sessionString)
	w := httptest.NewRecorder()
	srv.handleMainRooms(w, r)
	return w
}

func TestHandleMainRooms_Valid(t *testing.T) {
	rooms := []models.Room{
		{ID: 0, Name: "test_room_1"},
		{ID: 0, Name: "test_room_2"},
		{ID: 0, Name: "test_room_3"}}
	user := models.NewUser("email@test.com", "test_user", "123")
	userID, sessionString, roomIDs := createUserSessionRooms(t, user, rooms)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	defer tdb.DeleteRoomsByRoomIDs(t, roomIDs)
	defer tdb.DeleteRoomUsersByUserID(t, userID)
	defer tdb.DeleteSessionsByUserID(t, userID)
	w := handleRooms(t, &sessionString)
	wanted := wantedResult{
		Code:        http.StatusOK,
		Success:     true,
		ErrKind:     apperrors.KindErrNone,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleMainRooms_No_SessionString(t *testing.T) {
	w := handleRooms(t, nil)
	wanted := wantedResult{
		Code:        http.StatusUnauthorized,
		Success:     false,
		ErrKind:     apperrors.KindErrInvalidSessionString,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleMainRooms_No_Rooms(t *testing.T) {
	rooms := []models.Room{}
	user := models.NewUser("email@test.com", "test_user", "123")
	userID, sessionString, roomIDs := createUserSessionRooms(t, user, rooms)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	defer tdb.DeleteRoomsByRoomIDs(t, roomIDs)
	defer tdb.DeleteRoomUsersByUserID(t, userID)
	defer tdb.DeleteSessionsByUserID(t, userID)
	w := handleRooms(t, &sessionString)
	wanted := wantedResult{
		Code:        http.StatusOK,
		Success:     true,
		ErrKind:     apperrors.KindErrNone,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}
