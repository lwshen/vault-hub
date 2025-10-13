package api

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/internal/config"
	"github.com/lwshen/vault-hub/internal/email"
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
	if err := model.LogUserAction(model.ActionLoginUser, user.ID, model.SourceWeb, clientIP, userAgent); err != nil {
		slog.Error("Failed to create audit log for login", "error", err, "userID", user.ID)
	}

	resp := LoginResponse{
		Token: token,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// Signup handles user registration requests
// It validates input, creates the user account, and returns a JWT token
func (Server) Signup(c *fiber.Ctx) error {
	// Parse request body
	input, err := parseSignupRequest(c)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Extract client information for audit logging
	clientIP, userAgent := getClientInfo(c)

	// Validate and create user parameters
	createParams, err := buildUserCreateParams(input)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Create the user account
	user, err := createUser(createParams)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	slog.Info("User created", "email", user.Email, "name", user.Name)

	// Log successful registration
	logSignupAudit(user.ID, clientIP, userAgent)

	// Send email verification if SMTP is enabled
	if err := sendEmailVerification(user); err != nil {
		slog.Error("Failed to send verification email", "error", err, "userID", user.ID)
		// Don't fail signup if email sending fails
	}

	// Generate authentication token
	token, err := user.GenerateToken()
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := SignupResponse{
		Token: token,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// parseSignupRequest parses and validates the signup request body
func parseSignupRequest(c *fiber.Ctx) (SignupRequest, error) {
	var input SignupRequest
	if err := c.BodyParser(&input); err != nil {
		return input, err
	}
	return input, nil
}

// buildUserCreateParams creates and validates user creation parameters
func buildUserCreateParams(input SignupRequest) (model.CreateUserParams, error) {
	email, err := getEmail(input.Email)
	if err != nil {
		return model.CreateUserParams{}, err
	}

	// For email/password signup, password is required
	if input.Password == "" {
		return model.CreateUserParams{}, fmt.Errorf("password is required")
	}

	createUserParams := model.CreateUserParams{
		Email:    string(email),
		Password: &input.Password,
		Name:     input.Name,
	}

	errors := createUserParams.Validate()
	if len(errors) > 0 {
		// Convert map values to slice for joining
		var errorMsgs []string
		for _, msg := range errors {
			errorMsgs = append(errorMsgs, msg)
		}
		errorMsg := strings.Join(errorMsgs, "; ")
		return model.CreateUserParams{}, fmt.Errorf("%s", errorMsg)
	}

	return createUserParams, nil
}

// createUser creates a new user account
func createUser(params model.CreateUserParams) (*model.User, error) {
	return params.Create()
}

// logSignupAudit records the signup action in audit logs
func logSignupAudit(userID uint, clientIP, userAgent string) {
	if err := model.LogUserAction(model.ActionRegisterUser, userID, model.SourceWeb, clientIP, userAgent); err != nil {
		slog.Error("Failed to create audit log for signup", "error", err, "userID", userID)
	}
}

func (Server) Logout(c *fiber.Ctx) error {
	// Try to get user information from context (set by JWT middleware)
	// If there is no authentication information, this should not prevent logout operation
	user, ok := c.Locals("user").(*model.User)
	if ok && user != nil {
		clientIP, userAgent := getClientInfo(c)
		if err := model.LogUserAction(model.ActionLogoutUser, user.ID, model.SourceWeb, clientIP, userAgent); err != nil {
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

// sendEmailVerification sends an email verification to the user
func sendEmailVerification(user *model.User) error {
	if !config.SMTPEnabled {
		slog.Debug("SMTP not enabled, skipping email verification", "userID", user.ID)
		return nil
	}

	// Create email verification token
	token, err := model.CreateEmailVerificationToken(user.ID)
	if err != nil {
		return fmt.Errorf("failed to create verification token: %w", err)
	}

	// Build verification URL
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", config.BaseURL, token.Token)

	// Send verification email
	recipientName := "User"
	if user.Name != nil {
		recipientName = *user.Name
	}

	if err := email.SendEmailVerificationAsync(user.Email, recipientName, verificationURL); err != nil {
		return fmt.Errorf("failed to queue verification email: %w", err)
	}

	slog.Info("Verification email queued", "email", user.Email, "userID", user.ID)
	return nil
}

// getClientInfo extracts IP address and User-Agent from the request
func getClientInfo(c *fiber.Ctx) (string, string) {
	// Get IP address (check for forwarded headers first)
	ip := c.Get("X-Forwarded-For")
	if ip == "" {
		ip = c.Get("X-Real-IP")
	}
	if ip == "" {
		ip = c.IP()
	}
	// Get User-Agent
	userAgent := c.Get("User-Agent")
	return ip, userAgent
}
