package http

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/asesor/app"
	httpplatform "lsp/internal/platform/http"
)

type Handler struct {
	service *app.Service
}

func NewHandler(service *app.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Upsert(ctx context.Context, input *UpsertProfileInput) (*ProfileOutput, error) {
	profile, err := h.service.Upsert(ctx, app.UpsertCommand{
		UserID:               input.UserID,
		TUKID:                input.Body.TUKID,
		Province:             input.Body.Province,
		City:                 input.Body.City,
		DailyAssignmentLimit: input.Body.DailyAssignmentLimit,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &ProfileOutput{Body: toProfileDTO(profile)}, nil
}

func (h *Handler) UpdateProfile(ctx context.Context, input *UpdateProfileInput) (*ProfileOutput, error) {
	profile, err := h.service.PatchProfile(ctx, app.PatchProfileCommand{
		UserID:               input.UserID,
		TUKID:                input.Body.TUKID,
		Province:             input.Body.Province,
		City:                 input.Body.City,
		DailyAssignmentLimit: input.Body.DailyAssignmentLimit,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &ProfileOutput{Body: toProfileDTO(profile)}, nil
}

func (h *Handler) Get(ctx context.Context, input *GetProfileInput) (*ProfileOutput, error) {
	profile, err := h.service.Get(ctx, input.UserID)
	if err != nil {
		return nil, mapError(err)
	}
	return &ProfileOutput{Body: toProfileDTO(profile)}, nil
}

func (h *Handler) List(ctx context.Context, input *ListProfilesInput) (*httpplatform.PagedResponse[ProfileDTO], error) {
	limit := input.SafeLimit()
	result, err := h.service.List(ctx, limit, input.Offset())
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]ProfileDTO, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toListItemDTO(item))
	}
	return httpplatform.NewPagedResponse(items, input.Page, limit, result.Total), nil
}

func (h *Handler) ListCertifications(ctx context.Context, input *ListCertificationsInput) (*CertificationListOutput, error) {
	items, err := h.service.ListCertifications(ctx, input.UserID)
	if err != nil {
		return nil, mapError(err)
	}
	out := &CertificationListOutput{}
	out.Body.Data = make([]CertificationDTO, 0, len(items))
	for _, item := range items {
		out.Body.Data = append(out.Body.Data, toCertificationDTO(item))
	}
	return out, nil
}

func (h *Handler) CreateCertification(ctx context.Context, input *CreateCertificationInput) (*CertificationOutput, error) {
	issuedAt, expiredAt, err := certificationDates(input.Body)
	if err != nil {
		return nil, err
	}
	cert, err := h.service.CreateCertification(ctx, app.CreateCertificationCommand{
		UserID:    input.UserID,
		Name:      input.Body.Name,
		Issuer:    input.Body.Issuer,
		Number:    input.Body.Number,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &CertificationOutput{Body: toCertificationDTO(cert)}, nil
}

func (h *Handler) UpdateCertification(ctx context.Context, input *UpdateCertificationInput) (*CertificationOutput, error) {
	issuedAt, expiredAt, err := certificationPatchDates(input.Body)
	if err != nil {
		return nil, err
	}
	cert, err := h.service.PatchCertification(ctx, app.PatchCertificationCommand{
		ID:        input.CertID,
		Name:      input.Body.Name,
		Issuer:    input.Body.Issuer,
		Number:    input.Body.Number,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &CertificationOutput{Body: toCertificationDTO(cert)}, nil
}

func (h *Handler) DeleteCertification(ctx context.Context, input *DeleteCertificationInput) (*DeleteCertificationOutput, error) {
	if err := h.service.DeleteCertification(ctx, input.UserID, input.CertID); err != nil {
		return nil, mapError(err)
	}
	return &DeleteCertificationOutput{}, nil
}

// certificationDates memvalidasi & mengubah tanggal body menjadi time.Time.
func certificationDates(body CertificationBody) (*time.Time, *time.Time, error) {
	issuedAt, err := parseDate(body.IssuedAt)
	if err != nil {
		return nil, nil, huma.Error400BadRequest("issued_at must be a valid date (YYYY-MM-DD)")
	}
	expiredAt, err := parseDate(body.ExpiredAt)
	if err != nil {
		return nil, nil, huma.Error400BadRequest("expired_at must be a valid date (YYYY-MM-DD)")
	}
	return issuedAt, expiredAt, nil
}

func certificationPatchDates(body CertificationPatchBody) (*time.Time, *time.Time, error) {
	issuedAt, err := parseDate(body.IssuedAt)
	if err != nil {
		return nil, nil, huma.Error400BadRequest("issued_at must be a valid date (YYYY-MM-DD)")
	}
	expiredAt, err := parseDate(body.ExpiredAt)
	if err != nil {
		return nil, nil, huma.Error400BadRequest("expired_at must be a valid date (YYYY-MM-DD)")
	}
	return issuedAt, expiredAt, nil
}
