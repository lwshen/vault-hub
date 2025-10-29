package generated_handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/packages/api/generated_models"
	"net/http"
)

// GetVaultByAPIKey -
func (c *Container) GetVaultByAPIKey(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}

// GetVaultByNameAPIKey -
func (c *Container) GetVaultByNameAPIKey(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}

// GetVaultsByAPIKey -
func (c *Container) GetVaultsByAPIKey(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}
