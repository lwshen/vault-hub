package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/api"
)

func main() {
	server := api.NewServer()

	app := fiber.New()

	api.RegisterHandlers(app, server)

	app.Static("/", "./web/dist")

	app.Get("/*", func(c *fiber.Ctx) error {
		if err := c.SendFile("./web/dist/index.html"); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return nil
	})

	log.Fatal(app.Listen(":3000"))
}
