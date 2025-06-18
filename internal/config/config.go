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
	OidcEnabled      bool
	OidcClientId     string
	OidcClientSecret string
	OidcIssuer       string
)

func init() {
	AppPort = getEnv("APP_PORT", "3000")
	DatabaseType = DatabaseTypeEnum(getEnv("DATABASE_TYPE", "sqlite"))
	DatabaseUrl = getEnv("DATABASE_URL", "data.db")
	JwtSecret = getEnv("JWT_SECRET", "secret")

	OidcClientId = getEnv("OIDC_CLIENT_ID", "")
	OidcClientSecret = getEnv("OIDC_CLIENT_SECRET", "")
	OidcIssuer = getEnv("OIDC_ISSUER", "")
	OidcEnabled = OidcClientId != "" || OidcClientSecret != "" || OidcIssuer != ""

	printConfig()

	checkConfig()
}

func printConfig() {
	slog.Info("Config", "AppPort", AppPort)
	slog.Info("Config", "DatabaseType", DatabaseType)
	slog.Info("Config", "DatabaseUrl", DatabaseUrl)
	slog.Info("Config", "OidcEnabled", OidcEnabled)
	if OidcEnabled {
		slog.Info("Config", "OidcClientId", OidcClientId)
		slog.Info("Config", "OidcClientSecret", mask(OidcClientSecret))
		slog.Info("Config", "OidcIssuer", OidcIssuer)
	}
}

func checkConfig() {
	hasError := false
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
