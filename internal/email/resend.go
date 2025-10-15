package email

import (
	"fmt"
	"log/slog"

	"github.com/lwshen/vault-hub/internal/config"
	resend "github.com/resend/resend-go/v2"
)

type ResendSender struct {
	client *resend.Client
}

func NewResendSender() (*ResendSender, error) {
	if config.ResendAPIKey == "" {
		return nil, fmt.Errorf("resend api key is not configured")
	}
	client := resend.NewClient(config.ResendAPIKey)
	return &ResendSender{client: client}, nil
}

func (s *ResendSender) Send(to string, subject string, htmlBody string) error {
	if !config.ResendEnabled {
		slog.Warn("Resend is disabled; skipping email send", "to", to, "subject", subject)
		return nil
	}
	params := &resend.SendEmailRequest{
		From:    formatFrom(config.ResendFromName, config.ResendFromAddress),
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
	}
	if _, err := s.client.Emails.Send(params); err != nil {
		return fmt.Errorf("failed to send email via Resend: %w", err)
	}
	return nil
}
