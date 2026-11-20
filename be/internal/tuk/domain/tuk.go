package domain

import "time"

// TUK adalah master data Tempat Uji Kompetensi.
type TUK struct {
	ID         string
	Name       string
	Logo       *string
	PostalCode string
	Phone      string
	Email      string
	Province   string
	City       string
	Address    string
	MapsURL    string
	Letterhead *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
