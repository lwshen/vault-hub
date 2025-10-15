package email

import (
	"log/slog"

	"github.com/lwshen/vault-hub/internal/config"
)

type noopSender struct{}

func (noopSender) Send(string, string, string) error { return nil }

// NewSender returns the configured email sender. It prefers Resend when enabled,
// falling back to SMTP for backward compatibility.
func NewSender() Sender {
	if !config.EmailEnabled {
		return noopSender{}
	}
	if config.ResendEnabled {
		resendSender, err := NewResendSender()
		if err != nil {
			slog.Error("Failed to initialize Resend sender, falling back to SMTP", "error", err)
		} else {
			return resendSender
		}
	}
	return NewSMTPSender()
}
