package domain

import "errors"

var (
	ErrNotFound = errors.New("independent assessment not found")
	ErrConflict = errors.New("independent assessment already exists")
	ErrNotDraft = errors.New("independent assessment is not a draft")
	// ErrNotSubmitted: asesmen mandiri belum disubmit oleh asesi.
	ErrNotSubmitted = errors.New("independent assessment is not submitted")
	// ErrNotEligible: permohonan APL-01 belum berstatus SUBMITTED/APPROVED.
	ErrNotEligible = errors.New("application is not eligible for independent assessment")
	// ErrIncomplete: masih ada unit yang belum dijawab / belum ada bukti.
	ErrIncomplete   = errors.New("independent assessment is incomplete")
	ErrInvalidToken = errors.New("invalid or expired verification token")
	ErrTokenUsed    = errors.New("verification token already used")
)
