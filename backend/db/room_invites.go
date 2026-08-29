package db

import (
	"context"
	"fmt"

	"github.com/Tanzor-Disco/skaM/models"
)

func (db *DB) CreateRoomInvite(ctx context.Context, invite models.RoomInvite) error {
	_, err := db.pool.Exec(ctx,
		`
	INSERT INTO room_invites (room_id,token)
	VALUES($1,$2)
	`,
		invite.RoomID, invite.Token)
	if err != nil {
		return fmt.Errorf("CreateRoomInvite: %w", err)
	}
	return nil
}

func (db *DB) GetRoomInviteByroomID(ctx context.Context, roomID models.RoomID) (models.RoomInvite, error) {
	var roomInvite models.RoomInvite
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

func (db *DB) GetRoomInviteByInviteToken(ctx context.Context, inviteToken string) (models.RoomInvite, error) {
	var roomInvite models.RoomInvite
	err := db.pool.QueryRow(ctx,
		`
	SELECT * FROM room_invites
	WHERE token = $1
	`,
		inviteToken).Scan(&roomInvite.ID, &roomInvite.RoomID, &roomInvite.Token)
	if err != nil {
		return roomInvite, fmt.Errorf("GetRoomInviteByroomID: queryRow: %w", err)
	}
	return roomInvite, nil
}
