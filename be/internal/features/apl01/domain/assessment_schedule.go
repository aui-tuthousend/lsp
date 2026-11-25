package domain

import "time"

// AssessmentSchedule adalah jadwal asesmen yang diajukan asesi (relasi 1:1 ke
// permohonan).
type AssessmentSchedule struct {
	ID            string
	ApplicationID string
	Mode          Mode
	ScheduledDate time.Time
	TukID         *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
