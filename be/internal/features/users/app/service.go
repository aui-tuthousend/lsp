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

type PatchUserCommand struct {
	ID              string
	Name            *string
	Phone           *string
	NomorRegistrasi *string
	Password        *string
	IsValid         *bool
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

// PatchUser mengubah sebagian data user (dipakai admin). Bila password diisi,
// password akan di-hash ulang sebelum disimpan.
func (s *Service) PatchUser(ctx context.Context, cmd PatchUserCommand) (domain.User, error) {
	var passwordHash *string
	if cmd.Password != nil && *cmd.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*cmd.Password), bcrypt.DefaultCost)
		if err != nil {
			return domain.User{}, err
		}
		value := string(hash)
		passwordHash = &value
	}
	return s.repo.Patch(ctx, domain.UserPatch{
		ID:              cmd.ID,
		Name:            cmd.Name,
		Phone:           cmd.Phone,
		NomorRegistrasi: cmd.NomorRegistrasi,
		PasswordHash:    passwordHash,
		IsValid:         cmd.IsValid,
	})
}

// DeleteUser melakukan soft delete pada user (dipakai admin).
func (s *Service) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// GetByID mengambil satu user berdasarkan id (dipakai admin).
func (s *Service) GetByID(ctx context.Context, id string) (domain.User, error) {
	return s.repo.FindByID(ctx, id)
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

func (s *Service) List(ctx context.Context, limit, offset int, roles []string) (domain.ListResult, error) {
	return s.repo.List(ctx, domain.ListParams{Limit: limit, Offset: offset, Roles: roles})
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

// --- Master data roles ---

type CreateRoleCommand struct {
	Code        domain.Role
	Name        string
	Description *string
}

func (s *Service) CreateRole(ctx context.Context, cmd CreateRoleCommand) (domain.RoleRecord, error) {
	role := domain.RoleRecord{ID: id.New(), Code: cmd.Code, Name: cmd.Name, Description: cmd.Description}
	if err := s.roles.Create(ctx, role); err != nil {
		return domain.RoleRecord{}, err
	}
	return s.roles.FindByID(ctx, role.ID)
}

type PatchRoleCommand struct {
	ID          string
	Code        *domain.Role
	Name        *string
	Description *string
}

func (s *Service) PatchRole(ctx context.Context, cmd PatchRoleCommand) (domain.RoleRecord, error) {
	return s.roles.Patch(ctx, domain.RolePatch{
		ID:          cmd.ID,
		Code:        cmd.Code,
		Name:        cmd.Name,
		Description: cmd.Description,
	})
}

func (s *Service) GetRole(ctx context.Context, id string) (domain.RoleRecord, error) {
	return s.roles.FindByID(ctx, id)
}

func (s *Service) ListRoles(ctx context.Context) ([]domain.RoleRecord, error) {
	return s.roles.List(ctx)
}

func (s *Service) DeleteRole(ctx context.Context, id string) error {
	return s.roles.Delete(ctx, id)
}
