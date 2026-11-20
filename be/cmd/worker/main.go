package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"lsp/internal/platform/cache"
	"lsp/internal/platform/config"
	"lsp/internal/platform/database"
	"lsp/internal/platform/events"
	"lsp/internal/platform/mailer"
	"lsp/internal/platform/worker"
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

	var mailerClient mailer.Mailer = mailer.NewLogMailer(slog.Default())
	if cfg.MailerDriver == "mailgun" {
		mailerClient = mailer.NewMailgunMailer(cfg.MailgunDomain, cfg.MailgunAPIKey, cfg.MailgunSender)
	}
	bus := events.NewRedisBus(rds.Client)
	processor := worker.NewOutboxProcessor(database.NewTxManager(pool), database.NewOutboxRepository(pool), mailerClient, bus, slog.Default())
	if err := processor.ProcessOnce(ctx); err != nil {
		slog.Error("process outbox", "error", err)
		os.Exit(1)
	}
}
