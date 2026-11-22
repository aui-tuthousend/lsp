package app

import (
	"context"
	"time"

	asesiapp "lsp/internal/features/asesi/app"
	asesidomain "lsp/internal/features/asesi/domain"
	usersapp "lsp/internal/features/users/app"
	"lsp/internal/platform/database"
)

// Service mengorkestrasi pendaftaran mandiri asesi: membuat user (role asesi)
// dan profil asesi dalam satu transaksi database.
type Service struct {
	users *usersapp.Service
	asesi *asesiapp.Service
	tx    *database.TxManager
}

func NewService(users *usersapp.Service, asesi *asesiapp.Service, tx *database.TxManager) *Service {
	return &Service{users: users, asesi: asesi, tx: tx}
}

// RegisterCommand membawa data akun sekaligus data profil asesi.
type RegisterCommand struct {
	// Akun (dipakai untuk membuat user).
	Email    string
	Password string
	FullName string
	Phone    string

	// Profil asesi.
	NIK               string
	Gender            asesidomain.Gender
	BirthPlace        string
	BirthDate         time.Time
	Nationality       string
	Occupation        string
	EducationLevel    asesidomain.EducationLevel
	KTPAddress        string
	DomicileSameAsKTP bool
	DomicileAddress   *string
	Province          string
	City              string
	District          string
	Village           string
	PostalCode        *string
	Photo             *string
	KTPFile           *string
	CVFile            *string
	NPWPFile          *string
}

type RegisterResult struct {
	Token  string
	UserID string
	Email  string
	Name   string
	Role   string
}

func (s *Service) Register(ctx context.Context, cmd RegisterCommand) (RegisterResult, error) {
	var result RegisterResult
	err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		auth, err := s.users.Register(txCtx, usersapp.RegisterCommand{
			Email:    cmd.Email,
			Password: cmd.Password,
			Name:     cmd.FullName,
			Phone:    &cmd.Phone,
		})
		if err != nil {
			return err
		}
		if _, err := s.asesi.Create(txCtx, asesiapp.CreateCommand{
			UserID:            auth.User.ID,
			FullName:          cmd.FullName,
			NIK:               cmd.NIK,
			Phone:             cmd.Phone,
			Gender:            cmd.Gender,
			BirthPlace:        cmd.BirthPlace,
			BirthDate:         cmd.BirthDate,
			Nationality:       cmd.Nationality,
			Occupation:        cmd.Occupation,
			EducationLevel:    cmd.EducationLevel,
			KTPAddress:        cmd.KTPAddress,
			DomicileSameAsKTP: cmd.DomicileSameAsKTP,
			DomicileAddress:   cmd.DomicileAddress,
			Province:          cmd.Province,
			City:              cmd.City,
			District:          cmd.District,
			Village:           cmd.Village,
			PostalCode:        cmd.PostalCode,
			Photo:             cmd.Photo,
			KTPFile:           cmd.KTPFile,
			CVFile:            cmd.CVFile,
			NPWPFile:          cmd.NPWPFile,
		}); err != nil {
			return err
		}
		result = RegisterResult{
			Token:  auth.Token,
			UserID: auth.User.ID,
			Email:  auth.User.Email,
			Name:   auth.User.Name,
			Role:   string(auth.User.Role),
		}
		return nil
	})
	if err != nil {
		return RegisterResult{}, err
	}
	return result, nil
}
