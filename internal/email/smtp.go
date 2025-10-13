package email

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"

	"github.com/lwshen/vault-hub/internal/config"
)

// SMTPClient handles email sending via SMTP
type SMTPClient struct {
	host     string
	port     string
	username string
	password string
	from     string
	fromName string
}

// NewSMTPClient creates a new SMTP client from configuration
func NewSMTPClient() *SMTPClient {
	return &SMTPClient{
		host:     config.SMTPHost,
		port:     config.SMTPPort,
		username: config.SMTPUsername,
		password: config.SMTPPassword,
		from:     config.SMTPFromEmail,
		fromName: config.SMTPFromName,
	}
}

// SendEmail sends an email with the given parameters
func (s *SMTPClient) SendEmail(to, subject, htmlBody, plainBody string) error {
	if !config.SMTPEnabled {
		slog.Warn("SMTP is not enabled, skipping email send", "to", to, "subject", subject)
		return nil
	}

	// Build the email message
	msg := s.buildMessage(to, subject, htmlBody, plainBody)

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	// Setup authentication
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	// Create TLS config
	tlsConfig := &tls.Config{
		ServerName: s.host,
	}

	// Connect to server with TLS
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		slog.Error("Failed to connect to SMTP server", "error", err, "addr", addr)
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		slog.Error("Failed to create SMTP client", "error", err)
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Authenticate
	if err = client.Auth(auth); err != nil {
		slog.Error("Failed to authenticate with SMTP server", "error", err)
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set sender
	if err = client.Mail(s.from); err != nil {
		slog.Error("Failed to set sender", "error", err)
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipient
	if err = client.Rcpt(to); err != nil {
		slog.Error("Failed to set recipient", "error", err, "to", to)
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send the email body
	w, err := client.Data()
	if err != nil {
		slog.Error("Failed to get data writer", "error", err)
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		slog.Error("Failed to write email body", "error", err)
		return fmt.Errorf("failed to write email body: %w", err)
	}

	err = w.Close()
	if err != nil {
		slog.Error("Failed to close data writer", "error", err)
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	// Send QUIT
	err = client.Quit()
	if err != nil {
		slog.Error("Failed to quit SMTP connection", "error", err)
		return fmt.Errorf("failed to quit: %w", err)
	}

	slog.Info("Email sent successfully", "to", to, "subject", subject)
	return nil
}

// buildMessage constructs the email message with headers and multipart body
func (s *SMTPClient) buildMessage(to, subject, htmlBody, plainBody string) string {
	boundary := "----=_Part_0_1234567890.1234567890"

	from := s.from
	if s.fromName != "" {
		from = fmt.Sprintf("%s <%s>", s.fromName, s.from)
	}

	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"", boundary),
	}

	var parts []string
	parts = append(parts, strings.Join(headers, "\r\n"))
	parts = append(parts, "")

	// Plain text part
	if plainBody != "" {
		parts = append(parts, fmt.Sprintf("--%s", boundary))
		parts = append(parts, "Content-Type: text/plain; charset=UTF-8")
		parts = append(parts, "Content-Transfer-Encoding: 7bit")
		parts = append(parts, "")
		parts = append(parts, plainBody)
		parts = append(parts, "")
	}

	// HTML part
	if htmlBody != "" {
		parts = append(parts, fmt.Sprintf("--%s", boundary))
		parts = append(parts, "Content-Type: text/html; charset=UTF-8")
		parts = append(parts, "Content-Transfer-Encoding: 7bit")
		parts = append(parts, "")
		parts = append(parts, htmlBody)
		parts = append(parts, "")
	}

	parts = append(parts, fmt.Sprintf("--%s--", boundary))

	return strings.Join(parts, "\r\n")
}
