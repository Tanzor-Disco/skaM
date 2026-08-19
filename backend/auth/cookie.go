package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
)

// CreateSessionString generates a unique token that is used in cookies and stored in user_sessions db
func CreateSessionString() (string, error) {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		return "", fmt.Errorf("createSession: %w", err)
	}
	sessionString := hex.EncodeToString(buf)
	return sessionString, nil
}

func CreateSessionCookie(sessionString string) http.Cookie {
	return http.Cookie{
		Name:     "session_string",
		Value:    sessionString,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7,
		HttpOnly: true,
		Secure:   false, //TODO: change to secure in production
	}
}
