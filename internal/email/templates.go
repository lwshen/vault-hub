package email

import (
	"bytes"
	"fmt"
	"html/template"
)

// EmailTemplate represents an email with HTML and plain text versions
type EmailTemplate struct {
	Subject   string
	HTMLBody  string
	PlainBody string
}

// TemplateData holds common data for email templates
type TemplateData struct {
	AppName     string
	RecipientName string
	ActionURL   string
	Token       string
	ExpiryTime  string
}

const baseHTML = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f4f4f4;
        }
        .container {
            background-color: #ffffff;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
            padding-bottom: 20px;
            border-bottom: 2px solid #f0f0f0;
        }
        .header h1 {
            margin: 0;
            color: #2563eb;
            font-size: 28px;
        }
        .content {
            margin-bottom: 30px;
        }
        .button {
            display: inline-block;
            padding: 12px 30px;
            background-color: #2563eb;
            color: #ffffff !important;
            text-decoration: none;
            border-radius: 6px;
            font-weight: 600;
            text-align: center;
            margin: 20px 0;
        }
        .button:hover {
            background-color: #1d4ed8;
        }
        .footer {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 2px solid #f0f0f0;
            text-align: center;
            font-size: 12px;
            color: #666;
        }
        .warning {
            background-color: #fef3c7;
            border-left: 4px solid #f59e0b;
            padding: 12px;
            margin: 20px 0;
            border-radius: 4px;
        }
        .code {
            font-family: 'Courier New', monospace;
            background-color: #f3f4f6;
            padding: 2px 6px;
            border-radius: 3px;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.AppName}}</h1>
        </div>
        <div class="content">
            {{.Content}}
        </div>
        <div class="footer">
            <p>This is an automated email from {{.AppName}}. Please do not reply to this email.</p>
            <p>&copy; 2025 {{.AppName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

// GenerateEmailVerificationEmail creates an email verification email
func GenerateEmailVerificationEmail(data TemplateData) (*EmailTemplate, error) {
	subject := fmt.Sprintf("Verify your %s email address", data.AppName)

	htmlContent := fmt.Sprintf(`
		<h2>Welcome to %s, %s!</h2>
		<p>Thank you for signing up. To complete your registration, please verify your email address by clicking the button below:</p>
		<p style="text-align: center;">
			<a href="%s" class="button">Verify Email Address</a>
		</p>
		<p>Or copy and paste this link into your browser:</p>
		<p class="code">%s</p>
		<div class="warning">
			<p><strong>This link will expire in %s.</strong></p>
			<p>If you did not create an account with %s, please ignore this email.</p>
		</div>
	`, data.AppName, data.RecipientName, data.ActionURL, data.ActionURL, data.ExpiryTime, data.AppName)

	plainContent := fmt.Sprintf(`
Welcome to %s, %s!

Thank you for signing up. To complete your registration, please verify your email address by visiting:

%s

This link will expire in %s.

If you did not create an account with %s, please ignore this email.

---
This is an automated email from %s. Please do not reply to this email.
	`, data.AppName, data.RecipientName, data.ActionURL, data.ExpiryTime, data.AppName, data.AppName)

	html, err := renderHTMLTemplate(data.AppName, htmlContent)
	if err != nil {
		return nil, err
	}

	return &EmailTemplate{
		Subject:   subject,
		HTMLBody:  html,
		PlainBody: plainContent,
	}, nil
}

// GeneratePasswordResetEmail creates a password reset email
func GeneratePasswordResetEmail(data TemplateData) (*EmailTemplate, error) {
	subject := fmt.Sprintf("Reset your %s password", data.AppName)

	htmlContent := fmt.Sprintf(`
		<h2>Password Reset Request</h2>
		<p>Hi %s,</p>
		<p>We received a request to reset your password for your %s account. Click the button below to create a new password:</p>
		<p style="text-align: center;">
			<a href="%s" class="button">Reset Password</a>
		</p>
		<p>Or copy and paste this link into your browser:</p>
		<p class="code">%s</p>
		<div class="warning">
			<p><strong>This link will expire in %s.</strong></p>
			<p>If you did not request a password reset, please ignore this email. Your password will remain unchanged.</p>
		</div>
	`, data.RecipientName, data.AppName, data.ActionURL, data.ActionURL, data.ExpiryTime)

	plainContent := fmt.Sprintf(`
Password Reset Request

Hi %s,

We received a request to reset your password for your %s account. Visit the following link to create a new password:

%s

This link will expire in %s.

If you did not request a password reset, please ignore this email. Your password will remain unchanged.

---
This is an automated email from %s. Please do not reply to this email.
	`, data.RecipientName, data.AppName, data.ActionURL, data.ExpiryTime, data.AppName)

	html, err := renderHTMLTemplate(data.AppName, htmlContent)
	if err != nil {
		return nil, err
	}

	return &EmailTemplate{
		Subject:   subject,
		HTMLBody:  html,
		PlainBody: plainContent,
	}, nil
}

// GenerateMagicLinkEmail creates a magic link login email
func GenerateMagicLinkEmail(data TemplateData) (*EmailTemplate, error) {
	subject := fmt.Sprintf("Your %s login link", data.AppName)

	htmlContent := fmt.Sprintf(`
		<h2>Login to Your Account</h2>
		<p>Hi %s,</p>
		<p>Click the button below to securely log in to your %s account:</p>
		<p style="text-align: center;">
			<a href="%s" class="button">Log In to %s</a>
		</p>
		<p>Or copy and paste this link into your browser:</p>
		<p class="code">%s</p>
		<div class="warning">
			<p><strong>This link will expire in %s and can only be used once.</strong></p>
			<p>If you did not request this login link, please ignore this email and ensure your account is secure.</p>
		</div>
	`, data.RecipientName, data.AppName, data.ActionURL, data.AppName, data.ActionURL, data.ExpiryTime)

	plainContent := fmt.Sprintf(`
Login to Your Account

Hi %s,

Click the link below to securely log in to your %s account:

%s

This link will expire in %s and can only be used once.

If you did not request this login link, please ignore this email and ensure your account is secure.

---
This is an automated email from %s. Please do not reply to this email.
	`, data.RecipientName, data.AppName, data.ActionURL, data.ExpiryTime, data.AppName)

	html, err := renderHTMLTemplate(data.AppName, htmlContent)
	if err != nil {
		return nil, err
	}

	return &EmailTemplate{
		Subject:   subject,
		HTMLBody:  html,
		PlainBody: plainContent,
	}, nil
}

// renderHTMLTemplate renders the base HTML template with content
func renderHTMLTemplate(appName, content string) (string, error) {
	tmpl, err := template.New("email").Parse(baseHTML)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		AppName string
		Content template.HTML
	}{
		AppName: appName,
		Content: template.HTML(content),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
