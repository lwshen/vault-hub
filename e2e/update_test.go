package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestUpdate_ByName_WithValue tests updating vault by name with direct value
func TestUpdate_ByName_WithValue(t *testing.T) {
	server := StartTestServer(t)

	result := RunCLI(t,
		"update",
		"--name", server.VaultName,
		"--value", "new-secret-value",
		"--base-url", server.URL,
		"--api-key", server.APIKey,
	)

	result.MustSucceed(t)
	if !result.ContainsStdout(t, "Vault updated successfully") {
		t.Errorf("Expected success message in stdout, got: %s", result.Stdout)
	}
}

// TestUpdate_ByID_WithValue tests updating vault by ID with direct value
func TestUpdate_ByID_WithValue(t *testing.T) {
	server := StartTestServer(t)

	result := RunCLI(t,
		"update",
		"--id", server.VaultID,
		"--value", "new-secret-value",
		"--base-url", server.URL,
		"--api-key", server.APIKey,
	)

	result.MustSucceed(t)
	if !result.ContainsStdout(t, "Vault updated successfully") {
		t.Errorf("Expected success message in stdout, got: %s", result.Stdout)
	}
}

// TestUpdate_ByName_WithValueFile tests updating vault using value from file
func TestUpdate_ByName_WithValueFile(t *testing.T) {
	server := StartTestServer(t)

	// Create temp file with secret
	tempDir := t.TempDir()
	secretFile := filepath.Join(tempDir, "secret.txt")
	if err := os.WriteFile(secretFile, []byte("secret-from-file"), 0600); err != nil {
		t.Fatalf("Failed to create secret file: %v", err)
	}

	result := RunCLI(t,
		"update",
		"--name", server.VaultName,
		"--value-file", secretFile,
		"--base-url", server.URL,
		"--api-key", server.APIKey,
	)

	result.MustSucceed(t)
	if !result.ContainsStdout(t, "Vault updated successfully") {
		t.Errorf("Expected success message in stdout, got: %s", result.Stdout)
	}
}

// TestUpdate_WithClientEncryption tests that client-side encryption works by default
func TestUpdate_WithClientEncryption(t *testing.T) {
	server := StartTestServer(t)

	result := RunCLI(t,
		"update",
		"--name", server.VaultName,
		"--value", "encrypted-secret",
		"--base-url", server.URL,
		"--api-key", server.APIKey,
		"--debug",
	)

	result.MustSucceed(t)
	// Check debug output for encryption
	if !result.ContainsStderr(t, "Encrypting vault value") {
		t.Errorf("Expected encryption debug message in stderr, got: %s", result.Stderr)
	}
}

// TestUpdate_WithoutClientEncryption tests disabling client-side encryption
func TestUpdate_WithoutClientEncryption(t *testing.T) {
	server := StartTestServer(t)

	result := RunCLI(t,
		"update",
		"--name", server.VaultName,
		"--value", "plaintext-secret",
		"--base-url", server.URL,
		"--api-key", server.APIKey,
		"--no-client-encryption",
		"--debug",
	)

	result.MustSucceed(t)
	// Should not encrypt
	if result.ContainsStderr(t, "Encrypting vault value") {
		t.Errorf("Should not encrypt when --no-client-encryption is set")
	}
}

// TestUpdate_JSONOutput tests JSON output format
func TestUpdate_JSONOutput(t *testing.T) {
	server := StartTestServer(t)

	result := RunCLI(t,
		"update",
		"--name", server.VaultName,
		"--value", "json-test-secret",
		"--base-url", server.URL,
		"--api-key", server.APIKey,
		"--output", "json",
	)

	result.MustSucceed(t)
	// Check JSON output contains expected fields
	if !strings.Contains(result.Stdout, "\"uniqueId\"") {
		t.Errorf("Expected JSON with uniqueId field, got: %s", result.Stdout)
	}
	if !strings.Contains(result.Stdout, "\"name\"") {
		t.Errorf("Expected JSON with name field, got: %s", result.Stdout)
	}
}

