package app

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"lsp/internal/platform/auth"
	"lsp/internal/platform/id"
	"lsp/internal/users/domain"
)

type Service struct {
	repo domain.Repository
	jwt  *auth.JWTProvider
}

type RegisterCommand struct {
	Email           string
	Password        string
	Name            string
	Role            domain.Role
	Phone           *string
	NomorRegistrasi *string
}

type LoginCommand struct {
	Email    string
	Password string
}

type AuthResult struct {
	User  domain.User
	Token string
}

func NewService(repo domain.Repository, jwt *auth.JWTProvider) *Service {
	return &Service{repo: repo, jwt: jwt}
}

func (s *Service) Register(ctx context.Context, cmd RegisterCommand) (AuthResult, error) {
	role := cmd.Role
	if role == "" {
		role = domain.RoleAsesi
	}
	if !domain.IsValidRole(role) {
		return AuthResult{}, domain.ErrInvalidRole
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResult{}, err
	}
	user := domain.User{ID: id.New(), Email: cmd.Email, PasswordHash: string(hash), Name: cmd.Name, Role: role, IsValid: true, Phone: cmd.Phone, NomorRegistrasi: cmd.NomorRegistrasi}
	if err := s.repo.Create(ctx, user); err != nil {
		return AuthResult{}, err
	}
	return s.authResult(user)
}

func (s *Service) Login(ctx context.Context, cmd LoginCommand) (AuthResult, error) {
	user, err := s.repo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		if err == domain.ErrNotFound {
			return AuthResult{}, domain.ErrInvalidCredentials
		}
		return AuthResult{}, err
	}
	if !user.IsValid || user.DeletedAt != nil {
		return AuthResult{}, domain.ErrInactiveUser
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(cmd.Password)); err != nil {
		return AuthResult{}, domain.ErrInvalidCredentials
	}
	return s.authResult(user)
}

func (s *Service) GetMe(ctx context.Context, claims auth.UserClaims) (domain.User, error) {
	return s.repo.FindByID(ctx, claims.UserID)
}

func (s *Service) List(ctx context.Context, limit, offset int) (domain.ListResult, error) {
	return s.repo.List(ctx, domain.ListParams{Limit: limit, Offset: offset})
}

func (s *Service) authResult(user domain.User) (AuthResult, error) {
	token, err := s.jwt.GenerateToken(auth.UserClaims{UserID: user.ID, Email: user.Email, Role: string(user.Role), Name: user.Name})
	if err != nil {
		return AuthResult{}, err
	}
	return AuthResult{User: user, Token: token}, nil
}
