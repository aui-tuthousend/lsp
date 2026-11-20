package http

import (
	"time"

	"lsp/internal/asesor/domain"
	httpplatform "lsp/internal/platform/http"
)

type UpsertProfileInput struct {
	UserID string            `path:"user_id"`
	Body   UpsertProfileBody `json:"body"`
}

type UpsertProfileBody struct {
	TUKID                *string `json:"tuk_id,omitempty" doc:"ID TUK tempat asesor berafiliasi"`
	Province             *string `json:"province,omitempty" maxLength:"120" doc:"Domisili asesor (provinsi)"`
	City                 *string `json:"city,omitempty" maxLength:"120" doc:"Domisili asesor (kota)"`
	DailyAssignmentLimit int     `json:"daily_assignment_limit" minimum:"0" doc:"Batas penugasan harian (0 = tanpa batas)"`
}

type GetProfileInput struct {
	UserID string `path:"user_id"`
}

type ListProfilesInput struct {
	httpplatform.PaginationQuery
}

type ProfileOutput struct {
	Body ProfileDTO `json:"body"`
}

type ProfileDTO struct {
	UserID               string     `json:"user_id"`
	Name                 string     `json:"name,omitempty"`
	Email                string     `json:"email,omitempty"`
	TUKID                *string    `json:"tuk_id,omitempty"`
	TUKName              *string    `json:"tuk_name,omitempty"`
	Province             *string    `json:"province,omitempty"`
	City                 *string    `json:"city,omitempty"`
	DailyAssignmentLimit int        `json:"daily_assignment_limit"`
	CreatedAt            *time.Time `json:"created_at,omitempty"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty"`
}

func toProfileDTO(profile domain.Profile) ProfileDTO {
	createdAt := profile.CreatedAt
	updatedAt := profile.UpdatedAt
	return ProfileDTO{
		UserID:               profile.UserID,
		TUKID:                profile.TUKID,
		Province:             profile.Province,
		City:                 profile.City,
		DailyAssignmentLimit: profile.DailyAssignmentLimit,
		CreatedAt:            &createdAt,
		UpdatedAt:            &updatedAt,
	}
}

func toListItemDTO(item domain.ListItem) ProfileDTO {
	return ProfileDTO{
		UserID:               item.UserID,
		Name:                 item.Name,
		Email:                item.Email,
		TUKID:                item.TUKID,
		TUKName:              item.TUKName,
		Province:             item.Province,
		City:                 item.City,
		DailyAssignmentLimit: item.DailyAssignmentLimit,
	}
}
