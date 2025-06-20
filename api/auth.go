package api

import (
	"log/slog"
	"strings"

	"github.com/lwshen/vault-hub/model"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/gofiber/fiber/v2"
)

func (Server) Login(c *fiber.Ctx) error {
	var input LoginRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	email, err := getEmail(input.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user := model.User{
		Email: email,
	}
	if err := user.GetByEmail(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	if !user.ComparePassword(input.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	token, err := user.GenerateToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp := LoginResponse{
		Token: token,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (Server) Signup(c *fiber.Ctx) error {
	var input SignupRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	email, err := getEmail(input.Email)
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
		var errorMsgs []string
		for _, msg := range errors {
			errorMsgs = append(errorMsgs, msg)
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": strings.Join(errorMsgs, "; "),
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

	resp := SignupResponse{
		Token: token,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (Server) Logout(c *fiber.Ctx) error {
	return nil
}

func getEmail(email openapi_types.Email) (string, error) {
	return string(email), nil
}
