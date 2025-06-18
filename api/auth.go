package api

import (
	"github.com/gofiber/fiber/v2"
)

func (Server) Login(c *fiber.Ctx) error {
	var input LoginRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return nil
}

func (Server) Signup(c *fiber.Ctx) error {
	var input SignupRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return nil
}

func (Server) Logout(c *fiber.Ctx) error {
	return nil
}
