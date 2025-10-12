package handler

import (
	"log/slog"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/internal/auth"
	"github.com/lwshen/vault-hub/model"
)

// LoginResponse represents the response for successful login
type LoginResponse struct {
	Token string `json:"token"`
}

func LoginOidc(c *fiber.Ctx) error {
	baseUrl := c.BaseURL()
	url, err := auth.AuthCodeURL(c, baseUrl)
	if err != nil {
		slog.Error("Failed to get OIDC URL", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	slog.Debug("Login with OIDC", "url", url)
	return c.Redirect(url)
}

func LoginOidcCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")
	slog.Debug("Login with OIDC callback", "uri", c.Request().URI(), "code", code, "state", state)

	err := auth.VerifyState(c, state)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	token, err := auth.Verify(c.Context(), code)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	slog.Debug("Login with OIDC callback", "token", token)

	userInfo, err := auth.UserInfo(c.Context(), token)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	slog.Debug("Login with OIDC callback", "userInfo", userInfo)

	// Extract claims from userInfo
	var claims map[string]interface{}
	if err := userInfo.Claims(&claims); err != nil {
		slog.Error("Failed to extract OIDC claims", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Extract email from claims
	email, ok := claims["email"].(string)
	if !ok || email == "" {
		slog.Error("OIDC userInfo missing email claim", "claims", claims)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Look up user by email
	user := model.User{
		Email: email,
	}
	if err := user.GetByEmail(); err != nil {
		// User doesn't exist, create new user from OIDC data
		name := ""
		if nameClaim, ok := claims["name"].(string); ok {
			name = nameClaim
		}

		createParams := model.CreateUserParams{
			Email:    email,
			Password: nil, // OIDC users don't need passwords
			Name:     name,
		}

		newUser, createErr := createParams.Create()
		if createErr != nil {
			slog.Error("Failed to create user from OIDC", "error", createErr, "email", email)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		user = *newUser
		slog.Info("User created from OIDC", "email", email, "name", name)
	}

	// Generate JWT token for the user
	jwtToken, err := user.GenerateToken()
	if err != nil {
		slog.Error("Failed to generate token for OIDC user", "error", err, "userID", user.ID)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Record successful login audit log
	clientIP, userAgent := getClientInfo(c)
	if err := model.LogUserAction(model.ActionLoginUser, user.ID, model.SourceWeb, clientIP, userAgent); err != nil {
		slog.Error("Failed to create audit log for OIDC login", "error", err, "userID", user.ID)
	}

	// Redirect back to frontend with token in URL fragment (hash) for security
	// URL fragments are never sent to the server, preventing token leakage in logs, Referer headers, and browser history
	redirectUrl := "/login#token=" + url.QueryEscape(jwtToken) + "&source=oidc"
	return c.Redirect(redirectUrl)
}

// getClientInfo extracts IP address and User-Agent from the request
func getClientInfo(c *fiber.Ctx) (string, string) {
	// Get IP address (check for forwarded headers first)
	ip := c.Get("X-Forwarded-For")
	if ip == "" {
		ip = c.Get("X-Real-IP")
	}
	if ip == "" {
		ip = c.IP()
	}
	// Get User-Agent
	userAgent := c.Get("User-Agent")
	return ip, userAgent
}
