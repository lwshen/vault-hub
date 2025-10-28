package route

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lwshen/vault-hub/internal/config"
	"github.com/lwshen/vault-hub/model"
)

// SlogMiddleware creates an Echo middleware that uses slog for logging
func SlogMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			start := time.Now()

			// Log request details
			logger.Info("HTTP Request",
				"method", req.Method,
				"path", req.URL.Path,
				"query", req.URL.RawQuery,
				"remote_ip", c.RealIP(),
				"user_agent", req.UserAgent(),
			)

			// Process request
			err := next(c)

			// Log response details
			duration := time.Since(start)
			res := c.Response()
			logger.Info("HTTP Response",
				"method", req.Method,
				"path", req.URL.Path,
				"status", res.Status,
				"duration_ms", duration.Milliseconds(),
				"size", res.Size,
			)

			return err
		}
	}
}

// SecurityHeadersMiddleware adds security-related headers
func SecurityHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			res := c.Response()
			res.Header().Set("X-Content-Type-Options", "nosniff")
			res.Header().Set("X-Frame-Options", "DENY")
			res.Header().Set("X-XSS-Protection", "1; mode=block")
			res.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			return next(c)
		}
	}
}

// CORSMiddleware provides CORS support
func CORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Configure appropriately for production
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	})
}

// AuthMiddleware handles authentication based on route patterns
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path

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
}

// jwtOnlyMiddleware ensures non-API-key routes only accept JWT authentication
func jwtOnlyMiddleware(c echo.Context, next echo.HandlerFunc) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "JWT token required",
		})
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid authorization header",
		})
	}

	tokenString := tokenParts[1]

	// Reject API keys on JWT-only routes
	if strings.HasPrefix(tokenString, "vhub_") {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "JWT token required for this endpoint",
		})
	}

	return handleJWTAuth(c, tokenString, next)
}

// apiKeyOnlyMiddleware ensures API key routes only accept API key authentication
func apiKeyOnlyMiddleware(c echo.Context, next echo.HandlerFunc) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "API key required for this endpoint",
		})
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid authorization header",
		})
	}

	tokenString := tokenParts[1]

	// Ensure it's an API key (starts with "vhub_")
	if !strings.HasPrefix(tokenString, "vhub_") {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "API key required for this endpoint",
		})
	}

	return handleAPIKeyAuth(c, tokenString, next)
}

func handleAPIKeyAuth(c echo.Context, apiKey string, next echo.HandlerFunc) error {
	key, err := model.ValidateAPIKey(apiKey)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid API key",
		})
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
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid token claims",
		})
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid user ID in token",
		})
	}

	var user model.User
	if err := model.DB.First(&user, uint(userID)).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "user not found",
		})
	}

	user.Password = nil
	c.Set("user", &user)

	return next(c)
}
