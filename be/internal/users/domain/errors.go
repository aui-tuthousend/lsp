package domain

import "errors"

var (
	ErrNotFound           = errors.New("user not found")
	ErrConflict           = errors.New("user conflict")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInactiveUser       = errors.New("inactive user")
	ErrInvalidRole        = errors.New("invalid user role")
)
