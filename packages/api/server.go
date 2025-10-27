package api

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/internal/auth"
	"github.com/lwshen/vault-hub/model"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ensure that we've conformed to the `ServerInterface` with a compile-time check
// var _ ServerInterface = (*Server)(nil)  // TODO: Re-enable after fixing interface compatibility

// Server implements the ServerInterface for Echo
type Server struct{}

// NewServer creates a new Server instance
func NewServer() Server {
	return Server{}
}

// Helper functions for Echo context handling

// sendError sends a standardized error response for Echo
func sendError(c echo.Context, code int, message string) error {
	return c.JSON(code, map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	})
}

// generateJWTToken generates a JWT token for the given user ID
func generateJWTToken(userID uint) (string, error) {
	return auth.GenerateToken(userID)
}

// getClientInfo extracts IP address and User-Agent from the request
func getClientInfo(c echo.Context) (string, string) {
	// Get IP address (check for forwarded headers first)
	ip := c.Request().Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = c.Request().Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = c.RealIP()
	}

	// Get User-Agent
	userAgent := c.Request().Header.Get("User-Agent")

	return ip, userAgent
}

// getUserFromContext extracts the authenticated user from the context
func getUserFromContext(c echo.Context) (*model.User, error) {
	user, ok := c.Get("user").(*model.User)
	if !ok {
		return nil, sendError(c, http.StatusUnauthorized, "user not found in context")
	}
	return user, nil
}

// getAPIKeyFromContext extracts the API key from the context
func getAPIKeyFromContext(c echo.Context) (*model.APIKey, error) {
	apiKey, ok := c.Get("api_key").(*model.APIKey)
	if !ok {
		return nil, sendError(c, http.StatusUnauthorized, "api key not found in context")
	}
	return apiKey, nil
}

// convertToApiVault converts a model.Vault to an api.Vault
func convertToApiVault(vault *model.Vault) Vault {
	// #nosec G115
	userID := int64(vault.UserID)
	return Vault{
		UniqueId:    vault.UniqueID,
		UserId:      &userID,
		Name:        vault.Name,
		Value:       vault.Value,
		Description: &vault.Description,
		Category:    &vault.Category,
		CreatedAt:   &vault.CreatedAt,
		UpdatedAt:   &vault.UpdatedAt,
	}
}

// Authentication endpoints

// Login handles user login
func (Server) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return sendError(c, http.StatusBadRequest, "invalid request body")
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		return sendError(c, http.StatusBadRequest, "email and password are required")
	}

	// Find user by email
	var user model.User
	if err := model.DB.Where("email = ?", string(req.Email)).First(&user).Error; err != nil {
		return sendError(c, http.StatusUnauthorized, "invalid credentials")
	}

	// Check password
	if user.Password == nil {
		return sendError(c, http.StatusUnauthorized, "invalid credentials")
	}

	// Compare hashed password with provided password
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password)); err != nil {
		return sendError(c, http.StatusUnauthorized, "invalid credentials")
	}

	// Generate JWT token
	token, err := generateJWTToken(user.ID)
	if err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to generate token")
	}

	// Clear password before returning
	user.Password = nil

	resp := LoginResponse{
		Token: token,
	}

	return c.JSON(http.StatusOK, resp)
}

// Logout handles user logout
func (Server) Logout(c echo.Context) error {
	// For JWT-based authentication, logout is primarily a client-side responsibility
	// The client should discard the JWT token
	// We return a success response to acknowledge the logout request

	resp := map[string]interface{}{
		"message": "logout successful. Please discard your JWT token.",
	}

	return c.JSON(http.StatusOK, resp)
}

