package testutils

import (
	"context"
	"log"
	"testing"
)

func (tdb *TestDB) DeleteRoomUsersByUserID(t *testing.T, userID int) {
	_, err := tdb.pool.Exec(context.Background(),
		`
	DELETE FROM room_users WHERE user_id = $1
	`,
		userID)
	if err != nil {
		log.Printf("1231")
		t.Fatalf("DeleteRoomUsersByUserID: %v", err)
	}
}

func (tdb *TestDB) DeleteRoomByRoomName(t *testing.T, roomName string) {
	_, err := tdb.pool.Exec(context.Background(),
		`
	DELETE FROM rooms WHERE room_name = $1
	`,
		roomName)
	if err != nil {
		t.Fatalf("DeleteRoomByRoomName: %v", err)
	}
}
