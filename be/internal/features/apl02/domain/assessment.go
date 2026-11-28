package domain

import "time"

// Status adalah status siklus hidup asesmen mandiri (FR.APL.02).
type Status string

const (
	StatusDraft     Status = "DRAFT"
	StatusSubmitted Status = "SUBMITTED"
	StatusReviewed  Status = "REVIEWED"
)

// Decision adalah rekomendasi asesor atas asesmen mandiri.
type Decision string

const (
	DecisionRecommended    Decision = "RECOMMENDED"
	DecisionNotRecommended Decision = "NOT_RECOMMENDED"
)

// Assessment adalah header asesmen mandiri (1:1 dengan permohonan APL-01).
type Assessment struct {
	ID             string
	ApplicationID  string
	Code           string
	Status         Status
	AsesorNote     *string
	AsesorDecision *Decision
	SubmittedAt    *time.Time
	ReviewedBy     *string
	ReviewedAt     *time.Time
	ZoomLink       *string
	ZoomDate       *time.Time
	ZoomTime       *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	// SchemaName/SchemaCode hanya diisi pada daftar (hasil join ke permohonan
	// APL-01) untuk menampilkan skema yang dipilih; kosong pada pengambilan tunggal.
	SchemaName string
	SchemaCode string
	// ApplicationCode/AsesiName/ApplicationStatus hanya diisi pada daftar
	// staff (antrian verifikasi) hasil join ke permohonan & profil asesi.
	ApplicationCode   string
	AsesiName         string
	ApplicationStatus string
	// ApplicationNote adalah catatan verifikasi admin pada permohonan APL-01
	// (mis. alasan penolakan); hanya diisi pada daftar.
	ApplicationNote *string
	// AsesorID/AsesorName/AsesorAssignmentStatus hanya diisi pada daftar staff
	// (hasil join ke penugasan asesor aktif pada permohonan).
	AsesorID               *string
	AsesorName             string
	AsesorAssignmentStatus *string
}

// UnitAnswer adalah jawaban asesi untuk satu unit kompetensi, beserta berkas
// bukti. UnitCode/UnitTitle diisi dari hasil join ke competency_units.
type UnitAnswer struct {
	ID               string
	AssessmentID     string
	CompetencyUnitID string
	IsCompetent      *bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	UnitCode         string
	UnitTitle        string
	Files            []UnitFile
}

// UnitFile adalah satu berkas bukti untuk sebuah jawaban unit.
type UnitFile struct {
	ID               string
	AssessmentUnitID string
	FileURL          string
	CreatedAt        time.Time
}

// ApplicationInfo adalah info permohonan/skema yang ditampilkan pada laman
// asesmen mandiri (read-only).
type ApplicationInfo struct {
	SchemaID      string
	SchemaName    string
	SchemaCode    string
	Purpose       string
	DataCreatedAt time.Time
}

// AssessmentDetail menggabungkan header dengan info permohonan & seluruh jawaban unit.
type AssessmentDetail struct {
	Assessment  Assessment
	Application ApplicationInfo
	Units       []UnitAnswer
	// Locked menandai formulir tidak boleh diubah lagi: status bukan DRAFT,
	// atau QR sudah terbit (token asesi masih aktif/menunggu).
	Locked bool
}
