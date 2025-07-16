package encryption

import (
	"os"
	"testing"

	"github.com/lwshen/vault-hub/internal/config"
)

func TestMain(m *testing.M) {
	// Set up test encryption key
	os.Setenv("ENCRYPTION_KEY", "test-encryption-key-for-testing-purposes")

	// Re-initialize config for tests
	config.EncryptionKey = "test-encryption-key-for-testing-purposes"

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "simple text",
			plaintext: "hello world",
		},
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "special characters",
			plaintext: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
		},
		{
			name:      "unicode characters",
			plaintext: "🔐🌍🚀 Hello, 世界! 🎉",
		},
		{
			name:      "json data",
			plaintext: `{"api_key": "secret-key-123", "database_url": "postgres://user:pass@localhost/db"}`,
		},
		{
			name:      "long text",
			plaintext: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt the plaintext
			ciphertext, err := Encrypt(tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}

			// Verify that ciphertext is different from plaintext (unless empty)
			if tt.plaintext != "" && ciphertext == tt.plaintext {
				t.Error("Ciphertext should be different from plaintext")
			}

			// Decrypt the ciphertext
			decrypted, err := Decrypt(ciphertext)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}

			// Verify that decrypted text matches original plaintext
			if decrypted != tt.plaintext {
				t.Errorf("Decrypt() = %v, want %v", decrypted, tt.plaintext)
			}
		})
	}
}

func TestEncryptDecryptMultipleTimes(t *testing.T) {
	plaintext := "test-data-123"

	// Encrypt the same plaintext multiple times
	var ciphertexts []string
	for i := 0; i < 5; i++ {
		ciphertext, err := Encrypt(plaintext)
		if err != nil {
			t.Fatalf("Encrypt() error = %v", err)
		}
		ciphertexts = append(ciphertexts, ciphertext)
	}

	// Verify that each encryption produces different ciphertext (due to random nonce)
	for i := 0; i < len(ciphertexts); i++ {
		for j := i + 1; j < len(ciphertexts); j++ {
			if ciphertexts[i] == ciphertexts[j] {
				t.Error("Multiple encryptions of the same plaintext should produce different ciphertexts")
			}
		}
	}

	// Verify that all ciphertexts decrypt to the same plaintext
	for i, ciphertext := range ciphertexts {
		decrypted, err := Decrypt(ciphertext)
		if err != nil {
			t.Fatalf("Decrypt() iteration %d error = %v", i, err)
		}
		if decrypted != plaintext {
			t.Errorf("Decrypt() iteration %d = %v, want %v", i, decrypted, plaintext)
		}
	}
}

func TestDecryptInvalidData(t *testing.T) {
	tests := []struct {
		name       string
		ciphertext string
		wantError  bool
	}{
		{
			name:       "invalid base64",
			ciphertext: "invalid-base64-data!@#$",
			wantError:  true,
		},
		{
			name:       "too short data",
			ciphertext: "c2hvcnQ=", // "short" in base64, but too short for valid ciphertext
			wantError:  true,
		},
		{
			name:       "valid base64 but invalid ciphertext",
			ciphertext: "dGhpcyBpcyBub3QgdmFsaWQgY2lwaGVydGV4dA==", // "this is not valid ciphertext" in base64
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decrypt(tt.ciphertext)
			if (err != nil) != tt.wantError {
				t.Errorf("Decrypt() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestEmptyStrings(t *testing.T) {
	// Test encrypting empty string
	ciphertext, err := Encrypt("")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	if ciphertext != "" {
		t.Errorf("Encrypt('') = %v, want empty string", ciphertext)
	}

	// Test decrypting empty string
	plaintext, err := Decrypt("")
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if plaintext != "" {
		t.Errorf("Decrypt('') = %v, want empty string", plaintext)
	}
}
