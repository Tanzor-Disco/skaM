package validate

import (
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
)

func RoomNameLength(roomName string) error {
	if len(roomName) > 0 && len(roomName) < 21 {
		return nil
	} 
	return apperrors.ErrInvalidRoomNameLength
}
