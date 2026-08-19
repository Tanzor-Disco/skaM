package apperrors

import (
	"errors"
)

var ErrInvalidRoomNameLength = errors.New("invalid room name length")

// KindErrors are used as part of JSON sent to frontend from server to add extra context
const (
	KindErrInvalidRoomNameLength = "ERR_INVALID_ROOM_NAME_LENGTH"
)
