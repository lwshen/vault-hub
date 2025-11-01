package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/internal/email"
	"github.com/lwshen/vault-hub/model"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	PasswordResetTTL  = 30 * time.Minute
	MagicLinkTTL      = 15 * time.Minute
	EmailSendCooldown = time.Minute

	emailTokenCodeSent        = "email_token_sent"
	emailTokenCodeRateLimited = "email_token_rate_limited"
	emailTokenCodeFailed      = "email_token_failed"
)

// formatTTLForEmail formats a duration for email display (e.g., "30m", "2h")
func formatTTLForEmail(d time.Duration) string {
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%.0fm", d.Minutes())
	default:
		return fmt.Sprintf("%.1fh", d.Hours())
	}
}

// EmailTokenOutcome captures the standard JSON payload used for email-token
// based flows (password reset, magic link) so both Fiber and Echo handlers can
// share the same business logic.
type EmailTokenOutcome struct {
	Success          bool
	Code             string
	Status           int
	RetryAfterHeader *string
}

// LoginWithPassword validates credentials, emits audit logs, and issues a JWT
// for the authenticated user.
func LoginWithPassword(emailAddr openapi_types.Email, password string, clientInfo ClientInfo) (LoginResponse, *APIError) {
	emailStr, err := getEmail(emailAddr)
	if err != nil {
		return LoginResponse{}, newAPIError(http.StatusBadRequest, err.Error())
	}

	user := model.User{Email: emailStr}
	if err := user.GetByEmail(); err != nil {
		return LoginResponse{}, newAPIError(http.StatusBadRequest, "Invalid email or password")
	}

	if !user.ComparePassword(password) {
		return LoginResponse{}, newAPIError(http.StatusBadRequest, "Invalid email or password")
	}

	token, err := user.GenerateToken()
	if err != nil {
		return LoginResponse{}, newAPIError(http.StatusInternalServerError, err.Error())
	}

	if err := model.LogUserAction(model.ActionLoginUser, user.ID, model.SourceWeb, clientInfo.IP, clientInfo.UserAgent); err != nil {
		slog.Error("Failed to create audit log for login", "error", err, "userID", user.ID)
	}

	return LoginResponse{Token: token}, nil
}

// SignupWithPassword creates a new user and returns the bootstrap JWT.
func SignupWithPassword(input SignupRequest, clientInfo ClientInfo) (SignupResponse, *APIError) {
	createParams, err := buildUserCreateParams(input)
	if err != nil {
		return SignupResponse{}, newAPIError(http.StatusBadRequest, err.Error())
	}

	user, err := createUser(createParams)
	if err != nil {
		return SignupResponse{}, newAPIError(http.StatusInternalServerError, err.Error())
	}

	slog.Info("User created", "email", user.Email, "name", deref(user.Name))

	logSignupAudit(user.ID, clientInfo.IP, clientInfo.UserAgent)

	clientIP := clientInfo.IP
	userAgent := clientInfo.UserAgent
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

	token, err := user.GenerateToken()
	if err != nil {
		return SignupResponse{}, newAPIError(http.StatusInternalServerError, err.Error())
	}

	return SignupResponse{Token: token}, nil
}

// RecordLogoutAudit captures logout events for auditing when a principal is present.
func RecordLogoutAudit(user *model.User, clientInfo ClientInfo) {
	if user == nil {
		return
	}
	if err := model.LogUserAction(model.ActionLogoutUser, user.ID, model.SourceWeb, clientInfo.IP, clientInfo.UserAgent); err != nil {
		slog.Error("Failed to create audit log for logout", "error", err, "userID", user.ID)
	}
}

// RequestPasswordResetEmail issues a password-reset token and schedules email delivery.
func RequestPasswordResetEmail(emailAddr openapi_types.Email, baseURL string) (EmailTokenOutcome, *APIError) {
	emailStr, err := getEmail(emailAddr)
	if err != nil {
		return EmailTokenOutcome{}, newAPIError(http.StatusBadRequest, err.Error())
	}

	outcome := EmailTokenOutcome{
		Success: true,
		Code:    emailTokenCodeSent,
		Status:  http.StatusOK,
	}

	user := model.User{Email: emailStr}
	if err := user.GetByEmail(); err != nil {
		return outcome, nil
	}

	limited, retryAfter, rateErr := model.EmailTokenRateLimited(user.ID, model.TokenPurposeResetPassword, EmailSendCooldown)
	if rateErr != nil {
		slog.Error("Failed to check password reset rate limit", "error", rateErr, "userID", user.ID)
	} else if limited {
		retryAfterHeader := fmt.Sprintf("%.0f", retryAfter.Seconds())
		return EmailTokenOutcome{
			Success:          false,
			Code:             emailTokenCodeRateLimited,
			Status:           http.StatusTooManyRequests,
			RetryAfterHeader: &retryAfterHeader,
		}, nil
	}

	token, _, tokenErr := model.CreateEmailToken(user.ID, model.TokenPurposeResetPassword, PasswordResetTTL)
	if tokenErr != nil {
		slog.Error("Failed to create password reset token", "error", tokenErr, "userID", user.ID)
		return EmailTokenOutcome{
			Success: false,
			Code:    emailTokenCodeFailed,
			Status:  http.StatusInternalServerError,
		}, nil
	}

	actionURL := fmt.Sprintf("%s/reset?token=%s", baseURL, url.QueryEscape(token))
	go func(u model.User, url string) {
		sender := email.NewSender()
		svc := email.NewService(sender, "Vault Hub")
		name := ""
		if u.Name != nil {
			name = *u.Name
		}
		if err := svc.SendPasswordReset(u.Email, name, url, formatTTLForEmail(PasswordResetTTL)); err != nil {
			slog.Error("Failed to send password reset email", "error", err, "email", u.Email)
		}
	}(user, actionURL)

	return outcome, nil
}

