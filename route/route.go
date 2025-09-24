package route

import (
    "io/fs"
    "net/http"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/filesystem"
    "github.com/lwshen/vault-hub/apps"
    "github.com/lwshen/vault-hub/handler"
    openapi "github.com/lwshen/vault-hub/packages/api"
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

    // Web (embedded)
    app.Use("/", filesystem.New(filesystem.Config{
        Root: http.FS(apps.WebDistFS),
    }))

    // SPA fallback to index.html for client-side routes
    app.Get("/*", func(c *fiber.Ctx) error {
        index, err := fs.ReadFile(apps.WebDistFS, "index.html")
        if err != nil {
            return c.SendStatus(fiber.StatusInternalServerError)
        }
        c.Type("html")
        return c.Send(index)
    })
}
