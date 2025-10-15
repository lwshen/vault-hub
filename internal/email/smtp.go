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

	msg := buildEmailMessage(to, subject, htmlBody)

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
		return sendImplicitTLS(addr, host, auth, to, msg)

	case "starttls":
		return sendStartTLS(addr, host, port, auth, to, msg)

	case "plain":
		// Plain SMTP (not recommended)
		return smtp.SendMail(addr, auth, config.SmtpFromAddress, []string{to}, msg)

	default:
		return fmt.Errorf("unsupported SMTP mode: %s", mode)
	}
}

func buildEmailMessage(to string, subject string, htmlBody string) []byte {
	headers := map[string]string{
		"From":         sanitizeHeaderValue(formatFrom(config.SmtpFromName, config.SmtpFromAddress)),
		"To":           sanitizeHeaderValue(to),
		"Subject":      sanitizeHeaderValue(subject),
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
	}
	var msgBuilder strings.Builder
	for k, v := range headers {
		msgBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msgBuilder.WriteString("\r\n")
	msgBuilder.WriteString(htmlBody)
	return []byte(msgBuilder.String())
}

func sendImplicitTLS(addr, host string, auth smtp.Auth, to string, msg []byte) error {
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
}

func sanitizeHeaderValue(value string) string {
	return strings.NewReplacer("\r", "", "\n", "").Replace(value)
}

func sendStartTLS(addr, host, port string, auth smtp.Auth, to string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer func() {
		if err := c.Quit(); err != nil {
			slog.Debug("smtp client quit error", "error", err)
		}
	}()
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
	if ok, _ := c.Extension("AUTH"); ok {
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
}

func formatFrom(name, address string) string {
	if strings.TrimSpace(name) == "" {
		return address
	}
	return fmt.Sprintf("%s <%s>", name, address)
}
