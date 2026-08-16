package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
)

func CreateSessionID() (string, error) {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		return "", fmt.Errorf("createSession: %w", err)
	}
	sessionID := hex.EncodeToString(buf)
	return sessionID, nil
}

func CreateSessionCookie(sessionID string) http.Cookie {
	return http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7,
		HttpOnly: true,
		Secure:   false, //TODO: change to secure in production
	}
}
