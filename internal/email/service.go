package email

import (
	"fmt"
)

type Service struct {
	sender       Sender
	appName      string
	supportEmail string
}

func NewService(sender Sender, appName, supportEmail string) *Service {
	return &Service{sender: sender, appName: appName, supportEmail: supportEmail}
}

func (s *Service) SendSignupConfirmation(to, userName string) error {
	data := TemplateData{
		Subject:      fmt.Sprintf("Welcome to %s", s.appName),
		AppName:      s.appName,
		UserName:     userName,
		SupportEmail: s.supportEmail,
	}
	body, err := renderTemplate("signup_confirmation.html.tmpl", data)
	if err != nil {
		return err
	}
	return s.sender.Send(to, data.Subject, body)
}

func (s *Service) SendPasswordReset(to, userName, actionURL, ttl string) error {
	data := TemplateData{
		Subject:      fmt.Sprintf("Reset your %s password", s.appName),
		AppName:      s.appName,
		UserName:     userName,
		ActionURL:    actionURL,
		TTL:          ttl,
		SupportEmail: s.supportEmail,
	}
	body, err := renderTemplate("password_reset.html.tmpl", data)
	if err != nil {
		return err
	}
	return s.sender.Send(to, data.Subject, body)
}

func (s *Service) SendMagicLink(to, userName, actionURL, ttl string) error {
	data := TemplateData{
		Subject:      fmt.Sprintf("Your %s magic link", s.appName),
		AppName:      s.appName,
		UserName:     userName,
		ActionURL:    actionURL,
		TTL:          ttl,
		SupportEmail: s.supportEmail,
	}
	body, err := renderTemplate("magic_link_login.html.tmpl", data)
	if err != nil {
		return err
	}
	return s.sender.Send(to, data.Subject, body)
}
