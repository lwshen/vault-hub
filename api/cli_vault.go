package api

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/model"
	"golang.org/x/crypto/pbkdf2"
	"gorm.io/gorm"
)

// GetVaultsByAPIKey - Get all vaults for a given API key
func (s Server) GetVaultsByAPIKey(c *fiber.Ctx) error {
	apiKey, ok := c.Locals("api_key").(*model.APIKey)
	if !ok {
		return handler.SendError(c, fiber.StatusUnauthorized, "API key not found in context")
	}

	// Get all accessible vaults for this API key (encrypted)
	vaults, err := apiKey.GetAccessibleVaults()
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Convert to API VaultLite format (no decryption needed)
	apiVaults := make([]VaultLite, len(vaults))
	for i, vault := range vaults {
		apiVaults[i] = convertToApiVaultLite(&vault)
	}

	return c.Status(fiber.StatusOK).JSON(apiVaults)
}

// GetVaultByAPIKey - Get a single vault by unique ID for a given API key
func (s Server) GetVaultByAPIKey(c *fiber.Ctx, uniqueId string) error {
	return s.getVaultByAPIKey(c, uniqueId, func(apiKey *model.APIKey) (*model.Vault, error) {
		var vault model.Vault
		err := vault.GetByUniqueID(uniqueId, apiKey.UserID)
		return &vault, err
	})
}

// GetVaultByNameAPIKey - Get a single vault by name for a given API key
func (s Server) GetVaultByNameAPIKey(c *fiber.Ctx, name string) error {
	return s.getVaultByAPIKey(c, name, func(apiKey *model.APIKey) (*model.Vault, error) {
		var vault model.Vault
		err := vault.GetByName(name, apiKey.UserID)
		return &vault, err
	})
}

// getVaultByAPIKey - Common logic for getting a vault via API key
func (s Server) getVaultByAPIKey(c *fiber.Ctx, encryptSalt string, vaultGetter func(*model.APIKey) (*model.Vault, error)) error {
	apiKey, ok := c.Locals("api_key").(*model.APIKey)
	if !ok {
		return handler.SendError(c, fiber.StatusUnauthorized, "API key not found in context")
	}

	// Get the vault using the provided getter function
	vault, err := vaultGetter(apiKey)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return handler.SendError(c, fiber.StatusNotFound, "vault not found")
		}
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Check if the API key has access to this specific vault
	if !apiKey.HasVaultAccess(vault.ID) {
		return handler.SendError(c, fiber.StatusForbidden, "API key does not have access to this vault")
	}

	// Log read action (using the API key user ID)
	ip, userAgent := getClientInfo(c)
	if err := model.LogVaultAction(vault.ID, model.ActionReadVault, apiKey.UserID, ip, userAgent); err != nil {
		slog.Error("Failed to create audit log for read vault", "error", err, "vaultID", vault.ID)
	}

	// Enhanced security: Apply additional client-side encryption if requested
	enableClientEncryption := c.Get("X-Enable-Client-Encryption")
	if enableClientEncryption == "true" {
		// Get the original API key from the Authorization header to use for key derivation
		authHeader := c.Get("Authorization")
		originalAPIKey := authHeader[7:] // Remove "Bearer " prefix

		encryptedValue, err := encryptForClientWithDerivedKey(vault.Value, originalAPIKey, encryptSalt)
		if err != nil {
			return handler.SendError(c, fiber.StatusInternalServerError, "failed to encrypt value for client")
		}
		vault.Value = encryptedValue
	}

	return c.Status(fiber.StatusOK).JSON(convertToApiVault(vault))
}

// encryptForClientWithDerivedKey encrypts the vault value using a key derived from the API key
// This provides additional security without requiring key exchange
func encryptForClientWithDerivedKey(plaintext, apiKey, salt string) (string, error) {
	// Derive encryption key from API key + vault unique ID as salt
	// This ensures each vault gets a different encryption key even with same API key
	derivedKey := pbkdf2.Key([]byte(apiKey), []byte(salt), 100000, 32, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the plaintext
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Return base64 encoded result
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
