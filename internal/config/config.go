package config

import (
	"log/slog"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type DatabaseTypeEnum string

const (
	DatabaseTypeSQLite   DatabaseTypeEnum = "sqlite"
	DatabaseTypeMySQL    DatabaseTypeEnum = "mysql"
	DatabaseTypePostgres DatabaseTypeEnum = "postgres"
)

var (
	AppPort      string
	DatabaseType DatabaseTypeEnum
	DatabaseUrl  string
)

func init() {
	AppPort = getEnv("APP_PORT", "3000")
	DatabaseType = DatabaseTypeEnum(getEnv("DATABASE_TYPE", "sqlite"))
	DatabaseUrl = getEnv("DATABASE_URL", "data.db")

	printConfig()
}

func printConfig() {
	slog.Info("Config", "AppPort", AppPort)
	slog.Info("Config", "DatabaseType", DatabaseType)
	slog.Info("Config", "DatabaseUrl", DatabaseUrl)
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
