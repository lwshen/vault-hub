package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/model"
	"github.com/lwshen/vault-hub/packages/api/generated/models"
	"gorm.io/gorm"
)

// auditAPIKeyOperation creates an audit log entry for API key operations
func auditAPIKeyOperation(ctx echo.Context, action model.ActionType, userID uint, apiKeyID uint, apiKeyName string) {
	ip, userAgent := getClientInfoEcho(ctx)

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

// GetAPIKeys - Get API keys for the current user with pagination
func (c *Container) GetAPIKeys(ctx echo.Context) error {
	user, err := getUserFromEchoContext(ctx)
	if err != nil {
		return err
	}

	// Parse pagination parameters with defaults
	pageSize := 20
	if ps := ctx.QueryParam("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil {
			pageSize = v
		}
	}

	pageIndex := 1
	if pi := ctx.QueryParam("pageIndex"); pi != "" {
		if v, err := strconv.Atoi(pi); err == nil {
			pageIndex = v
		}
	}

	// Validate parameters
	if pageSize < 1 || pageSize > 1000 {
		return SendError(ctx, http.StatusBadRequest, "pageSize must be between 1 and 1000")
	}
	if pageIndex < 1 {
		return SendError(ctx, http.StatusBadRequest, "pageIndex must be at least 1")
	}

	// Get paginated API keys for the user
	apiKeys, totalCount, err := model.GetUserAPIKeysWithPagination(user.ID, pageSize, pageIndex)
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, "failed to get API keys")
	}

	// Convert to API format - initialize with empty slice to ensure [] in JSON
	apiKeyList := make([]models.VaultApiKey, 0)
	for _, apiKey := range apiKeys {
		apiAPIKey, err := convertToGeneratedAPIKeyWithVaults(&apiKey)
		if err != nil {
			return SendError(ctx, http.StatusInternalServerError, "failed to convert API key")
		}
		apiKeyList = append(apiKeyList, *apiAPIKey)
	}

	// #nosec G115
	response := models.ApiKeysResponse{
		ApiKeys:    apiKeyList,
		TotalCount: int32(totalCount),
		PageSize:   int32(pageSize),
		PageIndex:  int32(pageIndex),
	}

	return ctx.JSON(http.StatusOK, response)
}

// CreateAPIKey handles the creation of a new API key
func (c *Container) CreateAPIKey(ctx echo.Context) error {
	// Get authenticated user
	user, err := getUserFromEchoContext(ctx)
	if err != nil {
		return err
	}

	// Parse and validate request
	var req models.CreateApiKeyRequest
	if err := ctx.Bind(&req); err != nil {
		return SendError(ctx, http.StatusBadRequest, "invalid request body")
	}

	// Convert vault unique IDs to database IDs
	vaultIDs, err := convertVaultUniqueIDsToDBIDs(req.VaultUniqueIds, user.ID)
	if err != nil {
		return err
	}

	// Build API key creation parameters
	params := buildCreateAPIKeyParams(req, user.ID, vaultIDs)

	// Validate parameters
	if validationErrors := params.Validate(); len(validationErrors) > 0 {
		return SendError(ctx, http.StatusBadRequest, "validation failed")
	}

	// Create the API key
	apiKey, plainKey, err := params.Create()
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, "failed to create API key")
	}

	// Convert to API format
	apiAPIKey, err := convertToGeneratedAPIKeyWithVaults(apiKey)
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, "failed to convert API key")
	}

	// Log the creation for audit purposes
	auditAPIKeyOperation(ctx, model.ActionCreateAPIKey, user.ID, apiKey.ID, apiKey.Name)

	response := models.CreateApiKeyResponse{
		ApiKey: *apiAPIKey,
		Key:    plainKey,
	}

	return ctx.JSON(http.StatusCreated, response)
}

