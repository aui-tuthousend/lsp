package domain

import "errors"

var (
	ErrNotFound = errors.New("certification schema not found")
	ErrConflict = errors.New("certification schema already exists")
)
