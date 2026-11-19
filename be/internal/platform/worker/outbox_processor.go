package worker

import (
	"context"
	"log/slog"
	"time"

	"lsp/internal/platform/database"
	"lsp/internal/platform/mailer"
)

type OutboxProcessor struct {
	tx      *database.TxManager
	outbox  *database.OutboxRepository
	mailer  mailer.Mailer
	logger  *slog.Logger
	limit   int
	backoff time.Duration
}

func NewOutboxProcessor(tx *database.TxManager, outbox *database.OutboxRepository, mailer mailer.Mailer, logger *slog.Logger) *OutboxProcessor {
	return &OutboxProcessor{tx: tx, outbox: outbox, mailer: mailer, logger: logger, limit: 10, backoff: time.Minute}
}

func (p *OutboxProcessor) ProcessOnce(ctx context.Context) error {
	return p.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		events, err := p.outbox.FetchPendingForUpdate(txCtx, p.limit)
		if err != nil {
			return err
		}
		for _, event := range events {
			_ = p.mailer
			p.logger.InfoContext(txCtx, "outbox event pending", "id", event.ID, "type", event.EventType)
			if err := p.outbox.MarkProcessed(txCtx, event.ID); err != nil {
				return err
			}
		}
		return nil
	})
}
