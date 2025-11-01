package generated_handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/packages/api/generated_models"
	"net/http"
)

// GetAuditLogs -
func (c *Container) GetAuditLogs(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}

// GetAuditMetrics -
func (c *Container) GetAuditMetrics(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld{
		Message: "Hello World",
	})
}