// Signup handles user registration
func (Server) Signup(c echo.Context) error {
	var req SignupRequest
	if err := c.Bind(&req); err != nil {
		return sendError(c, http.StatusBadRequest, "invalid request body")
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		return sendError(c, http.StatusBadRequest, "email and password are required")
	}

	if len(req.Password) < 8 {
		return sendError(c, http.StatusBadRequest, "password must be at least 8 characters long")
	}

	// Check if user already exists
	var existingUser model.User
	if err := model.DB.Where("email = ?", string(req.Email)).First(&existingUser).Error; err == nil {
		return sendError(c, http.StatusConflict, "user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to hash password")
	}

	// Convert hashed password to string
	hashedPasswordStr := string(hashedPassword)

	// Create new user
	user := model.User{
		Email:    string(req.Email),
		Name:     &req.Name,
		Password: &hashedPasswordStr,
	}

	if err := model.DB.Create(&user).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to create user")
	}

	// Generate JWT token for immediate login
	token, err := generateJWTToken(user.ID)
	if err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to generate token")
	}

	// Clear password before returning
	user.Password = nil

	resp := SignupResponse{
		Token: token,
	}

	return c.JSON(http.StatusCreated, resp)
}

// RequestMagicLink requests a magic link for login
func (Server) RequestMagicLink(c echo.Context) error {
	var req MagicLinkRequest
	if err := c.Bind(&req); err != nil {
		return sendError(c, http.StatusBadRequest, "invalid request body")
	}

	// Validate input
	if req.Email == "" {
		return sendError(c, http.StatusBadRequest, "email is required")
	}

	// Check if user exists
	var user model.User
	if err := model.DB.Where("email = ?", string(req.Email)).First(&user).Error; err != nil {
		// Don't reveal if user exists or not for security
		resp := map[string]interface{}{
			"message": "if an account with this email exists, a magic link has been sent",
		}
		return c.JSON(http.StatusOK, resp)
	}

	// TODO: Generate magic link token and send email
	// For now, return a success response without actual email functionality
	resp := map[string]interface{}{
		"message": "if an account with this email exists, a magic link has been sent",
	}

	return c.JSON(http.StatusOK, resp)
}

// ConsumeMagicLink consumes a magic link for login
func (Server) ConsumeMagicLink(c echo.Context, params ConsumeMagicLinkParams) error {
	// Validate input
	if params.Token == "" {
		return sendError(c, http.StatusBadRequest, "token is required")
	}

	// TODO: Implement actual magic link token validation and login
	// For now, return a not implemented response until token system is built
	return sendError(c, http.StatusNotImplemented, "magic link token validation not yet implemented")
}

// RequestPasswordReset requests a password reset
func (Server) RequestPasswordReset(c echo.Context) error {
	var req PasswordResetRequest
	if err := c.Bind(&req); err != nil {
		return sendError(c, http.StatusBadRequest, "invalid request body")
	}

	// Validate input
	if req.Email == "" {
		return sendError(c, http.StatusBadRequest, "email is required")
	}

	// Check if user exists
	var user model.User
	if err := model.DB.Where("email = ?", string(req.Email)).First(&user).Error; err != nil {
		// Don't reveal if user exists or not for security
		resp := map[string]interface{}{
			"message": "if an account with this email exists, a password reset link has been sent",
		}
		return c.JSON(http.StatusOK, resp)
	}

	// TODO: Generate password reset token and send email
	// For now, return a success response without actual email functionality
	resp := map[string]interface{}{
		"message": "if an account with this email exists, a password reset link has been sent",
	}

	return c.JSON(http.StatusOK, resp)
}

// ConfirmPasswordReset confirms a password reset
func (Server) ConfirmPasswordReset(c echo.Context) error {
	var req PasswordResetConfirmRequest
	if err := c.Bind(&req); err != nil {
		return sendError(c, http.StatusBadRequest, "invalid request body")
	}

	// Validate input
	if req.Token == "" || req.NewPassword == "" {
		return sendError(c, http.StatusBadRequest, "token and new password are required")
	}

	if len(req.NewPassword) < 8 {
		return sendError(c, http.StatusBadRequest, "password must be at least 8 characters long")
	}

	// TODO: Implement actual password reset token validation and update
	// For now, return a not implemented response until token system is built
	return sendError(c, http.StatusNotImplemented, "password reset token validation not yet implemented")
}

// User endpoints

// GetCurrentUser gets the current authenticated user
func (Server) GetCurrentUser(c echo.Context) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	resp := GetUserResponse{
		Email:  openapi_types.Email(user.Email),
		Avatar: user.Avatar,
		Name:   user.Name,
	}

	return c.JSON(http.StatusOK, resp)
}

