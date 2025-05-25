package api

import "github.com/gofiber/fiber/v2"

func (Server) PostApiAuthLogin(c *fiber.Ctx) error {
	var input LoginRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return nil
}

func (Server) PostApiAuthSignup(c *fiber.Ctx) error {
	var input SignupRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return nil
}

func (Server) GetApiAuthLogout(c *fiber.Ctx) error {
	return nil
}
