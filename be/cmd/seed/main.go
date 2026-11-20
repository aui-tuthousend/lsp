package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"golang.org/x/crypto/bcrypt"
	"lsp/internal/platform/config"
	"lsp/internal/platform/database"
	"lsp/internal/platform/id"
	"lsp/internal/users/adapters/postgres"
	"lsp/internal/users/domain"
)

type roleSeed struct {
	Code        domain.Role
	Name        string
	Description *string
}

type userSeed struct {
	Name            string
	Email           string
	Password        string
	Role            domain.Role
	Phone           *string
	NomorRegistrasi *string
}

// Master data role. Idempoten: role dengan code yang sama akan dilewati.
var roleSeeds = []roleSeed{
	{Code: domain.RoleAdmin, Name: "Admin"},
	{Code: domain.RoleAsesi, Name: "Asesi"},
	{Code: domain.RoleAsesor, Name: "Asesor"},
	{Code: domain.RoleValidator, Name: "Validator"},
	{Code: domain.RoleManagerSertifikasi, Name: "Manager Sertifikasi"},
}

// Data user awal untuk pengembangan lokal. Idempoten: user dengan email yang
// sama akan dilewati, sehingga aman dijalankan berulang kali.
var userSeeds = []userSeed{
	{Name: "Admin LSP", Email: "admin@lsp.local", Password: "password", Role: domain.RoleAdmin, Phone: strPtr("081200000001")},
	{Name: "Asesi Contoh", Email: "asesi@lsp.local", Password: "password", Role: domain.RoleAsesi, Phone: strPtr("081200000002"), NomorRegistrasi: strPtr("REG-ASESI-0001")},
	{Name: "Asesor Contoh", Email: "asesor@lsp.local", Password: "password", Role: domain.RoleAsesor, Phone: strPtr("081200000003"), NomorRegistrasi: strPtr("REG-ASESOR-0001")},
	{Name: "Validator Contoh", Email: "validator@lsp.local", Password: "password", Role: domain.RoleValidator, Phone: strPtr("081200000004")},
	{Name: "Manager Sertifikasi", Email: "manager@lsp.local", Password: "password", Role: domain.RoleManagerSertifikasi, Phone: strPtr("081200000005")},
}

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}
	pool, err := database.NewPostgresPool(ctx, cfg.PostgresDSN)
	if err != nil {
		slog.Error("connect postgres", "error", err)
		os.Exit(1)
	}
	
	defer pool.Close()

	roles := postgres.NewRolesRepository(pool)
	users := postgres.NewRepository(pool)

	roleIDs, err := seedRoles(ctx, roles)
	if err != nil {
		slog.Error("seed roles", "error", err)
		os.Exit(1)
	}
	for _, s := range userSeeds {
		if err := seedUser(ctx, users, roleIDs, s); err != nil {
			slog.Error("seed user", "email", s.Email, "error", err)
			os.Exit(1)
		}
	}
	slog.Info("seeding selesai", "roles", len(roleSeeds), "users", len(userSeeds))
}

// seedRoles memastikan seluruh master role tersedia dan mengembalikan map
// code -> id role untuk dipakai saat seeding user.
func seedRoles(ctx context.Context, roles *postgres.RolesRepository) (map[domain.Role]string, error) {
	ids := make(map[domain.Role]string, len(roleSeeds))
	for _, s := range roleSeeds {
		existing, err := roles.FindByCode(ctx, s.Code)
		if err == nil {
			ids[s.Code] = existing.ID
			slog.Info("role sudah ada, dilewati", "code", s.Code)
			continue
		}
		if !errors.Is(err, domain.ErrNotFound) {
			return nil, err
		}
		record := domain.RoleRecord{ID: id.New(), Code: s.Code, Name: s.Name, Description: s.Description}
		if err := roles.Create(ctx, record); err != nil {
			return nil, err
		}
		slog.Info("role dibuat", "code", s.Code, "id", record.ID)
		ids[s.Code] = record.ID
	}
	return ids, nil
}

func seedUser(ctx context.Context, users *postgres.Repository, roleIDs map[domain.Role]string, s userSeed) error {
	_, err := users.FindByEmail(ctx, s.Email)
	if err == nil {
		slog.Info("user sudah ada, dilewati", "email", s.Email, "role", s.Role)
		return nil
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return err
	}
	roleID, ok := roleIDs[s.Role]
	if !ok {
		return errors.New("role id tidak ditemukan untuk " + string(s.Role))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := domain.User{
		ID:              id.New(),
		RoleID:          roleID,
		Role:            s.Role,
		IsValid:         true,
		Name:            s.Name,
		Email:           s.Email,
		Phone:           s.Phone,
		NomorRegistrasi: s.NomorRegistrasi,
		PasswordHash:    string(hash),
	}
	if err := users.Create(ctx, user); err != nil {
		return err
	}
	slog.Info("user dibuat", "email", s.Email, "role", s.Role, "id", user.ID)
	return nil
}

func strPtr(value string) *string {
	return &value
}
