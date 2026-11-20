package bootstrap

import (
	"context"
	"log/slog"
	stdhttp "net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	notificationshttp "lsp/internal/notifications/adapters/http"
	"lsp/internal/platform/auth"
	"lsp/internal/platform/config"
	"lsp/internal/platform/events"
	httpplatform "lsp/internal/platform/http"
	"lsp/internal/platform/storage"
	usershttp "lsp/internal/users/adapters/http"
	userspg "lsp/internal/users/adapters/postgres"
	"lsp/internal/users/app"
)

type App struct {
	Config        config.Config
	Router        *stdhttp.ServeMux
	Handler       stdhttp.Handler
	API           huma.API
	Users         *UsersModule
	Notifications *NotificationsModule
	Schema        *SchemaModule
	APL01         *APL01Module
	APL02         *APL02Module
	Portfolio     *PortfolioModule
	ObjectStorage storage.ObjectStorage
}

type UsersModule struct {
	Repository *userspg.Repository
	Roles      *userspg.RolesRepository
	Service    *app.Service
	Handler    *usershttp.Handler
}

func NewApp(ctx context.Context, cfg config.Config, pool *pgxpool.Pool, rdb *redis.Client, objStorage storage.ObjectStorage) (*App, error) {
	mux := stdhttp.NewServeMux()
	humaConfig := huma.DefaultConfig("LSP API", "0.1.0")
	humaConfig.OpenAPIPath = "/openapi"
	humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{"bearerAuth": {Type: "http", Scheme: "bearer", BearerFormat: "JWT"}}
	humaConfig.Security = []map[string][]string{{"bearerAuth": {}}}
	api := humago.New(mux, humaConfig)

	jwtProvider := auth.NewJWTProvider(cfg.JWTSecret, cfg.JWTExpiry)
	if cfg.DisableAuth {
		api.UseMiddleware(auth.NewDevAuthMiddleware())
	} else {
		api.UseMiddleware(auth.NewJWTAuthMiddleware(jwtProvider))
	}

	bus := events.NewRedisBus(rdb)

	usersModule := newUsersModule(pool, jwtProvider)
	usershttp.RegisterRoutes(api, usersModule.Handler)
	notificationsModule := newNotificationsModule(ctx, bus)
	notificationshttp.RegisterRoutes(api, notificationsModule.Handler)
	registerHealth(api, pool, rdb)

	app := &App{
		Config:        cfg,
		Router:        mux,
		API:           api,
		Users:         usersModule,
		Notifications: notificationsModule,
		Schema:        newSchemaModule(),
		APL01:         newAPL01Module(),
		APL02:         newAPL02Module(),
		Portfolio:     newPortfolioModule(objStorage),
		ObjectStorage: objStorage,
	}
	app.Handler = httpplatform.Chain(mux,
		httpplatform.RequestID,
		httpplatform.Logger(slog.Default()),
		httpplatform.Recover(slog.Default()),
		httpplatform.CORS(),
		httpplatform.MaxBodySize(1<<20),
		httpplatform.Timeout(30*time.Second, "/events"),
	)
	return app, nil
}

func newUsersModule(pool *pgxpool.Pool, jwtProvider *auth.JWTProvider) *UsersModule {
	repo := userspg.NewRepository(pool)
	roles := userspg.NewRolesRepository(pool)
	service := app.NewService(repo, roles, jwtProvider)
	handler := usershttp.NewHandler(service)
	return &UsersModule{Repository: repo, Roles: roles, Service: service, Handler: handler}
}

type HealthInput struct{}

type HealthOutput struct {
	Body struct {
		Status string `json:"status"`
	} `json:"body"`
}

func registerHealth(api huma.API, pool *pgxpool.Pool, rdb *redis.Client) {
	huma.Get(api, "/health", healthHandler(pool, rdb), auth.PublicOperation, httpplatform.WithDoc("health", "Health check", "Return service health status.", "System"))
}

func healthHandler(pool *pgxpool.Pool, rdb *redis.Client) func(context.Context, *HealthInput) (*HealthOutput, error) {
	return func(ctx context.Context, input *HealthInput) (*HealthOutput, error) {
		if err := pool.Ping(ctx); err != nil {
			return nil, err
		}
		if rdb != nil {
			if err := rdb.Ping(ctx).Err(); err != nil {
				return nil, err
			}
		}
		out := &HealthOutput{}
		out.Body.Status = "ok"
		return out, nil
	}
}
