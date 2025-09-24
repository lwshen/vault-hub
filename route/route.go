package route

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/lwshen/vault-hub/handler"
	openapi "github.com/lwshen/vault-hub/packages/api"
)

func SetupRoutes(app *fiber.App, webDistFS ...embed.FS) {
	app.Use(jwtMiddleware)

	server := openapi.NewServer()
	openapi.RegisterHandlers(app, server)

	api := app.Group("/api")

	// Auth
	auth := api.Group("/auth")
	auth.Get("/login/oidc", handler.LoginOidc)
	auth.Get("/callback/oidc", handler.LoginOidcCallback)

	// Web - Serve static files (embedded or filesystem)
	if len(webDistFS) > 0 {
		// Serve from embedded filesystem
		distFS, err := fs.Sub(webDistFS[0], "apps/web/dist")
		if err != nil {
			panic("Failed to create sub filesystem for web dist: " + err.Error())
		}

		app.Use("/", filesystem.New(filesystem.Config{
			Root:   http.FS(distFS),
			Index:  "index.html",
			Browse: false,
		}))

		// SPA fallback - serve index.html for all non-API routes
		app.Get("/*", func(c *fiber.Ctx) error {
			// Skip API routes
			if len(c.Path()) >= 4 && c.Path()[:4] == "/api" {
				return c.Next()
			}
			
			indexContent, err := fs.ReadFile(distFS, "index.html")
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			c.Set("Content-Type", "text/html")
			return c.Send(indexContent)
		})
	} else {
		// Fallback to filesystem serving (for development)
		app.Static("/", "./apps/web/dist")
		app.Get("/*", func(c *fiber.Ctx) error {
			if err := c.SendFile("./apps/web/dist/index.html"); err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			return nil
		})
	}
}
