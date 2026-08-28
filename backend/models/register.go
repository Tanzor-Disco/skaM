package models

import ()

// RegisterRequest represents data passed from user to server after registration form submit
type RegisterRequest struct {
	Email    string
	Username string
	Password string
}
