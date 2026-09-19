package apperrors

import (
	"errors"
)

var ErrInvalidRoomNameLength = errors.New("invalid room name length")
var ErrInvalidEmail = errors.New("invalid email format")
var ErrInvalidEmailLength = errors.New("invalid email length")
var ErrEmailTaken = errors.New("email already taken")
var ErrInvalidJSON = errors.New("invalid JSON")
var ErrInvalidPasswordChars = errors.New("forbidden password characters")
var ErrInvalidPasswordLength = errors.New("invalid password length")
var ErrInvalidUsernameLength = errors.New("invalid username length")
var ErrTokenCreation = errors.New("Failed to create a token")
var ErrUniqueViolation = errors.New("SQL unique constraints violation")

// KindErrors are used as part of JSON sent to frontend from server to add extra context
const (
	KindErrInternal               = "ERR_INTERNAL"
	KindErrWrongLoginData         = "ERR_WRONG_LOGIN_DATA"
	KindErrUnauthorized           = "ERR_UNAUTHORIZED"
	KindErrInvalidSessionString   = "ERR_INVALID_SESSION_STRING"
	KindErrInvalidRoomNameLength  = "ERR_INVALID_ROOM_NAME_LENGTH"
	KindErrEmailTaken             = "ERR_EMAIL_TAKEN"
	KindErrInvalidEmail           = "ERR_INVALID_EMAIL"
	KindErrInvalidEmailLength     = "ERR_INVALID_EMAIL_LENGTH"
	KindErrInvalidUsernameLength  = "ERR_INVALID_USERNAME_LENGTH"
	KindErrForbiddenPasswordChars = "ERR_FORBIDDEN_PASSWORD_CHARS"
	KindErrInvalidPasswordLength  = "ERR_INVALID_PASSWORD_LENGTH"
	KindErrNone                   = "ERR_NONE"
	KindErrInvalidQueryParam      = "ERR_INVALID_QUERY_PARAMETER"
	KindErrInvalidJSON            = "ERR_INVALID_JSON"
	KindErrInvalidInviteToken     = "ERR_INVALID_INVITE_TOKEN"
	KindErrUniqueViolation        = "ERR_UNIQUE_VIOLATION"
	KindErrWebsocketUpgrade       = "ERR_WEBSOCKET_FAILED_UPGRADE"
)
