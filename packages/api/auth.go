package api

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/lwshen/vault-hub/handler"
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

	slog.Info("User created", "email", user.Email, "name", *user.Name)

	// Log successful registration
	logSignupAudit(user.ID, clientIP, userAgent)

	// Fire-and-forget signup email (do not block response)
	go func(u *model.User) {
		sender := email.NewSMTPSender()
		svc := email.NewService(sender, "Vault Hub")
		name := ""
		if u.Name != nil {
			name = *u.Name
		}
		if err := svc.SendSignupConfirmation(u.Email, name); err != nil {
			slog.Error("Failed to send signup confirmation", "error", err, "email", u.Email)
		}
		_ = model.LogUserAction(model.ActionSendSignupEmail, u.ID, model.SourceWeb, clientIP, userAgent)
	}(user)

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

// Temporary request types until OpenAPI generation is updated
type PasswordResetRequest struct {
	Email openapi_types.Email `json:"email"`
}

type PasswordResetConfirmRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type MagicLinkRequest struct {
	Email openapi_types.Email `json:"email"`
}

// RequestPasswordReset creates a password reset token and sends email
func (Server) RequestPasswordReset(c *fiber.Ctx) error {
	var input PasswordResetRequest
	if err := c.BodyParser(&input); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	emailStr, err := getEmail(input.Email)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	// Always respond 200 to avoid user enumeration
	// Attempt to find user; if not found, still return success
	user := model.User{Email: emailStr}
	if err := user.GetByEmail(); err == nil {
		token, _, err := model.CreateEmailToken(user.ID, model.TokenPurposeResetPassword, 30*time.Minute)
		if err == nil {
			baseURL := c.BaseURL()
			actionURL := fmt.Sprintf("%s/reset?token=%s", baseURL, token)
			go func(u model.User, url string) {
				sender := email.NewSMTPSender()
				svc := email.NewService(sender, "Vault Hub")
				name := ""
				if u.Name != nil {
					name = *u.Name
				}
				if err := svc.SendPasswordReset(u.Email, name, url, "30m"); err != nil {
					slog.Error("Failed to send password reset email", "error", err, "email", u.Email)
				}
			}(user, actionURL)
		}
	}
	return c.SendStatus(fiber.StatusOK)
}

// ConfirmPasswordReset verifies token and updates password
func (Server) ConfirmPasswordReset(c *fiber.Ctx) error {
	var input PasswordResetConfirmRequest
	if err := c.BodyParser(&input); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	t, err := model.VerifyAndConsumeEmailToken(input.Token, model.TokenPurposeResetPassword)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "invalid or expired token")
	}
	var user model.User
	user.ID = t.UserID
	if err := model.DB.First(&user, user.ID).Error; err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "user not found")
	}
	// update password
	if input.NewPassword == "" {
		return handler.SendError(c, fiber.StatusBadRequest, "newPassword is required")
	}
	params := model.CreateUserParams{Email: user.Email, Password: &input.NewPassword, Name: deref(user.Name)}
	if errs := params.Validate(); len(errs) > 0 {
		return handler.SendError(c, fiber.StatusBadRequest, "password does not meet requirements")
	}
	hashed, err := model.HashPassword(input.NewPassword)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to hash password")
	}
	user.Password = &hashed
	if err := model.DB.Save(&user).Error; err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to update password")
	}
	return c.SendStatus(fiber.StatusOK)
}

// RequestMagicLink creates a magic link login token and emails it
func (Server) RequestMagicLink(c *fiber.Ctx) error {
	var input MagicLinkRequest
	if err := c.BodyParser(&input); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	emailStr, err := getEmail(input.Email)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	user := model.User{Email: emailStr}
	if err := user.GetByEmail(); err == nil {
		token, _, err := model.CreateEmailToken(user.ID, model.TokenPurposeMagicLink, 15*time.Minute)
		if err == nil {
			baseURL := c.BaseURL()
			actionURL := fmt.Sprintf("%s/api/auth/magic-link/consume?token=%s", baseURL, token)
			go func(u model.User, url string) {
				sender := email.NewSMTPSender()
				svc := email.NewService(sender, "Vault Hub")
				name := ""
				if u.Name != nil {
					name = *u.Name
				}
				if err := svc.SendMagicLink(u.Email, name, url, "15m"); err != nil {
					slog.Error("Failed to send magic link email", "error", err, "email", u.Email)
				}
			}(user, actionURL)
		}
	}
	return c.SendStatus(fiber.StatusOK)
}

// ConsumeMagicLink verifies token, generates JWT and redirects with fragment
func (Server) ConsumeMagicLink(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	t, err := model.VerifyAndConsumeEmailToken(token, model.TokenPurposeMagicLink)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var user model.User
	user.ID = t.UserID
	if err := model.DB.First(&user, user.ID).Error; err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	jwtToken, err := user.GenerateToken()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	redirectUrl := "/login#token=" + url.QueryEscape(jwtToken) + "&source=magic"
	return c.Redirect(redirectUrl)
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
