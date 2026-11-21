package app

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"lsp/internal/features/users/domain"
	"lsp/internal/platform/auth"
	"lsp/internal/platform/id"
)

type Service struct {
	repo       domain.Repository
	roles      domain.RolesRepository
	asesorInit domain.AsesorProfileInitializer
	jwt        *auth.JWTProvider
}

type RegisterCommand struct {
	Email           string
	Password        string
	Name            string
	Phone           *string
	NomorRegistrasi *string
}

type CreateUserCommand struct {
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

func NewService(repo domain.Repository, roles domain.RolesRepository, asesorInit domain.AsesorProfileInitializer, jwt *auth.JWTProvider) *Service {
	return &Service{repo: repo, roles: roles, asesorInit: asesorInit, jwt: jwt}
}

// Register adalah pendaftaran mandiri (endpoint publik) dan hanya membuat asesi.
func (s *Service) Register(ctx context.Context, cmd RegisterCommand) (AuthResult, error) {
	user, err := s.newUser(ctx, cmd.Email, cmd.Password, cmd.Name, domain.RoleAsesi, cmd.Phone, cmd.NomorRegistrasi)
	if err != nil {
		return AuthResult{}, err
	}
	return s.authResult(user)
}

// CreateUser dipakai admin untuk membuat user dengan role apa pun.
func (s *Service) CreateUser(ctx context.Context, cmd CreateUserCommand) (domain.User, error) {
	user, err := s.newUser(ctx, cmd.Email, cmd.Password, cmd.Name, cmd.Role, cmd.Phone, cmd.NomorRegistrasi)
	if err != nil {
		return domain.User{}, err
	}
	if user.Role == domain.RoleAsesor && s.asesorInit != nil {
		if err := s.asesorInit.EnsureProfile(ctx, user.ID); err != nil {
			return domain.User{}, err
		}
	}
	return user, nil
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

// newUser memvalidasi role, me-resolve role_id dari master data, meng-hash
// password, lalu menyimpan user.
func (s *Service) newUser(ctx context.Context, email, password, name string, roleCode domain.Role, phone, nomorRegistrasi *string) (domain.User, error) {
	if !domain.IsValidRole(roleCode) {
		return domain.User{}, domain.ErrInvalidRole
	}
	role, err := s.roles.FindByCode(ctx, roleCode)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.User{}, domain.ErrInvalidRole
		}
		return domain.User{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}
	user := domain.User{
		ID:              id.New(),
		RoleID:          role.ID,
		Role:            role.Code,
		Email:           email,
		PasswordHash:    string(hash),
		Name:            name,
		IsValid:         true,
		Phone:           phone,
		NomorRegistrasi: nomorRegistrasi,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *Service) authResult(user domain.User) (AuthResult, error) {
	token, err := s.jwt.GenerateToken(auth.UserClaims{UserID: user.ID, Email: user.Email, Role: string(user.Role), Name: user.Name})
	if err != nil {
		return AuthResult{}, err
	}
	return AuthResult{User: user, Token: token}, nil
}
