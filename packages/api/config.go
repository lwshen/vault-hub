package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/internal/config"
)

// GetConfig returns public configuration that requires no authentication
// This endpoint performs NO database operations and is safe for public access
func (s Server) GetConfig(ctx *fiber.Ctx) error {
	resp := ConfigResponse{
		OidcEnabled:  config.OidcEnabled,
		EmailEnabled: config.EmailEnabled,
	}

	return ctx.
		Status(http.StatusOK).
		JSON(resp)
}
