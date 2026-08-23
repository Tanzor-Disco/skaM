package testutils

import (
	"context"
	"testing"
	"log"
)

func (tdb *TestDB) DeleteRoomUsersByUserID(t *testing.T, userID int64) {
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


