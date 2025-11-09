package api

import (
	"fmt"
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

	err := model.LogAPIKeyAction(apiKeyID, action, userID, model.SourceWeb, ip, userAgent)

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
func convertToApiAPIKey(apiKey *model.APIKey) (*VaultAPIKey, error) {
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
	return &VaultAPIKey{
		Id:         id,
		Name:       apiKey.Name,
		Vaults:     &apiVaults,
		ExpiresAt:  expiresAt,
		LastUsedAt: lastUsedAt,
		IsActive:   !apiKey.DeletedAt.Valid,
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

	// Convert to API format - initialize with empty slice to ensure [] in JSON
	apiKeyList := make([]VaultAPIKey, 0)
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

// CreateAPIKey handles the creation of a new API key
// It validates the request, converts vault IDs, creates the key, and returns the response
func (s Server) CreateAPIKey(c *fiber.Ctx) error {
	// Get authenticated user
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Parse and validate request
	req, err := parseCreateAPIKeyRequest(c)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "invalid request body")
	}

	// Convert vault unique IDs to database IDs
	vaultIDs, err := convertVaultUniqueIDs(req.VaultUniqueIds, user.ID)
	if err != nil {
		return err
	}

	// Build API key creation parameters
	params := buildCreateAPIKeyParams(req, user.ID, vaultIDs)

	// Validate parameters
	if err := validateCreateAPIKeyParams(params); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "validation failed")
	}

	// Create the API key
	apiKey, plainKey, err := createAPIKey(params)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to create API key")
	}

	// Convert to API format
	apiAPIKey, err := convertToApiAPIKey(apiKey)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to convert API key")
	}

	// Log the creation for audit purposes
	auditAPIKeyOperation(c, model.ActionCreateAPIKey, user.ID, apiKey.ID, apiKey.Name)

	response := CreateAPIKeyResponse{
		ApiKey: *apiAPIKey,
		Key:    plainKey,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// parseCreateAPIKeyRequest parses the request body for API key creation
func parseCreateAPIKeyRequest(c *fiber.Ctx) (CreateAPIKeyRequest, error) {
	var req CreateAPIKeyRequest
	err := c.BodyParser(&req)
	return req, err
}

// convertVaultUniqueIDs converts vault unique IDs to database IDs
func convertVaultUniqueIDs(vaultUniqueIds *[]string, userID uint) ([]uint, error) {
	if vaultUniqueIds == nil || len(*vaultUniqueIds) == 0 {
		return []uint{}, nil
	}

	var vaultIDs []uint
	for _, uniqueID := range *vaultUniqueIds {
		vault, err := findVaultByUniqueID(uniqueID, userID)
		if err != nil {
			return nil, err
		}
		vaultIDs = append(vaultIDs, vault.ID)
	}
	return vaultIDs, nil
}

// findVaultByUniqueID finds a vault by unique ID and user ID
func findVaultByUniqueID(uniqueID string, userID uint) (*model.Vault, error) {
	var vault model.Vault
	err := model.DB.Where("unique_id = ? AND user_id = ?", uniqueID, userID).First(&vault).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("vault not found: %s", uniqueID)
		}
		return nil, fmt.Errorf("failed to validate vault %s: %v", uniqueID, err)
	}
	return &vault, nil
}

// buildCreateAPIKeyParams constructs API key creation parameters
func buildCreateAPIKeyParams(req CreateAPIKeyRequest, userID uint, vaultIDs []uint) model.CreateAPIKeyParams {
	params := model.CreateAPIKeyParams{
		UserID:   userID,
		Name:     req.Name,
		VaultIDs: vaultIDs,
	}

	if req.ExpiresAt != nil {
		params.ExpiresAt = req.ExpiresAt
	}

	return params
}

// validateCreateAPIKeyParams validates the API key creation parameters
func validateCreateAPIKeyParams(params model.CreateAPIKeyParams) error {
	if validationErrors := params.Validate(); len(validationErrors) > 0 {
		return fmt.Errorf("validation errors occurred")
	}
	return nil
}

// createAPIKey creates a new API key with the given parameters
func createAPIKey(params model.CreateAPIKeyParams) (*model.APIKey, string, error) {
	return params.Create()
}

