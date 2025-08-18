package api

import (
	"log/slog"
	"strings"

	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/model"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/gofiber/fiber/v2"
)

func (Server) Login(c *fiber.Ctx) error {
	var input LoginRequest
	if err := c.BodyParser(&input); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	email, err := getEmail(input.Email)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	clientIP, userAgent := getClientInfo(c)

	user := model.User{
		Email: email,
	}
	if err := user.GetByEmail(); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "Invalid email or password")
	}

	if !user.ComparePassword(input.Password) {
		return handler.SendError(c, fiber.StatusBadRequest, "Invalid email or password")
	}

	token, err := user.GenerateToken()
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Record successful login audit log
	if err := model.LogUserAction(model.ActionLoginUser, user.ID, clientIP, userAgent); err != nil {
		slog.Error("Failed to create audit log for login", "error", err, "userID", user.ID)
	}

	resp := LoginResponse{
		Token: token,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (Server) Signup(c *fiber.Ctx) error {
	var input SignupRequest
	if err := c.BodyParser(&input); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	email, err := getEmail(input.Email)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	clientIP, userAgent := getClientInfo(c)

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
		return handler.SendError(c, fiber.StatusBadRequest, strings.Join(errorMsgs, "; "))
	}

	user, err := createUserParams.Create()
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	slog.Info("User created", "email", user.Email, "name", user.Name)

	// Record successful registration audit log
	if err := model.LogUserAction(model.ActionRegisterUser, user.ID, clientIP, userAgent); err != nil {
		slog.Error("Failed to create audit log for signup", "error", err, "userID", user.ID)
	}

	token, err := user.GenerateToken()
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := SignupResponse{
		Token: token,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (Server) Logout(c *fiber.Ctx) error {
	// Try to get user information from context (set by JWT middleware)
	// If there is no authentication information, this should not prevent logout operation
	user, ok := c.Locals("user").(*model.User)
	if ok && user != nil {
		clientIP, userAgent := getClientInfo(c)
		if err := model.LogUserAction(model.ActionLogoutUser, user.ID, clientIP, userAgent); err != nil {
			slog.Error("Failed to create audit log for logout", "error", err, "userID", user.ID)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}

func getEmail(email openapi_types.Email) (string, error) {
	return string(email), nil
}
