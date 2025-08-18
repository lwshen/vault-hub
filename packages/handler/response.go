package handler

import "github.com/gofiber/fiber/v2"

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SendError sends a standardized error response
func SendError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{
		"error": ErrorResponse{
			Code:    code,
			Message: message,
		},
	})
}
