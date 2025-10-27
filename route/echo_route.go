package route

import (
	"io/fs"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/internal/embed"
)

// SetupEchoRoutes configures Echo router with all routes and middleware
func SetupEchoRoutes(e *echo.Echo) error {
	// TODO: Add middleware when needed
	// e.Use(SlogMiddleware(logger))
	// e.Use(SecurityHeadersMiddleware())
	// e.Use(CORSMiddleware())

	// Import and reference the handlers to avoid circular imports
	// TODO: Implement these handlers with proper business logic
	// For now, these return placeholder responses

	// Health and status endpoints (public)
	e.GET("/api/health", func(c echo.Context) error {
		response := map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now(),
		}
		return c.JSON(http.StatusOK, response)
	})

	e.GET("/api/status", func(c echo.Context) error {
		response := map[string]interface{}{
			"version":  "1.0.0-migration",
			"uptime":   "0h 0m",
			"database": map[string]string{"status": "healthy"},
			"system":   map[string]string{"status": "healthy"},
		}
		return c.JSON(http.StatusOK, response)
	})

	e.GET("/api/config", func(c echo.Context) error {
		response := map[string]interface{}{
			"isOidcEnabled":      false, // TODO: Implement proper config check
			"isEmailEnabled":     false, // TODO: Implement proper config check
			"passwordMinLength":  8,
			"isRegistrationOpen": true,
		}
		return c.JSON(http.StatusOK, response)
	})

	// Authentication endpoints (public with auth flow)
	auth := e.Group("/api/auth")
	auth.POST("/login", func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Login not yet implemented in Echo migration"})
	})
	auth.POST("/signup", func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Signup not yet implemented in Echo migration"})
	})
	auth.GET("/logout", func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Logout not yet implemented in Echo migration"})
	})
	auth.GET("/login/oidc", func(c echo.Context) error {
		// TODO: Implement OIDC login
		return c.JSON(http.StatusNotImplemented, map[string]string{"error": "OIDC login not yet implemented in Echo migration"})
	})

	// TODO: Add remaining endpoints with proper handlers
	// Vault endpoints, API key endpoints, etc.

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

// TODO: Implement OIDC handlers when needed
