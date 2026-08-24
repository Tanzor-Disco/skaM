package invite

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func CreateInviteToken() (string, error) {
	var token string

	raw := make([]byte, 32)
	_, err := rand.Read(raw)
	if err != nil {
		return token, fmt.Errorf("CreateInviteToken: rand.Read: %w", err)
	}
	token = base64.RawURLEncoding.EncodeToString(raw)
	return token, nil
}
