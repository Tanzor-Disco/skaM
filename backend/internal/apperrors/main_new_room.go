package apperrors

import (
	"errors"
)

var ErrInvalidRoomNameLength = errors.New("invalid room name length")

const (
	KindErrInvalidRoomNameLength = "ERR_INVALID_ROOM_NAME_LENGTH"
)
