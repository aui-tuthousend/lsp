package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"lsp/internal/features/apl01/domain"
	"lsp/internal/platform/database"
	"lsp/internal/platform/events"
	"lsp/internal/platform/id"
)

// signatureTTL adalah masa berlaku token QR verifikasi APL-01.
const signatureTTL = 24 * time.Hour

// eventApplicationVerified dipublikasikan ke topic user asesi setelah APL-01
// diverifikasi via QR (form_signatures -> dipakai).
const eventApplicationVerified = "apl01.verified"

type Service struct {
	repo domain.Repository
	tx   *database.TxManager
	bus  events.Bus
}

func NewService(repo domain.Repository, tx *database.TxManager, bus events.Bus) *Service {
	return &Service{repo: repo, tx: tx, bus: bus}
}

// CreateOrGetDraft mengambil draft APL-01 aktif milik user; bila belum ada,
// membuat draft baru. Dipakai saat asesi mulai mengisi formulir.
func (s *Service) CreateOrGetDraft(ctx context.Context, userID string) (domain.ApplicationDetail, error) {
	existing, err := s.repo.FindActiveDraft(ctx, userID)
	if err == nil {
		return s.Detail(ctx, existing.ID)
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return domain.ApplicationDetail{}, err
	}

	application := domain.Application{
		ID:          id.New(),
		UserID:      userID,
		Code:        newCode(),
		Status:      domain.StatusDraft,
		CurrentStep: 1,
	}
	if err := s.repo.CreateApplication(ctx, application); err != nil {
		return domain.ApplicationDetail{}, err
	}
	return s.Detail(ctx, application.ID)
}

// Detail mengambil satu permohonan beserta seluruh data tahapnya.
func (s *Service) Detail(ctx context.Context, applicationID string) (domain.ApplicationDetail, error) {
	application, err := s.repo.FindApplicationByID(ctx, applicationID)
	if err != nil {
		return domain.ApplicationDetail{}, err
	}
	detail := domain.ApplicationDetail{Application: application}

	data, err := s.repo.FindApplicationData(ctx, applicationID)
	switch {
	case err == nil:
		detail.Data = &data
	case errors.Is(err, domain.ErrNotFound):
	default:
		return domain.ApplicationDetail{}, err
	}

	documents, err := s.repo.ListApplicationDocuments(ctx, applicationID)
	if err != nil {
		return domain.ApplicationDetail{}, err
	}
	detail.Documents = documents

	schedule, err := s.repo.FindAssessmentSchedule(ctx, applicationID)
	switch {
	case err == nil:
		detail.Schedule = &schedule
	case errors.Is(err, domain.ErrNotFound):
	default:
		return domain.ApplicationDetail{}, err
	}

	signatures, err := s.repo.ListSignatures(ctx, applicationID)
	if err != nil {
		return domain.ApplicationDetail{}, err
	}
	detail.Signatures = signatures

	// Terkunci bila sudah bukan draft, atau QR sudah terbit (belum discan).
	pendingSignature, err := s.repo.HasPendingSignature(ctx, applicationID)
	if err != nil {
		return domain.ApplicationDetail{}, err
	}
	detail.Locked = application.Status != domain.StatusDraft || pendingSignature

	return detail, nil
}

// List mengambil daftar permohonan. userID kosong berarti semua asesi (admin).
func (s *Service) List(ctx context.Context, userID string, status *domain.Status, limit, offset int) (domain.ListResult, error) {
	return s.repo.ListApplications(ctx, domain.ListParams{UserID: userID, Status: status, Limit: limit, Offset: offset})
}

// UpdateStep menyimpan step terakhir yang dikerjakan asesi (untuk resume draft).
func (s *Service) UpdateStep(ctx context.Context, applicationID string, step int) (domain.ApplicationDetail, error) {
	if err := s.ensureDraft(ctx, applicationID); err != nil {
		return domain.ApplicationDetail{}, err
	}
	if _, err := s.repo.UpdateCurrentStep(ctx, applicationID, step); err != nil {
		return domain.ApplicationDetail{}, err
	}
	return s.Detail(ctx, applicationID)
}

