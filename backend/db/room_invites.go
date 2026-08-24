package db

import (
	"context"
	"fmt"
)

type RoomInvite struct {
	ID     int64
	RoomID int64
	Token  string
}

func NewRoomInvite(roomID int64, token string) RoomInvite {
	return RoomInvite{
		RoomID: roomID,
		Token:  token,
	}
}

func (db *DB) CreateRoomInvite(ctx context.Context, invite RoomInvite) error {
	_, err := db.pool.Exec(ctx,
		`
	INSERT INTO room_invites (room_id,token)
	VALUES($1,$2)
	`,
		invite.RoomID, invite.Token)
	if err != nil {
		return fmt.Errorf("createRoomInvite: %w", err)
	}
	return nil
}

func (db *DB) GetRoomInviteByroomID(ctx context.Context, roomID int64) (RoomInvite, error) {
	var roomInvite RoomInvite
	err := db.pool.QueryRow(ctx,
		`
	SELECT * FROM room_invites
	WHERE room_id = $1
	`,
		roomID).Scan(&roomInvite.ID, &roomInvite.RoomID, &roomInvite.Token)
	if err != nil {
		return roomInvite, fmt.Errorf("getRoomInviteByroomID: queryRow: %w", err)
	}
	return roomInvite, nil
}
