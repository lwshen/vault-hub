package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/internal/email"
	"github.com/lwshen/vault-hub/model"
	"github.com/lwshen/vault-hub/packages/api/generated/models"
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

// LoginEcho handles POST /api/auth/login
func (c *Container) Login(ctx echo.Context) error {
	var input models.LoginRequest
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

	resp := models.LoginResponse{
		Token: token,
	}

	return ctx.JSON(http.StatusOK, resp)
}

// SignupEcho handles POST /api/auth/signup
func (c *Container) Signup(ctx echo.Context) error {
	var input models.SignupRequest
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

	resp := models.SignupResponse{
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
func buildUserCreateParamsEcho(input models.SignupRequest) (model.CreateUserParams, error) {
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

// RequestPasswordReset creates a password reset token and sends email
func (c *Container) RequestPasswordReset(ctx echo.Context) error {
	var input models.PasswordResetRequest
	if err := ctx.Bind(&input); err != nil {
		return SendError(ctx, http.StatusBadRequest, err.Error())
	}

	// Always respond 200 to avoid user enumeration
	// Attempt to find user; if not found, still return success
	user := model.User{Email: input.Email}
	if err := user.GetByEmail(); err == nil {
		limited, retryAfter, rateErr := model.EmailTokenRateLimited(user.ID, model.TokenPurposeResetPassword, EmailSendCooldown)
		if rateErr != nil {
			slog.Error("Failed to check password reset rate limit", "error", rateErr, "userID", user.ID)
		} else if limited {
			slog.Warn("Password reset email rate limited", "userID", user.ID, "retryAfter", retryAfter)
			ctx.Response().Header().Set(echo.HeaderRetryAfter, fmt.Sprintf("%.0f", retryAfter.Seconds()))
			return ctx.JSON(http.StatusTooManyRequests, models.EmailTokenResponse{
				Success: false,
				Code:    emailTokenCodeRateLimited,
			})
		}

		token, _, err := model.CreateEmailToken(user.ID, model.TokenPurposeResetPassword, PasswordResetTTL)
		if err != nil {
			slog.Error("Failed to create password reset token", "error", err, "userID", user.ID)
			return ctx.JSON(http.StatusInternalServerError, models.EmailTokenResponse{
				Success: false,
				Code:    emailTokenCodeFailed,
			})
		}

		// Construct base URL from request
		scheme := ctx.Scheme()
		host := ctx.Request().Host
		baseURL := fmt.Sprintf("%s://%s", scheme, host)
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
	}

	return ctx.JSON(http.StatusOK, models.EmailTokenResponse{
		Success: true,
		Code:    emailTokenCodeSent,
	})
}

// ConfirmPasswordReset verifies token and updates password
func (c *Container) ConfirmPasswordReset(ctx echo.Context) error {
	var input models.PasswordResetConfirmRequest
	if err := ctx.Bind(&input); err != nil {
		return SendError(ctx, http.StatusBadRequest, err.Error())
	}

	t, err := model.VerifyAndConsumeEmailToken(input.Token, model.TokenPurposeResetPassword)
	if err != nil {
		return SendError(ctx, http.StatusBadRequest, "invalid or expired token")
	}

	var user model.User
	user.ID = t.UserID
	if err := model.DB.First(&user, user.ID).Error; err != nil {
		return SendError(ctx, http.StatusInternalServerError, "user not found")
	}

	// Validate password
	if input.NewPassword == "" {
		return SendError(ctx, http.StatusBadRequest, "newPassword is required")
	}

	// Use model validation to check password requirements
	name := ""
	if user.Name != nil {
		name = *user.Name
	}
	params := model.CreateUserParams{Email: user.Email, Password: &input.NewPassword, Name: name}
	if errs := params.Validate(); len(errs) > 0 {
		return SendError(ctx, http.StatusBadRequest, "password does not meet requirements")
	}

	// Hash and update password
	hashed, err := model.HashPassword(input.NewPassword)
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, "failed to hash password")
	}
	user.Password = &hashed
	if err := model.DB.Save(&user).Error; err != nil {
		return SendError(ctx, http.StatusInternalServerError, "failed to update password")
	}

	return ctx.NoContent(http.StatusOK)
}

// RequestMagicLink creates a magic link login token and emails it
func (c *Container) RequestMagicLink(ctx echo.Context) error {
	var input models.MagicLinkRequest
	if err := ctx.Bind(&input); err != nil {
		return SendError(ctx, http.StatusBadRequest, err.Error())
	}

	user := model.User{Email: input.Email}
	if err := user.GetByEmail(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Info("Magic link request user not found", "email", input.Email)
			return ctx.JSON(http.StatusOK, models.EmailTokenResponse{
				Success: false,
				Code:    emailTokenCodeFailed,
			})
		}
		slog.Error("Failed to look up user for magic link", "email", input.Email, "error", err)
		return SendError(ctx, http.StatusInternalServerError, "Unable to send a magic link right now.")
	}

	limited, retryAfter, rateErr := model.EmailTokenRateLimited(user.ID, model.TokenPurposeMagicLink, EmailSendCooldown)
	if rateErr != nil {
		slog.Error("Failed to check magic link rate limit", "error", rateErr, "userID", user.ID)
	} else if limited {
		slog.Warn("Magic link email rate limited", "userID", user.ID, "retryAfter", retryAfter)
		ctx.Response().Header().Set(echo.HeaderRetryAfter, fmt.Sprintf("%.0f", retryAfter.Seconds()))
		return ctx.JSON(http.StatusTooManyRequests, models.EmailTokenResponse{
			Success: false,
			Code:    emailTokenCodeRateLimited,
		})
	}

	token, _, tokenErr := model.CreateEmailToken(user.ID, model.TokenPurposeMagicLink, MagicLinkTTL)
	if tokenErr != nil {
		slog.Error("Failed to create magic link token", "error", tokenErr, "userID", user.ID)
		return ctx.JSON(http.StatusInternalServerError, models.EmailTokenResponse{
			Success: false,
			Code:    emailTokenCodeFailed,
		})
	}

	// Construct base URL from request
	scheme := ctx.Scheme()
	host := ctx.Request().Host
	baseURL := fmt.Sprintf("%s://%s", scheme, host)
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

	return ctx.JSON(http.StatusOK, models.EmailTokenResponse{
		Success: true,
		Code:    emailTokenCodeSent,
	})
}

