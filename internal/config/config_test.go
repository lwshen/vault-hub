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
	originalAppPort := os.Getenv("APP_PORT")
	originalDBType := os.Getenv("DATABASE_TYPE")
	originalDBUrl := os.Getenv("DATABASE_URL")

	// Clean up after test
	defer func() {
		if originalAppPort != "" {
			os.Setenv("APP_PORT", originalAppPort)
		} else {
			os.Unsetenv("APP_PORT")
		}
		if originalDBType != "" {
			os.Setenv("DATABASE_TYPE", originalDBType)
		} else {
			os.Unsetenv("DATABASE_TYPE")
		}
		if originalDBUrl != "" {
			os.Setenv("DATABASE_URL", originalDBUrl)
		} else {
			os.Unsetenv("DATABASE_URL")
		}
	}()

	// Test default values
	os.Unsetenv("APP_PORT")
	os.Unsetenv("DATABASE_TYPE")
	os.Unsetenv("DATABASE_URL")

	// The init() function will be called automatically when the package is imported
	// We can verify the default values
	if AppPort != "3000" {
		t.Errorf("AppPort = %v, want %v", AppPort, "3000")
	}
	if DatabaseType != DatabaseTypeSQLite {
		t.Errorf("DatabaseType = %v, want %v", DatabaseType, DatabaseTypeSQLite)
	}
	if DatabaseUrl != "data.db" {
		t.Errorf("DatabaseUrl = %v, want %v", DatabaseUrl, "data.db")
	}

	// Test custom values
	os.Setenv("APP_PORT", "8080")
	os.Setenv("DATABASE_TYPE", "postgres")
	os.Setenv("DATABASE_URL", "postgres://localhost:5432/mydb")

	// Since we can't call init() directly, we'll need to manually set the values
	AppPort = getEnv("APP_PORT", "3000")
	DatabaseType = DatabaseTypeEnum(getEnv("DATABASE_TYPE", "sqlite"))
	DatabaseUrl = getEnv("DATABASE_URL", "data.db")

	if AppPort != "8080" {
		t.Errorf("AppPort = %v, want %v", AppPort, "8080")
	}
	if DatabaseType != DatabaseTypePostgres {
		t.Errorf("DatabaseType = %v, want %v", DatabaseType, DatabaseTypePostgres)
	}
	if DatabaseUrl != "postgres://localhost:5432/mydb" {
		t.Errorf("DatabaseUrl = %v, want %v", DatabaseUrl, "postgres://localhost:5432/mydb")
	}
}
