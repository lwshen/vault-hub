package api

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/model"
	"golang.org/x/crypto/pbkdf2"
	"gorm.io/gorm"
)

// GetVaultsForAPIKey lists all vaults accessible to the provided API key.
func GetVaultsForAPIKey(apiKey *model.APIKey) ([]VaultLite, *APIError) {
	if apiKey == nil {
		return nil, newAPIError(http.StatusUnauthorized, "API key not found in context")
	}

	vaults, err := apiKey.GetAccessibleVaults()
	if err != nil {
		return nil, newAPIError(http.StatusInternalServerError, err.Error())
	}

	apiVaults := make([]VaultLite, len(vaults))
	for i := range vaults {
		apiVaults[i] = convertToApiVaultLite(&vaults[i])
	}

	return apiVaults, nil
}

// GetVaultByAPIKeyWithLookup centralizes vault retrieval, authorization, audit
// logging, and optional client-side encryption for API key traffic.
func GetVaultByAPIKeyWithLookup(apiKey *model.APIKey, lookup func(*model.APIKey) (*model.Vault, error), encryptSalt string, enableClientEncryption bool, clientInfo ClientInfo, authorizationHeader string, encryptionFlagProvided bool) (Vault, *APIError) {
	if apiKey == nil {
		return Vault{}, newAPIError(http.StatusUnauthorized, "API key not found in context")
	}

	vault, err := lookup(apiKey)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return Vault{}, newAPIError(http.StatusNotFound, "vault not found")
		}
		return Vault{}, newAPIError(http.StatusInternalServerError, err.Error())
	}

	if !apiKey.HasVaultAccess(vault.ID) {
		return Vault{}, newAPIError(http.StatusForbidden, "API key does not have access to this vault")
	}

	if err := model.LogVaultAction(vault.ID, model.ActionReadVault, apiKey.UserID, model.SourceCLI, &apiKey.ID, clientInfo.IP, clientInfo.UserAgent); err != nil {
		slog.Error("Failed to create audit log for read vault", "error", err, "vaultID", vault.ID)
	}

	if enableClientEncryption {
		plainKey, err := deriveAPIKeyFromAuthorization(authorizationHeader)
		if err != nil {
			return Vault{}, newAPIError(http.StatusBadRequest, err.Error())
		}

		originalValueLen := len(vault.Value)
		encryptedValue, err := encryptForClientWithDerivedKey(vault.Value, plainKey, encryptSalt)
		if err != nil {
			slog.Error("Failed to encrypt vault value for client", "error", err, "vaultID", vault.ID)
			return Vault{}, newAPIError(http.StatusInternalServerError, "failed to encrypt value for client")
		}
		vault.Value = encryptedValue
		slog.Debug("Vault value encrypted for client", "vaultID", vault.ID, "originalLen", originalValueLen, "encryptedLen", len(encryptedValue), "salt", encryptSalt)
	} else if encryptionFlagProvided {
		slog.Debug("Client-side encryption not enabled", "vaultID", vault.ID)
	} else {
		slog.Debug("No client-side encryption preference provided", "vaultID", vault.ID)
	}

	return convertToApiVault(vault), nil
}

func deriveAPIKeyFromAuthorization(header string) (string, error) {
	if header == "" {
		return "", fmt.Errorf("Authorization header required for client encryption")
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("invalid Authorization header format")
	}
	if !strings.HasPrefix(parts[1], "vhub_") {
		return "", fmt.Errorf("API key required for client encryption")
	}

	return parts[1], nil
}

// GetVaultsByAPIKey - Get all vaults for a given API key
func (s Server) GetVaultsByAPIKey(c *fiber.Ctx) error {
	apiKey, ok := c.Locals("api_key").(*model.APIKey)
	if !ok {
		return handler.SendError(c, fiber.StatusUnauthorized, "API key not found in context")
	}

	resp, apiErr := GetVaultsForAPIKey(apiKey)
	if apiErr != nil {
		return handler.SendError(c, apiErr.Status, apiErr.Message)
	}

	return c.Status(http.StatusOK).JSON(resp)
}

// GetVaultByAPIKey - Get a single vault by unique ID for a given API key
func (s Server) GetVaultByAPIKey(c *fiber.Ctx, uniqueId string, params GetVaultByAPIKeyParams) error {
	return s.getVaultByAPIKey(c, uniqueId, params.XEnableClientEncryption, func(apiKey *model.APIKey) (*model.Vault, error) {
		var vault model.Vault
		err := vault.GetByUniqueID(uniqueId, apiKey.UserID)
		return &vault, err
	})
}

// GetVaultByNameAPIKey - Get a single vault by name for a given API key
func (s Server) GetVaultByNameAPIKey(c *fiber.Ctx, name string, params GetVaultByNameAPIKeyParams) error {
	return s.getVaultByAPIKey(c, name, params.XEnableClientEncryption, func(apiKey *model.APIKey) (*model.Vault, error) {
		var vault model.Vault
		err := vault.GetByName(name, apiKey.UserID)
		return &vault, err
	})
}

// getVaultByAPIKey - Common logic for getting a vault via API key
func (s Server) getVaultByAPIKey(c *fiber.Ctx, encryptSalt string, enableClientEncryptionParam *string, vaultGetter func(*model.APIKey) (*model.Vault, error)) error {
	apiKey, ok := c.Locals("api_key").(*model.APIKey)
	if !ok {
		return handler.SendError(c, fiber.StatusUnauthorized, "API key not found in context")
	}

	clientInfo := getClientInfoDetails(c)
	enableClientEncryption := enableClientEncryptionParam != nil && *enableClientEncryptionParam == "true"
	resp, apiErr := GetVaultByAPIKeyWithLookup(apiKey, vaultGetter, encryptSalt, enableClientEncryption, clientInfo, c.Get("Authorization"), enableClientEncryptionParam != nil)
	if apiErr != nil {
		return handler.SendError(c, apiErr.Status, apiErr.Message)
	}

	return c.Status(http.StatusOK).JSON(resp)
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
