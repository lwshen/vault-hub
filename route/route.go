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

	// Web
	embedFS, err := embed.GetDistFS()
	if err != nil {
		panic(err)
	}
	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(embedFS),
		Browse:       false,
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))
}
