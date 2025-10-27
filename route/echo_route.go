package route

import (
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/internal/embed"
	"github.com/lwshen/vault-hub/packages/api"
)

// SetupEchoRoutes configures Echo router with all routes and middleware
func SetupEchoRoutes(e *echo.Echo) error {
	// Add middleware
	// TODO: Pass actual logger to SlogMiddleware
	// e.Use(SlogMiddleware(nil))
	e.Use(SecurityHeadersMiddleware())
	e.Use(CORSMiddleware())
	e.Use(AuthMiddleware())

	// Register only basic handlers for now
	// TODO: Fix interface compatibility and register full API
	registerBasicHandlers(e)

	// Setup static file serving from embedded filesystem
	distFS, err := embed.GetDistFS()
	if err != nil {
		return err
	}

	// Serve static assets (CSS, JS, images, etc.)
	e.GET("/assets/*", func(c echo.Context) error {
		return echo.WrapHandler(http.FileServer(http.FS(distFS)))(c)
	})

	// Serve other static files (fonts, icons, etc.)
	e.GET("/fonts/*", func(c echo.Context) error {
		return echo.WrapHandler(http.FileServer(http.FS(distFS)))(c)
	})
	e.GET("/icon.svg", func(c echo.Context) error {
		return echo.WrapHandler(http.FileServer(http.FS(distFS)))(c)
	})
	e.GET("/vite.svg", func(c echo.Context) error {
		return echo.WrapHandler(http.FileServer(http.FS(distFS)))(c)
	})

	// Serve index.html for all other routes (SPA routing)
	// This should be added last to act as a catch-all for client-side routes
	e.GET("/*", func(c echo.Context) error {
		indexHTML, err := fs.ReadFile(distFS, "index.html")
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Page not found"})
		}
		return c.HTMLBlob(http.StatusOK, indexHTML)
	})

	return nil
}

// registerBasicHandlers registers essential handlers for testing
func registerBasicHandlers(e *echo.Echo) {
	server := api.NewServer()

	// Authentication endpoints
	e.POST("/api/auth/login", server.Login)
	e.POST("/api/auth/signup", server.Signup)
	e.GET("/api/auth/logout", server.Logout)

	// User endpoints
	e.GET("/api/user", server.GetCurrentUser)

	// System endpoints
	e.GET("/api/config", server.GetConfig)
	e.GET("/api/health", server.Health)
	e.GET("/api/status", server.GetStatus)
}

// TODO: Implement OIDC handlers when needed
