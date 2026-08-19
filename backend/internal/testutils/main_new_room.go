package testutils

import (
	"context"
	"testing"
)

func (tdb *TestDB) DeleteRoomUsersByUserID(t *testing.T, userID int) {
	_, err := tdb.pool.Exec(context.Background(),
		`
	DELETE FROM room_users WHERE user_id = $1
	`,
		userID)
	if err != nil {
		t.Fatalf("DeleteRoomUsersByUserID: %v", err)
	}
}
