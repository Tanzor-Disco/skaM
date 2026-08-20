package db

import (
	"context"
	"fmt"
)

type RoomUser struct {
	RoomID int
	UserID int
}

func NewRoomUser(roomID, userID int) RoomUser {
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

func (db *DB) GetRoomIDSByUserID(ctx context.Context, userID int) ([]int, error) {
	rows, err := db.pool.Query(ctx,
		`
	SELECT room_id FROM room_users WHERE user_id = $1
	`,
		userID)
	if err != nil {
		return make([]int, 0), fmt.Errorf("GetRoomIDSByUserID: %w", err)
	}
	var roomIDS []int
	for rows.Next() {
		var roomID int
		err := rows.Scan(&roomID)
		if err != nil {
			return make([]int, 0), fmt.Errorf("GetRoomIDSByUserID: %w", err)
		}
		roomIDS = append(roomIDS, roomID)
	}
	return roomIDS, nil
}
