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

	// Routes starting with /api/api-key/ MUST use API key authentication
	if strings.HasPrefix(path, "/api/api-key/") {
		return apiKeyOnlyMiddleware(c)
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Next()
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return handler.SendError(c, fiber.StatusUnauthorized, "invalid authorization header")
	}

	tokenString := tokenParts[1]

	// Handle JWT token
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
