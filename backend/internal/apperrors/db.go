package apperrors

import (
	"errors"
)

var ErrEmailTaken = errors.New("email is already taken")
var ErrUserNotFound = errors.New("no rows found")
var ErrDB = errors.New("db error")