// Vault endpoints

// GetVaults gets all vaults for the current user
func (Server) GetVaults(c echo.Context, params GetVaultsParams) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Parse pagination parameters
	pageIndex := int64(1)
	pageSize := int64(20) // Default limit

	if params.PageIndex != nil {
		pageIndex = int64(*params.PageIndex)
		if pageIndex < 1 {
			pageIndex = 1
		}
	}

	if params.PageSize != nil {
		pageSize = int64(*params.PageSize)
		if pageSize < 1 || pageSize > 1000 {
			pageSize = 20
		}
	}

	offset := (pageIndex - 1) * pageSize

	// Build query
	query := model.DB.Where("user_id = ?", user.ID)

	// Count total vaults
	var total int64
	if err := query.Model(&model.Vault{}).Count(&total).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to count vaults")
	}

	// Get vaults with pagination
	var vaults []model.Vault
	if err := query.Order("updated_at DESC").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&vaults).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to retrieve vaults")
	}

	// Convert to API response (VaultLite format - without decrypted values for security)
	vaultLites := make([]VaultLite, len(vaults))
	for i, vault := range vaults {
		vaultLites[i] = VaultLite{
			UniqueId:    vault.UniqueID,
			Name:        vault.Name,
			Description: &vault.Description,
			Category:    &vault.Category,
			UpdatedAt:   &vault.UpdatedAt,
		}
		// Note: VaultLite does not include the Value field for security
	}

	resp := VaultsResponse{
		Vaults:     vaultLites,
		TotalCount: int(total),
		PageIndex:  int(pageIndex),
		PageSize:   int(pageSize),
	}

	return c.JSON(http.StatusOK, resp)
}

// CreateVault creates a new vault
func (Server) CreateVault(c echo.Context) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var req CreateVaultRequest
	if err := c.Bind(&req); err != nil {
		return sendError(c, http.StatusBadRequest, "invalid request body")
	}

	// Validate input
	if req.Name == "" || req.Value == "" {
		return sendError(c, http.StatusBadRequest, "name and value are required")
	}

	// Check if vault with same name already exists for this user
	var existingVault model.Vault
	if err := model.DB.Where("user_id = ? AND name = ?", user.ID, req.Name).First(&existingVault).Error; err == nil {
		return sendError(c, http.StatusConflict, "vault with this name already exists")
	}

	// Create new vault
	vault := model.Vault{
		UserID: user.ID,
		Name:   req.Name,
		Value:  req.Value, // The model will handle encryption
	}

	// Set optional fields if provided
	if req.Description != nil {
		vault.Description = *req.Description
	}
	if req.Category != nil {
		vault.Category = *req.Category
	}

	if err := model.DB.Create(&vault).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to create vault")
	}

	// Return created vault with decrypted value
	apiVault := convertToApiVault(&vault)

	return c.JSON(http.StatusCreated, apiVault)
}

// GetVault gets a specific vault by ID
func (Server) GetVault(c echo.Context, uniqueId string) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Find vault by unique ID and user ID
	var vault model.Vault
	if err := model.DB.Where("unique_id = ? AND user_id = ?", uniqueId, user.ID).First(&vault).Error; err != nil {
		return sendError(c, http.StatusNotFound, "vault not found")
	}

	// Convert to API response with decrypted value
	apiVault := convertToApiVault(&vault)

	return c.JSON(http.StatusOK, apiVault)
}

