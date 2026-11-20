package domain

import "errors"

var (
	ErrNotFound     = errors.New("asesor profile not found")
	ErrNotAsesor    = errors.New("user is not an asesor")
	ErrTUKNotFound  = errors.New("tuk not found")
	ErrInvalidLimit = errors.New("daily assignment limit must be >= 0")
)
