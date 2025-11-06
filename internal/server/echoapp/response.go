package echoapp

import "github.com/labstack/echo/v4"

// sendError mirrors the Fiber SendError helper while targeting Echo contexts.
func sendError(c echo.Context, status int, message string) error {
	return c.JSON(status, map[string]any{
		"error": map[string]any{
			"code":    status,
			"message": message,
		},
	})
}
