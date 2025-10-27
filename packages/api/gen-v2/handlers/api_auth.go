package handlers
import (
	"github.com/GIT_USER_ID/GIT_REPO_ID/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ConfirmPasswordReset - 
func (c *Container) ConfirmPasswordReset(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.HelloWorld {
		Message: "Hello World",
	})
}


// ConsumeMagicLink - 
func (c *Container) ConsumeMagicLink(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.HelloWorld {
		Message: "Hello World",
	})
}


// Login - 
func (c *Container) Login(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.HelloWorld {
		Message: "Hello World",
	})
}


// Logout - 
func (c *Container) Logout(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.HelloWorld {
		Message: "Hello World",
	})
}


// RequestMagicLink - 
func (c *Container) RequestMagicLink(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.HelloWorld {
		Message: "Hello World",
	})
}


// RequestPasswordReset - 
func (c *Container) RequestPasswordReset(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.HelloWorld {
		Message: "Hello World",
	})
}


// Signup - 
func (c *Container) Signup(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.HelloWorld {
		Message: "Hello World",
	})
}

