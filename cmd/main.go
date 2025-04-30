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

	log.Fatal(app.Listen(":3000"))
}
