package model

import (
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Initialize database for tests
	logger := slog.Default()
	if err := Open(logger); err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Exit with test result code
	os.Exit(code)
}

func TestGetAPIKeyUsageStats(t *testing.T) {
	// Create a test user with unique email
	password := "test"
	uniqueEmail := "test-api-key-usage-" + time.Now().Format("20060102150405") + "@example.com"
	user := User{
		Email:    uniqueEmail,
		Password: &password,
	}
	if err := DB.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	defer DB.Delete(&user)

	// Create a test vault
	vault := Vault{
		UserID:   user.ID,
		Name:     "Test Vault for API Key Stats",
		UniqueID: "test-vault-stats-" + time.Now().Format("20060102150405"),
		Value:    "test-value",
	}
	if err := DB.Create(&vault).Error; err != nil {
		t.Fatalf("Failed to create test vault: %v", err)
	}
	defer DB.Delete(&vault)

	// Create test API key
	params := CreateAPIKeyParams{
		UserID:   user.ID,
		Name:     "Test API Key for Stats",
		VaultIDs: []uint{vault.ID},
	}
	apiKey, _, err := params.Create()
	if err != nil {
		t.Fatalf("Failed to create test API key: %v", err)
	}
	defer DB.Delete(apiKey)

	// Create some audit logs to simulate usage
	now := time.Now()

	// Create logs at different times
	auditLogs := []AuditLog{
		{
			APIKeyID:  &apiKey.ID,
			VaultID:   &vault.ID,
			Action:    ActionReadVault,
			UserID:    user.ID,
			Source:    SourceCLI,
			IPAddress: "127.0.0.1",
			UserAgent: "test-cli",
		},
		{
			APIKeyID:  &apiKey.ID,
			VaultID:   &vault.ID,
			Action:    ActionReadVault,
			UserID:    user.ID,
			Source:    SourceCLI,
			IPAddress: "127.0.0.1",
			UserAgent: "test-cli",
		},
		{
			APIKeyID:  &apiKey.ID,
			VaultID:   &vault.ID,
			Action:    ActionReadVault,
			UserID:    user.ID,
			Source:    SourceCLI,
			IPAddress: "127.0.0.1",
			UserAgent: "test-cli",
		},
		{
			// Non-vault action (should still be counted in total)
			APIKeyID:  &apiKey.ID,
			Action:    ActionCreateAPIKey,
			UserID:    user.ID,
			Source:    SourceWeb,
			IPAddress: "127.0.0.1",
			UserAgent: "test-browser",
		},
	}

	// Create logs with different timestamps
	for i, log := range auditLogs {
		log.Model.CreatedAt = now.Add(-time.Duration(i) * time.Hour)
		if err := DB.Create(&log).Error; err != nil {
			t.Fatalf("Failed to create audit log %d: %v", i, err)
		}
		defer DB.Delete(&log)
	}

	// Update the API key's LastUsedAt timestamp
	lastUsedAt := now
	apiKey.LastUsedAt = &lastUsedAt
	if err := DB.Save(apiKey).Error; err != nil {
		t.Fatalf("Failed to update API key last used at: %v", err)
	}

	// Test GetAPIKeyUsageStats
	stats, err := GetAPIKeyUsageStats(apiKey.ID)
	if err != nil {
		t.Fatalf("GetAPIKeyUsageStats failed: %v", err)
	}

	// Verify basic statistics
	if stats == nil {
		t.Fatal("Expected stats to be non-nil")
	}

	// Check total requests (should be 4)
	if stats.TotalRequests != 4 {
		t.Errorf("Expected TotalRequests to be 4, got %d", stats.TotalRequests)
	}

	// Check vault access count (should be 3)
	if stats.VaultAccessCount != 3 {
		t.Errorf("Expected VaultAccessCount to be 3, got %d", stats.VaultAccessCount)
	}

	// Check last24Hours (all 4 logs are within 24 hours)
	if stats.Last24Hours != 4 {
		t.Errorf("Expected Last24Hours to be 4, got %d", stats.Last24Hours)
	}

	// Check last7Days (all logs should be counted)
	if stats.Last7Days != 4 {
		t.Errorf("Expected Last7Days to be 4, got %d", stats.Last7Days)
	}

	// Check last30Days (all logs should be counted)
	if stats.Last30Days != 4 {
		t.Errorf("Expected Last30Days to be 4, got %d", stats.Last30Days)
	}

	// Check LastUsedAt
	if stats.LastUsedAt == nil {
		t.Error("Expected LastUsedAt to be non-nil")
	} else if !stats.LastUsedAt.Equal(lastUsedAt) {
		t.Errorf("Expected LastUsedAt to be %v, got %v", lastUsedAt, *stats.LastUsedAt)
	}

	// Check vault breakdown (should have 1 vault with 3 accesses)
	if len(stats.VaultBreakdown) != 1 {
		t.Errorf("Expected VaultBreakdown to have 1 entry, got %d", len(stats.VaultBreakdown))
	} else {
		breakdown := stats.VaultBreakdown[0]
		if breakdown.VaultID != vault.ID {
			t.Errorf("Expected VaultID to be %d, got %d", vault.ID, breakdown.VaultID)
		}
		if breakdown.VaultName != vault.Name {
			t.Errorf("Expected VaultName to be %s, got %s", vault.Name, breakdown.VaultName)
		}
		if breakdown.VaultUniqueID != vault.UniqueID {
			t.Errorf("Expected VaultUniqueID to be %s, got %s", vault.UniqueID, breakdown.VaultUniqueID)
		}
		if breakdown.AccessCount != 3 {
			t.Errorf("Expected AccessCount to be 3, got %d", breakdown.AccessCount)
		}
	}
}

