package domain

import (
	"context"
	"time"
)

// Entity asesmen mandiri pada tabel bersama form_signatures (FR.APL.02).
const (
	SignatureEntity     = "asesmen_mandiri"
	SignatureRoleAsesi  = "asesi"
	SignatureRoleAsesor = "asesor"

	SignatureStatusPending = "menunggu"
	SignatureStatusUsed    = "dipakai"
)

type ListParams struct {
	// UserID membatasi hasil ke asesi pemilik; kosong = semua (staff).
	UserID string
	Limit  int
	Offset int
}

type ListResult struct {
	Items []Assessment
	Total int
}

// ApplicationContext adalah potongan data permohonan APL-01 yang dibutuhkan
// modul apl02 (pemilik, status, skema terpilih, asesor yang di-plotting) plus
// info skema/tujuan untuk ditampilkan pada laman asesmen mandiri.
type ApplicationContext struct {
	UserID        string
	Status        string
	SchemaID      string
	AsesorID      *string
	SchemaName    string
	SchemaCode    string
	Purpose       string
	DataCreatedAt time.Time
}

// Review membawa hasil tinjauan asesor atas asesmen mandiri.
type Review struct {
	AssessmentID string
	Decision     *Decision
	Note         *string
	ReviewedBy   string
	ReviewedAt   time.Time
	ZoomLink     *string
	ZoomDate     *time.Time
	ZoomTime     *string
}

// Signature adalah token verifikasi QR pada tabel bersama form_signatures.
type Signature struct {
	ID            string
	ApplicationID string
	Entity        string
	Role          string
	TokenHash     string
	Status        string
	ExpiresAt     time.Time
	UsedAt        *time.Time
	RevokedAt     *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Repository interface {
	// Konteks permohonan APL-01 terkait (pemilik, status, skema, asesor).
	FindApplicationContext(ctx context.Context, applicationID string) (ApplicationContext, error)

	// Header asesmen mandiri.
	CreateAssessment(ctx context.Context, assessment Assessment) error
	FindAssessmentByID(ctx context.Context, id string) (Assessment, error)
	FindAssessmentByApplicationID(ctx context.Context, applicationID string) (Assessment, error)
	ListAssessments(ctx context.Context, params ListParams) (ListResult, error)
	MarkSubmitted(ctx context.Context, id string, submittedAt time.Time) (Assessment, error)
	MarkReviewed(ctx context.Context, review Review) (Assessment, error)

	// Jawaban per unit kompetensi.
	ListUnitAnswers(ctx context.Context, assessmentID string) ([]UnitAnswer, error)
	// CreateUnitAnswers membuat baris jawaban (kerangka) untuk daftar unit.
	CreateUnitAnswers(ctx context.Context, answers []UnitAnswer) error
	// ReplaceUnitAnswers menyimpan jawaban (set-semantik) beserta berkas bukti.
	ReplaceUnitAnswers(ctx context.Context, assessmentID string, answers []UnitAnswer) error
	// CountUnansweredUnits menghitung unit yang belum dijawab (is_competent null).
	CountUnansweredUnits(ctx context.Context, assessmentID string) (int, error)
	// CountCompetentWithoutFiles menghitung unit kompeten yang belum punya bukti.
	CountCompetentWithoutFiles(ctx context.Context, assessmentID string) (int, error)
	// ListSchemaUnitIDs mengembalikan id unit kompetensi milik sebuah skema.
	ListSchemaUnitIDs(ctx context.Context, schemaID string) ([]string, error)

	// Token verifikasi QR (tabel bersama form_signatures).
	CreateSignature(ctx context.Context, signature Signature) error
	RevokeActiveSignatures(ctx context.Context, applicationID, entity, role string, revokedAt time.Time) error
	FindSignatureByTokenHash(ctx context.Context, tokenHash string) (Signature, error)
	MarkSignatureUsed(ctx context.Context, id string, usedAt time.Time, ip, userAgent *string) (Signature, error)
}
