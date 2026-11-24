package http

import (
	"context"

	"lsp/internal/features/documentrequirements/app"
	httpplatform "lsp/internal/platform/http"
)

type Handler struct {
	service *app.Service
}

func NewHandler(service *app.Service) *Handler {
	return &Handler{service: service}
}

// --- Master jenis dokumen ---

func (h *Handler) CreateDocumentType(ctx context.Context, input *CreateDocumentTypeInput) (*DocumentTypeOutput, error) {
	docType, err := h.service.CreateDocumentType(ctx, app.CreateDocumentTypeCommand{
		Code:  input.Body.Code,
		Label: input.Body.Label,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &DocumentTypeOutput{Body: toDocumentTypeDTO(docType)}, nil
}

func (h *Handler) UpdateDocumentType(ctx context.Context, input *UpdateDocumentTypeInput) (*DocumentTypeOutput, error) {
	docType, err := h.service.PatchDocumentType(ctx, app.PatchDocumentTypeCommand{
		ID:    input.ID,
		Code:  input.Body.Code,
		Label: input.Body.Label,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &DocumentTypeOutput{Body: toDocumentTypeDTO(docType)}, nil
}

func (h *Handler) GetDocumentType(ctx context.Context, input *GetDocumentTypeInput) (*DocumentTypeOutput, error) {
	docType, err := h.service.GetDocumentType(ctx, input.ID)
	if err != nil {
		return nil, mapError(err)
	}
	return &DocumentTypeOutput{Body: toDocumentTypeDTO(docType)}, nil
}

func (h *Handler) ListDocumentTypes(ctx context.Context, input *ListDocumentTypesInput) (*httpplatform.PagedResponse[DocumentTypeDTO], error) {
	limit := input.SafeLimit()
	result, err := h.service.ListDocumentTypes(ctx, limit, input.Offset())
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]DocumentTypeDTO, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toDocumentTypeDTO(item))
	}
	return httpplatform.NewPagedResponse(items, input.Page, limit, result.Total), nil
}

func (h *Handler) DeleteDocumentType(ctx context.Context, input *DeleteDocumentTypeInput) (*DeleteOutput, error) {
	if err := h.service.DeleteDocumentType(ctx, input.ID); err != nil {
		return nil, mapError(err)
	}
	return &DeleteOutput{}, nil
}

// --- Persyaratan dokumen per skema ---

func (h *Handler) ListRequirements(ctx context.Context, input *ListRequirementsInput) (*RequirementsOutput, error) {
	requirements, err := h.service.ListRequirements(ctx, input.SchemaID)
	if err != nil {
		return nil, mapError(err)
	}
	return &RequirementsOutput{Body: RequirementsBody{Data: toRequirementDTOs(requirements)}}, nil
}

func (h *Handler) ReplaceRequirements(ctx context.Context, input *ReplaceRequirementsInput) (*RequirementsOutput, error) {
	items := make([]app.ReplaceRequirementItem, 0, len(input.Body.Items))
	for _, item := range input.Body.Items {
		items = append(items, app.ReplaceRequirementItem{
			DocumentTypeID: item.DocumentTypeID,
			IsRequired:     item.IsRequired,
		})
	}
	requirements, err := h.service.ReplaceRequirements(ctx, app.ReplaceRequirementsCommand{
		CertificationSchemaID: input.SchemaID,
		Items:                 items,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &RequirementsOutput{Body: RequirementsBody{Data: toRequirementDTOs(requirements)}}, nil
}
