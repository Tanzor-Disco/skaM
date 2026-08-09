package apperrors

import (
	"errors"
)

var ErrInvalidEmail = errors.New("invalid email format")
var ErrInvalidJSON = errors.New("invalid JSON")
var ErrInvalidPasswordChars = errors.New("forbidden password characters")
var ErrInvalidPasswordLength = errors.New("invalid password length")
var ErrInvalidUsernameLength = errors.New("invalid username length")

//error kinds for http responses
const (
	KindErrHashingPassword = "ERR_HASHING_PASSWORD"
	KindErrEmailTaken = "ERR_EMAIL_TAKEN"
	KindErrDB = "ERR_DB"
	KindErrInvalidEmail = "ERR_INVALID_EMAIL"
	KindErrInvalidUsernameLength = "ERR_INVALID_USERNAME_LENGTH"
	KindErrForbiddenPasswordChars = "ERR_FORBIDDEN_PASSWORD_CHARS"
	KindErrInvalidPasswordLength = "ERR_INVALID_PASSWORD_LENGTH"
	KindErrNone = "ERR_NONE"
)