// convertVaultUniqueIDsToDBIDs converts vault unique IDs to database IDs
func convertVaultUniqueIDsToDBIDs(vaultUniqueIds []string, userID uint) ([]uint, error) {
	if len(vaultUniqueIds) == 0 {
		return []uint{}, nil
	}

	var vaultIDs []uint
	for _, uniqueID := range vaultUniqueIds {
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
func buildCreateAPIKeyParams(req models.CreateApiKeyRequest, userID uint, vaultIDs []uint) model.CreateAPIKeyParams {
	params := model.CreateAPIKeyParams{
		UserID:   userID,
		Name:     req.Name,
		VaultIDs: vaultIDs,
	}

	// Handle optional ExpiresAt field
	if req.ExpiresAt != nil {
		params.ExpiresAt = req.ExpiresAt
	}

	return params
}

// UpdateAPIKey handles updating an existing API key's properties
func (c *Container) UpdateAPIKey(ctx echo.Context) error {
	// Get authenticated user
	user, err := getUserFromEchoContext(ctx)
	if err != nil {
		return err
	}

	// Parse ID from path parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return SendError(ctx, http.StatusBadRequest, "invalid API key ID")
	}

	// Find and validate API key ownership
	apiKey, err := findAPIKeyByID(id, user.ID)
	if err != nil {
		return SendError(ctx, http.StatusNotFound, "API key not found")
	}

	// Parse update request
	var req models.UpdateApiKeyRequest
	if err := ctx.Bind(&req); err != nil {
		return SendError(ctx, http.StatusBadRequest, "invalid request body")
	}

	// Convert vault unique IDs to database IDs if provided
	var vaultIDs *[]uint
	if len(req.VaultUniqueIds) > 0 {
		ids, err := convertVaultUniqueIDsToDBIDs(req.VaultUniqueIds, user.ID)
		if err != nil {
			return err
		}
		vaultIDs = &ids
	}

	// Build update parameters
	updateParams := buildUpdateAPIKeyParams(req, vaultIDs)

	// Validate update parameters
	if validationErrors := updateParams.ValidateForUpdate(user.ID, apiKey.ID); len(validationErrors) > 0 {
		return SendError(ctx, http.StatusBadRequest, "validation failed")
	}

	// Perform the update
	if err := apiKey.Update(updateParams); err != nil {
		return SendError(ctx, http.StatusInternalServerError, "failed to update API key")
	}

	// Convert to API format
	apiAPIKey, err := convertToGeneratedAPIKeyWithVaults(apiKey)
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, "failed to convert API key")
	}

	// Log the update for audit purposes
	auditAPIKeyOperation(ctx, model.ActionUpdateAPIKey, user.ID, apiKey.ID, apiKey.Name)

	return ctx.JSON(http.StatusOK, *apiAPIKey)
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

// buildUpdateAPIKeyParams constructs API key update parameters
func buildUpdateAPIKeyParams(req models.UpdateApiKeyRequest, vaultIDs *[]uint) model.UpdateAPIKeyParams {
	params := model.UpdateAPIKeyParams{
		VaultIDs: vaultIDs,
	}

	// Handle optional Name field
	if req.Name != "" {
		params.Name = &req.Name
	}

	// Handle optional ExpiresAt field
	if req.ExpiresAt != nil {
		params.ExpiresAt = req.ExpiresAt
	}

	return params
}

// DeleteAPIKey - Delete an API key
func (c *Container) DeleteAPIKey(ctx echo.Context) error {
	user, err := getUserFromEchoContext(ctx)
	if err != nil {
		return err
	}

	// Parse ID from path parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return SendError(ctx, http.StatusBadRequest, "invalid API key ID")
	}

	// Find the API key
	var apiKey model.APIKey
	err = model.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&apiKey).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return SendError(ctx, http.StatusNotFound, "API key not found")
		}
		return SendError(ctx, http.StatusInternalServerError, "failed to get API key")
	}

	// Delete the API key (soft delete)
	err = apiKey.Delete()
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, "failed to delete API key")
	}

	// Create audit log for API key deletion
	auditAPIKeyOperation(ctx, model.ActionDeleteAPIKey, user.ID, apiKey.ID, apiKey.Name)

	return ctx.NoContent(http.StatusNoContent)
}
