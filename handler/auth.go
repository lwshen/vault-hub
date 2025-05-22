package handler

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/internal/auth"
	"github.com/lwshen/vault-hub/model"
)

func Signup(c *fiber.Ctx) error {
	var input model.CreateUserParams
	if err := c.BodyParser(&input); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if errors := input.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	return nil
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return nil
}

func Logout(c *fiber.Ctx) error {
	return nil
}

func LoginOidc(c *fiber.Ctx) error {
	baseUrl := c.BaseURL()
	url := auth.AuthCodeURL(c, baseUrl)
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
