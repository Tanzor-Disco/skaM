package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
	"github.com/jackc/pgx/v5/pgconn"
)

func (db *DB) AddUserToRoom(ctx context.Context, roomUser models.RoomUser) error {
	_, err := db.pool.Exec(ctx,
		`
	INSERT INTO room_users (user_id,room_id)
	VALUES ($1,$2)
	`,
		roomUser.UserID, roomUser.RoomID,
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return apperrors.ErrUniqueViolation
	}
	return err
}

func (db *DB) GetRoomIDSByUserID(ctx context.Context, userID models.UserID) ([]models.RoomID, error) {
	rows, err := db.pool.Query(ctx,
		`
	SELECT room_id FROM room_users WHERE user_id = $1
	`,
		userID)
	if err != nil {
		return make([]models.RoomID, 0), fmt.Errorf("GetRoomIDSByUserID: %w", err)
	}
	var roomIDS []models.RoomID
	for rows.Next() {
		var roomID models.RoomID
		err := rows.Scan(&roomID)
		if err != nil {
			return make([]models.RoomID, 0), fmt.Errorf("GetRoomIDSByUserID: %w", err)
		}
		roomIDS = append(roomIDS, roomID)
	}
	return roomIDS, nil
}

func (db *DB) ValidateRoomUser(ctx context.Context, userID models.UserID, roomID models.RoomID) error {
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
