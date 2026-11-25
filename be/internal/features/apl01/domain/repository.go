package domain

import (
	"context"
	"time"
)

type ListParams struct {
	// UserID membatasi hasil ke satu asesi; kosong = semua asesi (admin).
	UserID string
	Status *Status
	Limit  int
	Offset int
}

type ListResult struct {
	Items []Application
	Total int
}

type Repository interface {
	// Header permohonan.
	CreateApplication(ctx context.Context, application Application) error
	FindApplicationByID(ctx context.Context, id string) (Application, error)
	FindActiveDraft(ctx context.Context, userID string) (Application, error)
	ListApplications(ctx context.Context, params ListParams) (ListResult, error)
	UpdateCurrentStep(ctx context.Context, id string, step int) (Application, error)
	// MarkSubmitted menandai permohonan SUBMITTED (dipanggil saat asesi
	// memverifikasi via QR).
	MarkSubmitted(ctx context.Context, id string, submittedAt time.Time) (Application, error)
	ApplyDecision(ctx context.Context, decision Decision) (Application, error)

	// Data sertifikasi (step 2).
	UpsertApplicationData(ctx context.Context, data ApplicationData) (ApplicationData, error)
	FindApplicationData(ctx context.Context, applicationID string) (ApplicationData, error)

	// Bukti kelengkapan (step 3). Penggantian set bersifat transaksional:
	// hapus yang tidak lagi dikirim lalu upsert yang dikirim.
	ListApplicationDocuments(ctx context.Context, applicationID string) ([]ApplicationDocument, error)
	ReplaceApplicationDocuments(ctx context.Context, applicationID string, documents []ApplicationDocument) error
	// CountMissingRequiredDocuments menghitung dokumen wajib (sesuai skema) yang
	// belum memiliki file.
	CountMissingRequiredDocuments(ctx context.Context, applicationID string) (int, error)

	// Jadwal asesmen (step 4).
	UpsertAssessmentSchedule(ctx context.Context, schedule AssessmentSchedule) (AssessmentSchedule, error)
	FindAssessmentSchedule(ctx context.Context, applicationID string) (AssessmentSchedule, error)
	// HasScheduleConflict melaporkan apakah asesi pemilik permohonan sudah punya
	// permohonan lain (status selain REJECTED) pada tanggal jadwal yang sama.
	HasScheduleConflict(ctx context.Context, applicationID string, date time.Time) (bool, error)

	// Token verifikasi QR.
	CreateSignature(ctx context.Context, signature FormSignature) error
	RevokeActiveSignatures(ctx context.Context, applicationID string, entity FormEntity, role SignerRole, revokedAt time.Time) error
	FindSignatureByTokenHash(ctx context.Context, tokenHash string) (FormSignature, error)
	MarkSignatureUsed(ctx context.Context, id string, usedAt time.Time, ip, userAgent *string) (FormSignature, error)
	ListSignatures(ctx context.Context, applicationID string) ([]FormSignature, error)
}
