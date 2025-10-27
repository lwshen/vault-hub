package route

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lwshen/vault-hub/internal/embed"
	"github.com/lwshen/vault-hub/packages/api"
)

// SetupEchoRoutes configures Echo router with all routes and middleware
func SetupEchoRoutes(e *echo.Echo, logger *slog.Logger) error {
	// Global middleware
	e.Use(SlogMiddleware(logger))
	e.Use(SecurityHeadersMiddleware())
	e.Use(CORSMiddleware())

	// Import and reference the handlers to avoid circular imports
	// TODO: Implement these handlers with proper business logic
	// For now, these return placeholder responses

	// Health and status endpoints (public)
	e.GET("/api/health", func(c echo.Context) error {
		adapter := api.NewEchoAdapter()
		response := adapter.ConvertHealthCheckResponse("ok", time.Now())
		return c.JSON(http.StatusOK, response)
	})

	e.GET("/api/status", func(c echo.Context) error {
		adapter := api.NewEchoAdapter()
		response := api.StatusResponse{
			Version: "1.0.0-migration",
			Uptime: "0h 0m",
			Database: api.StatusResponseDatabase{Status: "healthy"},
			System: api.StatusResponseSystem{Status: "healthy"},
		}
		return c.JSON(http.StatusOK, response)
	})

	e.GET("/api/config", func(c echo.Context) error {
		adapter := api.NewEchoAdapter()
		response := adapter.ConvertConfigResponse()
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

	return nil
}

	// Static file serving (SPA support)
	embedFS, err := embed.GetDistFS()
	if err != nil {
		slog.Error("Failed to initialize embedded filesystem", "error", err)
		os.Exit(1)
	}

	// Serve static files
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "dist",
		Browse:      false,
		HTML5:       true,
		Index:       "index.html",
		IgnoreBase:  false,
		RootURL:     "/",
		Filesystem:  http.FS(embedFS),
	}))

	return nil
}

// handleOIDCLogin handles OIDC login redirect (custom implementation)
func handleOIDCLogin(c echo.Context) error {
	// TODO: Implement OIDC login logic similar to existing handler.LoginOidc
	return echo.NewHTTPError(http.StatusNotImplemented, "OIDC login not yet migrated to Echo")
}

// handleOIDCCallback handles OIDC callback (custom implementation)
func handleOIDCCallback(c echo.Context) error {
	// TODO: Implement OIDC callback logic similar to existing handler.LoginOidcCallback
	return echo.NewHTTPError(http.StatusNotImplemented, "OIDC callback not yet migrated to Echo")
}