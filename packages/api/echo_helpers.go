package api

import (
	"math"
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

// safeInt64ToInt32 safely converts int64 to int32, capping at int32 max value
func safeInt64ToInt32(val int64) int32 {
	if val > math.MaxInt32 {
		return math.MaxInt32
	}
	if val < math.MinInt32 {
		return math.MinInt32
	}
	return int32(val)
}

// getAPIKeyFromContext extracts the API key from the Echo context
func getAPIKeyFromEchoContext(ctx echo.Context) (*model.APIKey, error) {
	apiKey, ok := ctx.Get("api_key").(*model.APIKey)
	if !ok {
		return nil, SendError(ctx, http.StatusUnauthorized, "api_key not found in context")
	}
	return apiKey, nil
}
