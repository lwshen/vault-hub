package echoapp

import (
	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/handler"
)

// sendError mirrors the Fiber SendError helper while targeting Echo contexts.
func sendError(c echo.Context, status int, message string) error {
	return c.JSON(status, map[string]any{
		"error": handler.ErrorResponse{
			Code:    status,
			Message: message,
		},
	})
}
