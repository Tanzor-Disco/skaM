package testutils

import (
	"context"
	"testing"
)

func (tdb *TestDB) DeleteRoomByRoomID(t *testing.T, id int) {
	_, err := tdb.pool.Exec(context.Background(),
		`
	DELETE FROM rooms WHERE id = $1
	`,
		id)
	if err != nil {
		t.Fatalf("DeleteRoomByRoomID: %v", err)
	}
}

func (tdb *TestDB) DeleteRoomsByRoomIDs(t *testing.T, roomIDs []int) {
	for _, roomID := range roomIDs {
		tdb.DeleteRoomByRoomID(t, roomID)
	}
}
