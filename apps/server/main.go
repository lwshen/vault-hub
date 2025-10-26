package main

import (
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lwshen/vault-hub/internal/config"
	"github.com/lwshen/vault-hub/internal/version"
	"github.com/lwshen/vault-hub/model"
	"github.com/lwshen/vault-hub/route"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Info("Starting VaultHub Server", "version", version.Version, "commit", version.Commit)

	err := model.Open(logger)
	if err != nil {
		logger.Error("Failed to open database", "error", err)
		os.Exit(1)
	}

	e := echo.New()

	// Built-in middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Setup routes
	route.SetupRoutesEcho(e)

	// Start server
	e.Logger.Fatal(e.Start(":" + config.AppPort))
}
