package worker

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"lsp/internal/platform/database"
	"lsp/internal/platform/events"
	"lsp/internal/platform/mailer"
)

type OutboxProcessor struct {
	tx      *database.TxManager
	outbox  *database.OutboxRepository
	mailer  mailer.Mailer
	bus     events.Bus
	logger  *slog.Logger
	limit   int
	backoff time.Duration
}

func NewOutboxProcessor(tx *database.TxManager, outbox *database.OutboxRepository, mailer mailer.Mailer, bus events.Bus, logger *slog.Logger) *OutboxProcessor {
	return &OutboxProcessor{tx: tx, outbox: outbox, mailer: mailer, bus: bus, logger: logger, limit: 10, backoff: time.Minute}
}

func (p *OutboxProcessor) ProcessOnce(ctx context.Context) error {
	// Kumpulkan event yang sudah diproses di dalam transaksi, lalu publish ke
	// bus setelah commit agar tidak ada notifikasi terkirim saat rollback.
	var processed []database.OutboxEvent
	err := p.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		pending, err := p.outbox.FetchPendingForUpdate(txCtx, p.limit)
		if err != nil {
			return err
		}
		for _, event := range pending {
			_ = p.mailer
			p.logger.InfoContext(txCtx, "outbox event pending", "id", event.ID, "type", event.EventType)
			if err := p.outbox.MarkProcessed(txCtx, event.ID); err != nil {
				return err
			}
			processed = append(processed, event)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, event := range processed {
		p.publish(ctx, event)
	}
	return nil
}

// publish mengirim event ke topic SSE berdasarkan metadata di payload outbox
// (user_id dan/atau role). Bersifat best-effort: kegagalan hanya dicatat log.
func (p *OutboxProcessor) publish(ctx context.Context, event database.OutboxEvent) {
	if p.bus == nil {
		return
	}
	var meta struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
	_ = json.Unmarshal(event.Payload, &meta)

	var topics []string
	if meta.UserID != "" {
		topics = append(topics, events.UserTopic(meta.UserID))
	}
	if meta.Role != "" {
		topics = append(topics, events.RoleTopic(meta.Role))
	}
	for _, topic := range topics {
		if err := p.bus.Publish(ctx, events.Event{Topic: topic, Type: event.EventType, Data: event.Payload}); err != nil {
			p.logger.ErrorContext(ctx, "publish event", "topic", topic, "type", event.EventType, "error", err)
		}
	}
}
