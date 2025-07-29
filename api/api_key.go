package api

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/model"
	"gorm.io/gorm"
)

// auditAPIKeyOperation creates an audit log entry for API key operations
func auditAPIKeyOperation(c *fiber.Ctx, action model.ActionType, userID uint, apiKeyID uint, apiKeyName string) {
	ip, userAgent := getClientInfo(c)

	err := model.LogAPIKeyAction(apiKeyID, action, userID, ip, userAgent)

	if err != nil {
		// Log the audit error but don't fail the main operation
		slog.Error("Failed to create audit log for API key operation",
			"action", action,
			"user_id", userID,
			"api_key_id", apiKeyID,
			"api_key_name", apiKeyName,
			"error", err)
	}
}

// convertToApiAPIKey converts a model.APIKey to an api.APIKey
func convertToApiAPIKey(apiKey *model.APIKey) (*APIKey, error) {
	// Get accessible vaults for this API key
	vaults, err := apiKey.GetAccessibleVaults()
	if err != nil {
		return nil, err
	}

	// Convert vaults to VaultLite
	var apiVaults []VaultLite
	for _, vault := range vaults {
		apiVaults = append(apiVaults, convertToApiVaultLite(&vault))
	}

	// Convert timestamps
	var expiresAt, lastUsedAt *time.Time
	if apiKey.ExpiresAt != nil {
		expiresAt = apiKey.ExpiresAt
	}
	if apiKey.LastUsedAt != nil {
		lastUsedAt = apiKey.LastUsedAt
	}

	// #nosec G115
	id := int64(apiKey.ID)
	return &APIKey{
		Id:         id,
		Name:       apiKey.Name,
		Vaults:     &apiVaults,
		ExpiresAt:  expiresAt,
		LastUsedAt: lastUsedAt,
		IsActive:   apiKey.IsActive,
		CreatedAt:  apiKey.CreatedAt,
		UpdatedAt:  &apiKey.UpdatedAt,
	}, nil
}

