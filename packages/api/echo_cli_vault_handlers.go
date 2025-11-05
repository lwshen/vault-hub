package api

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/model"
	"github.com/lwshen/vault-hub/packages/api/generated/models"
	"golang.org/x/crypto/pbkdf2"
	"gorm.io/gorm"
)

// GetVaultsByAPIKey - Get all vaults for a given API key
func (c *Container) GetVaultsByAPIKey(ctx echo.Context) error {
	apiKey, err := getAPIKeyFromEchoContext(ctx)
	if err != nil {
		return err
	}

	// Get all accessible vaults for this API key (encrypted)
	vaults, err := apiKey.GetAccessibleVaults()
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, err.Error())
	}

	// Convert to API VaultLite format (no decryption needed)
	apiVaults := make([]models.VaultLite, 0, len(vaults))
	for i := range vaults {
		apiVaults = append(apiVaults, convertToGeneratedVaultLite(&vaults[i]))
	}

	return ctx.JSON(http.StatusOK, apiVaults)
}

// GetVaultByAPIKey - Get a single vault by unique ID for a given API key
func (c *Container) GetVaultByAPIKey(ctx echo.Context) error {
	uniqueID := ctx.Param("uniqueId")
	enableClientEncryptionHeader := ctx.Request().Header.Get("X-Enable-Client-Encryption")

	return c.getVaultByAPIKey(ctx, uniqueID, enableClientEncryptionHeader, func(apiKey *model.APIKey) (*model.Vault, error) {
		var vault model.Vault
		err := vault.GetByUniqueID(uniqueID, apiKey.UserID)
		return &vault, err
	})
}

// GetVaultByNameAPIKey - Get a single vault by name for a given API key
func (c *Container) GetVaultByNameAPIKey(ctx echo.Context) error {
	name := ctx.Param("name")
	enableClientEncryptionHeader := ctx.Request().Header.Get("X-Enable-Client-Encryption")

	return c.getVaultByAPIKey(ctx, name, enableClientEncryptionHeader, func(apiKey *model.APIKey) (*model.Vault, error) {
		var vault model.Vault
		err := vault.GetByName(name, apiKey.UserID)
		return &vault, err
	})
}

// getVaultByAPIKey - Common logic for getting a vault via API key
func (c *Container) getVaultByAPIKey(ctx echo.Context, encryptSalt string, enableClientEncryptionHeader string, vaultGetter func(*model.APIKey) (*model.Vault, error)) error {
	apiKey, err := getAPIKeyFromEchoContext(ctx)
	if err != nil {
		return err
	}

	// Get the vault using the provided getter function
	vault, err := vaultGetter(apiKey)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return SendError(ctx, http.StatusNotFound, "vault not found")
		}
		return SendError(ctx, http.StatusInternalServerError, err.Error())
	}

	// Check if the API key has access to this specific vault
	if !apiKey.HasVaultAccess(vault.ID) {
		return SendError(ctx, http.StatusForbidden, "API key does not have access to this vault")
	}

	// Log read action (using the API key user ID)
	ip, userAgent := getClientInfoEcho(ctx)
	if err := model.LogVaultAction(vault.ID, model.ActionReadVault, apiKey.UserID, model.SourceCLI, &apiKey.ID, ip, userAgent); err != nil {
		slog.Error("Failed to create audit log for read vault", "error", err, "vaultID", vault.ID)
	}

	// Enhanced security: Apply additional client-side encryption if requested
	enableClientEncryption := enableClientEncryptionHeader == "true"
	if enableClientEncryption {
		slog.Debug("Client-side encryption requested", "header", enableClientEncryptionHeader, "vaultID", vault.ID)

		// Get the original API key from the Authorization header to use for key derivation
		authHeader := ctx.Request().Header.Get("Authorization")
		originalAPIKey := strings.TrimPrefix(authHeader, "Bearer ")

		originalValueLen := len(vault.Value)
		encryptedValue, err := encryptForClientWithDerivedKey(vault.Value, originalAPIKey, encryptSalt)
		if err != nil {
			slog.Error("Failed to encrypt vault value for client", "error", err, "vaultID", vault.ID)
			return SendError(ctx, http.StatusInternalServerError, "failed to encrypt value for client")
		}
		vault.Value = encryptedValue
		slog.Debug("Vault value encrypted for client",
			"vaultID", vault.ID,
			"originalLen", originalValueLen,
			"encryptedLen", len(encryptedValue),
			"salt", encryptSalt)
	} else {
		if enableClientEncryptionHeader != "" {
			slog.Debug("Client-side encryption not enabled", "headerValue", enableClientEncryptionHeader, "vaultID", vault.ID)
		} else {
			slog.Debug("No client-side encryption header received", "vaultID", vault.ID)
		}
	}

	return ctx.JSON(http.StatusOK, convertToGeneratedVault(vault))
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
