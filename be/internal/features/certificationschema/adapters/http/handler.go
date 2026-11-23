package http

import (
	"context"

	"lsp/internal/features/certificationschema/app"
	httpplatform "lsp/internal/platform/http"
)

type Handler struct {
	service *app.Service
}

func NewHandler(service *app.Service) *Handler {
	return &Handler{service: service}
}

// --- Skema sertifikasi ---

func (h *Handler) CreateSchema(ctx context.Context, input *CreateSchemaInput) (*SchemaOutput, error) {
	schema, err := h.service.CreateSchema(ctx, app.CreateSchemaCommand{
		Code:       input.Body.Code,
		Name:       input.Body.Name,
		SchemaCode: input.Body.SchemaCode,
		SortOrder:  input.Body.SortOrder,
		IsActive:   input.Body.IsActive,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &SchemaOutput{Body: toSchemaDTO(schema)}, nil
}

func (h *Handler) UpdateSchema(ctx context.Context, input *UpdateSchemaInput) (*SchemaOutput, error) {
	schema, err := h.service.PatchSchema(ctx, app.PatchSchemaCommand{
		ID:         input.ID,
		Code:       input.Body.Code,
		Name:       input.Body.Name,
		SchemaCode: input.Body.SchemaCode,
		SortOrder:  input.Body.SortOrder,
		IsActive:   input.Body.IsActive,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &SchemaOutput{Body: toSchemaDTO(schema)}, nil
}

// GetSchema mengembalikan skema beserta struktur unit -> elemen -> kriteria.
func (h *Handler) GetSchema(ctx context.Context, input *GetSchemaInput) (*SchemaTreeOutput, error) {
	tree, err := h.service.GetSchemaTree(ctx, input.ID)
	if err != nil {
		return nil, mapError(err)
	}
	return &SchemaTreeOutput{Body: toSchemaTreeDTO(tree)}, nil
}

func (h *Handler) ListSchemas(ctx context.Context, input *ListSchemasInput) (*httpplatform.PagedResponse[SchemaDTO], error) {
	limit := input.SafeLimit()
	result, err := h.service.ListSchemas(ctx, limit, input.Offset())
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]SchemaDTO, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toSchemaDTO(item))
	}
	return httpplatform.NewPagedResponse(items, input.Page, limit, result.Total), nil
}

func (h *Handler) DeleteSchema(ctx context.Context, input *DeleteSchemaInput) (*DeleteOutput, error) {
	if err := h.service.DeleteSchema(ctx, input.ID); err != nil {
		return nil, mapError(err)
	}
	return &DeleteOutput{}, nil
}

// --- Unit kompetensi ---

func (h *Handler) CreateUnit(ctx context.Context, input *CreateUnitInput) (*UnitOutput, error) {
	unit, err := h.service.CreateUnit(ctx, app.CreateUnitCommand{
		CertificationSchemaID: input.SchemaID,
		Code:                  input.Body.Code,
		Title:                 input.Body.Title,
		SortOrder:             input.Body.SortOrder,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &UnitOutput{Body: toUnitDTO(unit)}, nil
}

func (h *Handler) UpdateUnit(ctx context.Context, input *UpdateUnitInput) (*UnitOutput, error) {
	unit, err := h.service.PatchUnit(ctx, app.PatchUnitCommand{
		ID:        input.ID,
		Code:      input.Body.Code,
		Title:     input.Body.Title,
		SortOrder: input.Body.SortOrder,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &UnitOutput{Body: toUnitDTO(unit)}, nil
}

func (h *Handler) ListUnits(ctx context.Context, input *ListUnitsInput) (*httpplatform.PagedResponse[UnitDTO], error) {
	limit := input.SafeLimit()
	result, err := h.service.ListUnits(ctx, input.SchemaID, limit, input.Offset())
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]UnitDTO, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toUnitDTO(item))
	}
	return httpplatform.NewPagedResponse(items, input.Page, limit, result.Total), nil
}

func (h *Handler) DeleteUnit(ctx context.Context, input *DeleteUnitInput) (*DeleteOutput, error) {
	if err := h.service.DeleteUnit(ctx, input.ID); err != nil {
		return nil, mapError(err)
	}
	return &DeleteOutput{}, nil
}

// --- Elemen unit kompetensi ---

func (h *Handler) CreateElement(ctx context.Context, input *CreateElementInput) (*ElementOutput, error) {
	element, err := h.service.CreateElement(ctx, app.CreateElementCommand{
		CompetencyUnitID: input.UnitID,
		Title:            input.Body.Title,
		SortOrder:        input.Body.SortOrder,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &ElementOutput{Body: toElementDTO(element)}, nil
}

func (h *Handler) UpdateElement(ctx context.Context, input *UpdateElementInput) (*ElementOutput, error) {
	element, err := h.service.PatchElement(ctx, app.PatchElementCommand{
		ID:        input.ID,
		Title:     input.Body.Title,
		SortOrder: input.Body.SortOrder,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &ElementOutput{Body: toElementDTO(element)}, nil
}

func (h *Handler) ListElements(ctx context.Context, input *ListElementsInput) (*httpplatform.PagedResponse[ElementDTO], error) {
	limit := input.SafeLimit()
	result, err := h.service.ListElements(ctx, input.UnitID, limit, input.Offset())
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]ElementDTO, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toElementDTO(item))
	}
	return httpplatform.NewPagedResponse(items, input.Page, limit, result.Total), nil
}

func (h *Handler) DeleteElement(ctx context.Context, input *DeleteElementInput) (*DeleteOutput, error) {
	if err := h.service.DeleteElement(ctx, input.ID); err != nil {
		return nil, mapError(err)
	}
	return &DeleteOutput{}, nil
}

// --- Kriteria unjuk kerja (KUK) ---

func (h *Handler) CreateCriteria(ctx context.Context, input *CreateCriteriaInput) (*CriteriaOutput, error) {
	criteria, err := h.service.CreateCriteria(ctx, app.CreateCriteriaCommand{
		CompetencyUnitElementID: input.ElementID,
		Content:                 input.Body.Content,
		SortOrder:               input.Body.SortOrder,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &CriteriaOutput{Body: toCriteriaDTO(criteria)}, nil
}

func (h *Handler) UpdateCriteria(ctx context.Context, input *UpdateCriteriaInput) (*CriteriaOutput, error) {
	criteria, err := h.service.PatchCriteria(ctx, app.PatchCriteriaCommand{
		ID:        input.ID,
		Content:   input.Body.Content,
		SortOrder: input.Body.SortOrder,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &CriteriaOutput{Body: toCriteriaDTO(criteria)}, nil
}

func (h *Handler) ListCriteria(ctx context.Context, input *ListCriteriaInput) (*httpplatform.PagedResponse[CriteriaDTO], error) {
	limit := input.SafeLimit()
	result, err := h.service.ListCriteria(ctx, input.ElementID, limit, input.Offset())
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]CriteriaDTO, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toCriteriaDTO(item))
	}
	return httpplatform.NewPagedResponse(items, input.Page, limit, result.Total), nil
}

func (h *Handler) DeleteCriteria(ctx context.Context, input *DeleteCriteriaInput) (*DeleteOutput, error) {
	if err := h.service.DeleteCriteria(ctx, input.ID); err != nil {
		return nil, mapError(err)
	}
	return &DeleteOutput{}, nil
}
