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
	DemoEnabled       bool
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

type validation struct {
	ok  bool
	msg string
}

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

	DemoEnabled = getEnv("DEMO_ENABLED", "false") == "true"

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
	slog.Info("Config", "DemoEnabled", DemoEnabled)
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
	validations := make([]validation, 0, 16)
	validations = append(validations, baseValidations()...)
	validations = append(validations, oidcValidations()...)
	validations = append(validations, emailValidations()...)
	validations = append(validations, smtpValidations()...)
	validations = append(validations, resendValidations()...)

	if logValidationErrors(validations) {
		slog.Error("Config is invalid, exiting")
		os.Exit(1)
	}
}

func baseValidations() []validation {
	return []validation{
		{ok: JwtSecret != "", msg: "JwtSecret is not set"},
		{ok: EncryptionKey != "", msg: "EncryptionKey is not set"},
	}
}

func oidcValidations() []validation {
	if !OidcEnabled {
		return nil
	}
	return []validation{
		{ok: OidcClientId != "", msg: "OidcClientId is not set"},
		{ok: OidcClientSecret != "", msg: "OidcClientSecret is not set"},
		{ok: OidcIssuer != "", msg: "OidcIssuer is not set"},
	}
}

func emailValidations() []validation {
	if !EmailEnabled {
		return nil
	}
	return []validation{
		{ok: isValidEmailType(EmailType), msg: "Email type is invalid (EMAIL_TYPE). Use SMTP|RESEND"},
	}
}

func smtpValidations() []validation {
	if !SmtpEnabled {
		return nil
	}
	lowerMode := strings.ToLower(SmtpMode)
	return []validation{
		{ok: SmtpHost != "", msg: "SMTP host is not set (SMTP_HOST)"},
		{ok: SmtpPort != "", msg: "SMTP port is not set (SMTP_PORT)"},
		{ok: isValidSmtpMode(lowerMode), msg: "SMTP mode is invalid (SMTP_MODE). Use auto|starttls|implicit|plain"},
		{ok: SmtpFromAddress != "", msg: "SMTP from address is not set (SMTP_FROM_ADDRESS)"},
		{ok: SmtpUsername != "", msg: "SMTP username is not set (SMTP_USERNAME)"},
		{ok: SmtpPassword != "", msg: "SMTP password is not set (SMTP_PASSWORD)"},
	}
}

func resendValidations() []validation {
	if !ResendEnabled {
		return nil
	}
	return []validation{
		{ok: ResendAPIKey != "", msg: "Resend API key is not set (RESEND_API_KEY)"},
		{ok: ResendFromAddress != "", msg: "Resend from address is not set (RESEND_FROM_ADDRESS)"},
	}
}

func logValidationErrors(validations []validation) bool {
	hasError := false
	for _, v := range validations {
		if !v.ok {
			slog.Error(v.msg)
			hasError = true
		}
	}
	return hasError
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
