package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/model"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SendError sends a standardized error response for Echo
func SendError(ctx echo.Context, code int, message string) error {
	return ctx.JSON(code, map[string]interface{}{
		"error": ErrorResponse{
			Code:    code,
			Message: message,
		},
	})
}

// getClientInfo extracts IP address and User-Agent from the Echo request
func getClientInfoEcho(ctx echo.Context) (string, string) {
	// Get IP address (check for forwarded headers first)
	ip := ctx.Request().Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = ctx.Request().Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = ctx.RealIP()
	}

	// Get User-Agent
	userAgent := ctx.Request().Header.Get("User-Agent")

	return ip, userAgent
}

// getUserFromContext extracts the authenticated user from the Echo context
func getUserFromEchoContext(ctx echo.Context) (*model.User, error) {
	user, ok := ctx.Get("user").(*model.User)
	if !ok {
		return nil, SendError(ctx, http.StatusUnauthorized, "user not found in context")
	}
	return user, nil
}

// getUserIDFromContext extracts the user ID from the Echo context (for API key auth)
func getUserIDFromEchoContext(ctx echo.Context) (*uint, error) {
	userID, ok := ctx.Get("user_id").(*uint)
	if !ok {
		return nil, SendError(ctx, http.StatusUnauthorized, "user_id not found in context")
	}
	return userID, nil
}

// getAPIKeyFromContext extracts the API key from the Echo context
func getAPIKeyFromEchoContext(ctx echo.Context) (*model.APIKey, error) {
	apiKey, ok := ctx.Get("api_key").(*model.APIKey)
	if !ok {
		return nil, SendError(ctx, http.StatusUnauthorized, "api_key not found in context")
	}
	return apiKey, nil
}