// GetAPIKeys - Get API keys for the current user with pagination
func (s Server) GetAPIKeys(c *fiber.Ctx, params GetAPIKeysParams) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Parse pagination parameters with defaults
	pageSize := params.PageSize
	pageIndex := params.PageIndex

	// Apply defaults if parameters are zero
	if pageSize == 0 {
		pageSize = 20
	}
	if pageIndex == 0 {
		pageIndex = 1
	}

	// Validate parameters
	if pageSize < 1 || pageSize > 1000 {
		return handler.SendError(c, fiber.StatusBadRequest, "pageSize must be between 1 and 1000")
	}
	if pageIndex < 1 {
		return handler.SendError(c, fiber.StatusBadRequest, "pageIndex must be at least 1")
	}

	// Get paginated API keys for the user
	apiKeys, totalCount, err := model.GetUserAPIKeysWithPagination(user.ID, pageSize, pageIndex)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to get API keys")
	}

	// Convert to API format
	var apiKeyList []APIKey
	for _, apiKey := range apiKeys {
		apiAPIKey, err := convertToApiAPIKey(&apiKey)
		if err != nil {
			return handler.SendError(c, fiber.StatusInternalServerError, "failed to convert API key")
		}
		apiKeyList = append(apiKeyList, *apiAPIKey)
	}

	// #nosec G115
	response := APIKeysResponse{
		ApiKeys:    apiKeyList,
		TotalCount: int(totalCount),
		PageSize:   pageSize,
		PageIndex:  pageIndex,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// CreateAPIKey - Create a new API key
func (s Server) CreateAPIKey(c *fiber.Ctx) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var req CreateAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "invalid request body")
	}

	// Convert vault unique IDs to database IDs
	var vaultIDs []uint
	if req.VaultUniqueIds != nil && len(*req.VaultUniqueIds) > 0 {
		for _, uniqueID := range *req.VaultUniqueIds {
			var vault model.Vault
			err := model.DB.Where("unique_id = ? AND user_id = ?", uniqueID, user.ID).First(&vault).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return handler.SendError(c, fiber.StatusBadRequest, "vault not found: "+uniqueID)
				}
				return handler.SendError(c, fiber.StatusInternalServerError, "failed to validate vault")
			}
			vaultIDs = append(vaultIDs, vault.ID)
		}
	}

	// Create API key parameters
	params := model.CreateAPIKeyParams{
		UserID:   user.ID,
		Name:     req.Name,
		VaultIDs: vaultIDs,
	}

	if req.ExpiresAt != nil {
		params.ExpiresAt = req.ExpiresAt
	}

	// Validate parameters
	if validationErrors := params.Validate(); len(validationErrors) > 0 {
		return handler.SendError(c, fiber.StatusBadRequest, "validation failed")
	}

	// Create the API key
	apiKey, plainKey, err := params.Create()
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to create API key")
	}

	// Convert to API format
	apiAPIKey, err := convertToApiAPIKey(apiKey)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to convert API key")
	}

	// Create audit log for API key creation
	auditAPIKeyOperation(c, model.ActionCreateAPIKey, user.ID, apiKey.ID, apiKey.Name)

	response := CreateAPIKeyResponse{
		ApiKey: *apiAPIKey,
		Key:    plainKey,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// UpdateAPIKey - Update an API key (enable/disable or modify properties)
func (s Server) UpdateAPIKey(c *fiber.Ctx, id int64) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Find the API key
	var apiKey model.APIKey
	err = model.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&apiKey).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return handler.SendError(c, fiber.StatusNotFound, "API key not found")
		}
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to get API key")
	}

	var req UpdateAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "invalid request body")
	}

	// Convert vault unique IDs to database IDs if provided
	var vaultIDs *[]uint
	if req.VaultUniqueIds != nil {
		var ids []uint
		for _, uniqueID := range *req.VaultUniqueIds {
			var vault model.Vault
			err := model.DB.Where("unique_id = ? AND user_id = ?", uniqueID, user.ID).First(&vault).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return handler.SendError(c, fiber.StatusBadRequest, "vault not found: "+uniqueID)
				}
				return handler.SendError(c, fiber.StatusInternalServerError, "failed to validate vault")
			}
			ids = append(ids, vault.ID)
		}
		vaultIDs = &ids
	}

	// Update parameters
	updateParams := model.UpdateAPIKeyParams{
		Name:      req.Name,
		VaultIDs:  vaultIDs,
		ExpiresAt: req.ExpiresAt,
		IsActive:  req.IsActive,
	}

	// Validate parameters
	if validationErrors := updateParams.ValidateForUpdate(user.ID, apiKey.ID); len(validationErrors) > 0 {
		return handler.SendError(c, fiber.StatusBadRequest, "validation failed")
	}

	// Update the API key
	err = apiKey.Update(updateParams)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to update API key")
	}

	// Convert to API format
	apiAPIKey, err := convertToApiAPIKey(&apiKey)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to convert API key")
	}

	// Create audit log for API key update
	auditAPIKeyOperation(c, model.ActionUpdateAPIKey, user.ID, apiKey.ID, apiKey.Name)

	return c.Status(fiber.StatusOK).JSON(*apiAPIKey)
}

// DeleteAPIKey - Delete an API key
func (s Server) DeleteAPIKey(c *fiber.Ctx, id int64) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Find the API key
	var apiKey model.APIKey
	err = model.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&apiKey).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return handler.SendError(c, fiber.StatusNotFound, "API key not found")
		}
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to get API key")
	}

	// Delete the API key (soft delete)
	err = apiKey.Delete()
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to delete API key")
	}

	// Create audit log for API key deletion
	auditAPIKeyOperation(c, model.ActionDeleteAPIKey, user.ID, apiKey.ID, apiKey.Name)

	return c.SendStatus(fiber.StatusNoContent)
}
