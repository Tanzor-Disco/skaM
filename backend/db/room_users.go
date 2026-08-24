package db

import (
	"context"
	"fmt"
)

type RoomUser struct {
	RoomID int64
	UserID int64
}

func NewRoomUser(roomID, userID int64) RoomUser {
	return RoomUser{
		RoomID: roomID,
		UserID: userID,
	}
}

func (db *DB) AddUserToRoom(ctx context.Context, roomUser RoomUser) error {
	_, err := db.pool.Exec(ctx,
		`
	INSERT INTO room_users (user_id,room_id)
	VALUES ($1,$2)
	`,
		roomUser.UserID, roomUser.RoomID,
	)
	return err
}

func (db *DB) GetRoomIDSByUserID(ctx context.Context, userID int64) ([]int64, error) {
	rows, err := db.pool.Query(ctx,
		`
	SELECT room_id FROM room_users WHERE user_id = $1
	`,
		userID)
	if err != nil {
		return make([]int64, 0), fmt.Errorf("GetRoomIDSByUserID: %w", err)
	}
	var roomIDS []int64
	for rows.Next() {
		var roomID int64
		err := rows.Scan(&roomID)
		if err != nil {
			return make([]int64, 0), fmt.Errorf("GetRoomIDSByUserID: %w", err)
		}
		roomIDS = append(roomIDS, roomID)
	}
	return roomIDS, nil
}

func (db *DB) ValidateRoomUser(ctx context.Context, userID, roomID int64) error {
	var exists int
	err := db.pool.QueryRow(ctx,
		`
	SELECT 1 FROM room_users 
	WHERE user_id = $1 AND room_id = $2
	`,
		userID, roomID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ValidateRoomUser: %w", err)
	}
	return nil
}
