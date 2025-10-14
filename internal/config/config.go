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
	DatabaseType     DatabaseTypeEnum
	DatabaseUrl      string
	JwtSecret        string
	EncryptionKey    string
	OidcEnabled      bool
	OidcClientId     string
	OidcClientSecret string
	OidcIssuer       string
	// SMTP / Email
	SmtpEnabled     bool
	SmtpHost        string
	SmtpPort        string
	SmtpMode        string
	SmtpUsername    string
	SmtpPassword    string
	SmtpFromAddress string
	SmtpFromName    string
	SmtpTLS         bool
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
	SmtpEnabled = getEnv("SMTP_ENABLED", "false") == "true"
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
	if SmtpEnabled {
		if SmtpHost == "" {
			slog.Error("SMTP host is not set (SMTP_HOST)")
			hasError = true
		}
		if SmtpPort == "" {
			slog.Error("SMTP port is not set (SMTP_PORT)")
			hasError = true
		}
		mode := strings.ToLower(SmtpMode)
		if mode != "auto" && mode != "starttls" && mode != "implicit" && mode != "plain" {
			slog.Error("SMTP mode is invalid (SMTP_MODE). Use auto|starttls|implicit|plain", "SmtpMode", SmtpMode)
			hasError = true
		}
		if SmtpFromAddress == "" {
			slog.Error("SMTP from address is not set (SMTP_FROM_ADDRESS)")
			hasError = true
		}
		// For most providers, auth is required
		if SmtpUsername == "" {
			slog.Error("SMTP username is not set (SMTP_USERNAME)")
			hasError = true
		}
		if SmtpPassword == "" {
			slog.Error("SMTP password is not set (SMTP_PASSWORD)")
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
