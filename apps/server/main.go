package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/internal/config"
	"github.com/lwshen/vault-hub/internal/email"
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

	// Initialize email service
	if config.SMTPEnabled {
		email.InitService()
		logger.Info("Email service initialized successfully")
	} else {
		logger.Warn("Email service disabled - SMTP not configured")
	}

	app := fiber.New()

	app.Use(slogfiber.New(logger))

	route.SetupRoutes(app)

	log.Fatal(app.Listen(":" + config.AppPort))
}
