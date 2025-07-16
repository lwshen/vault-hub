package api

import (
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

// GetConfigurations handles GET /api/configurations
func (Server) GetConfigurations(c *fiber.Ctx, params GetConfigurationsParams) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	category := getStringValue(params.Category)

	configs, err := model.GetConfigurationsByUser(user.ID, category)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Log read action for each configuration
	ip, userAgent := getClientInfo(c)
	for _, config := range configs {
		_ = model.LogConfigurationAction(config.ID, model.ActionReadConfig, user.ID, ip, userAgent)
	}

	return c.Status(fiber.StatusOK).JSON(configs)
}

// GetConfiguration handles GET /api/configurations/{id}
func (Server) GetConfiguration(c *fiber.Ctx, id int64) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var config model.Configuration
	err = config.GetByID(uint(id), user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return handler.SendError(c, fiber.StatusNotFound, "configuration not found")
		}
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Log read action
	ip, userAgent := getClientInfo(c)
	_ = model.LogConfigurationAction(config.ID, model.ActionReadConfig, user.ID, ip, userAgent)

	return c.Status(fiber.StatusOK).JSON(config)
}

// CreateConfiguration handles POST /api/configurations
func (Server) CreateConfiguration(c *fiber.Ctx) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var input CreateConfigurationRequest
	if err := c.BodyParser(&input); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Create parameters
	params := model.CreateConfigurationParams{
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

	// Create configuration
	config, err := params.Create()
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Log create action
	ip, userAgent := getClientInfo(c)
	_ = model.LogConfigurationAction(config.ID, model.ActionCreateConfig, user.ID, ip, userAgent)

	return c.Status(fiber.StatusCreated).JSON(config)
}

// UpdateConfiguration handles PUT /api/configurations/{id}
func (Server) UpdateConfiguration(c *fiber.Ctx, id int64) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var config model.Configuration
	err = config.GetByID(uint(id), user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return handler.SendError(c, fiber.StatusNotFound, "configuration not found")
		}
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	var input UpdateConfigurationRequest
	if err := c.BodyParser(&input); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Create update parameters
	params := model.UpdateConfigurationParams{
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

	// Update configuration
	err = config.Update(&params)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Log update action
	ip, userAgent := getClientInfo(c)
	_ = model.LogConfigurationAction(config.ID, model.ActionUpdateConfig, user.ID, ip, userAgent)

	return c.Status(fiber.StatusOK).JSON(config)
}

// DeleteConfiguration handles DELETE /api/configurations/{id}
func (Server) DeleteConfiguration(c *fiber.Ctx, id int64) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var config model.Configuration
	err = config.GetByID(uint(id), user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return handler.SendError(c, fiber.StatusNotFound, "configuration not found")
		}
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Delete configuration
	err = config.Delete()
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Log delete action
	ip, userAgent := getClientInfo(c)
	_ = model.LogConfigurationAction(config.ID, model.ActionDeleteConfig, user.ID, ip, userAgent)

	return c.SendStatus(fiber.StatusNoContent)
}

// getStringValue safely gets string value from pointer, returns empty string if nil
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
