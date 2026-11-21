package domain

import "errors"

var (
	ErrNotFound        = errors.New("asesi profile not found")
	ErrConflict        = errors.New("asesi profile already exists")
	ErrNotAsesi        = errors.New("user is not an asesi")
	ErrInvalidDomicile = errors.New("domicile_address is required when domicile_same_as_ktp is false")
)
