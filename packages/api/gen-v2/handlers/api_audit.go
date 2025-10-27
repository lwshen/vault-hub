package handlers
import (
	"github.com/GIT_USER_ID/GIT_REPO_ID/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetAuditLogs - 
func (c *Container) GetAuditLogs(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.HelloWorld {
		Message: "Hello World",
	})
}


// GetAuditMetrics - 
func (c *Container) GetAuditMetrics(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.HelloWorld {
		Message: "Hello World",
	})
}

