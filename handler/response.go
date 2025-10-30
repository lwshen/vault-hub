package handler

import "github.com/labstack/echo/v4"

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SendError sends a standardized error response
func SendError(c echo.Context, code int, message string) error {
	return c.JSON(code, map[string]interface{}{
		"error": ErrorResponse{
			Code:    code,
			Message: message,
		},
	})
}
