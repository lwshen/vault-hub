package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/internal/config"
	"github.com/lwshen/vault-hub/internal/version"
	"github.com/lwshen/vault-hub/model"
	"github.com/lwshen/vault-hub/route"
	slogfiber "github.com/samber/slog-fiber"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Info("Starting VaultHub Server", "version", version.Version, "commit", version.Commit)

	err := model.Open(logger)
	if err != nil {
		logger.Error("Failed to open database", "error", err)
		os.Exit(1)
	}

	// Ensure demo user exists when demo mode is enabled
	if config.DemoEnabled {
		logger.Info("Demo mode enabled, ensuring demo user exists")
		if err := model.EnsureDemoUser(); err != nil {
			logger.Error("Failed to ensure demo user", "error", err)
			os.Exit(1)
		}
		logger.Info("Demo user verified successfully", "email", model.DemoUserEmail)
	}

	app := fiber.New()

	app.Use(slogfiber.New(logger))

	route.SetupRoutes(app)

	log.Fatal(app.Listen(":" + config.AppPort))
}
