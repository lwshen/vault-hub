package route

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/internal/embed"
	openapi "github.com/lwshen/vault-hub/packages/api"
)

func SetupRoutes(e *echo.Echo) {
	e.Use(jwtMiddleware)

	server := openapi.NewServer()
	openapi.RegisterHandlers(e, server)

	api := e.Group("/api")

	// Auth
	auth := api.Group("/auth")
	auth.GET("/login/oidc", handler.LoginOidc)
	auth.GET("/callback/oidc", handler.LoginOidcCallback)
	// Magic link consume endpoint exposed via API namespace
	auth.GET("/magic-link/token", func(c echo.Context) error {
		token := c.QueryParam("token")
		return server.ConsumeMagicLink(c, openapi.ConsumeMagicLinkParams{Token: token})
	})

	// Web - Static file serving with SPA fallback
	embedFS, err := embed.GetDistFS()
	if err != nil {
		slog.Error("Failed to initialize embedded filesystem", "error", err)
		os.Exit(1)
	}

	// Serve static files with SPA fallback
	// Use Echo's static file handler
	assetHandler := echo.WrapHandler(http.FileServer(http.FS(embedFS)))
	e.GET("/*", func(c echo.Context) error {
		path := c.Request().URL.Path

		// Skip API routes
		if embedFSIsAPI(path) {
			return c.NoContent(http.StatusNotFound)
		}

		// Try to serve the file
		file, err := embedFS.Open(path)
		if err == nil {
			file.Close()
			return assetHandler(c)
		}

		// For SPA routing, serve index.html for any non-API route
		indexFile, err := embedFS.Open("index.html")
		if err == nil {
			indexFile.Close()
			// Use Echo's File method which handles embedded files correctly
			return c.File("index.html")
		}

		return assetHandler(c)
	})
}

// embedFSIsAPI checks if a path is an API route
func embedFSIsAPI(path string) bool {
	return len(path) >= 5 && path[:5] == "/api/"
}