// ConfirmPasswordResetToken updates the user's password after token verification.
func ConfirmPasswordResetToken(token string, newPassword string) *APIError {
	t, err := model.VerifyAndConsumeEmailToken(token, model.TokenPurposeResetPassword)
	if err != nil {
		return newAPIError(http.StatusBadRequest, "invalid or expired token")
	}

	var user model.User
	user.ID = t.UserID
	if err := model.DB.First(&user, user.ID).Error; err != nil {
		return newAPIError(http.StatusInternalServerError, "user not found")
	}

	if newPassword == "" {
		return newAPIError(http.StatusBadRequest, "newPassword is required")
	}

	params := model.CreateUserParams{Email: user.Email, Password: &newPassword, Name: deref(user.Name)}
	if errs := params.Validate(); len(errs) > 0 {
		return newAPIError(http.StatusBadRequest, "password does not meet requirements")
	}

	hashed, err := model.HashPassword(newPassword)
	if err != nil {
		return newAPIError(http.StatusInternalServerError, "failed to hash password")
	}

	user.Password = &hashed
	if err := model.DB.Save(&user).Error; err != nil {
		return newAPIError(http.StatusInternalServerError, "failed to update password")
	}

	return nil
}

// RequestMagicLinkEmail issues a magic-link email when the user exists.
func RequestMagicLinkEmail(emailAddr openapi_types.Email, baseURL string) (EmailTokenOutcome, *APIError) {
	emailStr, err := getEmail(emailAddr)
	if err != nil {
		return EmailTokenOutcome{}, newAPIError(http.StatusBadRequest, err.Error())
	}

	user := model.User{Email: emailStr}
	if err := user.GetByEmail(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EmailTokenOutcome{
				Success: false,
				Code:    emailTokenCodeFailed,
				Status:  http.StatusOK,
			}, nil
		}
		slog.Error("Failed to look up user for magic link", "email", emailStr, "error", err)
		return EmailTokenOutcome{}, newAPIError(http.StatusInternalServerError, "Unable to send a magic link right now.")
	}

	limited, retryAfter, rateErr := model.EmailTokenRateLimited(user.ID, model.TokenPurposeMagicLink, EmailSendCooldown)
	if rateErr != nil {
		slog.Error("Failed to check magic link rate limit", "error", rateErr, "userID", user.ID)
	} else if limited {
		retryAfterHeader := fmt.Sprintf("%.0f", retryAfter.Seconds())
		return EmailTokenOutcome{
			Success:          false,
			Code:             emailTokenCodeRateLimited,
			Status:           http.StatusTooManyRequests,
			RetryAfterHeader: &retryAfterHeader,
		}, nil
	}

	token, _, tokenErr := model.CreateEmailToken(user.ID, model.TokenPurposeMagicLink, MagicLinkTTL)
	if tokenErr != nil {
		slog.Error("Failed to create magic link token", "error", tokenErr, "userID", user.ID)
		return EmailTokenOutcome{
			Success: false,
			Code:    emailTokenCodeFailed,
			Status:  http.StatusInternalServerError,
		}, nil
	}

	actionURL := fmt.Sprintf("%s/login/magic-link?token=%s", baseURL, url.QueryEscape(token))
	go func(u model.User, url string) {
		sender := email.NewSender()
		svc := email.NewService(sender, "Vault Hub")
		name := ""
		if u.Name != nil {
			name = *u.Name
		}
		if err := svc.SendMagicLink(u.Email, name, url, formatTTLForEmail(MagicLinkTTL)); err != nil {
			slog.Error("Failed to send magic link email", "error", err, "email", u.Email)
		}
	}(user, actionURL)

	return EmailTokenOutcome{
		Success: true,
		Code:    emailTokenCodeSent,
		Status:  http.StatusOK,
	}, nil
}

