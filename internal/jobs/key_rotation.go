package jobs

import (
	"log/slog"
	"sync"
	"time"

	"github.com/lwshen/vault-hub/internal/encryption"
	"github.com/lwshen/vault-hub/model"
)

// KeyRotationJob handles batch re-encryption of vault values
type KeyRotationJob struct {
	mu              sync.Mutex
	running         bool
	totalVaults     int
	processedVaults int
	failedVaults    int
	lastError       error
}

// NewKeyRotationJob creates a new key rotation job
func NewKeyRotationJob() *KeyRotationJob {
	return &KeyRotationJob{}
}

// Status returns the current status of the key rotation job
func (j *KeyRotationJob) Status() map[string]interface{} {
	j.mu.Lock()
	defer j.mu.Unlock()

	return map[string]interface{}{
		"running":         j.running,
		"totalVaults":     j.totalVaults,
		"processedVaults": j.processedVaults,
		"failedVaults":    j.failedVaults,
		"progress":        j.progress(),
		"lastError":       j.lastError,
	}
}

// progress returns the progress percentage
func (j *KeyRotationJob) progress() float64 {
	if j.totalVaults == 0 {
		return 0
	}
	return float64(j.processedVaults) / float64(j.totalVaults) * 100
}

// Run executes the key rotation job
func (j *KeyRotationJob) Run(newEncryptionKey string) error {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return nil // Already running
	}
	j.running = true
	j.lastError = nil
	j.mu.Unlock()

	slog.Info("Starting key rotation job")

	// Get all vaults
	vaults, err := model.GetAllVaults()
	if err != nil {
		j.mu.Lock()
		j.running = false
		j.lastError = err
		j.mu.Unlock()
		return err
	}

	j.mu.Lock()
	j.totalVaults = len(vaults)
	j.processedVaults = 0
	j.failedVaults = 0
	j.mu.Unlock()

	// Create key manager with new key
	keyManager := encryption.NewKeyManager(newEncryptionKey)
	// Set current key as previous key for decryption during migration
	keyManager.SetPreviousKey(encryption.CurrentKey())

	for _, vault := range vaults {
		if err := j.rotateVaultKey(vault, keyManager); err != nil {
			slog.Error("Failed to rotate key for vault", "vaultID", vault.ID, "error", err)
			j.mu.Lock()
			j.failedVaults++
			j.mu.Unlock()
		}

		j.mu.Lock()
		j.processedVaults++
		j.mu.Unlock()

		// Small delay to avoid overwhelming the database
		time.Sleep(10 * time.Millisecond)
	}

	slog.Info("Key rotation job completed",
		"total", j.totalVaults,
		"processed", j.processedVaults,
		"failed", j.failedVaults)

	j.mu.Lock()
	j.running = false
	j.mu.Unlock()

	return nil
}

// rotateVaultKey rotates the encryption key for a single vault
func (j *KeyRotationJob) rotateVaultKey(vault model.Vault, keyManager *encryption.KeyManager) error {
	// Get the current decrypted value
	decrypted, err := encryption.Decrypt(vault.Value)
	if err != nil {
		return err
	}

	// Re-encrypt with new key
	encrypted, err := keyManager.Encrypt(decrypted)
	if err != nil {
		return err
	}

	// Update the vault
	return model.DB.Model(&vault).Update("value", encrypted).Error
}

// IsRunning returns whether the job is currently running
func (j *KeyRotationJob) IsRunning() bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.running
}
