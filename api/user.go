package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/model"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (Server) GetCurrentUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*model.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found in context"})
	}

	resp := GetUserResponse{
		Email:  openapi_types.Email(user.Email),
		Avatar: user.Avatar,
		Name:   user.Name,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
