package domain

import "errors"

var (
	ErrNotFound                   = errors.New("asesor profile not found")
	ErrNotAsesor                  = errors.New("user is not an asesor")
	ErrTUKNotFound                = errors.New("tuk not found")
	ErrInvalidLimit               = errors.New("daily assignment limit must be >= 0")
	ErrCertificationNotFound      = errors.New("asesor certification not found")
	ErrInvalidCertificationPeriod = errors.New("certification expired_at must be after issued_at")
)