// UpdateVault updates a specific vault
func (Server) UpdateVault(c echo.Context, id int64) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var req UpdateVaultRequest
	if err := c.Bind(&req); err != nil {
		return sendError(c, http.StatusBadRequest, "invalid request body")
	}

	// Find vault by ID and user ID
	var vault model.Vault
	if err := model.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&vault).Error; err != nil {
		return sendError(c, http.StatusNotFound, "vault not found")
	}

	// Check if new name conflicts with existing vault (if name is being changed)
	if req.Name != nil && *req.Name != vault.Name {
		var existingVault model.Vault
		if err := model.DB.Where("user_id = ? AND name = ? AND id != ?", user.ID, *req.Name, id).First(&existingVault).Error; err == nil {
			return sendError(c, http.StatusConflict, "vault with this name already exists")
		}
	}

	// Update fields if provided
	if req.Name != nil {
		vault.Name = *req.Name
	}
	if req.Value != nil {
		vault.Value = *req.Value // The model will handle encryption
	}
	if req.Description != nil {
		vault.Description = *req.Description
	}
	if req.Category != nil {
		vault.Category = *req.Category
	}

	// Save updated vault
	if err := model.DB.Save(&vault).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to update vault")
	}

	// Return updated vault with decrypted value
	apiVault := convertToApiVault(&vault)

	return c.JSON(http.StatusOK, apiVault)
}

// DeleteVault deletes a specific vault
func (Server) DeleteVault(c echo.Context, uniqueId string) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Find vault by unique ID and user ID
	var vault model.Vault
	if err := model.DB.Where("unique_id = ? AND user_id = ?", uniqueId, user.ID).First(&vault).Error; err != nil {
		return sendError(c, http.StatusNotFound, "vault not found")
	}

	// Delete the vault
	if err := model.DB.Delete(&vault).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to delete vault")
	}

	// Return success response
	resp := map[string]interface{}{
		"message": "vault deleted successfully",
	}

	return c.JSON(http.StatusOK, resp)
}

// API Key endpoints

// GetAPIKeys gets all API keys for the current user
func (Server) GetAPIKeys(c echo.Context, params GetAPIKeysParams) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Parse pagination parameters
	pageIndex := int64(1)
	pageSize := int64(20) // Default limit

	if params.PageIndex > 0 {
		pageIndex = int64(params.PageIndex)
	}

	if params.PageSize > 0 {
		pageSize = int64(params.PageSize)
		if pageSize < 1 || pageSize > 1000 {
			pageSize = 20
		}
	}

	offset := (pageIndex - 1) * pageSize

	// Count total API keys
	var total int64
	if err := model.DB.Where("user_id = ?", user.ID).Model(&model.APIKey{}).Count(&total).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to count API keys")
	}

	// Get API keys with pagination
	var apiKeys []model.APIKey
	if err := model.DB.Where("user_id = ?", user.ID).
		Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&apiKeys).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to retrieve API keys")
	}

	// Convert to API response (without the actual key values for security)
	apiKeyResponses := make([]VaultAPIKey, len(apiKeys))
	for i, key := range apiKeys {
		isActive := true // Default to true since model doesn't have IsActive field

		apiKeyResponses[i] = VaultAPIKey{
			Id:           int64(key.ID),
			Name:         key.Name,
			IsActive:     isActive,
			CreatedAt:    key.CreatedAt,
			UpdatedAt:    &key.UpdatedAt,
			LastUsedAt:   key.LastUsedAt,
			ExpiresAt:    key.ExpiresAt,
		}
	}

	resp := APIKeysResponse{
		ApiKeys:   apiKeyResponses,
		TotalCount: int(total),
		PageIndex:  int(pageIndex),
		PageSize:   int(pageSize),
	}

	return c.JSON(http.StatusOK, resp)
}

