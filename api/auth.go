package api

import (
	"log/slog"

	"github.com/lwshen/vault-hub/model"

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

	email, err := input.Email.MarshalJSON()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	createUserParams := model.CreateUserParams{
		Email:    string(email),
		Password: input.Password,
		Name:     input.Name,
	}

	errors := createUserParams.Validate()
	if len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors,
		})
	}

	user, err := createUserParams.Create()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	slog.Info("User created", "email", user.Email, "name", user.Name)

	token, err := user.GenerateToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

func (Server) Logout(c *fiber.Ctx) error {
	return nil
}
