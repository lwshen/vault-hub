package route

import (
	"github.com/gofiber/fiber/v2"
	openapi "github.com/lwshen/vault-hub/packages/api"
	"github.com/lwshen/vault-hub/packages/handler"
)

func SetupRoutes(app *fiber.App) {
	app.Use(jwtMiddleware)

	server := openapi.NewServer()
	openapi.RegisterHandlers(app, server)

	api := app.Group("/api")

	// Auth
	auth := api.Group("/auth")
	auth.Get("/login/oidc", handler.LoginOidc)
	auth.Get("/callback/oidc", handler.LoginOidcCallback)

	// Web
	app.Static("/", "./apps/web/dist")
	app.Get("/*", func(c *fiber.Ctx) error {
		if err := c.SendFile("./apps/web/dist/index.html"); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return nil
	})
}
