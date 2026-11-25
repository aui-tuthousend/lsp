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
	apl01http "lsp/internal/features/apl01/adapters/http"
	asesihttp "lsp/internal/features/asesi/adapters/http"
	asesorhttp "lsp/internal/features/asesor/adapters/http"
	asesorpg "lsp/internal/features/asesor/adapters/postgres"
	schemahttp "lsp/internal/features/certificationschema/adapters/http"
	docreqhttp "lsp/internal/features/documentrequirements/adapters/http"
	notificationshttp "lsp/internal/features/notifications/adapters/http"
	registrationhttp "lsp/internal/features/registration/adapters/http"
	tukhttp "lsp/internal/features/tuk/adapters/http"
	uploadshttp "lsp/internal/features/uploads/adapters/http"
	usershttp "lsp/internal/features/users/adapters/http"
	"lsp/internal/platform/auth"
	"lsp/internal/platform/config"
	"lsp/internal/platform/database"
	"lsp/internal/platform/events"
	httpplatform "lsp/internal/platform/http"
	"lsp/internal/platform/storage"
)

type App struct {
	Config          config.Config
	Router          *stdhttp.ServeMux
	Handler         stdhttp.Handler
	API             huma.API
	Users           *UsersModule
	Notifications   *NotificationsModule
	Registration    *RegistrationModule
	Asesi           *AsesiModule
	Uploads         *UploadsModule
	Asesor          *AsesorModule
	TUK             *TUKModule
	Schema          *SchemaModule
	DocRequirements *DocumentRequirementsModule
	APL01           *APL01Module
	APL02           *APL02Module
	Portfolio       *PortfolioModule
	ObjectStorage   storage.ObjectStorage
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
	txManager := database.NewTxManager(pool)

	asesorRepo := asesorpg.NewRepository(pool)
	usersModule := newUsersModule(pool, asesorRepo, jwtProvider)
	usershttp.RegisterRoutes(api, usersModule.Handler)
	asesorModule := newAsesorModule(asesorRepo, usersModule.Repository)
	asesorhttp.RegisterRoutes(api, asesorModule.Handler)
	asesiModule := newAsesiModule(pool, usersModule.Repository, objStorage)
	asesihttp.RegisterRoutes(api, asesiModule.Handler)
	registrationModule := newRegistrationModule(usersModule.Service, asesiModule.Service, txManager)
	registrationhttp.RegisterRoutes(api, registrationModule.Handler)
	docReqModule := newDocumentRequirementsModule(pool, txManager)
	docreqhttp.RegisterRoutes(api, docReqModule.Handler)
	uploadsModule := newUploadsModule(objStorage, docReqModule.Repository)
	uploadshttp.RegisterRoutes(api, uploadsModule.Handler)
	tukModule := newTUKModule(pool)
	tukhttp.RegisterRoutes(api, tukModule.Handler)
	schemaModule := newSchemaModule(pool)
	schemahttp.RegisterRoutes(api, schemaModule.Handler)
	apl01Module := newAPL01Module(pool, txManager, bus)
	apl01http.RegisterRoutes(api, apl01Module.Handler)
	notificationsModule := newNotificationsModule(ctx, bus)
	notificationshttp.RegisterRoutes(api, notificationsModule.Handler)
	registerHealth(api, pool, rdb)

	app := &App{
		Config:          cfg,
		Router:          mux,
		API:             api,
		Users:           usersModule,
		Notifications:   notificationsModule,
		Registration:    registrationModule,
		Asesi:           asesiModule,
		Uploads:         uploadsModule,
		Asesor:          asesorModule,
		TUK:             tukModule,
		Schema:          schemaModule,
		DocRequirements: docReqModule,
		APL01:           apl01Module,
		APL02:           newAPL02Module(),
		Portfolio:       newPortfolioModule(objStorage),
		ObjectStorage:   objStorage,
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