// PutData menyimpan data sertifikasi (step 2): skema + tujuan asesmen.
func (s *Service) PutData(ctx context.Context, cmd PutDataCommand) (domain.ApplicationDetail, error) {
	if err := s.ensureDraft(ctx, cmd.ApplicationID); err != nil {
		return domain.ApplicationDetail{}, err
	}
	if _, err := s.repo.UpsertApplicationData(ctx, domain.ApplicationData{
		ID:                    id.New(),
		ApplicationID:         cmd.ApplicationID,
		CertificationSchemaID: cmd.CertificationSchemaID,
		Purpose:               cmd.Purpose,
	}); err != nil {
		return domain.ApplicationDetail{}, err
	}
	return s.Detail(ctx, cmd.ApplicationID)
}

// PutDocuments mengganti seluruh bukti kelengkapan (step 3) dengan daftar yang
// dikirim (set-semantik, idempoten).
func (s *Service) PutDocuments(ctx context.Context, cmd PutDocumentsCommand) (domain.ApplicationDetail, error) {
	if err := s.ensureDraft(ctx, cmd.ApplicationID); err != nil {
		return domain.ApplicationDetail{}, err
	}
	documents := make([]domain.ApplicationDocument, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		status := domain.DocumentIncomplete
		if item.FileURL != nil && *item.FileURL != "" {
			status = domain.DocumentMeets
		}
		documents = append(documents, domain.ApplicationDocument{
			ID:             id.New(),
			ApplicationID:  cmd.ApplicationID,
			DocumentTypeID: item.DocumentTypeID,
			FileURL:        item.FileURL,
			Status:         status,
		})
	}
	err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.repo.ReplaceApplicationDocuments(txCtx, cmd.ApplicationID, documents)
	})
	if err != nil {
		return domain.ApplicationDetail{}, err
	}
	return s.Detail(ctx, cmd.ApplicationID)
}

// PutSchedule menyimpan jadwal asesmen (step 4). Menolak bila asesi sudah punya
// pengajuan lain pada tanggal yang sama.
func (s *Service) PutSchedule(ctx context.Context, cmd PutScheduleCommand) (domain.ApplicationDetail, error) {
	if err := s.ensureDraft(ctx, cmd.ApplicationID); err != nil {
		return domain.ApplicationDetail{}, err
	}
	conflict, err := s.repo.HasScheduleConflict(ctx, cmd.ApplicationID, cmd.ScheduledDate)
	if err != nil {
		return domain.ApplicationDetail{}, err
	}
	if conflict {
		return domain.ApplicationDetail{}, domain.ErrScheduleConflict
	}
	if _, err := s.repo.UpsertAssessmentSchedule(ctx, domain.AssessmentSchedule{
		ID:            id.New(),
		ApplicationID: cmd.ApplicationID,
		Mode:          cmd.Mode,
		ScheduledDate: cmd.ScheduledDate,
		TukID:         cmd.TukID,
	}); err != nil {
		return domain.ApplicationDetail{}, err
	}
	return s.Detail(ctx, cmd.ApplicationID)
}

