package http

import (
	"time"

	httpplatform "lsp/internal/platform/http"
	"lsp/internal/tuk/domain"
)

type CreateInput struct {
	Body TUKBody `json:"body"`
}

type UpdateInput struct {
	ID   string  `path:"id"`
	Body TUKBody `json:"body"`
}

type GetInput struct {
	ID string `path:"id"`
}

type DeleteInput struct {
	ID string `path:"id"`
}

type ListInput struct {
	httpplatform.PaginationQuery
}

type TUKBody struct {
	Name       string  `json:"name" minLength:"2" maxLength:"255"`
	Logo       *string `json:"logo,omitempty" doc:"URL/path logo TUK"`
	PostalCode string  `json:"postal_code" minLength:"3" maxLength:"10"`
	Phone      string  `json:"phone" minLength:"5" maxLength:"20"`
	Email      string  `json:"email" format:"email" maxLength:"255"`
	Province   string  `json:"province" minLength:"2" maxLength:"120"`
	City       string  `json:"city" minLength:"2" maxLength:"120"`
	Address    string  `json:"address" minLength:"5"`
	MapsURL    string  `json:"maps_url" minLength:"5" doc:"Link Google Maps atau provider maps lain"`
	Letterhead *string `json:"letterhead,omitempty" doc:"URL/path kop surat TUK"`
}

type TUKOutput struct {
	Body TUKDTO `json:"body"`
}

type TUKDTO struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Logo       *string    `json:"logo,omitempty"`
	PostalCode string     `json:"postal_code"`
	Phone      string     `json:"phone"`
	Email      string     `json:"email"`
	Province   string     `json:"province"`
	City       string     `json:"city"`
	Address    string     `json:"address"`
	MapsURL    string     `json:"maps_url"`
	Letterhead *string    `json:"letterhead,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

type DeleteOutput struct{}

func toCreateCommand(body TUKBody) createCommandDTO {
	return createCommandDTO{
		Name:       body.Name,
		Logo:       body.Logo,
		PostalCode: body.PostalCode,
		Phone:      body.Phone,
		Email:      body.Email,
		Province:   body.Province,
		City:       body.City,
		Address:    body.Address,
		MapsURL:    body.MapsURL,
		Letterhead: body.Letterhead,
	}
}

type createCommandDTO struct {
	Name       string
	Logo       *string
	PostalCode string
	Phone      string
	Email      string
	Province   string
	City       string
	Address    string
	MapsURL    string
	Letterhead *string
}

func toDTO(tuk domain.TUK) TUKDTO {
	return TUKDTO{
		ID:         tuk.ID,
		Name:       tuk.Name,
		Logo:       tuk.Logo,
		PostalCode: tuk.PostalCode,
		Phone:      tuk.Phone,
		Email:      tuk.Email,
		Province:   tuk.Province,
		City:       tuk.City,
		Address:    tuk.Address,
		MapsURL:    tuk.MapsURL,
		Letterhead: tuk.Letterhead,
		CreatedAt:  tuk.CreatedAt,
		UpdatedAt:  tuk.UpdatedAt,
		DeletedAt:  tuk.DeletedAt,
	}
}
