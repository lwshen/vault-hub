package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/model"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// BuildCurrentUserResponse converts the authenticated user to API shape.
func BuildCurrentUserResponse(user *model.User) (GetUserResponse, *APIError) {
	if user == nil {
		return GetUserResponse{}, newAPIError(http.StatusUnauthorized, "user not found in context")
	}

	resp := GetUserResponse{
		Email:  openapi_types.Email(user.Email),
		Avatar: user.Avatar,
		Name:   user.Name,
	}
	return resp, nil
}

func (Server) GetCurrentUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*model.User)
	if !ok {
		return handler.SendError(c, fiber.StatusUnauthorized, "user not found in context")
	}

	resp, apiErr := BuildCurrentUserResponse(user)
	if apiErr != nil {
		return handler.SendError(c, apiErr.Status, apiErr.Message)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
