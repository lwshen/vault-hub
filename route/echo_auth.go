package route

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/internal/config"
	"github.com/lwshen/vault-hub/model"
)

// isPublicRoute checks if a route is public and doesn't need authentication
func isPublicRoute(path string) bool {
	publicRoutes := []string{
		"/api/health",
		"/api/version",
		"/api/config",
		"/api/auth/login",
		"/api/auth/signup",
		"/api/auth/login/oidc",
		"/api/auth/callback/oidc",
		"/api/auth/password/reset/request",
		"/api/auth/password/reset/confirm",
		"/api/auth/magic-link/request",
		"/api/auth/magic-link/token",
	}

	for _, route := range publicRoutes {
		if strings.HasPrefix(path, route) {
			return true
		}
	}
	return false
}

// EchoJWTMiddleware handles JWT and API key authentication for Echo
func EchoJWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path

			// Public routes that don't need authentication
			if isPublicRoute(path) {
				return next(c)
			}

			// Routes starting with /api/cli/ MUST use API key authentication
			if strings.HasPrefix(path, "/api/cli/") {
				return echoAPIKeyOnlyMiddleware(c, next)
			}

			// All other /api/ routes require JWT authentication
			if strings.HasPrefix(path, "/api/") {
				return echoJWTOnlyMiddleware(c, next)
			}

			// Non-API routes (web assets, etc.) don't need auth
			return next(c)
		}
	}
}

// echoJWTOnlyMiddleware ensures non-API-key routes only accept JWT authentication
func echoJWTOnlyMiddleware(c echo.Context, next echo.HandlerFunc) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "JWT token required")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header")
	}

	tokenString := tokenParts[1]

	// Reject API keys on JWT-only routes
	if strings.HasPrefix(tokenString, "vhub_") {
		return echo.NewHTTPError(http.StatusUnauthorized, "JWT token required for this endpoint")
	}

	return handleEchoJWTAuth(c, next, tokenString)
}

// echoAPIKeyOnlyMiddleware ensures API key routes only accept API key authentication
func echoAPIKeyOnlyMiddleware(c echo.Context, next echo.HandlerFunc) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "API key required for this endpoint")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header")
	}

	tokenString := tokenParts[1]

	// Ensure it's an API key (starts with "vhub_")
	if !strings.HasPrefix(tokenString, "vhub_") {
		return echo.NewHTTPError(http.StatusUnauthorized, "API key required for this endpoint")
	}

	return handleEchoAPIKeyAuth(c, next, tokenString)
}

// handleEchoAPIKeyAuth validates API key and sets context
func handleEchoAPIKeyAuth(c echo.Context, next echo.HandlerFunc, apiKey string) error {
	key, err := model.ValidateAPIKey(apiKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid API key")
	}

	// Set context variables for handlers
	c.Set("user_id", &key.UserID)
	c.Set("api_key", key)

	return next(c)
}

// handleEchoJWTAuth validates JWT token and sets context
func handleEchoJWTAuth(c echo.Context, next echo.HandlerFunc, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
	}

	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user ID in token")
	}

	userID := uint(userIDFloat)

	var user model.User
	if err := model.DB.First(&user, userID).Error; err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
	}

	// Clear password for security
	user.Password = nil

	// Set user in context
	c.Set("user", &user)

	return next(c)
}

// Helper function to get current user from Echo context
func GetCurrentUserFromEcho(c echo.Context) (*model.User, error) {
	user, ok := c.Get("user").(*model.User)
	if !ok || user == nil {
		return nil, fmt.Errorf("user not found in context")
	}
	return user, nil
}

// Helper function to get user ID from Echo context
func GetUserIDFromEcho(c echo.Context) (*uint, error) {
	if userID, ok := c.Get("user_id").(*uint); ok && userID != nil {
		return userID, nil
	}
	return nil, fmt.Errorf("user_id not found in context")
}

// Helper function to get API key from Echo context
func GetAPIKeyFromEcho(c echo.Context) (*model.APIKey, error) {
	apiKey, ok := c.Get("api_key").(*model.APIKey)
	if !ok || apiKey == nil {
		return nil, fmt.Errorf("api_key not found in context")
	}
	return apiKey, nil
}
