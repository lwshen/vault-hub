package echoapp

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/internal/config"
	"github.com/lwshen/vault-hub/model"
)

// SecurityMiddleware mirrors the Fiber auth gate by routing requests through the
// appropriate JWT or API key authentication handlers.
func SecurityMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()

			if isPublicRoute(path) {
				return next(c)
			}

			if strings.HasPrefix(path, "/api/cli/") {
				return apiKeyOnlyMiddleware(next)(c)
			}

			if strings.HasPrefix(path, "/api/") {
				return jwtOnlyMiddleware(next)(c)
			}

			return next(c)
		}
	}
}

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

func jwtOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return sendError(c, http.StatusUnauthorized, "JWT token required")
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return sendError(c, http.StatusUnauthorized, "invalid authorization header")
		}

		tokenString := tokenParts[1]

		if strings.HasPrefix(tokenString, "vhub_") {
			return sendError(c, http.StatusUnauthorized, "JWT token required for this endpoint")
		}

		if err := handleJWTAuth(c, tokenString); err != nil {
			return err
		}

		return next(c)
	}
}

func apiKeyOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return sendError(c, http.StatusUnauthorized, "API key required for this endpoint")
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return sendError(c, http.StatusUnauthorized, "invalid authorization header")
		}

		tokenString := tokenParts[1]

		if !strings.HasPrefix(tokenString, "vhub_") {
			return sendError(c, http.StatusUnauthorized, "API key required for this endpoint")
		}

		if err := handleAPIKeyAuth(c, tokenString); err != nil {
			return err
		}

		return next(c)
	}
}

func handleAPIKeyAuth(c echo.Context, apiKey string) error {
	key, err := model.ValidateAPIKey(apiKey)
	if err != nil {
		return sendError(c, http.StatusUnauthorized, "invalid API key")
	}

	c.Set("user_id", &key.UserID)
	c.Set("api_key", key)

	return nil
}

func handleJWTAuth(c echo.Context, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return sendError(c, http.StatusUnauthorized, "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return sendError(c, http.StatusUnauthorized, "invalid token claims")
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		return sendError(c, http.StatusUnauthorized, "invalid user ID in token")
	}

	var user model.User
	if err := model.DB.First(&user, uint(userID)).Error; err != nil {
		return sendError(c, http.StatusUnauthorized, "user not found")
	}

	user.Password = nil
	c.Set("user", &user)

	return nil
}
