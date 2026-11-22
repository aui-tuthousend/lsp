package http

import (
	"context"

	"lsp/internal/features/uploads/app"
)

type Handler struct {
	service *app.Service
}

func NewHandler(service *app.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Presign(ctx context.Context, input *PresignInput) (*PresignOutput, error) {
	result, err := h.service.Presign(ctx, app.PresignCommand{
		Purpose:     input.Body.Purpose,
		ContentType: input.Body.ContentType,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &PresignOutput{Body: PresignResponse{
		ObjectKey: result.ObjectKey,
		UploadURL: result.UploadURL,
		ExpiresIn: result.ExpiresIn,
	}}, nil
}
