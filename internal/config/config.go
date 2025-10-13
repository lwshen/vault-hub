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

var (
	AppPort          string
	AppName          string
	BaseURL          string
	DatabaseType     DatabaseTypeEnum
	DatabaseUrl      string
	JwtSecret        string
	EncryptionKey    string
	OidcEnabled      bool
	OidcClientId     string
	OidcClientSecret string
	OidcIssuer       string
	// SMTP configuration
	SMTPEnabled               bool
	SMTPHost                  string
	SMTPPort                  string
	SMTPUsername              string
	SMTPPassword              string
	SMTPFromEmail             string
	SMTPFromName              string
	EmailVerificationRequired bool
)

func init() {
	AppPort = getEnv("APP_PORT", "3000")
	AppName = getEnv("APP_NAME", "VaultHub")
	BaseURL = getEnv("BASE_URL", "http://localhost:3000")
	JwtSecret = getEnv("JWT_SECRET", "")
	EncryptionKey = getEnv("ENCRYPTION_KEY", "")
	DatabaseType = DatabaseTypeEnum(getEnv("DATABASE_TYPE", "sqlite"))
	DatabaseUrl = getEnv("DATABASE_URL", "data.db")

	OidcClientId = getEnv("OIDC_CLIENT_ID", "")
	OidcClientSecret = getEnv("OIDC_CLIENT_SECRET", "")
	OidcIssuer = getEnv("OIDC_ISSUER", "")
	OidcEnabled = OidcClientId != "" || OidcClientSecret != "" || OidcIssuer != ""

	// SMTP configuration
	SMTPHost = getEnv("SMTP_HOST", "")
	SMTPPort = getEnv("SMTP_PORT", "587")
	SMTPUsername = getEnv("SMTP_USERNAME", "")
	SMTPPassword = getEnv("SMTP_PASSWORD", "")
	SMTPFromEmail = getEnv("SMTP_FROM_EMAIL", "")
	SMTPFromName = getEnv("SMTP_FROM_NAME", AppName)
	SMTPEnabled = SMTPHost != "" && SMTPUsername != "" && SMTPPassword != "" && SMTPFromEmail != ""
	EmailVerificationRequired = getEnv("EMAIL_VERIFICATION_REQUIRED", "false") == "true"

	printConfig()

	checkConfig()
}

func printConfig() {
	slog.Info("Config", "AppPort", AppPort)
	slog.Info("Config", "AppName", AppName)
	slog.Info("Config", "BaseURL", BaseURL)
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
	slog.Info("Config", "SMTPEnabled", SMTPEnabled)
	if SMTPEnabled {
		slog.Info("Config", "SMTPHost", SMTPHost)
		slog.Info("Config", "SMTPPort", SMTPPort)
		slog.Info("Config", "SMTPUsername", SMTPUsername)
		slog.Info("Config", "SMTPPassword", mask(SMTPPassword))
		slog.Info("Config", "SMTPFromEmail", SMTPFromEmail)
		slog.Info("Config", "SMTPFromName", SMTPFromName)
	}
	slog.Info("Config", "EmailVerificationRequired", EmailVerificationRequired)
}

func checkConfig() {
	hasError := false
	if JwtSecret == "" {
		slog.Error("JwtSecret is not set")
		hasError = true
	}
	if EncryptionKey == "" {
		slog.Error("EncryptionKey is not set")
		hasError = true
	}
	if OidcEnabled {
		if OidcClientId == "" {
			slog.Error("OidcClientId is not set")
			hasError = true
		}
		if OidcClientSecret == "" {
			slog.Error("OidcClientSecret is not set")
			hasError = true
		}
		if OidcIssuer == "" {
			slog.Error("OidcIssuer is not set")
			hasError = true
		}
	}
	if hasError {
		slog.Error("Config is invalid, exiting")
		os.Exit(1)
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
