package api

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/model"
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

// GetVaultByAPIKey - Get a single vault for a given API key
func (s Server) GetVaultByAPIKey(c *fiber.Ctx, uniqueId string) error {
	apiKey, ok := c.Locals("api_key").(*model.APIKey)
	if !ok {
		return handler.SendError(c, fiber.StatusUnauthorized, "API key not found in context")
	}

	// First, find the vault by unique ID and user ID (who owns the API key)
	var vault model.Vault
	err := vault.GetByUniqueID(uniqueId, apiKey.UserID)
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

	return c.Status(fiber.StatusOK).JSON(convertToApiVault(&vault))
}
