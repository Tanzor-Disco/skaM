package verify

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"encoding/hex"
	"github.com/Tanzor-Disco/skaM/models"
	"net/smtp"
	"text/template"
)

func CreateToken() (string, string, error) {
	var token string
	var tokenHashed string

	raw := make([]byte, 32)
	_, err := rand.Read(raw)
	if err != nil {
		return token, tokenHashed, err
	}
	token = base64.RawURLEncoding.EncodeToString(raw)
	sum := sha256.Sum256([]byte(token))
	tokenHashed = hex.EncodeToString(sum[:])
	return token, tokenHashed, nil
}

//go:embed templates/body.html
var bodyEmbed string

func getEmailBody(baseURL string, token string) (string, error) {
	var body string
	template, err := template.New("body").Parse(bodyEmbed)
	if err != nil {
		return body, err
	}
	apiDest := baseURL + "/api/verify/email?token=" + token
	var buf bytes.Buffer
	err = template.Execute(&buf, apiDest)
	if err != nil {
		return body, err
	}
	body = buf.String()
	return body, nil
}

func SendEmail(baseURL, token, to string, data models.SMTPData) error {
	var auth smtp.Auth
	if data.Username != "" && data.Password != "" {
		auth = smtp.PlainAuth("", data.Username, data.Password, data.Host)
	}
	body, err := getEmailBody(baseURL, token)
	if err != nil {
		return err
	}
	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: Validate skaM Registration\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" + "\r\n" +
			body)
	err = smtp.SendMail(data.Addr, auth, data.From, []string{to}, msg)
	return err
}
