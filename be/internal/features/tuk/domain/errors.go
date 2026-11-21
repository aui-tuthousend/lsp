package domain

import "errors"

var (
	ErrNotFound = errors.New("tuk not found")
	ErrConflict = errors.New("tuk already exists")
)
