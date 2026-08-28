package testutils

import (
	"context"
	"testing"

	"github.com/Tanzor-Disco/skaM/models"
)

func (tdb *TestDB) DeleteRoomInvitesByRoomID(t *testing.T, RoomID models.RoomID) {
	_, err := tdb.pool.Exec(context.Background(),
		`
	DELETE FROM room_invites WHERE room_id = $1
	`,
		RoomID)
	if err != nil {
		t.Fatalf("DeleteRoomInvitesByUserID: %v", err)
	}
}
