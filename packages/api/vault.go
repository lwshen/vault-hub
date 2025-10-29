package api

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/model"
	"gorm.io/gorm"
)

// getClientInfo extracts IP address and User-Agent from the request.
func getClientInfo(c *fiber.Ctx) (string, string) {
	info := ExtractClientInfo(func(key string) string {
		return c.Get(key)
	}, c.IP())
	return info.IP, info.UserAgent
}

func getClientInfoDetails(c *fiber.Ctx) ClientInfo {
	return ExtractClientInfo(func(key string) string {
		return c.Get(key)
	}, c.IP())
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
		UpdatedAt:   &vault.UpdatedAt,
	}
}

// GetVaultsForUser handles listing vaults with pagination logic.
func GetVaultsForUser(user *model.User, params GetVaultsParams) (VaultsResponse, *APIError) {
	if user == nil {
		return VaultsResponse{}, newAPIError(http.StatusUnauthorized, "user not found in context")
	}

	pageSize := 20
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}
	pageIndex := 1
	if params.PageIndex != nil {
		pageIndex = *params.PageIndex
	}

	if pageSize < 1 || pageSize > 1000 {
		return VaultsResponse{}, newAPIError(http.StatusBadRequest, "pageSize must be between 1 and 1000")
	}
	if pageIndex < 1 {
		return VaultsResponse{}, newAPIError(http.StatusBadRequest, "pageIndex must be at least 1")
	}

	vaults, totalCount, err := model.GetUserVaultsWithPagination(user.ID, pageSize, pageIndex)
	if err != nil {
		return VaultsResponse{}, newAPIError(http.StatusInternalServerError, err.Error())
	}

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

	return response, nil
}

// GetVaultForUser retrieves a single vault and logs access.
func GetVaultForUser(user *model.User, uniqueID string, clientInfo ClientInfo) (Vault, *APIError) {
	if user == nil {
		return Vault{}, newAPIError(http.StatusUnauthorized, "user not found in context")
	}

	var vault model.Vault
	if err := vault.GetByUniqueID(uniqueID, user.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return Vault{}, newAPIError(http.StatusNotFound, "vault not found")
		}
		return Vault{}, newAPIError(http.StatusInternalServerError, err.Error())
	}

	if err := model.LogVaultAction(vault.ID, model.ActionReadVault, user.ID, model.SourceWeb, nil, clientInfo.IP, clientInfo.UserAgent); err != nil {
		slog.Error("Failed to create audit log for read vault", "error", err, "vaultID", vault.ID)
	}

	return convertToApiVault(&vault), nil
}

// CreateVaultForUser validates and persists a new vault.
func CreateVaultForUser(user *model.User, input CreateVaultRequest, clientInfo ClientInfo) (Vault, *APIError) {
	if user == nil {
		return Vault{}, newAPIError(http.StatusUnauthorized, "user not found in context")
	}

	uniqueID, err := uuid.NewV7()
	if err != nil {
		return Vault{}, newAPIError(http.StatusInternalServerError, err.Error())
	}

	params := model.CreateVaultParams{
		UniqueID:    uniqueID.String(),
		UserID:      user.ID,
		Name:        input.Name,
		Value:       input.Value,
		Description: getStringValue(input.Description),
		Category:    getStringValue(input.Category),
	}

	if errors := params.Validate(); len(errors) > 0 {
		var errorMsgs []string
		for _, msg := range errors {
			errorMsgs = append(errorMsgs, msg)
		}
		return Vault{}, newAPIError(http.StatusBadRequest, strings.Join(errorMsgs, "; "))
	}

	vault, err := params.Create()
	if err != nil {
		return Vault{}, newAPIError(http.StatusInternalServerError, err.Error())
	}

	if err := model.LogVaultAction(vault.ID, model.ActionCreateVault, user.ID, model.SourceWeb, nil, clientInfo.IP, clientInfo.UserAgent); err != nil {
		slog.Error("Failed to create audit log for create vault", "error", err, "vaultID", vault.ID)
	}

	return convertToApiVault(vault), nil
}

