package email

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/smtp"
	"strings"

	"github.com/lwshen/vault-hub/internal/config"
)

type Sender interface {
	Send(to string, subject string, htmlBody string) error
}

type SMTPSender struct{}

func NewSMTPSender() *SMTPSender { return &SMTPSender{} }

func (s *SMTPSender) Send(to string, subject string, htmlBody string) error {
	if !config.SmtpEnabled {
		slog.Warn("SMTP is disabled; skipping email send", "to", to, "subject", subject)
		return nil
	}

	host := config.SmtpHost
	port := config.SmtpPort
	addr := net.JoinHostPort(host, port)

	auth := smtp.PlainAuth("", config.SmtpUsername, config.SmtpPassword, host)

	headers := map[string]string{
		"From":         formatFrom(config.SmtpFromName, config.SmtpFromAddress),
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
	}
	var msgBuilder strings.Builder
	for k, v := range headers {
		msgBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msgBuilder.WriteString("\r\n")
	msgBuilder.WriteString(htmlBody)
	msg := []byte(msgBuilder.String())

	// Try STARTTLS when TLS is requested
	if config.SmtpTLS {
		tlsConfig := &tls.Config{
			ServerName: host,
			MinVersion: tls.VersionTLS12,
		}
		// Dial TCP first
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return err
		}
		c, err := smtp.NewClient(conn, host)
		if err != nil {
			return err
		}
		defer func() {
			if err := c.Quit(); err != nil {
				slog.Debug("smtp client quit error", "error", err)
			}
		}()
		if err := c.Auth(auth); err != nil {
			return err
		}
		if err := c.Mail(config.SmtpFromAddress); err != nil {
			return err
		}
		if err := c.Rcpt(to); err != nil {
			return err
		}
		w, err := c.Data()
		if err != nil {
			return err
		}
		if _, err := w.Write(msg); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
		return nil
	}

	// Plain SMTP without TLS (not recommended)
	return smtp.SendMail(addr, auth, config.SmtpFromAddress, []string{to}, msg)
}

func formatFrom(name, address string) string {
	if strings.TrimSpace(name) == "" {
		return address
	}
	return fmt.Sprintf("%s <%s>", name, address)
}
