package generated_handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/packages/api/generated_models"
	"net/http"
)

// ConfirmPasswordReset -
func (c *Container) ConfirmPasswordReset(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}

// ConsumeMagicLink -
func (c *Container) ConsumeMagicLink(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}

// Login -
func (c *Container) Login(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}

// Logout -
func (c *Container) Logout(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}

// RequestMagicLink -
func (c *Container) RequestMagicLink(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}

// RequestPasswordReset -
func (c *Container) RequestPasswordReset(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}

// Signup -
func (c *Container) Signup(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}