// UpdateVaultForUser applies updates to an existing vault.
func UpdateVaultForUser(user *model.User, uniqueID string, input UpdateVaultRequest, clientInfo ClientInfo) (Vault, *APIError) {
	if user == nil {
		return Vault{}, newAPIError(http.StatusUnauthorized, "user not found in context")
	}

	var vault model.Vault
	if err := vault.GetByUniqueID(uniqueID, user.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return Vault{}, newAPIError(http.StatusNotFound, "vault not found")
		}
		return Vault{}, newAPIError(http.StatusInternalServerError, err.Error())
	}

	params := model.UpdateVaultParams{
		Name:        input.Name,
		Value:       input.Value,
		Description: input.Description,
		Category:    input.Category,
	}

	if errors := params.Validate(); len(errors) > 0 {
		var errorMsgs []string
		for _, msg := range errors {
			errorMsgs = append(errorMsgs, msg)
		}
		return Vault{}, newAPIError(http.StatusBadRequest, strings.Join(errorMsgs, "; "))
	}

	if err := vault.Update(&params); err != nil {
		return Vault{}, newAPIError(http.StatusInternalServerError, err.Error())
	}

	if err := model.LogVaultAction(vault.ID, model.ActionUpdateVault, user.ID, model.SourceWeb, nil, clientInfo.IP, clientInfo.UserAgent); err != nil {
		slog.Error("Failed to create audit log for update vault", "error", err, "vaultID", vault.ID)
	}

	return convertToApiVault(&vault), nil
}

// DeleteVaultForUser removes a vault owned by the authenticated user.
func DeleteVaultForUser(user *model.User, uniqueID string, clientInfo ClientInfo) *APIError {
	if user == nil {
		return newAPIError(http.StatusUnauthorized, "user not found in context")
	}

	var vault model.Vault
	if err := vault.GetByUniqueID(uniqueID, user.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return newAPIError(http.StatusNotFound, "vault not found")
		}
		return newAPIError(http.StatusInternalServerError, err.Error())
	}

	if err := vault.Delete(); err != nil {
		return newAPIError(http.StatusInternalServerError, err.Error())
	}

	if err := model.LogVaultAction(vault.ID, model.ActionDeleteVault, user.ID, model.SourceWeb, nil, clientInfo.IP, clientInfo.UserAgent); err != nil {
		slog.Error("Failed to create audit log for delete vault", "error", err, "vaultID", vault.ID)
	}

	return nil
}

// GetVaults handles GET /api/vaults with pagination
func (Server) GetVaults(c *fiber.Ctx, params GetVaultsParams) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	response, apiErr := GetVaultsForUser(user, params)
	if apiErr != nil {
		return handler.SendError(c, apiErr.Status, apiErr.Message)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetVault handles GET /api/vaults/{unique_id}
func (Server) GetVault(c *fiber.Ctx, uniqueID string) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	resp, apiErr := GetVaultForUser(user, uniqueID, getClientInfoDetails(c))
	if apiErr != nil {
		return handler.SendError(c, apiErr.Status, apiErr.Message)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
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

	resp, apiErr := CreateVaultForUser(user, input, getClientInfoDetails(c))
	if apiErr != nil {
		return handler.SendError(c, apiErr.Status, apiErr.Message)
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
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

	resp, apiErr := UpdateVaultForUser(user, uniqueID, input, getClientInfoDetails(c))
	if apiErr != nil {
		return handler.SendError(c, apiErr.Status, apiErr.Message)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// DeleteVault handles DELETE /api/vaults/{unique_id}
func (Server) DeleteVault(c *fiber.Ctx, uniqueID string) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	if apiErr := DeleteVaultForUser(user, uniqueID, getClientInfoDetails(c)); apiErr != nil {
		return handler.SendError(c, apiErr.Status, apiErr.Message)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// getStringValue safely gets string value from pointer, returns empty string if nil
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
