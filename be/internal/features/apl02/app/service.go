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

	"lsp/internal/features/apl02/domain"
	"lsp/internal/platform/database"
	"lsp/internal/platform/events"
	"lsp/internal/platform/id"
)

// signatureTTL adalah masa berlaku token QR verifikasi asesmen mandiri.
const signatureTTL = 24 * time.Hour

// eventAssessmentSubmitted dipublikasikan ke topic user asesi setelah asesmen
// mandiri diverifikasi via QR (form_signatures -> dipakai).
const eventAssessmentSubmitted = "apl02.submitted"

type Service struct {
	repo domain.Repository
	tx   *database.TxManager
	bus  events.Bus
}

func NewService(repo domain.Repository, tx *database.TxManager, bus events.Bus) *Service {
	return &Service{repo: repo, tx: tx, bus: bus}
}

// CreateOrGet mengambil asesmen mandiri milik sebuah permohonan; bila belum ada,
// membuat header baru beserta kerangka jawaban seluruh unit kompetensi skema.
// Hanya boleh diakses pemilik permohonan dan setelah permohonan SUBMITTED/APPROVED.
func (s *Service) CreateOrGet(ctx context.Context, userID, applicationID string) (domain.AssessmentDetail, error) {
	app, err := s.repo.FindApplicationContext(ctx, applicationID)
	if err != nil {
		return domain.AssessmentDetail{}, err
	}
	if app.UserID != userID {
		return domain.AssessmentDetail{}, domain.ErrNotFound
	}
	if app.Status != "SUBMITTED" && app.Status != "APPROVED" {
		return domain.AssessmentDetail{}, domain.ErrNotEligible
	}

	existing, err := s.repo.FindAssessmentByApplicationID(ctx, applicationID)
	if err == nil {
		return s.Detail(ctx, existing.ID)
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return domain.AssessmentDetail{}, err
	}

	unitIDs, err := s.repo.ListSchemaUnitIDs(ctx, app.SchemaID)
	if err != nil {
		return domain.AssessmentDetail{}, err
	}

	assessment := domain.Assessment{
		ID:            id.New(),
		ApplicationID: applicationID,
		Code:          newCode(),
		Status:        domain.StatusDraft,
	}
	answers := make([]domain.UnitAnswer, 0, len(unitIDs))
	for _, unitID := range unitIDs {
		answers = append(answers, domain.UnitAnswer{
			ID:               id.New(),
			AssessmentID:     assessment.ID,
			CompetencyUnitID: unitID,
		})
	}

	err = s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.CreateAssessment(txCtx, assessment); err != nil {
			return err
		}
		if len(answers) == 0 {
			return nil
		}
		return s.repo.CreateUnitAnswers(txCtx, answers)
	})
	if err != nil {
		return domain.AssessmentDetail{}, err
	}
	return s.Detail(ctx, assessment.ID)
}

// Detail mengambil satu asesmen mandiri beserta info permohonan & jawaban unitnya.
func (s *Service) Detail(ctx context.Context, id string) (domain.AssessmentDetail, error) {
	assessment, err := s.repo.FindAssessmentByID(ctx, id)
	if err != nil {
		return domain.AssessmentDetail{}, err
	}
	units, err := s.repo.ListUnitAnswers(ctx, id)
	if err != nil {
		return domain.AssessmentDetail{}, err
	}
	appCtx, err := s.repo.FindApplicationContext(ctx, assessment.ApplicationID)
	if err != nil {
		return domain.AssessmentDetail{}, err
	}
	// Terkunci bila sudah bukan draft, atau QR sudah terbit (belum discan).
	pendingSignature, err := s.repo.HasPendingSignature(ctx, assessment.ApplicationID)
	if err != nil {
		return domain.AssessmentDetail{}, err
	}
	return domain.AssessmentDetail{
		Assessment: assessment,
		Application: domain.ApplicationInfo{
			SchemaID:      appCtx.SchemaID,
			SchemaName:    appCtx.SchemaName,
			SchemaCode:    appCtx.SchemaCode,
			Purpose:       appCtx.Purpose,
			DataCreatedAt: appCtx.DataCreatedAt,
		},
		Units:  units,
		Locked: assessment.Status != domain.StatusDraft || pendingSignature,
	}, nil
}

// List mengambil daftar asesmen mandiri sesuai filter. UserID kosong berarti
// semua asesi (staff); ExcludeDraft dipakai antrian verifikasi agar hanya
// asesmen yang sudah diajukan yang tampil.
func (s *Service) List(ctx context.Context, params domain.ListParams) (domain.ListResult, error) {
	return s.repo.ListAssessments(ctx, params)
}

// ApplicationContext mengembalikan konteks permohonan APL-01 pemilik sebuah
// asesmen mandiri (dipakai middleware otorisasi).
func (s *Service) ApplicationContext(ctx context.Context, assessmentID string) (domain.ApplicationContext, error) {
	assessment, err := s.repo.FindAssessmentByID(ctx, assessmentID)
	if err != nil {
		return domain.ApplicationContext{}, err
	}
	return s.repo.FindApplicationContext(ctx, assessment.ApplicationID)
}

