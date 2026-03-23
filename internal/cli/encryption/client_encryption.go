package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

// EncryptForClient encrypts the vault value using a key derived from the API key
// This matches the server-side encryption in packages/api/cli_vault.go
func EncryptForClient(plaintext, apiKey, salt string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// Derive encryption key from API key + salt (vault unique ID or name)
	// This ensures each vault gets a different encryption key even with same API key
	derivedKey := pbkdf2.Key([]byte(apiKey), []byte(salt), 100000, 32, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonceSize := gcm.NonceSize()
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the plaintext
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return encoded, nil
}

// DecryptForClient decrypts the vault value using a key derived from the API key
// This matches the server-side encryption in packages/api/cli_vault.go
func DecryptForClient(encryptedValue, apiKey, salt string) (string, error) {
	if encryptedValue == "" {
		return "", nil
	}

	// Derive encryption key from API key + salt (vault unique ID or name)
	// This ensures each vault gets a different encryption key even with same API key
	derivedKey := pbkdf2.Key([]byte(apiKey), []byte(salt), 100000, 32, sha256.New)

	// Decode base64 encoded ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce from the beginning of the ciphertext
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the ciphertext
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}
