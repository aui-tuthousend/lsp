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

type seedUser struct {
	Name            string
	Email           string
	Password        string
	Role            domain.Role
	Phone           *string
	NomorRegistrasi *string
}

// Data user awal untuk pengembangan lokal. Idempoten: user dengan email yang
// sama akan dilewati, sehingga aman dijalankan berulang kali.
var seeds = []seedUser{
	{Name: "Admin LSP", Email: "admin@lsp.local", Password: "password", Role: domain.RoleAdmin, Phone: strPtr("081200000001")},
	{Name: "Asesi Contoh", Email: "asesi@lsp.local", Password: "password", Role: domain.RoleAsesi, Phone: strPtr("081200000002"), NomorRegistrasi: strPtr("REG-ASESI-0001")},
	{Name: "Asesor Contoh", Email: "asesor@lsp.local", Password: "password", Role: domain.RoleAsesor, Phone: strPtr("081200000003"), NomorRegistrasi: strPtr("REG-ASESOR-0001")},
	{Name: "Validator Contoh", Email: "validator@lsp.local", Password: "password", Role: domain.RoleValidator, Phone: strPtr("081200000004")},
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

	repo := postgres.NewRepository(pool)
	for _, s := range seeds {
		if err := seedOne(ctx, repo, s); err != nil {
			slog.Error("seed user", "email", s.Email, "error", err)
			os.Exit(1)
		}
	}
	slog.Info("seeding selesai", "total", len(seeds))
}

func seedOne(ctx context.Context, repo *postgres.Repository, s seedUser) error {
	_, err := repo.FindByEmail(ctx, s.Email)
	if err == nil {
		slog.Info("user sudah ada, dilewati", "email", s.Email, "role", s.Role)
		return nil
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := domain.User{
		ID:              id.New(),
		Role:            s.Role,
		IsValid:         true,
		Name:            s.Name,
		Email:           s.Email,
		Phone:           s.Phone,
		NomorRegistrasi: s.NomorRegistrasi,
		PasswordHash:    string(hash),
	}
	if err := repo.Create(ctx, user); err != nil {
		return err
	}
	slog.Info("user dibuat", "email", s.Email, "role", s.Role, "id", user.ID)
	return nil
}

func strPtr(value string) *string {
	return &value
}