// PutAnswers menyimpan jawaban per unit (set-semantik) beserta berkas bukti.
func (s *Service) PutAnswers(ctx context.Context, cmd PutAnswersCommand) (domain.AssessmentDetail, error) {
	if err := s.ensureDraft(ctx, cmd.AssessmentID); err != nil {
		return domain.AssessmentDetail{}, err
	}
	answers := make([]domain.UnitAnswer, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		files := make([]domain.UnitFile, 0)
		// Berkas hanya disimpan bila asesi menyatakan kompeten.
		if item.IsCompetent != nil && *item.IsCompetent {
			for _, url := range item.FileURLs {
				if strings.TrimSpace(url) == "" {
					continue
				}
				files = append(files, domain.UnitFile{ID: id.New(), FileURL: url})
			}
		}
		answers = append(answers, domain.UnitAnswer{
			CompetencyUnitID: item.CompetencyUnitID,
			IsCompetent:      item.IsCompetent,
			Files:            files,
		})
	}
	err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.repo.ReplaceUnitAnswers(txCtx, cmd.AssessmentID, answers)
	})
	if err != nil {
		return domain.AssessmentDetail{}, err
	}
	return s.Detail(ctx, cmd.AssessmentID)
}

// Submit memvalidasi kelengkapan dan menerbitkan token QR verifikasi APL-02.
// Status asesmen tetap DRAFT sampai asesi memverifikasi via QR (lihat Verify).
// Idempoten bila `existingToken` masih aktif untuk permohonan ini.
func (s *Service) Submit(ctx context.Context, assessmentID, existingToken string) (SubmitResult, error) {
	assessment, err := s.repo.FindAssessmentByID(ctx, assessmentID)
	if err != nil {
		return SubmitResult{}, err
	}
	if assessment.Status != domain.StatusDraft {
		return SubmitResult{}, domain.ErrNotDraft
	}
	if err := s.validateComplete(ctx, assessmentID); err != nil {
		return SubmitResult{}, err
	}

	now := time.Now().UTC()

	if existingToken != "" {
		if sig, err := s.repo.FindSignatureByTokenHash(ctx, hashToken(existingToken)); err == nil &&
			sig.ApplicationID == assessment.ApplicationID &&
			sig.Entity == domain.SignatureEntity &&
			sig.Status == domain.SignatureStatusPending &&
			now.Before(sig.ExpiresAt) {
			detail, err := s.Detail(ctx, assessmentID)
			if err != nil {
				return SubmitResult{}, err
			}
			return SubmitResult{Assessment: detail, VerifyToken: existingToken}, nil
		}
	}

	rawToken, tokenHash, err := newToken()
	if err != nil {
		return SubmitResult{}, err
	}
	err = s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.RevokeActiveSignatures(txCtx, assessment.ApplicationID, domain.SignatureEntity, domain.SignatureRoleAsesi, now); err != nil {
			return err
		}
		return s.repo.CreateSignature(txCtx, domain.Signature{
			ID:            id.New(),
			ApplicationID: assessment.ApplicationID,
			Entity:        domain.SignatureEntity,
			Role:          domain.SignatureRoleAsesi,
			TokenHash:     tokenHash,
			Status:        domain.SignatureStatusPending,
			ExpiresAt:     now.Add(signatureTTL),
		})
	})
	if err != nil {
		return SubmitResult{}, err
	}

	detail, err := s.Detail(ctx, assessmentID)
	if err != nil {
		return SubmitResult{}, err
	}
	return SubmitResult{Assessment: detail, VerifyToken: rawToken}, nil
}

// Verify mengonfirmasi hasil scan QR (menandai token terpakai & asesmen SUBMITTED).
func (s *Service) Verify(ctx context.Context, userID, token string, ip, userAgent *string) (domain.AssessmentDetail, error) {
	signature, err := s.repo.FindSignatureByTokenHash(ctx, hashToken(token))
	if err != nil {
		return domain.AssessmentDetail{}, err
	}
	if signature.Entity != domain.SignatureEntity {
		return domain.AssessmentDetail{}, domain.ErrInvalidToken
	}
	if signature.Status != domain.SignatureStatusPending {
		if signature.Status == domain.SignatureStatusUsed {
			return domain.AssessmentDetail{}, domain.ErrTokenUsed
		}
		return domain.AssessmentDetail{}, domain.ErrInvalidToken
	}
	now := time.Now().UTC()
	if now.After(signature.ExpiresAt) {
		return domain.AssessmentDetail{}, domain.ErrInvalidToken
	}

	app, err := s.repo.FindApplicationContext(ctx, signature.ApplicationID)
	if err != nil {
		return domain.AssessmentDetail{}, err
	}
	if app.UserID != userID {
		return domain.AssessmentDetail{}, domain.ErrInvalidToken
	}
	assessment, err := s.repo.FindAssessmentByApplicationID(ctx, signature.ApplicationID)
	if err != nil {
		return domain.AssessmentDetail{}, err
	}

	err = s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if _, err := s.repo.MarkSignatureUsed(txCtx, signature.ID, now, ip, userAgent); err != nil {
			return err
		}
		_, err := s.repo.MarkSubmitted(txCtx, assessment.ID, now)
		return err
	})
	if err != nil {
		return domain.AssessmentDetail{}, err
	}

	s.publishSubmitted(ctx, app.UserID, assessment)
	return s.Detail(ctx, assessment.ID)
}

