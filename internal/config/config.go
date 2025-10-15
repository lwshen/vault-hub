package config

import (
	"log/slog"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

type DatabaseTypeEnum string

const (
	DatabaseTypeSQLite   DatabaseTypeEnum = "sqlite"
	DatabaseTypeMySQL    DatabaseTypeEnum = "mysql"
	DatabaseTypePostgres DatabaseTypeEnum = "postgres"
)

const (
	EmailTypeSMTP   = "SMTP"
	EmailTypeResend = "RESEND"
)

var (
	AppPort           string
	DatabaseType      DatabaseTypeEnum
	DatabaseUrl       string
	JwtSecret         string
	EncryptionKey     string
	OidcEnabled       bool
	OidcClientId      string
	OidcClientSecret  string
	OidcIssuer        string
	EmailEnabled      bool
	EmailType         string
	SmtpEnabled       bool
	SmtpHost          string
	SmtpPort          string
	SmtpMode          string
	SmtpUsername      string
	SmtpPassword      string
	SmtpFromAddress   string
	SmtpFromName      string
	SmtpTLS           bool
	ResendEnabled     bool
	ResendAPIKey      string
	ResendFromAddress string
	ResendFromName    string
)

func init() {
	AppPort = getEnv("APP_PORT", "3000")
	JwtSecret = getEnv("JWT_SECRET", "")
	EncryptionKey = getEnv("ENCRYPTION_KEY", "")
	DatabaseType = DatabaseTypeEnum(getEnv("DATABASE_TYPE", "sqlite"))
	DatabaseUrl = getEnv("DATABASE_URL", "data.db")

	OidcClientId = getEnv("OIDC_CLIENT_ID", "")
	OidcClientSecret = getEnv("OIDC_CLIENT_SECRET", "")
	OidcIssuer = getEnv("OIDC_ISSUER", "")
	OidcEnabled = OidcClientId != "" || OidcClientSecret != "" || OidcIssuer != ""

	// SMTP
	rawEmailType := strings.ToUpper(strings.TrimSpace(getEnv("EMAIL_TYPE", "")))
	switch rawEmailType {
	case EmailTypeSMTP, EmailTypeResend:
		EmailType = rawEmailType
	default:
		if getEnv("RESEND_ENABLED", "false") == "true" {
			EmailType = EmailTypeResend
		} else {
			EmailType = EmailTypeSMTP
		}
	}

	rawEmailEnabled := strings.TrimSpace(getEnv("EMAIL_ENABLED", ""))
	if rawEmailEnabled != "" {
		EmailEnabled = strings.EqualFold(rawEmailEnabled, "true")
	} else {
		switch EmailType {
		case EmailTypeResend:
			EmailEnabled = getEnv("RESEND_ENABLED", "false") == "true"
		default:
			EmailEnabled = getEnv("SMTP_ENABLED", "false") == "true"
		}
	}

	SmtpHost = getEnv("SMTP_HOST", "")
	SmtpPort = getEnv("SMTP_PORT", "587")
	// SMTP_MODE controls TLS behavior: auto|starttls|implicit|plain
	// auto: choose by port (465=implicit, 587=starttls, else try STARTTLS then plain)
	SmtpMode = getEnv("SMTP_MODE", "auto")
	SmtpUsername = getEnv("SMTP_USERNAME", "")
	SmtpPassword = getEnv("SMTP_PASSWORD", "")
	SmtpFromAddress = getEnv("SMTP_FROM_ADDRESS", "")
	SmtpFromName = getEnv("SMTP_FROM_NAME", "Vault Hub")
	SmtpTLS = getEnv("SMTP_TLS", "true") == "true"
	ResendEnabled = EmailEnabled && EmailType == EmailTypeResend
	SmtpEnabled = EmailEnabled && EmailType == EmailTypeSMTP
	ResendAPIKey = getEnv("RESEND_API_KEY", "")
	ResendFromAddress = getEnv("RESEND_FROM_ADDRESS", SmtpFromAddress)
	ResendFromName = getEnv("RESEND_FROM_NAME", SmtpFromName)

	printConfig()

	checkConfig()
}

func printConfig() {
	slog.Info("Config", "AppPort", AppPort)
	slog.Info("Config", "JwtSecret", mask(JwtSecret))
	slog.Info("Config", "EncryptionKey", mask(EncryptionKey))
	slog.Info("Config", "DatabaseType", DatabaseType)
	slog.Info("Config", "DatabaseUrl", DatabaseUrl)
	slog.Info("Config", "OidcEnabled", OidcEnabled)
	if OidcEnabled {
		slog.Info("Config", "OidcClientId", OidcClientId)
		slog.Info("Config", "OidcClientSecret", mask(OidcClientSecret))
		slog.Info("Config", "OidcIssuer", OidcIssuer)
	}
	slog.Info("Config", "EmailEnabled", EmailEnabled)
	slog.Info("Config", "EmailType", EmailType)
	slog.Info("Config", "SmtpEnabled", SmtpEnabled)
	if SmtpEnabled {
		slog.Info("Config", "SmtpHost", SmtpHost)
		slog.Info("Config", "SmtpPort", SmtpPort)
		slog.Info("Config", "SmtpMode", SmtpMode)
		slog.Info("Config", "SmtpUsername", SmtpUsername)
		slog.Info("Config", "SmtpPassword", mask(SmtpPassword))
		slog.Info("Config", "SmtpFromAddress", SmtpFromAddress)
		slog.Info("Config", "SmtpFromName", SmtpFromName)
		slog.Info("Config", "SmtpTLS", SmtpTLS)
	}
	slog.Info("Config", "ResendEnabled", ResendEnabled)
	if ResendEnabled {
		slog.Info("Config", "ResendFromAddress", ResendFromAddress)
		slog.Info("Config", "ResendFromName", ResendFromName)
	}
}

func checkConfig() {
	type validation struct {
		ok  bool
		msg string
	}

	lowerMode := strings.ToLower(SmtpMode)

	validations := []validation{
		{JwtSecret != "", "JwtSecret is not set"},
		{EncryptionKey != "", "EncryptionKey is not set"},

		{!EmailEnabled || isValidEmailType(EmailType), "Email type is invalid (EMAIL_TYPE). Use SMTP|RESEND"},

		// OIDC checks are required only when OIDC is enabled
		{!OidcEnabled || OidcClientId != "", "OidcClientId is not set"},
		{!OidcEnabled || OidcClientSecret != "", "OidcClientSecret is not set"},
		{!OidcEnabled || OidcIssuer != "", "OidcIssuer is not set"},

		// SMTP checks are required only when SMTP is enabled
		{!SmtpEnabled || SmtpHost != "", "SMTP host is not set (SMTP_HOST)"},
		{!SmtpEnabled || SmtpPort != "", "SMTP port is not set (SMTP_PORT)"},
		{!SmtpEnabled || isValidSmtpMode(lowerMode), "SMTP mode is invalid (SMTP_MODE). Use auto|starttls|implicit|plain"},
		{!SmtpEnabled || SmtpFromAddress != "", "SMTP from address is not set (SMTP_FROM_ADDRESS)"},
		{!SmtpEnabled || SmtpUsername != "", "SMTP username is not set (SMTP_USERNAME)"},
		{!SmtpEnabled || SmtpPassword != "", "SMTP password is not set (SMTP_PASSWORD)"},
		{!ResendEnabled || ResendAPIKey != "", "Resend API key is not set (RESEND_API_KEY)"},
		{!ResendEnabled || ResendFromAddress != "", "Resend from address is not set (RESEND_FROM_ADDRESS)"},
	}

	hasError := false
	for _, v := range validations {
		if !v.ok {
			slog.Error(v.msg)
			hasError = true
		}
	}

	if hasError {
		slog.Error("Config is invalid, exiting")
		os.Exit(1)
	}
}

func isValidEmailType(emailType string) bool {
	switch strings.ToUpper(emailType) {
	case EmailTypeSMTP, EmailTypeResend:
		return true
	default:
		return false
	}
}

func isValidSmtpMode(mode string) bool {
	switch mode {
	case "auto", "starttls", "implicit", "plain":
		return true
	default:
		return false
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func mask(value string) string {
	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}
	return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}