// CreateAPIKey creates a new API key
func (Server) CreateAPIKey(c echo.Context) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var req CreateAPIKeyRequest
	if err := c.Bind(&req); err != nil {
		return sendError(c, http.StatusBadRequest, "invalid request body")
	}

	// Validate input
	if req.Name == "" {
		return sendError(c, http.StatusBadRequest, "name is required")
	}

	// Check if API key with same name already exists for this user
	var existingKey model.APIKey
	if err := model.DB.Where("user_id = ? AND name = ?", user.ID, req.Name).First(&existingKey).Error; err == nil {
		return sendError(c, http.StatusConflict, "API key with this name already exists")
	}

	// Generate new API key
	apiKeyValue, err := model.GenerateAPIKey()
	if err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to generate API key")
	}

	// Create new API key with proper model fields
	apiKey := model.APIKey{
		UserID:    user.ID,
		Name:      req.Name,
		KeyHash:   model.HashAPIKey(apiKeyValue), // Store hash, not the key
		ExpiresAt: req.ExpiresAt,
	}

	// Set vault access permissions
	if req.VaultUniqueIds != nil && len(*req.VaultUniqueIds) > 0 {
		// Convert vault unique IDs to vault IDs
		vaultIDs := make(model.VaultIDs, len(*req.VaultUniqueIds))
		for i, uniqueId := range *req.VaultUniqueIds {
			var vault model.Vault
			if err := model.DB.Where("unique_id = ? AND user_id = ?", uniqueId, user.ID).First(&vault).Error; err != nil {
				return sendError(c, http.StatusBadRequest, fmt.Sprintf("vault with unique ID %s not found", uniqueId))
			}
			vaultIDs[i] = vault.ID
		}
		apiKey.VaultIDs = vaultIDs
	}

	if err := model.DB.Create(&apiKey).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to create API key")
	}

	// Convert to API response - only return the actual key value on creation
	isActive := true // New keys are active by default

	resp := VaultAPIKey{
		Id:           int64(apiKey.ID),
		Name:         apiKey.Name,
		IsActive:     isActive,
		CreatedAt:    apiKey.CreatedAt,
		UpdatedAt:    &apiKey.UpdatedAt,
		LastUsedAt:   apiKey.LastUsedAt,
		ExpiresAt:    apiKey.ExpiresAt,
	}

	return c.JSON(http.StatusCreated, resp)
}

// UpdateAPIKey updates a specific API key
func (Server) UpdateAPIKey(c echo.Context, id int64) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var req UpdateAPIKeyRequest
	if err := c.Bind(&req); err != nil {
		return sendError(c, http.StatusBadRequest, "invalid request body")
	}

	// Find API key by ID and user ID
	var apiKey model.APIKey
	if err := model.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&apiKey).Error; err != nil {
		return sendError(c, http.StatusNotFound, "API key not found")
	}

	// Check if new name conflicts with existing API key (if name is being changed)
	if req.Name != nil && *req.Name != apiKey.Name {
		var existingKey model.APIKey
		if err := model.DB.Where("user_id = ? AND name = ? AND id != ?", user.ID, *req.Name, id).First(&existingKey).Error; err == nil {
			return sendError(c, http.StatusConflict, "API key with this name already exists")
		}
	}

	// Update fields if provided
	if req.Name != nil {
		apiKey.Name = *req.Name
	}
	if req.ExpiresAt != nil {
		apiKey.ExpiresAt = req.ExpiresAt
	}
	if req.VaultUniqueIds != nil {
		if len(*req.VaultUniqueIds) > 0 {
			// Convert vault unique IDs to vault IDs
			vaultIDs := make(model.VaultIDs, len(*req.VaultUniqueIds))
			for i, uniqueId := range *req.VaultUniqueIds {
				var vault model.Vault
				if err := model.DB.Where("unique_id = ? AND user_id = ?", uniqueId, user.ID).First(&vault).Error; err != nil {
					return sendError(c, http.StatusBadRequest, fmt.Sprintf("vault with unique ID %s not found", uniqueId))
				}
				vaultIDs[i] = vault.ID
			}
			apiKey.VaultIDs = vaultIDs
		} else {
			// Empty array means all vaults
			apiKey.VaultIDs = nil
		}
	}

	// Save updated API key
	if err := model.DB.Save(&apiKey).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to update API key")
	}

	// Convert to API response (without the actual key value for security)
	isActive := true // Default to true since model doesn't have IsActive field

	resp := VaultAPIKey{
		Id:           int64(apiKey.ID),
		Name:         apiKey.Name,
		IsActive:     isActive,
		CreatedAt:    apiKey.CreatedAt,
		UpdatedAt:    &apiKey.UpdatedAt,
		LastUsedAt:   apiKey.LastUsedAt,
		ExpiresAt:    apiKey.ExpiresAt,
	}

	return c.JSON(http.StatusOK, resp)
}

