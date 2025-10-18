package users

import "errors"

var (
	ErrDuplicateEmail     = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
)