// ConsumeMagicLinkToken validates the provided token and returns a JWT when successful.
func ConsumeMagicLinkToken(token string) (string, *APIError) {
	if token == "" {
		return "", newAPIError(http.StatusBadRequest, "missing token")
	}

	t, err := model.VerifyAndConsumeEmailToken(token, model.TokenPurposeMagicLink)
	if err != nil {
		return "", newAPIError(http.StatusBadRequest, "invalid or expired token")
	}

	var user model.User
	user.ID = t.UserID
	if err := model.DB.First(&user, user.ID).Error; err != nil {
		return "", newAPIError(http.StatusInternalServerError, "user not found")
	}

	jwToken, err := user.GenerateToken()
	if err != nil {
		return "", newAPIError(http.StatusInternalServerError, "failed to generate token")
	}

	return jwToken, nil
}

func (Server) Login(c *fiber.Ctx) error {
	var input LoginRequest
	if err := c.BodyParser(&input); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	clientInfo := getClientInfoDetails(c)

	resp, apiErr := LoginWithPassword(input.Email, input.Password, clientInfo)
	if apiErr != nil {
		return handler.SendError(c, apiErr.Status, apiErr.Message)
	}

	return c.Status(http.StatusOK).JSON(resp)
}

// Signup handles user registration requests
// It validates input, creates the user account, and returns a JWT token
func (Server) Signup(c *fiber.Ctx) error {
	// Parse request body
	input, err := parseSignupRequest(c)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	clientInfo := getClientInfoDetails(c)

	resp, apiErr := SignupWithPassword(input, clientInfo)
	if apiErr != nil {
		return handler.SendError(c, apiErr.Status, apiErr.Message)
	}

	return c.Status(http.StatusOK).JSON(resp)
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
		clientInfo := getClientInfoDetails(c)
		RecordLogoutAudit(user, clientInfo)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}

func getEmail(email openapi_types.Email) (string, error) {
	return string(email), nil
}

// RequestPasswordReset creates a password reset token and sends email
func (Server) RequestPasswordReset(c *fiber.Ctx) error {
	var input PasswordResetRequest
	if err := c.BodyParser(&input); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	outcome, apiErr := RequestPasswordResetEmail(input.Email, c.BaseURL())
	if apiErr != nil {
		return handler.SendError(c, apiErr.Status, apiErr.Message)
	}
	if outcome.RetryAfterHeader != nil {
		c.Set(fiber.HeaderRetryAfter, *outcome.RetryAfterHeader)
	}
	return c.Status(outcome.Status).JSON(fiber.Map{
		"success": outcome.Success,
		"code":    outcome.Code,
	})
}

// ConfirmPasswordReset verifies token and updates password
func (Server) ConfirmPasswordReset(c *fiber.Ctx) error {
	var input PasswordResetConfirmRequest
	if err := c.BodyParser(&input); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	if apiErr := ConfirmPasswordResetToken(input.Token, input.NewPassword); apiErr != nil {
		return handler.SendError(c, apiErr.Status, apiErr.Message)
	}
	return c.SendStatus(http.StatusOK)
}

// RequestMagicLink creates a magic link login token and emails it
func (Server) RequestMagicLink(c *fiber.Ctx) error {
	var input MagicLinkRequest
	if err := c.BodyParser(&input); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	outcome, apiErr := RequestMagicLinkEmail(input.Email, c.BaseURL())
	if apiErr != nil {
		return handler.SendError(c, apiErr.Status, apiErr.Message)
	}
	if outcome.RetryAfterHeader != nil {
		c.Set(fiber.HeaderRetryAfter, *outcome.RetryAfterHeader)
	}
	return c.Status(outcome.Status).JSON(fiber.Map{
		"success": outcome.Success,
		"code":    outcome.Code,
	})
}

// ConsumeMagicLink verifies token, generates JWT and redirects with fragment
func (Server) ConsumeMagicLink(c *fiber.Ctx, params ConsumeMagicLinkParams) error {
	token := params.Token
	acceptsJSON := strings.Contains(c.Get(fiber.HeaderAccept), fiber.MIMEApplicationJSON)
	jwtToken, apiErr := ConsumeMagicLinkToken(token)
	if apiErr != nil {
		if acceptsJSON {
			return c.Status(apiErr.Status).JSON(fiber.Map{
				"error": apiErr.Message,
				"code":  emailTokenCodeFailed,
			})
		}
		return c.SendStatus(apiErr.Status)
	}
	redirectFragment := "/login#token=" + url.QueryEscape(jwtToken) + "&source=magic"

	if acceptsJSON {
		return c.JSON(fiber.Map{
			"token":       jwtToken,
			"redirectUrl": fmt.Sprintf("%s/dashboard", c.BaseURL()),
			"code":        emailTokenCodeSent,
			"success":     true,
		})
	}

	return c.Redirect(redirectFragment)
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
