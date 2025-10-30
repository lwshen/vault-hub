package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/internal/config"
)

// GetConfig returns public configuration that requires no authentication
// This endpoint performs NO database operations and is safe for public access
func (s Server) GetConfig(ctx echo.Context) error {
	resp := ConfigResponse{
		OidcEnabled:  config.OidcEnabled,
		EmailEnabled: config.EmailEnabled,
	}

	return ctx.JSON(http.StatusOK, resp)
}
