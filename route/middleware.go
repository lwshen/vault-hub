package route

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/internal/config"
	"github.com/lwshen/vault-hub/model"
)

func jwtMiddleware(c *fiber.Ctx) error {
	path := c.Path()

	// Public routes that don't need authentication
	if isPublicRoute(path) {
		return c.Next()
	}

	// Routes starting with /api/api-key/ MUST use API key authentication
	if strings.HasPrefix(path, "/api/api-key/") {
		return apiKeyOnlyMiddleware(c)
	}

	// All other /api/ routes require JWT authentication
	if strings.HasPrefix(path, "/api/") {
		return jwtOnlyMiddleware(c)
	}

	// Non-API routes (web assets, etc.) don't need auth
	return c.Next()
}

// isPublicRoute checks if a route is public and doesn't need authentication
func isPublicRoute(path string) bool {
	publicRoutes := []string{
		"/api/auth/login",
		"/api/auth/register", 
		"/api/auth/login/oidc",
		"/api/auth/callback/oidc",
	}
	
	for _, route := range publicRoutes {
		if strings.HasPrefix(path, route) {
			return true
		}
	}
	return false
}

// jwtOnlyMiddleware ensures non-API-key routes only accept JWT authentication
func jwtOnlyMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return handler.SendError(c, fiber.StatusUnauthorized, "JWT token required")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return handler.SendError(c, fiber.StatusUnauthorized, "invalid authorization header")
	}

	tokenString := tokenParts[1]

	// Reject API keys on JWT-only routes
	if strings.HasPrefix(tokenString, "vhub_") {
		return handler.SendError(c, fiber.StatusUnauthorized, "JWT token required for this endpoint")
	}

	return handleJWTAuth(c, tokenString)
}

// apiKeyOnlyMiddleware ensures API key routes only accept API key authentication
func apiKeyOnlyMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return handler.SendError(c, fiber.StatusUnauthorized, "API key required for this endpoint")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return handler.SendError(c, fiber.StatusUnauthorized, "invalid authorization header")
	}

	tokenString := tokenParts[1]

	// Ensure it's an API key (starts with "vhub_")
	if !strings.HasPrefix(tokenString, "vhub_") {
		return handler.SendError(c, fiber.StatusUnauthorized, "API key required for this endpoint")
	}

	return handleAPIKeyAuth(c, tokenString)
}

func handleAPIKeyAuth(c *fiber.Ctx, apiKey string) error {
	key, err := model.ValidateAPIKey(apiKey)
	if err != nil {
		return handler.SendError(c, fiber.StatusUnauthorized, "invalid API key")
	}

	c.Locals("user_id", &key.UserID)
	c.Locals("api_key", key)

	return c.Next()
}

func handleJWTAuth(c *fiber.Ctx, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return handler.SendError(c, fiber.StatusUnauthorized, "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return handler.SendError(c, fiber.StatusUnauthorized, "invalid token claims")
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		return handler.SendError(c, fiber.StatusUnauthorized, "invalid user ID in token")
	}

	var user model.User
	if err := model.DB.First(&user, uint(userID)).Error; err != nil {
		return handler.SendError(c, fiber.StatusUnauthorized, "user not found")
	}

	user.Password = nil
	c.Locals("user", &user)

	return c.Next()
}
