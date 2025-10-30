package route

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/internal/config"
	"github.com/lwshen/vault-hub/model"
)

// SlogMiddleware creates an Echo middleware that logs requests using slog
func SlogMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			logger.Info("request",
				"method", req.Method,
				"path", req.URL.Path,
				"status", res.Status,
				"ip", c.RealIP(),
			)

			return err
		}
	}
}

func jwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Path()

		// Public routes that don't need authentication
		if isPublicRoute(path) {
			return next(c)
		}

		// Routes starting with /api/cli/ MUST use API key authentication
		if strings.HasPrefix(path, "/api/cli/") {
			return apiKeyOnlyMiddleware(c, next)
		}

		// All other /api/ routes require JWT authentication
		if strings.HasPrefix(path, "/api/") {
			return jwtOnlyMiddleware(c, next)
		}

		// Non-API routes (web assets, etc.) don't need auth
		return next(c)
	}
}

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

// jwtOnlyMiddleware ensures non-API-key routes only accept JWT authentication
func jwtOnlyMiddleware(c echo.Context, next echo.HandlerFunc) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return handler.SendError(c, http.StatusUnauthorized, "JWT token required")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return handler.SendError(c, http.StatusUnauthorized, "invalid authorization header")
	}

	tokenString := tokenParts[1]

	// Reject API keys on JWT-only routes
	if strings.HasPrefix(tokenString, "vhub_") {
		return handler.SendError(c, http.StatusUnauthorized, "JWT token required for this endpoint")
	}

	return handleJWTAuth(c, tokenString, next)
}

// apiKeyOnlyMiddleware ensures API key routes only accept API key authentication
func apiKeyOnlyMiddleware(c echo.Context, next echo.HandlerFunc) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return handler.SendError(c, http.StatusUnauthorized, "API key required for this endpoint")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return handler.SendError(c, http.StatusUnauthorized, "invalid authorization header")
	}

	tokenString := tokenParts[1]

	// Ensure it's an API key (starts with "vhub_")
	if !strings.HasPrefix(tokenString, "vhub_") {
		return handler.SendError(c, http.StatusUnauthorized, "API key required for this endpoint")
	}

	return handleAPIKeyAuth(c, tokenString, next)
}

func handleAPIKeyAuth(c echo.Context, apiKey string, next echo.HandlerFunc) error {
	key, err := model.ValidateAPIKey(apiKey)
	if err != nil {
		return handler.SendError(c, http.StatusUnauthorized, "invalid API key")
	}

	c.Set("user_id", &key.UserID)
	c.Set("api_key", key)

	return next(c)
}

func handleJWTAuth(c echo.Context, tokenString string, next echo.HandlerFunc) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return handler.SendError(c, http.StatusUnauthorized, "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return handler.SendError(c, http.StatusUnauthorized, "invalid token claims")
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		return handler.SendError(c, http.StatusUnauthorized, "invalid user ID in token")
	}

	var user model.User
	if err := model.DB.First(&user, uint(userID)).Error; err != nil {
		return handler.SendError(c, http.StatusUnauthorized, "user not found")
	}

	user.Password = nil
	c.Set("user", &user)

	return next(c)
}
