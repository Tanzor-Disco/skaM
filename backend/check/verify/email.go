package verify

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/Tanzor-Disco/skaM/models"
)

// CreateToken generates a unique token for email confirmation.
// It returns the raw token for the confirmation URL and its SHA-256 hash for database storage in pending_users.
func CreateToken() (string, string, error) {
	var token string
	var tokenHashed string

	raw := make([]byte, 32)
	_, err := rand.Read(raw)
	if err != nil {
		return token, tokenHashed, fmt.Errorf("CreateToken: %w", err)
	}
	token = base64.RawURLEncoding.EncodeToString(raw)
	sum := sha256.Sum256([]byte(token))
	tokenHashed = hex.EncodeToString(sum[:])
	return token, tokenHashed, nil
}

//go:embed templates/body.html
var bodyEmbed string

// getEmailBody uses embedded body.html to create a body for go SMTP module
// It replaces {{.}} in body.html with a link to an endpoint that contains a unique token
func getEmailBody(baseURL string, token string) (string, error) {
	var body string
	template, err := template.New("body").Parse(bodyEmbed)
	if err != nil {
		return body, fmt.Errorf("getEmailBody: template.New: %w", err)
	}
	apiDest := baseURL + "/api/verify/email?token=" + token
	var buf bytes.Buffer
	err = template.Execute(&buf, apiDest)
	if err != nil {
		return body, fmt.Errorf("getEmailBody: template.Execute: %w", err)
	}
	body = buf.String()
	return body, nil
}

// SendEmail uses go SMTP module to send a verification email to a specified user
func SendEmail(token, to, baseURL string, data models.SMTPData) error {
	var auth smtp.Auth
	if data.Username != "" && data.Password != "" {
		auth = smtp.PlainAuth("", data.Username, data.Password, data.Host)
	}
	body, err := getEmailBody(baseURL, token)
	if err != nil {
		return fmt.Errorf("sendEmail: %w", err)
	}
	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: Validate skaM Registration\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" + "\r\n" +
			body)
	err = smtp.SendMail(data.Addr, auth, data.From, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("sendEmail: %w", err)
	}
	return nil
}
