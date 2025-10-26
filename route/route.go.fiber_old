package route

import (
	"log/slog"
	"net/http"
	"os"

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
	// Magic link consume endpoint exposed via API namespace
	auth.Get("/magic-link/token", func(c *fiber.Ctx) error {
		token := c.Query("token")
		return server.ConsumeMagicLink(c, openapi.ConsumeMagicLinkParams{Token: token})
	})

	// Web
	embedFS, err := embed.GetDistFS()
	if err != nil {
		slog.Error("Failed to initialize embedded filesystem", "error", err)
		os.Exit(1)
	}
	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(embedFS),
		Browse:       false,
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))
}
