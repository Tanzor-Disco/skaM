package apperrors

import (
	"errors"
)

var ErrEmailTaken = errors.New("Error: email is already taken")

