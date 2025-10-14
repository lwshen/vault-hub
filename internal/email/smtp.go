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

	mode := strings.ToLower(config.SmtpMode)
	// Backward compatibility: if SMTP_MODE isn't set (auto) then honor SmtpTLS + port hints
	if mode == "auto" {
		if config.SmtpTLS {
			if port == "465" {
				mode = "implicit"
			} else {
				mode = "starttls"
			}
		} else {
			mode = "plain"
		}
	}

	switch mode {
	case "implicit":
		// SMTPS on 465: connect with TLS immediately
		tlsConfig := &tls.Config{
			ServerName: host,
			MinVersion: tls.VersionTLS12,
		}
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

	case "starttls":
		// SMTP on 587: connect plain, then upgrade via STARTTLS
		c, err := smtp.Dial(addr)
		if err != nil {
			return err
		}
		defer func() {
			if err := c.Quit(); err != nil {
				slog.Debug("smtp client quit error", "error", err)
			}
		}()
		// Attempt STARTTLS if supported
		tlsConfig := &tls.Config{
			ServerName: host,
			MinVersion: tls.VersionTLS12,
		}
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(tlsConfig); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("SMTP server does not support STARTTLS: %s:%s", host, port)
		}
		if ok, _ := c.Extension("AUTH"); ok || true {
			// Many servers require TLS before AUTH PLAIN/LOGIN; now safe to auth
			if err := c.Auth(auth); err != nil {
				return err
			}
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

	case "plain":
		// Plain SMTP (not recommended)
		return smtp.SendMail(addr, auth, config.SmtpFromAddress, []string{to}, msg)

	default:
		return fmt.Errorf("unsupported SMTP mode: %s", mode)
	}
}

func formatFrom(name, address string) string {
	if strings.TrimSpace(name) == "" {
		return address
	}
	return fmt.Sprintf("%s <%s>", name, address)
}
