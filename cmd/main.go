package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/api"
	"github.com/lwshen/vault-hub/internal/config"
	"github.com/lwshen/vault-hub/internal/db"
	slogfiber "github.com/samber/slog-fiber"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db.Open(logger)

	server := api.NewServer()

	app := fiber.New()

	app.Use(slogfiber.New(logger))

	api.RegisterHandlers(app, server)

	app.Static("/", "./web/dist")

	app.Get("/*", func(c *fiber.Ctx) error {
		if err := c.SendFile("./web/dist/index.html"); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return nil
	})

	log.Fatal(app.Listen(":" + config.AppPort))
}
