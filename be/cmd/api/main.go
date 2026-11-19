package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lsp/internal/bootstrap"
	"lsp/internal/platform/cache"
	"lsp/internal/platform/config"
	"lsp/internal/platform/database"
	"lsp/internal/platform/storage"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	rds := cache.NewRedisClient(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)
	defer rds.Close()
	if err := rds.Ping(ctx); err != nil {
		slog.Error("connect redis", "error", err)
		os.Exit(1)
	}

	var objStorage storage.ObjectStorage
	if cfg.HasR2Config() {
		objStorage, err = storage.NewS3Storage(ctx, cfg)
		if err != nil {
			slog.Error("init storage", "error", err)
			os.Exit(1)
		}
	}

	app, err := bootstrap.NewApp(cfg, pool, rds.Client, objStorage)
	if err != nil {
		slog.Error("bootstrap app", "error", err)
		os.Exit(1)
	}
	server := &http.Server{Addr: cfg.HTTPAddr, Handler: app.Handler}
	go func() {
		slog.Info("api listening", "addr", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server", "error", err)
			os.Exit(1)
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown", "error", err)
	}
}
