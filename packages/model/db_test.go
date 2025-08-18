package model

import (
	"log/slog"
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	logger := slog.Default()

	// Test database connection
	err := Open(logger)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Verify connection is working
	err = checkConnection()
	if err != nil {
		t.Fatalf("Database connection check failed: %v", err)
	}
}
