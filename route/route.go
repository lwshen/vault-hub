package route

import (
	"github.com/gofiber/fiber/v2"
	openapi "github.com/lwshen/vault-hub/api"
	"github.com/lwshen/vault-hub/handler"
)

func SetupRoutes(app *fiber.App) {
	server := openapi.NewServer()
	openapi.RegisterHandlers(app, server)

	api := app.Group("/api")

	// Auth
	auth := api.Group("/auth")
	auth.Post("/signup", handler.Signup)
	auth.Post("/login", handler.Login)
	auth.Post("/logout", handler.Logout)
	auth.Get("/login/oidc", handler.LoginOidc)
	auth.Get("/callback/oidc", handler.LoginOidcCallback)

	// Web
	app.Static("/", "./web/dist")
	app.Get("/*", func(c *fiber.Ctx) error {
		if err := c.SendFile("./web/dist/index.html"); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return nil
	})
}
