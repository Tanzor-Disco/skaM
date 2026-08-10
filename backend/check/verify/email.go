package verify

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func CreateToken() (string,string,error) {
	var token string 
	var tokenHashed string 

	raw := make([]byte,32)
	_,err := rand.Read(raw)
	if err != nil {
		return token,tokenHashed,err
	}
	token = base64.RawURLEncoding.EncodeToString(raw)
	sum := sha256.Sum256([]byte(token))
	tokenHashed = hex.EncodeToString(sum[:])
	return token,tokenHashed,nil
}

