package api

import (
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/model"
	"gorm.io/gorm"
)

// getClientInfo extracts IP address and User-Agent from the request
func getClientInfo(c *fiber.Ctx) (string, string) {
	// Get IP address (check for forwarded headers first)
	ip := c.Get("X-Forwarded-For")
	if ip == "" {
		ip = c.Get("X-Real-IP")
	}
	if ip == "" {
		ip = c.IP()
	}

	// Get User-Agent
	userAgent := c.Get("User-Agent")

	return ip, userAgent
}

// getUserFromContext extracts the authenticated user from the context
func getUserFromContext(c *fiber.Ctx) (*model.User, error) {
	user, ok := c.Locals("user").(*model.User)
	if !ok {
		return nil, handler.SendError(c, fiber.StatusUnauthorized, "user not found in context")
	}
	return user, nil
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
		Favourite:   &vault.Favourite,
		CreatedAt:   &vault.CreatedAt,
		UpdatedAt:   &vault.UpdatedAt,
	}
}

// convertToApiVaultLite converts a model.Vault to an api.VaultLite
func convertToApiVaultLite(vault *model.Vault) VaultLite {
	return VaultLite{
		UniqueId:    vault.UniqueID,
		Name:        vault.Name,
		Description: &vault.Description,
		Category:    &vault.Category,
		Favourite:   &vault.Favourite,
		UpdatedAt:   &vault.UpdatedAt,
	}
}

// GetVaults handles GET /api/vaults with pagination
func (Server) GetVaults(c *fiber.Ctx, params GetVaultsParams) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Apply defaults if not provided
	pageSize := 20
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}
	pageIndex := 1
	if params.PageIndex != nil {
		pageIndex = *params.PageIndex
	}

	// Validate bounds
	if pageSize < 1 || pageSize > 1000 {
		return handler.SendError(c, fiber.StatusBadRequest, "pageSize must be between 1 and 1000")
	}
	if pageIndex < 1 {
		return handler.SendError(c, fiber.StatusBadRequest, "pageIndex must be at least 1")
	}

	// Query paginated vaults for current user via model
	vaults, totalCount, err := model.GetUserVaultsWithPagination(user.ID, pageSize, pageIndex)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Convert to API VaultLite slice
	apiVaults := make([]VaultLite, 0, len(vaults))
	for i := range vaults {
		apiVaults = append(apiVaults, convertToApiVaultLite(&vaults[i]))
	}

	response := VaultsResponse{
		Vaults:     apiVaults,
		TotalCount: int(totalCount),
		PageSize:   pageSize,
		PageIndex:  pageIndex,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetVault handles GET /api/vaults/{unique_id}
func (Server) GetVault(c *fiber.Ctx, uniqueID string) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var vault model.Vault
	err = vault.GetByUniqueID(uniqueID, user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return handler.SendError(c, fiber.StatusNotFound, "vault not found")
		}
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Log read action
	ip, userAgent := getClientInfo(c)
	if err := model.LogVaultAction(vault.ID, model.ActionReadVault, user.ID, model.SourceWeb, nil, ip, userAgent); err != nil {
		slog.Error("Failed to create audit log for read vault", "error", err, "vaultID", vault.ID)
	}

	return c.Status(fiber.StatusOK).JSON(convertToApiVault(&vault))
}

// CreateVault handles POST /api/vaults
func (Server) CreateVault(c *fiber.Ctx) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var input CreateVaultRequest
	if err := c.BodyParser(&input); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Create parameters
	uniqueID, err := uuid.NewV7()
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	params := model.CreateVaultParams{
		UniqueID:    uniqueID.String(),
		UserID:      user.ID,
		Name:        input.Name,
		Value:       input.Value,
		Description: getStringValue(input.Description),
		Category:    getStringValue(input.Category),
	}

	// Validate parameters
	errors := params.Validate()
	if len(errors) > 0 {
		var errorMsgs []string
		for _, msg := range errors {
			errorMsgs = append(errorMsgs, msg)
		}
		return handler.SendError(c, fiber.StatusBadRequest, strings.Join(errorMsgs, "; "))
	}

	// Create vault
	vault, err := params.Create()
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Log create action
	ip, userAgent := getClientInfo(c)
	if err := model.LogVaultAction(vault.ID, model.ActionCreateVault, user.ID, model.SourceWeb, nil, ip, userAgent); err != nil {
		slog.Error("Failed to create audit log for create vault", "error", err, "vaultID", vault.ID)
	}

	return c.Status(fiber.StatusCreated).JSON(convertToApiVault(vault))
}

// UpdateVault handles PUT /api/vaults/{unique_id}
func (Server) UpdateVault(c *fiber.Ctx, uniqueID string) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var vault model.Vault
	err = vault.GetByUniqueID(uniqueID, user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return handler.SendError(c, fiber.StatusNotFound, "vault not found")
		}
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	var input UpdateVaultRequest
	if err := c.BodyParser(&input); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Create update parameters
	params := model.UpdateVaultParams{
		Name:        input.Name,
		Value:       input.Value,
		Description: input.Description,
		Category:    input.Category,
	}

	// Validate parameters
	errors := params.Validate()
	if len(errors) > 0 {
		var errorMsgs []string
		for _, msg := range errors {
			errorMsgs = append(errorMsgs, msg)
		}
		return handler.SendError(c, fiber.StatusBadRequest, strings.Join(errorMsgs, "; "))
	}

	// Update vault
	err = vault.Update(&params)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Log update action
	ip, userAgent := getClientInfo(c)
	_ = model.LogVaultAction(vault.ID, model.ActionUpdateVault, user.ID, model.SourceWeb, nil, ip, userAgent)

	return c.Status(fiber.StatusOK).JSON(convertToApiVault(&vault))
}

// DeleteVault handles DELETE /api/vaults/{unique_id}
func (Server) DeleteVault(c *fiber.Ctx, uniqueID string) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var vault model.Vault
	err = vault.GetByUniqueID(uniqueID, user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return handler.SendError(c, fiber.StatusNotFound, "vault not found")
		}
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Delete vault
	err = vault.Delete()
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Log delete action
	ip, userAgent := getClientInfo(c)
	_ = model.LogVaultAction(vault.ID, model.ActionDeleteVault, user.ID, model.SourceWeb, nil, ip, userAgent)

	return c.SendStatus(fiber.StatusNoContent)
}

// getStringValue safely gets string value from pointer, returns empty string if nil
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// convertToApiVaultFilterOption converts a model.Vault to an api.VaultFilterOption
func convertToApiVaultFilterOption(vault *model.Vault) VaultFilterOption {
	return VaultFilterOption{
		UniqueId: vault.UniqueID,
		Name:     vault.Name,
	}
}

// GetVaultFilterOptions handles GET /api/vaults/filter-options
// Returns a minimal list of vaults (uniqueId and name only) for filter dropdowns
func (Server) GetVaultFilterOptions(c *fiber.Ctx) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Query all vaults for the user, selecting only needed fields
	var vaults []model.Vault
	err = model.DB.Where("user_id = ? AND deleted_at IS NULL", user.ID).
		Select("unique_id", "name").
		Order("name ASC").
		Find(&vaults).Error
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Convert to API VaultFilterOption slice
	apiVaults := make([]VaultFilterOption, 0, len(vaults))
	for i := range vaults {
		apiVaults = append(apiVaults, convertToApiVaultFilterOption(&vaults[i]))
	}

	response := VaultFilterOptionsResponse{
		Vaults: apiVaults,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
