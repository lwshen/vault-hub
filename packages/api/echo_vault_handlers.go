package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/model"
	"github.com/lwshen/vault-hub/packages/api/generated/models"
	"gorm.io/gorm"
)

// GetVaults handles GET /api/vaults with pagination
func (c *Container) GetVaults(ctx echo.Context) error {
	user, err := getUserFromEchoContext(ctx)
	if err != nil {
		return err
	}

	// Parse query parameters
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

	// Validate bounds
	if pageSize < 1 || pageSize > 1000 {
		return SendError(ctx, http.StatusBadRequest, "pageSize must be between 1 and 1000")
	}
	if pageIndex < 1 {
		return SendError(ctx, http.StatusBadRequest, "pageIndex must be at least 1")
	}

	// Query paginated vaults for current user via model
	vaults, totalCount, err := model.GetUserVaultsWithPagination(user.ID, pageSize, pageIndex)
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, err.Error())
	}

	// Convert to API VaultLite slice
	apiVaults := make([]models.VaultLite, 0, len(vaults))
	for i := range vaults {
		apiVaults = append(apiVaults, convertToGeneratedVaultLite(&vaults[i]))
	}

	// Safe conversion with bounds checking (these values are already validated)
	// pageSize is max 1000, pageIndex is validated >= 1, totalCount from DB
	response := models.VaultsResponse{
		Vaults:     apiVaults,
		TotalCount: safeInt64ToInt32(totalCount),
		PageSize:   int32(pageSize),  // #nosec G115 -- validated max 1000
		PageIndex:  int32(pageIndex), // #nosec G115 -- validated >= 1
	}

	return ctx.JSON(http.StatusOK, response)
}

// GetVault handles GET /api/vaults/{uniqueId}
func (c *Container) GetVault(ctx echo.Context) error {
	user, err := getUserFromEchoContext(ctx)
	if err != nil {
		return err
	}

	uniqueID := ctx.Param("uniqueId")

	var vault model.Vault
	err = vault.GetByUniqueID(uniqueID, user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return SendError(ctx, http.StatusNotFound, "vault not found")
		}
		return SendError(ctx, http.StatusInternalServerError, err.Error())
	}

	// Log read action
	ip, userAgent := getClientInfoEcho(ctx)
	if err := model.LogVaultAction(vault.ID, model.ActionReadVault, user.ID, model.SourceWeb, nil, ip, userAgent); err != nil {
		slog.Error("Failed to create audit log for read vault", "error", err, "vaultID", vault.ID)
	}

	return ctx.JSON(http.StatusOK, convertToGeneratedVault(&vault))
}

// CreateVault handles POST /api/vaults
func (c *Container) CreateVault(ctx echo.Context) error {
	user, err := getUserFromEchoContext(ctx)
	if err != nil {
		return err
	}

	var input models.CreateVaultRequest
	if err := ctx.Bind(&input); err != nil {
		return SendError(ctx, http.StatusBadRequest, err.Error())
	}

	// Validate required fields
	if input.Name == "" {
		return SendError(ctx, http.StatusBadRequest, "name is required")
	}
	if input.Value == "" {
		return SendError(ctx, http.StatusBadRequest, "value is required")
	}

	// Create vault parameters
	params := model.CreateVaultParams{
		UniqueID:    uuid.New().String(),
		UserID:      user.ID,
		Name:        input.Name,
		Value:       input.Value,
		Description: input.Description,
		Category:    input.Category,
	}

	vault, err := params.Create()
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, err.Error())
	}

	// Log create action
	ip, userAgent := getClientInfoEcho(ctx)
	if err := model.LogVaultAction(vault.ID, model.ActionCreateVault, user.ID, model.SourceWeb, nil, ip, userAgent); err != nil {
		slog.Error("Failed to create audit log for create vault", "error", err, "vaultID", vault.ID)
	}

	return ctx.JSON(http.StatusOK, convertToGeneratedVault(vault))
}

// UpdateVault handles PUT /api/vaults/{uniqueId}
func (c *Container) UpdateVault(ctx echo.Context) error {
	user, err := getUserFromEchoContext(ctx)
	if err != nil {
		return err
	}

	uniqueID := ctx.Param("uniqueId")

	var vault model.Vault
	err = vault.GetByUniqueID(uniqueID, user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return SendError(ctx, http.StatusNotFound, "vault not found")
		}
		return SendError(ctx, http.StatusInternalServerError, err.Error())
	}

	var input models.UpdateVaultRequest
	if err := ctx.Bind(&input); err != nil {
		return SendError(ctx, http.StatusBadRequest, err.Error())
	}

	// Create update parameters - only set fields that are provided
	params := model.UpdateVaultParams{}

	if input.Name != "" {
		params.Name = &input.Name
	}

	if input.Value != "" {
		params.Value = &input.Value
	}

	if input.Description != "" {
		params.Description = &input.Description
	}

	if input.Category != "" {
		params.Category = &input.Category
	}

	// Validate parameters
	errors := params.Validate()
	if len(errors) > 0 {
		var errorMsgs []string
		for _, msg := range errors {
			errorMsgs = append(errorMsgs, msg)
		}
		return SendError(ctx, http.StatusBadRequest, strings.Join(errorMsgs, "; "))
	}

	// Update vault
	err = vault.Update(&params)
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, err.Error())
	}

	// Log update action
	ip, userAgent := getClientInfoEcho(ctx)
	if err := model.LogVaultAction(vault.ID, model.ActionUpdateVault, user.ID, model.SourceWeb, nil, ip, userAgent); err != nil {
		slog.Error("Failed to create audit log for update vault", "error", err, "vaultID", vault.ID)
	}

	return ctx.JSON(http.StatusOK, convertToGeneratedVault(&vault))
}

// DeleteVault handles DELETE /api/vaults/{uniqueId}
func (c *Container) DeleteVault(ctx echo.Context) error {
	user, err := getUserFromEchoContext(ctx)
	if err != nil {
		return err
	}

	uniqueID := ctx.Param("uniqueId")

	var vault model.Vault
	err = vault.GetByUniqueID(uniqueID, user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return SendError(ctx, http.StatusNotFound, "vault not found")
		}
		return SendError(ctx, http.StatusInternalServerError, err.Error())
	}

	// Delete vault (soft delete)
	err = vault.Delete()
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, err.Error())
	}

	// Log delete action
	ip, userAgent := getClientInfoEcho(ctx)
	if err := model.LogVaultAction(vault.ID, model.ActionDeleteVault, user.ID, model.SourceWeb, nil, ip, userAgent); err != nil {
		slog.Error("Failed to create audit log for delete vault", "error", err, "vaultID", vault.ID)
	}

	return ctx.NoContent(http.StatusNoContent)
}