// ConsumeMagicLink verifies token, generates JWT and redirects with fragment
func (c *Container) ConsumeMagicLink(ctx echo.Context) error {
	token := ctx.QueryParam("token")
	acceptsJSON := strings.Contains(ctx.Request().Header.Get(echo.HeaderAccept), echo.MIMEApplicationJSON)

	if token == "" {
		if acceptsJSON {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "missing token",
				"code":  emailTokenCodeFailed,
			})
		}
		return ctx.NoContent(http.StatusBadRequest)
	}

	t, err := model.VerifyAndConsumeEmailToken(token, model.TokenPurposeMagicLink)
	if err != nil {
		if acceptsJSON {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid or expired token",
				"code":  emailTokenCodeFailed,
			})
		}
		return ctx.NoContent(http.StatusBadRequest)
	}

	var user model.User
	user.ID = t.UserID
	if err := model.DB.First(&user, user.ID).Error; err != nil {
		if acceptsJSON {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "user not found",
				"code":  emailTokenCodeFailed,
			})
		}
		return ctx.NoContent(http.StatusInternalServerError)
	}

	jwtToken, err := user.GenerateToken()
	if err != nil {
		if acceptsJSON {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to generate token",
				"code":  emailTokenCodeFailed,
			})
		}
		return ctx.NoContent(http.StatusInternalServerError)
	}

	redirectFragment := "/login#token=" + url.QueryEscape(jwtToken) + "&source=magic"

	if acceptsJSON {
		// Construct base URL from request
		scheme := ctx.Scheme()
		host := ctx.Request().Host
		baseURL := fmt.Sprintf("%s://%s", scheme, host)

		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"token":       jwtToken,
			"redirectUrl": fmt.Sprintf("%s/dashboard", baseURL),
			"code":        emailTokenCodeSent,
			"success":     true,
		})
	}

	return ctx.Redirect(http.StatusFound, redirectFragment)
}