// Submit memvalidasi kelengkapan dan menerbitkan token QR verifikasi APL-01.
// Status permohonan tetap DRAFT sampai asesi memverifikasi via QR (lihat Verify).
// Token mentah dikembalikan sekali (untuk ditampilkan sebagai QR); yang
// tersimpan hanya hash-nya.
//
// Idempoten bila `existingToken` (token yang dikirim klien) masih aktif
// (menunggu & belum kedaluwarsa) untuk permohonan ini: token yang sama
// dikembalikan tanpa menerbitkan token baru.
func (s *Service) Submit(ctx context.Context, applicationID, existingToken string) (SubmitResult, error) {
	application, err := s.repo.FindApplicationByID(ctx, applicationID)
	if err != nil {
		return SubmitResult{}, err
	}
	if application.Status != domain.StatusDraft {
		return SubmitResult{}, domain.ErrNotDraft
	}
	if err := s.validateComplete(ctx, applicationID); err != nil {
		return SubmitResult{}, err
	}

	now := time.Now().UTC()

	// Pakai ulang token yang sama bila klien masih memegang token aktif.
	if existingToken != "" {
		if sig, err := s.repo.FindSignatureByTokenHash(ctx, hashToken(existingToken)); err == nil &&
			sig.ApplicationID == applicationID &&
			sig.Status == domain.SignaturePending &&
			now.Before(sig.ExpiresAt) {
			detail, err := s.Detail(ctx, applicationID)
			if err != nil {
				return SubmitResult{}, err
			}
			return SubmitResult{Application: detail, VerifyToken: existingToken}, nil
		}
	}

	rawToken, tokenHash, err := newToken()
	if err != nil {
		return SubmitResult{}, err
	}
	err = s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.RevokeActiveSignatures(txCtx, applicationID, domain.EntityPermohonan, domain.SignerAsesi, now); err != nil {
			return err
		}
		return s.repo.CreateSignature(txCtx, domain.FormSignature{
			ID:            id.New(),
			ApplicationID: applicationID,
			Entity:        domain.EntityPermohonan,
			Role:          domain.SignerAsesi,
			TokenHash:     tokenHash,
			Status:        domain.SignaturePending,
			ExpiresAt:     now.Add(signatureTTL),
		})
	})
	if err != nil {
		return SubmitResult{}, err
	}

	detail, err := s.Detail(ctx, applicationID)
	if err != nil {
		return SubmitResult{}, err
	}
	return SubmitResult{Application: detail, VerifyToken: rawToken}, nil
}

// Verify mengonfirmasi hasil scan QR berdasarkan token (permohonan ditentukan
// oleh token): menandai token terpakai (single-use), mengubah status permohonan
// menjadi SUBMITTED, lalu mem-publish event SSE ke topic user asesi.
func (s *Service) Verify(ctx context.Context, userID, token string, ip, userAgent *string) (domain.ApplicationDetail, error) {
	signature, err := s.repo.FindSignatureByTokenHash(ctx, hashToken(token))
	if err != nil {
		return domain.ApplicationDetail{}, err
	}
	// Token harus milik formulir APL-01 (tabel form_signatures dipakai lintas tahap).
	if signature.Entity != domain.EntityPermohonan {
		return domain.ApplicationDetail{}, domain.ErrInvalidToken
	}
	if signature.Status != domain.SignaturePending {
		if signature.Status == domain.SignatureUsed {
			return domain.ApplicationDetail{}, domain.ErrTokenUsed
		}
		return domain.ApplicationDetail{}, domain.ErrInvalidToken
	}
	now := time.Now().UTC()
	if now.After(signature.ExpiresAt) {
		return domain.ApplicationDetail{}, domain.ErrInvalidToken
	}

	application, err := s.repo.FindApplicationByID(ctx, signature.ApplicationID)
	if err != nil {
		return domain.ApplicationDetail{}, err
	}
	if application.UserID != userID {
		return domain.ApplicationDetail{}, domain.ErrInvalidToken
	}

	err = s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if _, err := s.repo.MarkSignatureUsed(txCtx, signature.ID, now, ip, userAgent); err != nil {
			return err
		}
		_, err := s.repo.MarkSubmitted(txCtx, signature.ApplicationID, now)
		return err
	})
	if err != nil {
		return domain.ApplicationDetail{}, err
	}

	// Publish setelah commit, best-effort (dipakai FE untuk toast realtime).
	s.publishVerified(ctx, application)
	return s.Detail(ctx, signature.ApplicationID)
}

// VerifyStatus mengembalikan status token verifikasi tanpa memakainya. Dipakai
// laman verifikasi untuk menampilkan "sudah diverifikasi" tanpa tombol.
func (s *Service) VerifyStatus(ctx context.Context, userID, token string) (domain.FormSignature, domain.Application, error) {
	signature, err := s.repo.FindSignatureByTokenHash(ctx, hashToken(token))
	if err != nil {
		return domain.FormSignature{}, domain.Application{}, err
	}
	// Token harus milik formulir APL-01 (tabel form_signatures dipakai lintas tahap).
	if signature.Entity != domain.EntityPermohonan {
		return domain.FormSignature{}, domain.Application{}, domain.ErrNotFound
	}
	application, err := s.repo.FindApplicationByID(ctx, signature.ApplicationID)
	if err != nil {
		return domain.FormSignature{}, domain.Application{}, err
	}
	if application.UserID != userID {
		return domain.FormSignature{}, domain.Application{}, domain.ErrNotFound
	}
	return signature, application, nil
}