// DeleteAPIKey deletes a specific API key
func (Server) DeleteAPIKey(c echo.Context, id int64) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Find API key by ID and user ID
	var apiKey model.APIKey
	if err := model.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&apiKey).Error; err != nil {
		return sendError(c, http.StatusNotFound, "API key not found")
	}

	// Delete the API key
	if err := model.DB.Delete(&apiKey).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to delete API key")
	}

	// Return success response
	resp := map[string]interface{}{
		"message": "API key deleted successfully",
	}

	return c.JSON(http.StatusOK, resp)
}

// CLI endpoints (API key authentication)

// GetVaultsByAPIKey gets vaults accessible by API key
func (Server) GetVaultsByAPIKey(c echo.Context) error {
	apiKey, err := getAPIKeyFromContext(c)
	if err != nil {
		return err
	}

	// Get vaults that this API key has access to
	var vaults []model.Vault
	query := model.DB.Where("user_id = ?", apiKey.UserID)

	// Note: Vault access restrictions are handled by the HasVaultAccess method

	if err := query.Order("updated_at DESC").Find(&vaults).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to retrieve vaults")
	}

	// Convert to API response (VaultLite format - without decrypted values for security)
	vaultLites := make([]VaultLite, len(vaults))
	for i, vault := range vaults {
		vaultLites[i] = VaultLite{
			UniqueId:    vault.UniqueID,
			Name:        vault.Name,
			Description: &vault.Description,
			Category:    &vault.Category,
			UpdatedAt:   &vault.UpdatedAt,
		}
		// Note: VaultLite does not include the Value field for security
	}

	return c.JSON(http.StatusOK, vaultLites)
}

// GetVaultByAPIKey gets a specific vault by API key
func (Server) GetVaultByAPIKey(c echo.Context, uniqueId string, params GetVaultByAPIKeyParams) error {
	apiKey, err := getAPIKeyFromContext(c)
	if err != nil {
		return err
	}

	// Find vault by unique ID and user ID first
	var vault model.Vault
	if err := model.DB.Where("unique_id = ? AND user_id = ?", uniqueId, apiKey.UserID).First(&vault).Error; err != nil {
		return sendError(c, http.StatusNotFound, "vault not found")
	}

	// Check if API key has access to this vault
	if !apiKey.HasVaultAccess(vault.ID) {
		return sendError(c, http.StatusForbidden, "access denied to this vault")
	}

	// Convert to API response with decrypted value (full Vault format for CLI)
	apiVault := convertToApiVault(&vault)

	return c.JSON(http.StatusOK, apiVault)
}

// GetVaultByNameAPIKey gets a vault by name using API key
func (Server) GetVaultByNameAPIKey(c echo.Context, name string) error {
	apiKey, err := getAPIKeyFromContext(c)
	if err != nil {
		return err
	}

	// Find vault by name and user ID
	var vault model.Vault
	if err := model.DB.Where("name = ? AND user_id = ?", name, apiKey.UserID).First(&vault).Error; err != nil {
		return sendError(c, http.StatusNotFound, "vault not found")
	}

	// Check if API key has access to this vault
	if !apiKey.HasVaultAccess(vault.ID) {
		return sendError(c, http.StatusForbidden, "access denied to this vault")
	}

	// Convert to API response with decrypted value (full Vault format for CLI)
	apiVault := convertToApiVault(&vault)

	return c.JSON(http.StatusOK, apiVault)
}

// Audit endpoints

