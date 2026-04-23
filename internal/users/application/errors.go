package application

import "errors"

var (
	ErrInvalidEmail     = errors.New("email must be a valid email address")
	ErrInvalidPassword  = errors.New("password must be at least 6 characters long")
	ErrEmailAlreadyUsed = errors.New("an account with this email already exists")
)