// UpdateAPIKey handles updating an existing API key's properties
// It validates ownership, processes the update request, and returns the updated key
func (s Server) UpdateAPIKey(c *fiber.Ctx, id int64) error {
	// Get authenticated user
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Find and validate API key ownership
	apiKey, err := findAPIKeyByID(id, user.ID)
	if err != nil {
		return err
	}

	// Parse update request
	req, err := parseUpdateAPIKeyRequest(c)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "invalid request body")
	}

	// Convert vault unique IDs to database IDs if provided
	vaultIDs, err := convertOptionalVaultUniqueIDs(req.VaultUniqueIds, user.ID)
	if err != nil {
		return err
	}

	// Build update parameters
	updateParams := buildUpdateAPIKeyParams(req, vaultIDs)

	// Validate update parameters
	if err := validateUpdateAPIKeyParams(updateParams, user.ID, apiKey.ID); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "validation failed")
	}

	// Perform the update
	if err := updateAPIKey(apiKey, updateParams); err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to update API key")
	}

	// Convert to API format
	apiAPIKey, err := convertToApiAPIKey(apiKey)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to convert API key")
	}

	// Log the update for audit purposes
	auditAPIKeyOperation(c, model.ActionUpdateAPIKey, user.ID, apiKey.ID, apiKey.Name)

	return c.Status(fiber.StatusOK).JSON(*apiAPIKey)
}

// findAPIKeyByID finds an API key by ID and validates user ownership
func findAPIKeyByID(id int64, userID uint) (*model.APIKey, error) {
	var apiKey model.APIKey
	err := model.DB.Where("id = ? AND user_id = ?", id, userID).First(&apiKey).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("API key not found")
		}
		return nil, fmt.Errorf("failed to get API key: %v", err)
	}
	return &apiKey, nil
}

// parseUpdateAPIKeyRequest parses the request body for API key updates
func parseUpdateAPIKeyRequest(c *fiber.Ctx) (UpdateAPIKeyRequest, error) {
	var req UpdateAPIKeyRequest
	err := c.BodyParser(&req)
	return req, err
}

// convertOptionalVaultUniqueIDs converts optional vault unique IDs to database IDs
func convertOptionalVaultUniqueIDs(vaultUniqueIds *[]string, userID uint) (*[]uint, error) {
	if vaultUniqueIds == nil {
		return nil, nil
	}

	var ids []uint
	for _, uniqueID := range *vaultUniqueIds {
		vault, err := findVaultByUniqueID(uniqueID, userID)
		if err != nil {
			return nil, err
		}
		ids = append(ids, vault.ID)
	}
	return &ids, nil
}

// buildUpdateAPIKeyParams constructs API key update parameters
func buildUpdateAPIKeyParams(req UpdateAPIKeyRequest, vaultIDs *[]uint) model.UpdateAPIKeyParams {
	return model.UpdateAPIKeyParams{
		Name:      req.Name,
		VaultIDs:  vaultIDs,
		ExpiresAt: req.ExpiresAt,
	}
}

// validateUpdateAPIKeyParams validates the API key update parameters
func validateUpdateAPIKeyParams(params model.UpdateAPIKeyParams, userID uint, apiKeyID uint) error {
	if validationErrors := params.ValidateForUpdate(userID, apiKeyID); len(validationErrors) > 0 {
		return fmt.Errorf("validation errors occurred")
	}
	return nil
}

// updateAPIKey performs the actual API key update
func updateAPIKey(apiKey *model.APIKey, params model.UpdateAPIKeyParams) error {
	return apiKey.Update(params)
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

// GetAPIKeyUsage - Get usage statistics for a specific API key
func (s Server) GetAPIKeyUsage(c *fiber.Ctx, id int64) error {
	// Get authenticated user
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Find and validate API key ownership
	apiKey, err := findAPIKeyByID(id, user.ID)
	if err != nil {
		if err.Error() == "API key not found" {
			return handler.SendError(c, fiber.StatusNotFound, "API key not found")
		}
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to get API key")
	}

	// Get usage statistics
	stats, err := model.GetAPIKeyUsageStats(apiKey.ID)
	if err != nil {
		slog.Error("Failed to get API key usage statistics",
			"api_key_id", apiKey.ID,
			"user_id", user.ID,
			"error", err)
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to get usage statistics")
	}

	// Convert vault breakdown to API format
	vaultBreakdown := make([]VaultUsageBreakdown, 0)
	for _, vb := range stats.VaultBreakdown {
		// #nosec G115
		vaultBreakdown = append(vaultBreakdown, VaultUsageBreakdown{
			VaultId:       int64(vb.VaultID),
			VaultName:     vb.VaultName,
			VaultUniqueId: vb.VaultUniqueID,
			AccessCount:   vb.AccessCount,
		})
	}

	response := APIKeyUsageResponse{
		TotalRequests:    stats.TotalRequests,
		Last24Hours:      stats.Last24Hours,
		Last7Days:        stats.Last7Days,
		Last30Days:       stats.Last30Days,
		LastUsedAt:       stats.LastUsedAt,
		VaultAccessCount: stats.VaultAccessCount,
		VaultBreakdown:   vaultBreakdown,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
