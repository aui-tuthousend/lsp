package domain

import "errors"

var (
	ErrNotFound = errors.New("certification application not found")
	ErrConflict = errors.New("certification application already exists")
	ErrNotDraft = errors.New("certification application is not a draft")
	// ErrLocked: formulir terkunci karena QR sudah terbit (meski belum discan).
	ErrLocked       = errors.New("certification application is locked after submission")
	ErrIncomplete   = errors.New("certification application is incomplete")
	ErrInvalidToken = errors.New("invalid or expired verification token")
	ErrTokenUsed    = errors.New("verification token already used")
	// ErrScheduleConflict: asesi sudah punya pengajuan lain pada tanggal jadwal
	// yang sama.
	ErrScheduleConflict = errors.New("already have an application on this date")
	// ErrTukRequired: jadwal wajib menyertakan TUK, baik untuk mode offline
	// maupun online (TUK menentukan lokasi asesor yang ditugaskan).
	ErrTukRequired = errors.New("tuk_id is required")
	// ErrNotApproved: plotting asesor hanya boleh untuk permohonan yang sudah
	// disetujui admin LSP.
	ErrNotApproved = errors.New("certification application is not approved yet")
)
