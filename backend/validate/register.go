package validate

import (
	"net/mail"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
)

func validateEmail(address string) error {
	addr,err := mail.ParseAddress(address)
	if err == nil  && addr.Address == address {
		return nil
	}
	return apperrors.ErrInvalidEmail
}

func isLatin(char rune) bool {
	if (char >= 'a' && char <= 'z') || (char >='A' && char <= 'Z') {
		return true
	}
	return false
}

func isNumber(char rune) bool {
	if char > '0' && char <= '9' {
		return true
	}
	return false
}

func checkPasswordFormat(password string) error {
	for _,char := range password {
		if isLatin(char) || isNumber(char) {
			continue
		}
		return apperrors.ErrInvalidPasswordChars
	}
	return nil
}

func validatePassword(password string) error {
	//bcrypt limits the size of the string to 72 bytes; with latin character and numbers it's 72 characters
	//took a little less just to be safe
	if len(password) == 0 || len(password) > 70 {
		return apperrors.ErrInvalidPasswordLength
	}

	err := checkPasswordFormat(password)
	if err != nil {
		return err
	}

	return nil
}

func validateUsername(username string) error {
	//database username has type varchar 20 which limits the max length to 30
	if len(username) >= 0 && len(username) <= 20 {
		return nil
	} 
	return apperrors.ErrInvalidUsernameLength
}

func RegisterRequest(request models.RegisterRequest) error {
	err := validateEmail(request.Email) 
	if err != nil {
		return err
	}
	err = validatePassword(request.Password)
	if err != nil {
		return err
	}
	err = validateUsername(request.Username)
	if err != nil {
		return err
	}
	return nil
}