// publishVerified mengirim event SSE ke topic user asesi (best-effort).
func (s *Service) publishVerified(ctx context.Context, application domain.Application) {
	if s.bus == nil {
		return
	}
	payload, err := json.Marshal(map[string]any{
		"application_id": application.ID,
		"code":           application.Code,
		"status":         string(domain.StatusSubmitted),
	})
	if err != nil {
		return
	}
	_ = s.bus.Publish(ctx, events.Event{
		Topic: events.UserTopic(application.UserID),
		Type:  eventApplicationVerified,
		Data:  payload,
	})
}

// Decide menerapkan hasil verifikasi admin (APPROVED/REJECTED + plotting asesor).
func (s *Service) Decide(ctx context.Context, cmd DecideCommand) (domain.ApplicationDetail, error) {
	now := time.Now().UTC()
	var assignedAt *time.Time
	if cmd.AsesorID != nil {
		assignedAt = &now
	}
	if _, err := s.repo.ApplyDecision(ctx, domain.Decision{
		ApplicationID:    cmd.ApplicationID,
		Status:           cmd.Decision,
		VerifiedBy:       cmd.VerifiedBy,
		VerifiedAt:       now,
		Note:             cmd.Note,
		AsesorID:         cmd.AsesorID,
		AsesorAssignedAt: assignedAt,
	}); err != nil {
		return domain.ApplicationDetail{}, err
	}
	return s.Detail(ctx, cmd.ApplicationID)
}

func (s *Service) ensureDraft(ctx context.Context, applicationID string) error {
	application, err := s.repo.FindApplicationByID(ctx, applicationID)
	if err != nil {
		return err
	}
	if application.Status != domain.StatusDraft {
		return domain.ErrNotDraft
	}
	// Terkunci begitu QR terbit (meski belum discan): data tidak boleh berubah.
	pending, err := s.repo.HasPendingSignature(ctx, applicationID)
	if err != nil {
		return err
	}
	if pending {
		return domain.ErrLocked
	}
	return nil
}

func (s *Service) validateComplete(ctx context.Context, applicationID string) error {
	if _, err := s.repo.FindApplicationData(ctx, applicationID); err != nil {
		return domain.ErrIncomplete
	}
	if _, err := s.repo.FindAssessmentSchedule(ctx, applicationID); err != nil {
		return domain.ErrIncomplete
	}
	missing, err := s.repo.CountMissingRequiredDocuments(ctx, applicationID)
	if err != nil {
		return err
	}
	if missing > 0 {
		return domain.ErrIncomplete
	}
	return nil
}

// --- Command types ---

type PutDataCommand struct {
	ApplicationID         string
	CertificationSchemaID string
	Purpose               domain.Purpose
}

type DocumentItem struct {
	DocumentTypeID string
	FileURL        *string
}

type PutDocumentsCommand struct {
	ApplicationID string
	Items         []DocumentItem
}

type PutScheduleCommand struct {
	ApplicationID string
	Mode          domain.Mode
	ScheduledDate time.Time
	TukID         *string
}

type DecideCommand struct {
	ApplicationID string
	Decision      domain.Status
	Note          *string
	AsesorID      *string
	VerifiedBy    string
}

type SubmitResult struct {
	Application domain.ApplicationDetail
	VerifyToken string
}

// --- Helpers ---

func newCode() string {
	return fmt.Sprintf("APL01-%s-%s", time.Now().Format("20060102"), strings.ToUpper(id.New()[:6]))
}

// newToken menghasilkan token acak (untuk QR) beserta hash sha256-nya.
func newToken() (raw string, hash string, err error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	raw = base64.RawURLEncoding.EncodeToString(buf)
	return raw, hashToken(raw), nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
