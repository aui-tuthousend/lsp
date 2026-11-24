package domain

import "time"

// DocumentType adalah master jenis dokumen (mis. 'ktp', 'ijazah').
type DocumentType struct {
	ID        string
	Code      string
	Label     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DocumentTypePatch struct {
	ID    string
	Code  *string
	Label *string
}
