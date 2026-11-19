package mailer

import (
	"context"

	mailgun "github.com/mailgun/mailgun-go/v4"
)

type MailgunMailer struct {
	mg     *mailgun.MailgunImpl
	sender string
}

func NewMailgunMailer(domain, apiKey, sender string) *MailgunMailer {
	return &MailgunMailer{mg: mailgun.NewMailgun(domain, apiKey), sender: sender}
}

func (m *MailgunMailer) Send(ctx context.Context, message EmailMessage) error {
	msg := mailgun.NewMessage(m.sender, message.Subject, message.TextBody, message.To...)
	if message.HTMLBody != "" {
		msg.SetHTML(message.HTMLBody)
	}
	_, _, err := m.mg.Send(ctx, msg)
	return err
}
