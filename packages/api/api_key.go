package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/lwshen/vault-hub/model"
	"gorm.io/gorm"
)

// auditAPIKeyOperation creates an audit log entry for API key operations
func auditAPIKeyOperation(clientInfo ClientInfo, action model.ActionType, userID uint, apiKeyID uint, apiKeyName string) {
	err := model.LogAPIKeyAction(apiKeyID, action, userID, model.SourceWeb, clientInfo.IP, clientInfo.UserAgent)

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

// GetAPIKeysForUser lists API keys for the provided user.
func GetAPIKeysForUser(user *model.User, params GetAPIKeysParams) (APIKeysResponse, *APIError) {
	if user == nil {
		return APIKeysResponse{}, newAPIError(http.StatusUnauthorized, "user not found in context")
	}

	pageSize := params.PageSize
	pageIndex := params.PageIndex

	if pageSize == 0 {
		pageSize = 20
	}
	if pageIndex == 0 {
		pageIndex = 1
	}

	if pageSize < 1 || pageSize > 1000 {
		return APIKeysResponse{}, newAPIError(http.StatusBadRequest, "pageSize must be between 1 and 1000")
	}
	if pageIndex < 1 {
		return APIKeysResponse{}, newAPIError(http.StatusBadRequest, "pageIndex must be at least 1")
	}

	apiKeys, totalCount, err := model.GetUserAPIKeysWithPagination(user.ID, pageSize, pageIndex)
	if err != nil {
		return APIKeysResponse{}, newAPIError(http.StatusInternalServerError, "failed to get API keys")
	}

	apiKeyList := make([]VaultAPIKey, 0)
	for _, apiKey := range apiKeys {
		apiAPIKey, err := convertToApiAPIKey(&apiKey)
		if err != nil {
			return APIKeysResponse{}, newAPIError(http.StatusInternalServerError, "failed to convert API key")
		}
		apiKeyList = append(apiKeyList, *apiAPIKey)
	}

	response := APIKeysResponse{
		ApiKeys:    apiKeyList,
		TotalCount: int(totalCount),
		PageSize:   pageSize,
		PageIndex:  pageIndex,
	}

	return response, nil
}

// CreateAPIKeyForUser provisions a new API key for the given user.
func CreateAPIKeyForUser(user *model.User, req CreateAPIKeyRequest, clientInfo ClientInfo) (CreateAPIKeyResponse, *APIError) {
	if user == nil {
		return CreateAPIKeyResponse{}, newAPIError(http.StatusUnauthorized, "user not found in context")
	}

	vaultIDs, err := convertVaultUniqueIDs(req.VaultUniqueIds, user.ID)
	if err != nil {
		return CreateAPIKeyResponse{}, newAPIError(http.StatusBadRequest, err.Error())
	}

	params := buildCreateAPIKeyParams(req, user.ID, vaultIDs)

	if err := validateCreateAPIKeyParams(params); err != nil {
		return CreateAPIKeyResponse{}, newAPIError(http.StatusBadRequest, "validation failed")
	}

	apiKey, plainKey, err := createAPIKey(params)
	if err != nil {
		return CreateAPIKeyResponse{}, newAPIError(http.StatusInternalServerError, "failed to create API key")
	}

	apiAPIKey, err := convertToApiAPIKey(apiKey)
	if err != nil {
		return CreateAPIKeyResponse{}, newAPIError(http.StatusInternalServerError, "failed to convert API key")
	}

	auditAPIKeyOperation(clientInfo, model.ActionCreateAPIKey, user.ID, apiKey.ID, apiKey.Name)

	response := CreateAPIKeyResponse{
		ApiKey: *apiAPIKey,
		Key:    plainKey,
	}

	return response, nil
}

// UpdateAPIKeyForUser updates an API key if owned by the provided user.
func UpdateAPIKeyForUser(user *model.User, id int64, req UpdateAPIKeyRequest, clientInfo ClientInfo) (VaultAPIKey, *APIError) {
	if user == nil {
		return VaultAPIKey{}, newAPIError(http.StatusUnauthorized, "user not found in context")
	}

	apiKey, apiErr := findAPIKeyByID(id, user.ID)
	if apiErr != nil {
		return VaultAPIKey{}, apiErr
	}

	vaultIDs, err := convertOptionalVaultUniqueIDs(req.VaultUniqueIds, user.ID)
	if err != nil {
		return VaultAPIKey{}, newAPIError(http.StatusBadRequest, err.Error())
	}

	updateParams := buildUpdateAPIKeyParams(req, vaultIDs)

	if err := validateUpdateAPIKeyParams(updateParams, user.ID, apiKey.ID); err != nil {
		return VaultAPIKey{}, newAPIError(http.StatusBadRequest, "validation failed")
	}

	if err := updateAPIKey(apiKey, updateParams); err != nil {
		return VaultAPIKey{}, newAPIError(http.StatusInternalServerError, "failed to update API key")
	}

	apiAPIKey, err := convertToApiAPIKey(apiKey)
	if err != nil {
		return VaultAPIKey{}, newAPIError(http.StatusInternalServerError, "failed to convert API key")
	}

	auditAPIKeyOperation(clientInfo, model.ActionUpdateAPIKey, user.ID, apiKey.ID, apiKey.Name)

	return *apiAPIKey, nil
}

// DeleteAPIKeyForUser removes an API key owned by the provided user.
func DeleteAPIKeyForUser(user *model.User, id int64, clientInfo ClientInfo) *APIError {
	if user == nil {
		return newAPIError(http.StatusUnauthorized, "user not found in context")
	}

	apiKey, apiErr := findAPIKeyByID(id, user.ID)
	if apiErr != nil {
		return apiErr
	}

	if err := apiKey.Delete(); err != nil {
		return newAPIError(http.StatusInternalServerError, "failed to delete API key")
	}

	auditAPIKeyOperation(clientInfo, model.ActionDeleteAPIKey, user.ID, apiKey.ID, apiKey.Name)

	return nil
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

// findAPIKeyByID finds an API key by ID and validates user ownership
func findAPIKeyByID(id int64, userID uint) (*model.APIKey, *APIError) {
	var apiKey model.APIKey
	err := model.DB.Where("id = ? AND user_id = ?", id, userID).First(&apiKey).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, newAPIError(http.StatusNotFound, "API key not found")
		}
		return nil, newAPIError(http.StatusInternalServerError, "failed to get API key")
	}
	return &apiKey, nil
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
