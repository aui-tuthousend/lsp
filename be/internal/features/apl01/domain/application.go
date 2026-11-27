package domain

import "time"

// Status adalah status siklus hidup permohonan APL-01.
type Status string

const (
	StatusDraft     Status = "DRAFT"
	StatusSubmitted Status = "SUBMITTED"
	StatusApproved  Status = "APPROVED"
	StatusRejected  Status = "REJECTED"
)

// Purpose adalah tujuan asesmen (istilah regulasi BNSP).
type Purpose string

const (
	PurposeCertification     Purpose = "SERTIFIKASI"
	PurposeRecertification   Purpose = "SERTIFIKASI_ULANG"
	PurposeRPL               Purpose = "REKOGNISI_PEMBELAJARAN_LAMPAU"
	PurposeCurrentCompetency Purpose = "PENGAKUAN_KOMPETENSI_TERKINI"
)

// Mode adalah metode pelaksanaan asesmen.
type Mode string

const (
	ModeOnline  Mode = "ONLINE"
	ModeOffline Mode = "OFFLINE"
)

// DocumentStatus adalah status kelengkapan satu bukti dasar.
type DocumentStatus string

const (
	DocumentIncomplete DocumentStatus = "BELUM_LENGKAP"
	DocumentMeets      DocumentStatus = "MEMENUHI_SYARAT"
	DocumentNotMeets   DocumentStatus = "TIDAK_MEMENUHI_SYARAT"
)

// Application adalah header permohonan APL-01 (1 baris = 1 pengajuan asesi).
type Application struct {
	ID               string
	UserID           string
	Code             string
	Status           Status
	CurrentStep      int
	SubmittedAt      *time.Time
	VerifiedBy       *string
	VerifiedAt       *time.Time
	VerificationNote *string
	AsesorID         *string
	AsesorAssignedAt *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

// ApplicationDetail menggabungkan seluruh data tahap APL-01: data sertifikasi,
// bukti kelengkapan, jadwal, dan token verifikasi.
type ApplicationDetail struct {
	Application Application
	Data        *ApplicationData
	Documents   []ApplicationDocument
	Schedule    *AssessmentSchedule
	Signatures  []FormSignature
	// Locked menandai formulir tidak boleh diubah lagi: status bukan DRAFT,
	// atau QR sudah terbit (token asesi masih aktif/menunggu).
	Locked bool
}

// Decision membawa hasil verifikasi admin atas sebuah permohonan.
type Decision struct {
	ApplicationID    string
	Status           Status
	VerifiedBy       string
	VerifiedAt       time.Time
	Note             *string
	AsesorID         *string
	AsesorAssignedAt *time.Time
}
