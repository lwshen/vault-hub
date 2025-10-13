package api

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/internal/config"
	"github.com/lwshen/vault-hub/internal/email"
	"github.com/lwshen/vault-hub/model"
)

// RequestPasswordReset handles password reset requests
func (Server) RequestPasswordReset(c *fiber.Ctx) error {
	var input RequestPasswordResetRequest
	if err := c.BodyParser(&input); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	emailStr, err := getEmail(input.Email)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Always return success to prevent user enumeration
	// Find user by email
	user := model.User{Email: emailStr}
	if err := user.GetByEmail(); err != nil {
		// User doesn't exist, but we still return success for security
		slog.Debug("Password reset requested for non-existent email", "email", emailStr)
		return c.Status(fiber.StatusOK).JSON(RequestPasswordResetResponse{
			Message: "If an account with that email exists, a password reset link has been sent",
		})
	}

	// Delete any existing password reset tokens for this user
	if err := model.DeleteUserTokensByType(user.ID, model.TokenTypePasswordReset); err != nil {
		slog.Error("Failed to delete existing password reset tokens", "error", err, "userID", user.ID)
	}

	// Create password reset token
	token, err := model.CreatePasswordResetToken(user.ID)
	if err != nil {
		slog.Error("Failed to create password reset token", "error", err, "userID", user.ID)
		return c.Status(fiber.StatusOK).JSON(RequestPasswordResetResponse{
			Message: "If an account with that email exists, a password reset link has been sent",
		})
	}

	// Build reset URL
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", config.BaseURL, token.Token)

	// Send password reset email
	recipientName := "User"
	if user.Name != nil {
		recipientName = *user.Name
	}

	if err := email.SendPasswordResetAsync(user.Email, recipientName, resetURL); err != nil {
		slog.Error("Failed to queue password reset email", "error", err, "email", user.Email)
	}

	slog.Info("Password reset email queued", "email", user.Email, "userID", user.ID)

	return c.Status(fiber.StatusOK).JSON(RequestPasswordResetResponse{
		Message: "If an account with that email exists, a password reset link has been sent",
	})
}

// ResetPassword handles password reset with token
func (Server) ResetPassword(c *fiber.Ctx) error {
	var input ResetPasswordRequest
	if err := c.BodyParser(&input); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Validate password
	if ok, msg := model.IsPasswordValid(input.NewPassword); !ok {
		return handler.SendError(c, fiber.StatusBadRequest, msg)
	}

	// Find and validate token
	token, err := model.FindValidToken(input.Token, model.TokenTypePasswordReset)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "Invalid or expired reset token")
	}

	// Update user password
	if err := token.User.UpdatePassword(input.NewPassword); err != nil {
		slog.Error("Failed to update password", "error", err, "userID", token.UserID)
		return handler.SendError(c, fiber.StatusInternalServerError, "Failed to reset password")
	}

	// Mark token as used
	if err := token.MarkAsUsed(); err != nil {
		slog.Error("Failed to mark token as used", "error", err, "tokenID", token.ID)
	}

	slog.Info("Password reset successful", "userID", token.UserID)

	return c.Status(fiber.StatusOK).JSON(ResetPasswordResponse{
		Message: "Password has been reset successfully",
	})
}

// RequestMagicLink handles magic link login requests
func (Server) RequestMagicLink(c *fiber.Ctx) error {
	var input RequestMagicLinkRequest
	if err := c.BodyParser(&input); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	emailStr, err := getEmail(input.Email)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Always return success to prevent user enumeration
	user := model.User{Email: emailStr}
	if err := user.GetByEmail(); err != nil {
		slog.Debug("Magic link requested for non-existent email", "email", emailStr)
		return c.Status(fiber.StatusOK).JSON(RequestMagicLinkResponse{
			Message: "If an account with that email exists, a login link has been sent",
		})
	}

	// Delete any existing magic link tokens for this user
	if err := model.DeleteUserTokensByType(user.ID, model.TokenTypeMagicLink); err != nil {
		slog.Error("Failed to delete existing magic link tokens", "error", err, "userID", user.ID)
	}

	// Create magic link token
	token, err := model.CreateMagicLinkToken(user.ID)
	if err != nil {
		slog.Error("Failed to create magic link token", "error", err, "userID", user.ID)
		return c.Status(fiber.StatusOK).JSON(RequestMagicLinkResponse{
			Message: "If an account with that email exists, a login link has been sent",
		})
	}

	// Build magic link URL
	magicLinkURL := fmt.Sprintf("%s/api/auth/magic-link-callback?token=%s", config.BaseURL, token.Token)

	// Send magic link email
	recipientName := "User"
	if user.Name != nil {
		recipientName = *user.Name
	}

	if err := email.SendMagicLinkAsync(user.Email, recipientName, magicLinkURL); err != nil {
		slog.Error("Failed to queue magic link email", "error", err, "email", user.Email)
	}

	slog.Info("Magic link email queued", "email", user.Email, "userID", user.ID)

	return c.Status(fiber.StatusOK).JSON(RequestMagicLinkResponse{
		Message: "If an account with that email exists, a login link has been sent",
	})
}