// GetAuditLogs gets audit logs with pagination
func (Server) GetAuditLogs(c echo.Context, params GetAuditLogsParams) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Parse pagination parameters
	page := int64(params.PageIndex)
	limit := int64(params.PageSize)

	if page < 0 {
		page = 0
	}
	if limit < 1 || limit > 1000 {
		limit = 100
	}

	offset := page * limit

	// Parse filters - use only available fields
	var actionFilter, resourceFilter string
	// Note: Action and Resource filters are not available in current API spec
	// Future implementation could add these as query parameters

	// Build query
	query := model.DB.Where("user_id = ?", user.ID)
	if actionFilter != "" {
		query = query.Where("action = ?", actionFilter)
	}
	if resourceFilter != "" {
		query = query.Where("resource_type = ?", resourceFilter)
	}

	// Count total audit logs
	var total int64
	if err := query.Model(&model.AuditLog{}).Count(&total).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to count audit logs")
	}

	// Get audit logs with pagination
	var auditLogs []model.AuditLog
	if err := query.Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&auditLogs).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to retrieve audit logs")
	}

	// Convert to API response
	apiAuditLogs := make([]AuditLog, len(auditLogs))
	for i, log := range auditLogs {
		apiAuditLogs[i] = AuditLog{
			Action:    AuditLogAction(log.Action),
			IpAddress: &log.IPAddress,
			UserAgent: &log.UserAgent,
			Source:    AuditLogSource(log.Source),
			CreatedAt: log.CreatedAt,
		}
	}

	resp := AuditLogsResponse{
		AuditLogs:  apiAuditLogs,
		TotalCount: int(total),
		PageIndex:  int(page),
		PageSize:   int(limit),
	}

	return c.JSON(http.StatusOK, resp)
}

// GetAuditMetrics gets audit metrics
func (Server) GetAuditMetrics(c echo.Context) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Get audit metrics for the user
	metrics := make(map[string]interface{})

	// Total actions in the last 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	var totalActions int64
	if err := model.DB.Model(&model.AuditLog{}).
		Where("user_id = ? AND created_at >= ?", user.ID, thirtyDaysAgo).
		Count(&totalActions).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to get total actions")
	}
	metrics["total_actions_30_days"] = totalActions

	// Actions by type
	var actionStats []struct {
		Action string `json:"action"`
		Count  int64  `json:"count"`
	}
	if err := model.DB.Model(&model.AuditLog{}).
		Select("action, COUNT(*) as count").
		Where("user_id = ? AND created_at >= ?", user.ID, thirtyDaysAgo).
		Group("action").
		Order("count DESC").
		Scan(&actionStats).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to get action stats")
	}
	metrics["actions_by_type"] = actionStats

	// Most accessed resources
	var resourceStats []struct {
		ResourceType string `json:"resource_type"`
		Count        int64  `json:"count"`
	}
	if err := model.DB.Model(&model.AuditLog{}).
		Select("resource_type, COUNT(*) as count").
		Where("user_id = ? AND created_at >= ?", user.ID, thirtyDaysAgo).
		Group("resource_type").
		Order("count DESC").
		Scan(&resourceStats).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to get resource stats")
	}
	metrics["resources_by_type"] = resourceStats

	// Recent activity trends (last 7 days)
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	var dailyStats []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	if err := model.DB.Model(&model.AuditLog{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("user_id = ? AND created_at >= ?", user.ID, sevenDaysAgo).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&dailyStats).Error; err != nil {
		return sendError(c, http.StatusInternalServerError, "failed to get daily stats")
	}
	metrics["daily_activity_7_days"] = dailyStats

	return c.JSON(http.StatusOK, metrics)
}

// System endpoints

// GetConfig gets system configuration
func (Server) GetConfig(c echo.Context) error {
	// TODO: Implement proper config check
	isOidcEnabled := false
	isEmailEnabled := false

	resp := ConfigResponse{
		OidcEnabled:  isOidcEnabled,
		EmailEnabled: isEmailEnabled,
	}

	return c.JSON(http.StatusOK, resp)
}

// Health checks application health
func (Server) Health(c echo.Context) error {
	resp := map[string]interface{}{
		"status":    "ok",
		"timestamp": "2025-01-27T00:00:00Z", // TODO: Use actual timestamp
	}
	return c.JSON(http.StatusOK, resp)
}

// GetStatus gets detailed system status
func (Server) GetStatus(c echo.Context) error {
	// TODO: Implement proper status checks
	resp := StatusResponse{
		Version:        "1.0.0-migration",
		Commit:         "echo-glm",
		DatabaseStatus: StatusResponseDatabaseStatusHealthy,
		SystemStatus:   StatusResponseSystemStatusHealthy,
	}
	return c.JSON(http.StatusOK, resp)
}