// VerifyStatus mengembalikan status token verifikasi tanpa memakainya.
func (s *Service) VerifyStatus(ctx context.Context, userID, token string) (domain.Signature, domain.Assessment, error) {
	signature, err := s.repo.FindSignatureByTokenHash(ctx, hashToken(token))
	if err != nil {
		return domain.Signature{}, domain.Assessment{}, err
	}
	if signature.Entity != domain.SignatureEntity {
		return domain.Signature{}, domain.Assessment{}, domain.ErrNotFound
	}
	app, err := s.repo.FindApplicationContext(ctx, signature.ApplicationID)
	if err != nil {
		return domain.Signature{}, domain.Assessment{}, err
	}
	if app.UserID != userID {
		return domain.Signature{}, domain.Assessment{}, domain.ErrNotFound
	}
	assessment, err := s.repo.FindAssessmentByApplicationID(ctx, signature.ApplicationID)
	if err != nil {
		return domain.Signature{}, domain.Assessment{}, err
	}
	return signature, assessment, nil
}

// Review menerapkan hasil tinjauan asesor (rekomendasi, catatan, jadwal zoom).
func (s *Service) Review(ctx context.Context, cmd ReviewCommand) (domain.AssessmentDetail, error) {
	assessment, err := s.repo.FindAssessmentByID(ctx, cmd.AssessmentID)
	if err != nil {
		return domain.AssessmentDetail{}, err
	}
	if assessment.Status != domain.StatusSubmitted {
		return domain.AssessmentDetail{}, domain.ErrNotSubmitted
	}
	if _, err := s.repo.MarkReviewed(ctx, domain.Review{
		AssessmentID: cmd.AssessmentID,
		Decision:     cmd.Decision,
		Note:         cmd.Note,
		ReviewedBy:   cmd.ReviewedBy,
		ReviewedAt:   time.Now().UTC(),
		ZoomLink:     cmd.ZoomLink,
		ZoomDate:     cmd.ZoomDate,
		ZoomTime:     cmd.ZoomTime,
	}); err != nil {
		return domain.AssessmentDetail{}, err
	}
	return s.Detail(ctx, cmd.AssessmentID)
}

func (s *Service) publishSubmitted(ctx context.Context, userID string, assessment domain.Assessment) {
	if s.bus == nil {
		return
	}
	payload, err := json.Marshal(map[string]any{
		"assessment_id":  assessment.ID,
		"application_id": assessment.ApplicationID,
		"code":           assessment.Code,
		"status":         string(domain.StatusSubmitted),
	})
	if err != nil {
		return
	}
	_ = s.bus.Publish(ctx, events.Event{
		Topic: events.UserTopic(userID),
		Type:  eventAssessmentSubmitted,
		Data:  payload,
	})
}

func (s *Service) ensureDraft(ctx context.Context, assessmentID string) error {
	assessment, err := s.repo.FindAssessmentByID(ctx, assessmentID)
	if err != nil {
		return err
	}
	if assessment.Status != domain.StatusDraft {
		return domain.ErrNotDraft
	}
	// Terkunci begitu QR terbit (meski belum discan): jawaban tidak boleh berubah.
	pending, err := s.repo.HasPendingSignature(ctx, assessment.ApplicationID)
	if err != nil {
		return err
	}
	if pending {
		return domain.ErrLocked
	}
	return nil
}

func (s *Service) validateComplete(ctx context.Context, assessmentID string) error {
	unanswered, err := s.repo.CountUnansweredUnits(ctx, assessmentID)
	if err != nil {
		return err
	}
	if unanswered > 0 {
		return domain.ErrIncomplete
	}
	missing, err := s.repo.CountCompetentWithoutFiles(ctx, assessmentID)
	if err != nil {
		return err
	}
	if missing > 0 {
		return domain.ErrIncomplete
	}
	return nil
}

// --- Command types ---

type UnitAnswerItem struct {
	CompetencyUnitID string
	IsCompetent      *bool
	FileURLs         []string
}

type PutAnswersCommand struct {
	AssessmentID string
	Items        []UnitAnswerItem
}

type ReviewCommand struct {
	AssessmentID string
	Decision     *domain.Decision
	Note         *string
	ZoomLink     *string
	ZoomDate     *time.Time
	ZoomTime     *string
	ReviewedBy   string
}

type SubmitResult struct {
	Assessment  domain.AssessmentDetail
	VerifyToken string
}

// --- Helpers ---

func newCode() string {
	return fmt.Sprintf("APL02-%s-%s", time.Now().Format("20060102"), strings.ToUpper(id.New()[:6]))
}

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
