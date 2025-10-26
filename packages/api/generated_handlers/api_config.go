package generated_handlers
import (
	"github.com/lwshen/vault-hub/packages/api/generated_models"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetConfig - Get public configuration
func (c *Container) GetConfig(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld {
		Message: "Hello World",
	})
}

