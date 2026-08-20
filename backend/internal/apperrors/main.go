package apperrors

import ()

// KindErrors are used as part of JSON sent to frontend from server to add extra context
const (
	KindErrUnauthorized    = "ERR_UNAUTHORIZED"
	KindErrNoSessionString = "ERR_NO_SESSION_STRING"
)
