package database

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"lsp/internal/platform/id"
)

type OutboxEvent struct {
	ID            string
	AggregateType string
	AggregateID   string
	EventType     string
	Payload       json.RawMessage
	Status        string
	Attempts      int
	AvailableAt   time.Time
	ProcessedAt   *time.Time
	Error         *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type OutboxRepository struct {
	pool *pgxpool.Pool
}

func NewOutboxRepository(pool *pgxpool.Pool) *OutboxRepository {
	return &OutboxRepository{pool: pool}
}

func (r *OutboxRepository) Insert(ctx context.Context, event OutboxEvent) (string, error) {
	db := GetExecutor(ctx, r.pool)
	if event.ID == "" {
		event.ID = id.New()
	}
	if event.AvailableAt.IsZero() {
		event.AvailableAt = time.Now().UTC()
	}
	_, err := db.Exec(ctx, `
		INSERT INTO outbox_events (id, aggregate_type, aggregate_id, event_type, payload, available_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, event.ID, event.AggregateType, event.AggregateID, event.EventType, event.Payload, event.AvailableAt)
	return event.ID, err
}

func (r *OutboxRepository) FetchPendingForUpdate(ctx context.Context, limit int) ([]OutboxEvent, error) {
	db := GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT id, aggregate_type, aggregate_id, event_type, payload, status, attempts, available_at, processed_at, error, created_at, updated_at
		FROM outbox_events
		WHERE status = 'pending' AND available_at <= now()
		ORDER BY available_at ASC, created_at ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]OutboxEvent, 0)
	for rows.Next() {
		var event OutboxEvent
		if err := rows.Scan(&event.ID, &event.AggregateType, &event.AggregateID, &event.EventType, &event.Payload, &event.Status, &event.Attempts, &event.AvailableAt, &event.ProcessedAt, &event.Error, &event.CreatedAt, &event.UpdatedAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

func (r *OutboxRepository) MarkProcessed(ctx context.Context, eventID string) error {
	db := GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		UPDATE outbox_events
		SET status = 'processed', processed_at = now(), error = NULL
		WHERE id = $1
	`, eventID)
	return err
}

func (r *OutboxRepository) MarkFailed(ctx context.Context, eventID string, reason string, retryAt time.Time) error {
	db := GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		UPDATE outbox_events
		SET status = CASE WHEN attempts + 1 >= 5 THEN 'failed' ELSE 'pending' END,
		    attempts = attempts + 1,
		    available_at = $2,
		    error = $3
		WHERE id = $1
	`, eventID, retryAt, reason)
	return err
}
