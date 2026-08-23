package testutils

import (
	"context"
	"testing"
)

func (tdb *TestDB) DeleteRoomByRoomID(t *testing.T, id int64) {
	_, err := tdb.pool.Exec(context.Background(),
		`
	DELETE FROM rooms WHERE id = $1
	`,
		id)
	if err != nil {
		t.Fatalf("DeleteRoomByRoomID: %v", err)
	}
}

func (tdb *TestDB) DeleteRoomsByRoomIDs(t *testing.T, roomIDs []int64) {
	for _, roomID := range roomIDs {
		tdb.DeleteRoomByRoomID(t, roomID)
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
