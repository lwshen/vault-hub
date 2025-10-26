package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/internal/email"
	"github.com/lwshen/vault-hub/model"
	"github.com/lwshen/vault-hub/packages/api/generated_models"
)

// LoginEcho handles POST /api/auth/login
func (c *Container) Login(ctx echo.Context) error {
	var input generated_models.LoginRequest
	if err := ctx.Bind(&input); err != nil {
		return SendError(ctx, http.StatusBadRequest, err.Error())
	}

	clientIP, userAgent := getClientInfoEcho(ctx)

	user := model.User{
		Email: input.Email,
	}
	if err := user.GetByEmail(); err != nil {
		return SendError(ctx, http.StatusBadRequest, "Invalid email or password")
	}

	if !user.ComparePassword(input.Password) {
		return SendError(ctx, http.StatusBadRequest, "Invalid email or password")
	}

	token, err := user.GenerateToken()
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, err.Error())
	}

	// Record successful login audit log
	if err := model.LogUserAction(model.ActionLoginUser, user.ID, model.SourceWeb, clientIP, userAgent); err != nil {
		slog.Error("Failed to create audit log for login", "error", err, "userID", user.ID)
	}

	resp := generated_models.LoginResponse{
		Token: token,
	}

	return ctx.JSON(http.StatusOK, resp)
}

// SignupEcho handles POST /api/auth/signup
func (c *Container) Signup(ctx echo.Context) error {
	var input generated_models.SignupRequest
	if err := ctx.Bind(&input); err != nil {
		return SendError(ctx, http.StatusBadRequest, err.Error())
	}

	// Extract client information for audit logging
	clientIP, userAgent := getClientInfoEcho(ctx)

	// Validate and create user parameters
	createParams, err := buildUserCreateParamsEcho(input)
	if err != nil {
		return SendError(ctx, http.StatusBadRequest, err.Error())
	}

	// Create the user account
	user, err := createParams.Create()
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, err.Error())
	}

	slog.Info("User created", "email", user.Email, "name", *user.Name)

	// Log successful registration
	if err := model.LogUserAction(model.ActionRegisterUser, user.ID, model.SourceWeb, clientIP, userAgent); err != nil {
		slog.Error("Failed to create audit log for signup", "error", err, "userID", user.ID)
	}

	// Fire-and-forget signup email (do not block response)
	go func(u *model.User) {
		sender := email.NewSender()
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
		return SendError(ctx, http.StatusInternalServerError, err.Error())
	}

	resp := generated_models.SignupResponse{
		Token: token,
	}

	return ctx.JSON(http.StatusOK, resp)
}

// LogoutEcho handles GET /api/auth/logout
func (c *Container) Logout(ctx echo.Context) error {
	// Try to get user information from context (set by JWT middleware)
	user, ok := ctx.Get("user").(*model.User)
	if ok && user != nil {
		clientIP, userAgent := getClientInfoEcho(ctx)
		if err := model.LogUserAction(model.ActionLogoutUser, user.ID, model.SourceWeb, clientIP, userAgent); err != nil {
			slog.Error("Failed to create audit log for logout", "error", err, "userID", user.ID)
		}
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Successfully logged out",
	})
}

// buildUserCreateParamsEcho creates and validates user creation parameters for Echo
func buildUserCreateParamsEcho(input generated_models.SignupRequest) (model.CreateUserParams, error) {
	// For email/password signup, password is required
	if input.Password == "" {
		return model.CreateUserParams{}, fmt.Errorf("password is required")
	}

	createUserParams := model.CreateUserParams{
		Email:    input.Email,
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

// Stub implementations for other auth endpoints (to be implemented later)

func (c *Container) RequestPasswordReset(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Password reset not yet implemented",
	})
}

func (c *Container) ConfirmPasswordReset(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Password reset confirm not yet implemented",
	})
}

func (c *Container) RequestMagicLink(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Magic link not yet implemented",
	})
}

func (c *Container) ConsumeMagicLink(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Magic link consume not yet implemented",
	})
}