// TestUpdate_MissingNameAndID tests error when neither name nor ID provided
func TestUpdate_MissingNameAndID(t *testing.T) {
	server := StartTestServer(t)

	result := RunCLI(t,
		"update",
		"--value", "test-value",
		"--base-url", server.URL,
		"--api-key", server.APIKey,
	)

	result.MustFail(t, 1)
	if !result.ContainsStderr(t, "either --name or --id must be provided") {
		t.Errorf("Expected error message about missing name/id, got: %s", result.Stderr)
	}
}

// TestUpdate_BothNameAndID tests error when both name and ID provided
func TestUpdate_BothNameAndID(t *testing.T) {
	server := StartTestServer(t)

	result := RunCLI(t,
		"update",
		"--name", server.VaultName,
		"--id", server.VaultID,
		"--value", "test-value",
		"--base-url", server.URL,
		"--api-key", server.APIKey,
	)

	result.MustFail(t, 1)
	if !result.ContainsStderr(t, "cannot specify both --name and --id") {
		t.Errorf("Expected error message about using both name and id, got: %s", result.Stderr)
	}
}

// TestUpdate_MissingValue tests error when neither value nor value-file provided
func TestUpdate_MissingValue(t *testing.T) {
	server := StartTestServer(t)

	result := RunCLI(t,
		"update",
		"--name", server.VaultName,
		"--base-url", server.URL,
		"--api-key", server.APIKey,
	)

	result.MustFail(t, 1)
	if !result.ContainsStderr(t, "either --value or --value-file must be provided") {
		t.Errorf("Expected error message about missing value, got: %s", result.Stderr)
	}
}

// TestUpdate_EmptyValueFile tests error when value-file is empty
func TestUpdate_EmptyValueFile(t *testing.T) {
	server := StartTestServer(t)

	// Create empty temp file
	tempDir := t.TempDir()
	emptyFile := filepath.Join(tempDir, "empty.txt")
	if err := os.WriteFile(emptyFile, []byte{}, 0600); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	result := RunCLI(t,
		"update",
		"--name", server.VaultName,
		"--value-file", emptyFile,
		"--base-url", server.URL,
		"--api-key", server.APIKey,
	)

	result.MustFail(t, 1)
	if !result.ContainsStderr(t, "Value file is empty") {
		t.Errorf("Expected error message about empty file, got: %s", result.Stderr)
	}
}

// TestUpdate_NonExistentVault tests error when vault doesn't exist
func TestUpdate_NonExistentVault(t *testing.T) {
	server := StartTestServer(t)

	result := RunCLI(t,
		"update",
		"--name", "non-existent-vault-12345",
		"--value", "test-value",
		"--base-url", server.URL,
		"--api-key", server.APIKey,
	)

	result.MustFail(t, 1)
	// Should get 404 or similar error
	if !result.ContainsStderr(t, "not found") && !result.ContainsStderr(t, "404") {
		t.Errorf("Expected error message about vault not found, got: %s", result.Stderr)
	}
}

// TestUpdate_InvalidValueFile tests error when value-file doesn't exist
func TestUpdate_InvalidValueFile(t *testing.T) {
	server := StartTestServer(t)

	result := RunCLI(t,
		"update",
		"--name", server.VaultName,
		"--value-file", "/nonexistent/path/secret.txt",
		"--base-url", server.URL,
		"--api-key", server.APIKey,
	)

	result.MustFail(t, 1)
	if !result.ContainsStderr(t, "Failed to read value file") {
		t.Errorf("Expected error message about file not found, got: %s", result.Stderr)
	}
}

// TestUpdate_SpecialCharacters tests updating with special characters
func TestUpdate_SpecialCharacters(t *testing.T) {
	server := StartTestServer(t)

	specialValue := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
	result := RunCLI(t,
		"update",
		"--name", server.VaultName,
		"--value", specialValue,
		"--base-url", server.URL,
		"--api-key", server.APIKey,
	)

	result.MustSucceed(t)
}

// TestUpdate_MultilineValue tests updating with multiline value
func TestUpdate_MultilineValue(t *testing.T) {
	server := StartTestServer(t)

	multilineValue := "line1\nline2\nline3"
	result := RunCLI(t,
		"update",
		"--name", server.VaultName,
		"--value", multilineValue,
		"--base-url", server.URL,
		"--api-key", server.APIKey,
	)

	result.MustSucceed(t)
}