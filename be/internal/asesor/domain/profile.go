// Package domain berisi model dan port untuk profil asesor.
package domain

import "time"

// Profile adalah konfigurasi milik seorang asesor: afiliasi TUK, domisili,
// dan batas penugasan harian.
type Profile struct {
	UserID               string
	TUKID                *string // TUK tempat asesor berafiliasi (opsional)
	Province             *string // domisili asesor
	City                 *string // domisili asesor
	DailyAssignmentLimit int     // 0 = tanpa batas
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// ListItem adalah profil asesor yang dilengkapi data user & TUK untuk tampilan admin.
type ListItem struct {
	UserID               string
	Name                 string
	Email                string
	TUKID                *string
	TUKName              *string
	Province             *string
	City                 *string
	DailyAssignmentLimit int
}
