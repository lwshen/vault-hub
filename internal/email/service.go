package email

import (
	"fmt"
	"log/slog"

	"github.com/lwshen/vault-hub/internal/config"
)

// Service handles email operations with async queue support
type Service struct {
	client    *SMTPClient
	queue     chan *EmailJob
	queueSize int
}

// EmailJob represents an email to be sent
type EmailJob struct {
	To        string
	Subject   string
	HTMLBody  string
	PlainBody string
}

var (
	// GlobalService is the singleton email service instance
	GlobalService *Service
)

// InitService initializes the global email service
func InitService() {
	GlobalService = NewService(100) // Queue size of 100
	GlobalService.Start()
	slog.Info("Email service initialized", "queueSize", GlobalService.queueSize)
}

// NewService creates a new email service with the specified queue size
func NewService(queueSize int) *Service {
	return &Service{
		client:    NewSMTPClient(),
		queue:     make(chan *EmailJob, queueSize),
		queueSize: queueSize,
	}
}

// Start starts the email worker goroutines
func (s *Service) Start() {
	// Start multiple workers for processing emails
	workerCount := 3
	for i := 0; i < workerCount; i++ {
		go s.worker(i)
	}
	slog.Info("Started email workers", "count", workerCount)
}

// worker processes email jobs from the queue
func (s *Service) worker(id int) {
	slog.Debug("Email worker started", "id", id)
	for job := range s.queue {
		if err := s.client.SendEmail(job.To, job.Subject, job.HTMLBody, job.PlainBody); err != nil {
			slog.Error("Failed to send email", "error", err, "to", job.To, "subject", job.Subject, "worker", id)
		} else {
			slog.Debug("Email sent by worker", "worker", id, "to", job.To)
		}
	}
}

// QueueEmail adds an email to the send queue (non-blocking)
func (s *Service) QueueEmail(to, subject, htmlBody, plainBody string) error {
	job := &EmailJob{
		To:        to,
		Subject:   subject,
		HTMLBody:  htmlBody,
		PlainBody: plainBody,
	}

	select {
	case s.queue <- job:
		slog.Debug("Email queued", "to", to, "subject", subject)
		return nil
	default:
		slog.Error("Email queue is full", "to", to, "subject", subject)
		return fmt.Errorf("email queue is full")
	}
}

// SendEmailVerification sends an email verification email
func (s *Service) SendEmailVerification(to, recipientName, verificationURL string) error {
	data := TemplateData{
		AppName:       config.AppName,
		RecipientName: recipientName,
		ActionURL:     verificationURL,
		ExpiryTime:    "24 hours",
	}

	email, err := GenerateEmailVerificationEmail(data)
	if err != nil {
		return fmt.Errorf("failed to generate email: %w", err)
	}

	return s.QueueEmail(to, email.Subject, email.HTMLBody, email.PlainBody)
}

// SendPasswordReset sends a password reset email
func (s *Service) SendPasswordReset(to, recipientName, resetURL string) error {
	data := TemplateData{
		AppName:       config.AppName,
		RecipientName: recipientName,
		ActionURL:     resetURL,
		ExpiryTime:    "1 hour",
	}

	email, err := GeneratePasswordResetEmail(data)
	if err != nil {
		return fmt.Errorf("failed to generate email: %w", err)
	}

	return s.QueueEmail(to, email.Subject, email.HTMLBody, email.PlainBody)
}

// SendMagicLink sends a magic link login email
func (s *Service) SendMagicLink(to, recipientName, magicLinkURL string) error {
	data := TemplateData{
		AppName:       config.AppName,
		RecipientName: recipientName,
		ActionURL:     magicLinkURL,
		ExpiryTime:    "15 minutes",
	}

	email, err := GenerateMagicLinkEmail(data)
	if err != nil {
		return fmt.Errorf("failed to generate email: %w", err)
	}

	return s.QueueEmail(to, email.Subject, email.HTMLBody, email.PlainBody)
}

// Helper functions for convenience

// SendEmailVerificationAsync sends an email verification email using the global service
func SendEmailVerificationAsync(to, recipientName, verificationURL string) error {
	if GlobalService == nil {
		return fmt.Errorf("email service not initialized")
	}
	return GlobalService.SendEmailVerification(to, recipientName, verificationURL)
}

// SendPasswordResetAsync sends a password reset email using the global service
func SendPasswordResetAsync(to, recipientName, resetURL string) error {
	if GlobalService == nil {
		return fmt.Errorf("email service not initialized")
	}
	return GlobalService.SendPasswordReset(to, recipientName, resetURL)
}

// SendMagicLinkAsync sends a magic link email using the global service
func SendMagicLinkAsync(to, recipientName, magicLinkURL string) error {
	if GlobalService == nil {
		return fmt.Errorf("email service not initialized")
	}
	return GlobalService.SendMagicLink(to, recipientName, magicLinkURL)
}
