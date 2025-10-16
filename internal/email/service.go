package email

import (
	"fmt"
	"strings"
)

type Service struct {
	sender  Sender
	appName string
	limiter RateLimiter
}

const (
	rateLimitKindSignup        = "signup_confirmation"
	rateLimitKindPasswordReset = "password_reset"
	rateLimitKindMagicLink     = "magic_link"
)

func NewService(sender Sender, appName string) *Service {
	return &Service{sender: sender, appName: appName, limiter: DefaultRateLimiter()}
}

func (s *Service) allow(kind, to string) bool {
	if s.limiter == nil {
		return true
	}
	key := buildRateLimitKey(kind, to)
	return s.limiter.Allow(key)
}

func buildRateLimitKey(kind, to string) string {
	return kind + ":" + strings.ToLower(strings.TrimSpace(to))
}

func (s *Service) SendSignupConfirmation(to, userName string) error {
	if !s.allow(rateLimitKindSignup, to) {
		return ErrRateLimited
	}
	data := TemplateData{
		Subject:  fmt.Sprintf("Welcome to %s", s.appName),
		AppName:  s.appName,
		UserName: userName,
	}
	body, err := renderTemplate("signup_confirmation.html.tmpl", data)
	if err != nil {
		return err
	}
	return s.sender.Send(to, data.Subject, body)
}

func (s *Service) SendPasswordReset(to, userName, actionURL, ttl string) error {
	if !s.allow(rateLimitKindPasswordReset, to) {
		return ErrRateLimited
	}
	data := TemplateData{
		Subject:   fmt.Sprintf("Reset your %s password", s.appName),
		AppName:   s.appName,
		UserName:  userName,
		ActionURL: actionURL,
		TTL:       ttl,
	}
	body, err := renderTemplate("password_reset.html.tmpl", data)
	if err != nil {
		return err
	}
	return s.sender.Send(to, data.Subject, body)
}

func (s *Service) SendMagicLink(to, userName, actionURL, ttl string) error {
	if !s.allow(rateLimitKindMagicLink, to) {
		return ErrRateLimited
	}
	data := TemplateData{
		Subject:   fmt.Sprintf("Your %s magic link", s.appName),
		AppName:   s.appName,
		UserName:  userName,
		ActionURL: actionURL,
		TTL:       ttl,
	}
	body, err := renderTemplate("magic_link_login.html.tmpl", data)
	if err != nil {
		return err
	}
	return s.sender.Send(to, data.Subject, body)
}
