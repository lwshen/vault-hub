package api

import (
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
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

// GetVaults handles GET /api/vaults
func (Server) GetVaults(c *fiber.Ctx, params GetVaultsParams) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	category := getStringValue(params.Category)

	vaults, err := model.GetVaultsByUser(user.ID, category)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Log read action for each vault
	ip, userAgent := getClientInfo(c)
	for _, vault := range vaults {
		if err := model.LogVaultAction(vault.ID, model.ActionReadVault, user.ID, ip, userAgent); err != nil {
			slog.Error("Failed to create audit log for read vault", "error", err, "vaultID", vault.ID)
		}
	}

	return c.Status(fiber.StatusOK).JSON(vaults)
}

// GetVault handles GET /api/vaults/{id}
func (Server) GetVault(c *fiber.Ctx, id int32) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var vault model.Vault
	err = vault.GetByID(uint(id), user.ID) // #nosec G115
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return handler.SendError(c, fiber.StatusNotFound, "vault not found")
		}
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Log read action
	ip, userAgent := getClientInfo(c)
	if err := model.LogVaultAction(vault.ID, model.ActionReadVault, user.ID, ip, userAgent); err != nil {
		slog.Error("Failed to create audit log for read vault", "error", err, "vaultID", vault.ID)
	}

	return c.Status(fiber.StatusOK).JSON(vault)
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
	params := model.CreateVaultParams{
		UniqueID:    input.UniqueId,
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
	if err := model.LogVaultAction(vault.ID, model.ActionCreateVault, user.ID, ip, userAgent); err != nil {
		slog.Error("Failed to create audit log for create vault", "error", err, "vaultID", vault.ID)
	}

	return c.Status(fiber.StatusCreated).JSON(vault)
}

// UpdateVault handles PUT /api/vaults/{id}
func (Server) UpdateVault(c *fiber.Ctx, id int32) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var vault model.Vault
	err = vault.GetByID(uint(id), user.ID) // #nosec G115
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
	_ = model.LogVaultAction(vault.ID, model.ActionUpdateVault, user.ID, ip, userAgent)

	return c.Status(fiber.StatusOK).JSON(vault)
}

// DeleteVault handles DELETE /api/vaults/{id}
func (Server) DeleteVault(c *fiber.Ctx, id int32) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var vault model.Vault
	err = vault.GetByID(uint(id), user.ID) // #nosec G115
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
	_ = model.LogVaultAction(vault.ID, model.ActionDeleteVault, user.ID, ip, userAgent)

	return c.SendStatus(fiber.StatusNoContent)
}

// getStringValue safely gets string value from pointer, returns empty string if nil
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
