package generated_handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/packages/api/generated_models"
	"net/http"
)

// CreateVault -
func (c *Container) CreateVault(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}

// DeleteVault -
func (c *Container) DeleteVault(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}

// GetVault -
func (c *Container) GetVault(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}

// GetVaults -
func (c *Container) GetVaults(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}

// UpdateVault -
func (c *Container) UpdateVault(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}
