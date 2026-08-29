package db

import (
	"context"
	"fmt"

	"github.com/Tanzor-Disco/skaM/models"
)

func (db *DB) CreateRoom(ctx context.Context, room models.Room) (models.RoomID, error) {
	var roomID models.RoomID
	err := db.pool.QueryRow(ctx,
		`
	INSERT INTO rooms (room_name)
	VALUES ($1)
	RETURNING id
	`,
		room.Name,
	).Scan(&roomID)

	if err != nil {
		return roomID, fmt.Errorf("CreateRoom: %w", err)
	}
	return roomID, nil
}

func (db *DB) GetRoomByRoomID(ctx context.Context, roomID models.RoomID) (models.Room, error) {
	var room models.Room
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

// GetRoomsByRoomIDS recieves a slice of room ids
// And returns a slice of all the rooms in rooms table that match any of the ids
func (db *DB) GetRoomsByRoomIDS(ctx context.Context, roomIDS []models.RoomID) ([]models.Room, error) {
	var rooms []models.Room
	for _, roomID := range roomIDS {
		room, err := db.GetRoomByRoomID(ctx, roomID)
		if err != nil {
			return make([]models.Room, 0), fmt.Errorf("GetRoomsByRoomIDS: %w", err)
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}
