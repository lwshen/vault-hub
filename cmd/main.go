package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/internal/config"
	"github.com/lwshen/vault-hub/model"
	"github.com/lwshen/vault-hub/route"
	slogfiber "github.com/samber/slog-fiber"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := model.Open(logger)
	if err != nil {
		logger.Error("Failed to open database", "error", err)
		os.Exit(1)
	}

	app := fiber.New()

	app.Use(slogfiber.New(logger))

	route.SetupRoutes(app)

	log.Fatal(app.Listen(":" + config.AppPort))
}