func TestGetAPIKeyUsageStatsNoUsage(t *testing.T) {
	// Create a test user with unique email
	password := "test"
	uniqueEmail := "test-api-key-no-usage-" + time.Now().Format("20060102150405") + "@example.com"
	user := User{
		Email:    uniqueEmail,
		Password: &password,
	}
	if err := DB.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	defer DB.Delete(&user)

	// Create test API key with no usage
	params := CreateAPIKeyParams{
		UserID:   user.ID,
		Name:     "Test API Key with No Usage",
		VaultIDs: []uint{},
	}
	apiKey, _, err := params.Create()
	if err != nil {
		t.Fatalf("Failed to create test API key: %v", err)
	}
	defer DB.Delete(apiKey)

	// Test GetAPIKeyUsageStats for unused key
	stats, err := GetAPIKeyUsageStats(apiKey.ID)
	if err != nil {
		t.Fatalf("GetAPIKeyUsageStats failed: %v", err)
	}

	// Verify all counts are zero
	if stats.TotalRequests != 0 {
		t.Errorf("Expected TotalRequests to be 0, got %d", stats.TotalRequests)
	}
	if stats.Last24Hours != 0 {
		t.Errorf("Expected Last24Hours to be 0, got %d", stats.Last24Hours)
	}
	if stats.Last7Days != 0 {
		t.Errorf("Expected Last7Days to be 0, got %d", stats.Last7Days)
	}
	if stats.Last30Days != 0 {
		t.Errorf("Expected Last30Days to be 0, got %d", stats.Last30Days)
	}
	if stats.VaultAccessCount != 0 {
		t.Errorf("Expected VaultAccessCount to be 0, got %d", stats.VaultAccessCount)
	}
	if stats.LastUsedAt != nil {
		t.Errorf("Expected LastUsedAt to be nil, got %v", *stats.LastUsedAt)
	}
	if len(stats.VaultBreakdown) != 0 {
		t.Errorf("Expected VaultBreakdown to be empty, got %d entries", len(stats.VaultBreakdown))
	}
}

func TestGetAPIKeyUsageStatsNonExistent(t *testing.T) {
	// Test with non-existent API key ID
	_, err := GetAPIKeyUsageStats(999999)
	if err == nil {
		t.Error("Expected error for non-existent API key, got nil")
	}
}
