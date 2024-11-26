package auth

import "errors"

var (
	ErrInvalidPassword    = errors.New("invalid password")
	ErrEmailExist         = errors.New("email already exist")
	ErrFailedHash         = errors.New("failed to hash password")
	ErrFailedToken        = errors.New("failed to generate token")
	ErrFailedCreate       = errors.New("failed to create user")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
