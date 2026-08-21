package apperrors

import (
	"errors"
)

var ErrUserNotFound = errors.New("no rows found")
var ErrDB = errors.New("db error")
