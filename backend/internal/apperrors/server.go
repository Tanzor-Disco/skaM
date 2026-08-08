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
	KindErrRequestBodyRead = "ERR_REQUEST_BODY_READ"
	KindErrInvalidJSON = "ERR_INVALID_JSON"
	KindErrHashingPassword = "ERR_HASHING_PASSWORD"
	KindErrEmailTaken = "ERR_EMAIL_TAKEN"
	KindErrDB = "ERR_DB"
	KindErrNone = "ERR_NONE"
	KindErrInvalidEmail = "ERR_INVALID_EMAIL"
)
