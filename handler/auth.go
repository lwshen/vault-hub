package handler

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/internal/auth"
)

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

	return nil
}
