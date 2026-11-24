package domain

import "errors"

var (
	ErrNotFound = errors.New("document requirement not found")
	ErrConflict = errors.New("document requirement already exists")
)
