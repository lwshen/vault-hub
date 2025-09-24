package route

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/internal/embed"
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

	// Web - serve from embedded filesystem
	distFS, err := embed.GetDistFS()
	if err != nil {
		// Fallback to serving from filesystem if embed fails
		app.Static("/", "./apps/web/dist")
		app.Get("/*", func(c *fiber.Ctx) error {
			if err := c.SendFile("./apps/web/dist/index.html"); err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			return nil
		})
	} else {
		// Serve from embedded filesystem
		app.Use("/", filesystem.New(filesystem.Config{
			Root:         http.FS(distFS),
			Browse:       false,
			Index:        "index.html",
			NotFoundFile: "index.html",
		}))
	}
}
