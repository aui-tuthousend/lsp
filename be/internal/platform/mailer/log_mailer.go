package mailer

import (
	"context"
	"log/slog"
)

type LogMailer struct {
	logger *slog.Logger
}

func NewLogMailer(logger *slog.Logger) *LogMailer {
	return &LogMailer{logger: logger}
}

func (m *LogMailer) Send(ctx context.Context, message EmailMessage) error {
	m.logger.InfoContext(ctx, "email logged", "to", message.To, "subject", message.Subject, "text", message.TextBody)
	return nil
}
