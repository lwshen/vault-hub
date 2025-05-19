package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		fallback string
		want     string
	}{
		{
			name:     "existing environment variable",
			key:      "TEST_KEY",
			value:    "test_value",
			fallback: "fallback_value",
			want:     "test_value",
		},
		{
			name:     "non-existing environment variable",
			key:      "NON_EXISTING_KEY",
			value:    "",
			fallback: "fallback_value",
			want:     "fallback_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}

			if got := getEnv(tt.key, tt.fallback); got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigValues(t *testing.T) {
	// Save original environment variables
	originalAppPort, appPortSet := os.LookupEnv("APP_PORT")
	originalDBType, dbTypeSet := os.LookupEnv("DATABASE_TYPE")
	originalDBUrl, dbUrlSet := os.LookupEnv("DATABASE_URL")

	// Clean up after test
	defer func() {
		if appPortSet {
			os.Setenv("APP_PORT", originalAppPort)
		} else {
			os.Unsetenv("APP_PORT")
		}
		if dbTypeSet {
			os.Setenv("DATABASE_TYPE", originalDBType)
		} else {
			os.Unsetenv("DATABASE_TYPE")
		}
		if dbUrlSet {
			os.Setenv("DATABASE_URL", originalDBUrl)
		} else {
			os.Unsetenv("DATABASE_URL")
		}

		// Re-populate global config variables to reflect the restored environment.
		// This is important if any test modified these global vars, ensuring they are reset.
		AppPort = getEnv("APP_PORT", "3000")
		DatabaseType = DatabaseTypeEnum(getEnv("DATABASE_TYPE", "sqlite"))
		DatabaseUrl = getEnv("DATABASE_URL", "data.db")
	}()

	t.Run("Defaults", func(t *testing.T) {
		// Unset environment variables to test default fetching logic
		os.Unsetenv("APP_PORT")
		os.Unsetenv("DATABASE_TYPE")
		os.Unsetenv("DATABASE_URL")

		// Test what getEnv would return for defaults.
		// This is because the package's init() has already run with the initial environment.
		// We are testing the logic that init() *would* use if these env vars were unset.
		if port := getEnv("APP_PORT", "3000"); port != "3000" {
			t.Errorf("Default AppPort: getEnv() = %v, want %v", port, "3000")
		}
		if dbType := DatabaseTypeEnum(getEnv("DATABASE_TYPE", "sqlite")); dbType != DatabaseTypeSQLite {
			t.Errorf("Default DatabaseType: getEnv() = %v, want %v", dbType, DatabaseTypeSQLite)
		}
		if dbUrl := getEnv("DATABASE_URL", "data.db"); dbUrl != "data.db" {
			t.Errorf("Default DatabaseUrl: getEnv() = %v, want %v", dbUrl, "data.db")
		}
	})

	t.Run("CustomValues", func(t *testing.T) {
		customAppPort := "8080"
		customDBType := DatabaseTypePostgres
		customDBUrl := "postgres://localhost:5432/mydb"

		os.Setenv("APP_PORT", customAppPort)
		os.Setenv("DATABASE_TYPE", string(customDBType))
		os.Setenv("DATABASE_URL", customDBUrl)

		// Manually update the global config variables using getEnv,
		// simulating init() running with these custom environment variables.
		AppPort = getEnv("APP_PORT", "3000") // Default fallbacks are from original config.go
		DatabaseType = DatabaseTypeEnum(getEnv("DATABASE_TYPE", "sqlite"))
		DatabaseUrl = getEnv("DATABASE_URL", "data.db")

		if AppPort != customAppPort {
			t.Errorf("Custom AppPort = %v, want %v", AppPort, customAppPort)
		}
		if DatabaseType != customDBType {
			t.Errorf("Custom DatabaseType = %v, want %v", DatabaseType, customDBType)
		}
		if DatabaseUrl != customDBUrl {
			t.Errorf("Custom DatabaseUrl = %v, want %v", DatabaseUrl, customDBUrl)
		}
	})
}
