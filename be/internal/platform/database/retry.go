package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

func RunWithRetry(ctx context.Context, maxAttempts int, fn func(context.Context) error) error {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := fn(ctx); err != nil {
			lastErr = err
			if !isRetryablePgError(err) || attempt == maxAttempts {
				return err
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(attempt) * 50 * time.Millisecond):
			}
			continue
		}
		return nil
	}
	return lastErr
}

func isRetryablePgError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	return pgErr.Code == "40P01" || pgErr.Code == "40001"
}