// MagicLinkCallback handles magic link callback
func (Server) MagicLinkCallback(c *fiber.Ctx, params MagicLinkCallbackParams) error {
	// Find and validate token
	token, err := model.FindValidToken(params.Token, model.TokenTypeMagicLink)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "Invalid or expired magic link")
	}

	// Mark token as used
	if err := token.MarkAsUsed(); err != nil {
		slog.Error("Failed to mark token as used", "error", err, "tokenID", token.ID)
	}

	// Generate JWT token
	jwtToken, err := token.User.GenerateToken()
	if err != nil {
		slog.Error("Failed to generate JWT token", "error", err, "userID", token.UserID)
		return handler.SendError(c, fiber.StatusInternalServerError, "Failed to generate authentication token")
	}

	// Record successful login audit log
	clientIP, userAgent := getClientInfo(c)
	if err := model.LogUserAction(model.ActionLoginUser, token.UserID, model.SourceWeb, clientIP, userAgent); err != nil {
		slog.Error("Failed to create audit log for magic link login", "error", err, "userID", token.UserID)
	}

	slog.Info("Magic link login successful", "userID", token.UserID)

	return c.Status(fiber.StatusOK).JSON(MagicLinkCallbackResponse{
		Token: jwtToken,
	})
}

// VerifyEmail handles email verification
func (Server) VerifyEmail(c *fiber.Ctx, params VerifyEmailParams) error {
	// Find and validate token
	token, err := model.FindValidToken(params.Token, model.TokenTypeEmailVerify)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "Invalid or expired verification token")
	}

	// Mark email as verified
	if err := token.User.MarkEmailAsVerified(); err != nil {
		slog.Error("Failed to mark email as verified", "error", err, "userID", token.UserID)
		return handler.SendError(c, fiber.StatusInternalServerError, "Failed to verify email")
	}

	// Mark token as used
	if err := token.MarkAsUsed(); err != nil {
		slog.Error("Failed to mark token as used", "error", err, "tokenID", token.ID)
	}

	// Generate JWT token for automatic login
	jwtToken, err := token.User.GenerateToken()
	if err != nil {
		slog.Error("Failed to generate JWT token", "error", err, "userID", token.UserID)
		return handler.SendError(c, fiber.StatusInternalServerError, "Failed to generate authentication token")
	}

	slog.Info("Email verified successfully", "userID", token.UserID, "email", token.User.Email)

	return c.Status(fiber.StatusOK).JSON(VerifyEmailResponse{
		Message: "Email verified successfully",
		Token:   jwtToken,
	})
}

// ResendVerification handles resending verification email
func (Server) ResendVerification(c *fiber.Ctx) error {
	// Get user from context (set by JWT middleware)
	user, ok := c.Locals("user").(*model.User)
	if !ok || user == nil {
		return handler.SendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	// Check if email is already verified
	if user.EmailVerified {
		return handler.SendError(c, fiber.StatusBadRequest, "Email is already verified")
	}

	// Delete any existing email verification tokens for this user
	if err := model.DeleteUserTokensByType(user.ID, model.TokenTypeEmailVerify); err != nil {
		slog.Error("Failed to delete existing email verification tokens", "error", err, "userID", user.ID)
	}

	// Create new verification token
	token, err := model.CreateEmailVerificationToken(user.ID)
	if err != nil {
		slog.Error("Failed to create email verification token", "error", err, "userID", user.ID)
		return handler.SendError(c, fiber.StatusInternalServerError, "Failed to send verification email")
	}

	// Build verification URL
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", config.BaseURL, token.Token)

	// Send verification email
	recipientName := "User"
	if user.Name != nil {
		recipientName = *user.Name
	}

	if err := email.SendEmailVerificationAsync(user.Email, recipientName, verificationURL); err != nil {
		slog.Error("Failed to queue verification email", "error", err, "email", user.Email)
		return handler.SendError(c, fiber.StatusInternalServerError, "Failed to send verification email")
	}

	slog.Info("Verification email resent", "userID", user.ID, "email", user.Email)

	return c.Status(fiber.StatusOK).JSON(ResendVerificationResponse{
		Message: "Verification email has been sent",
	})
}
