package domain

import "errors"

var (
	ErrNotFound     = errors.New("certification application not found")
	ErrConflict     = errors.New("certification application already exists")
	ErrNotDraft     = errors.New("certification application is not a draft")
	ErrIncomplete   = errors.New("certification application is incomplete")
	ErrInvalidToken = errors.New("invalid or expired verification token")
	ErrTokenUsed    = errors.New("verification token already used")
)
