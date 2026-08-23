package db

import (
	"context"
	"fmt"
)

// Room represents a row in rooms table
// It can be ised both for inserting or retrieving info from rooms table
type Room struct {
	ID   int
	Name string
}

func NewRoom(name string) Room {
	return Room{
		Name: name,
	}
}

func (db *DB) CreateRoom(ctx context.Context, room Room) (int64, error) {
	var roomID int64
	err := db.pool.QueryRow(context.Background(),
		`
	INSERT INTO rooms (room_name)
	VALUES ($1)
	RETURNING id
	`,
		room.Name,
	).Scan(&roomID)
	return roomID, err
}

func (db *DB) getRoomByRoomID(ctx context.Context, roomID int64) (Room, error) {
	var room Room
	err := db.pool.QueryRow(ctx,
		`
	SELECT * FROM rooms WHERE id = $1
	`,
		roomID).Scan(&room.ID, &room.Name)
	if err != nil {
		return room, fmt.Errorf("getRoomByRoomID: %w", err)
	}
	return room, nil
}

// GetRoomsByRoomIDS recieves a slice of integers
// And returns a slice of all the rooms in rooms table that match any of the ids
func (db *DB) GetRoomsByRoomIDS(ctx context.Context, roomIDS []int64) ([]Room, error) {
	var rooms []Room
	for _, roomID := range roomIDS {
		room, err := db.getRoomByRoomID(ctx, roomID)
		if err != nil {
			return make([]Room, 0), fmt.Errorf("GetRoomsByRoomIDS: %w", err)
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}
