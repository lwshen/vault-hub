package generated_handlers
import (
	"github.com/lwshen/vault-hub/packages/api/generated_models"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetCurrentUser - 
func (c *Container) GetCurrentUser(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated_models.HelloWorld {
		Message: "Hello World",
	})
}

