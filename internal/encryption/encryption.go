package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/lwshen/vault-hub/internal/config"
)

// KeyManager manages encryption keys including rotation
type KeyManager struct {
	currentKey  []byte
	previousKey []byte // Key for decryption during migration
}

// NewKeyManager creates a new key manager
func NewKeyManager(encryptionKey string) *KeyManager {
	return &KeyManager{
		currentKey:  deriveKey(encryptionKey),
		previousKey: nil,
	}
}

// SetPreviousKey sets the previous key for decryption during migration
func (km *KeyManager) SetPreviousKey(encryptionKey string) {
	km.previousKey = deriveKey(encryptionKey)
}

// Encrypt encrypts plaintext using AES-256-GCM with the current key
func (km *KeyManager) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(km.currentKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts using the current key, or falls back to previous key
func (km *KeyManager) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// Try current key first
	plaintext, err := km.decryptWithKey(ciphertext, km.currentKey)
	if err == nil {
		return plaintext, nil
	}

	// Try previous key if available
	if km.previousKey != nil {
		plaintext, err = km.decryptWithKey(ciphertext, km.previousKey)
		if err == nil {
			return plaintext, nil
		}
	}

	return "", fmt.Errorf("failed to decrypt with any key")
}

// decryptWithKey decrypts using the specified key
func (km *KeyManager) decryptWithKey(ciphertext string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// deriveKey derives a 32-byte key from the encryption key using SHA-256
func deriveKey(encryptionKey string) []byte {
	hash := sha256.Sum256([]byte(encryptionKey))
	return hash[:]
}

// Global key manager instance
var km *KeyManager

// Init initializes the global key manager
func Init(encryptionKey string) {
	km = NewKeyManager(encryptionKey)
}

// SetMigrationKey sets the previous key for migration
func SetMigrationKey(encryptionKey string) {
	if km == nil {
		Init(encryptionKey)
		return
	}
	km.SetPreviousKey(encryptionKey)
}

// Encrypt encrypts plaintext using AES-256-GCM (uses global key manager)
func Encrypt(plaintext string) (string, error) {
	if km == nil {
		return encryptLegacy(plaintext)
	}
	return km.Encrypt(plaintext)
}

// Decrypt decrypts base64 encoded ciphertext (uses global key manager)
func Decrypt(ciphertext string) (string, error) {
	if km == nil {
		return decryptLegacy(ciphertext)
	}
	return km.Decrypt(ciphertext)
}

// encryptLegacy uses the original single-key encryption
func encryptLegacy(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	key := deriveKey(config.EncryptionKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptLegacy uses the original single-key decryption
func decryptLegacy(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	key := deriveKey(config.EncryptionKey)

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// CurrentKey returns the current encryption key from config
func CurrentKey() string {
	return config.EncryptionKey
}
