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

// JWTMiddlewareEcho creates an Echo middleware for JWT and API key authentication
func JWTMiddlewareEcho() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			path := ctx.Path()

			// Public routes that don't need authentication
			if isPublicRouteEcho(path) {
				return next(ctx)
			}

			// Routes starting with /api/cli/ MUST use API key authentication
			if strings.HasPrefix(path, "/api/cli/") {
				return apiKeyOnlyMiddlewareEcho(ctx, next)
			}

			// All other /api/ routes require JWT authentication
			if strings.HasPrefix(path, "/api/") {
				return jwtOnlyMiddlewareEcho(ctx, next)
			}

			// Non-API routes (web assets, etc.) don't need auth
			return next(ctx)
		}
	}
}

// isPublicRouteEcho checks if a route is public and doesn't need authentication
func isPublicRouteEcho(path string) bool {
	publicRoutes := []string{
		"/api/health",
		"/api/status",
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

// jwtOnlyMiddlewareEcho ensures non-API-key routes only accept JWT authentication
func jwtOnlyMiddlewareEcho(ctx echo.Context, next echo.HandlerFunc) error {
	authHeader := ctx.Request().Header.Get("Authorization")
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

	return handleJWTAuthEcho(ctx, next, tokenString)
}

// apiKeyOnlyMiddlewareEcho ensures API key routes only accept API key authentication
func apiKeyOnlyMiddlewareEcho(ctx echo.Context, next echo.HandlerFunc) error {
	authHeader := ctx.Request().Header.Get("Authorization")
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

	return handleAPIKeyAuthEcho(ctx, next, tokenString)
}

// handleAPIKeyAuthEcho validates API key and sets context
func handleAPIKeyAuthEcho(ctx echo.Context, next echo.HandlerFunc, apiKey string) error {
	key, err := model.ValidateAPIKey(apiKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid API key")
	}

	ctx.Set("user_id", &key.UserID)
	ctx.Set("api_key", key)

	return next(ctx)
}

// handleJWTAuthEcho validates JWT and sets user in context
func handleJWTAuthEcho(ctx echo.Context, next echo.HandlerFunc, tokenString string) error {
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

	userID, ok := claims["sub"].(float64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user ID in token")
	}

	var user model.User
	if err := model.DB.First(&user, uint(userID)).Error; err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
	}

	user.Password = nil
	ctx.Set("user", &user)

	return next(ctx)
}
