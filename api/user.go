package api

import "github.com/gofiber/fiber/v2"

func (Server) GetCurrentUser(c *fiber.Ctx) error {
	resp := GetUserResponse{
		Email: "test@test.com",
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
