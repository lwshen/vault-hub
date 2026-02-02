package api

import (
	"encoding/json"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/internal/auth"
	"github.com/lwshen/vault-hub/model"
)

// MFASetupResponse represents the MFA setup response
type MFASetupResponse struct {
	Secret        string   `json:"secret"`
	QRCodeURL     string   `json:"qrCodeUrl"`
	RecoveryCodes []string `json:"recoveryCodes"`
}

// MFAVerifyRequest represents the MFA verification request
type MFAVerifyRequest struct {
	Code string `json:"code"`
}

// MFASetupRequest represents the MFA setup request (empty for initiating setup)
type MFASetupRequest struct{}

// GetMFASetup initiates MFA setup for the current user
func (s Server) GetMFASetup(c *fiber.Ctx) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Check if MFA is already enabled
	existingSettings, err := model.GetMFASettings(user.ID)
	if err == nil && existingSettings.Status == model.MFAStatusEnabled {
		return handler.SendError(c, fiber.StatusBadRequest, "MFA is already enabled")
	}

	// Generate a new TOTP secret
	secret, err := auth.GenerateTOTPSecret()
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to generate MFA secret")
	}

	// Generate recovery codes
	recoveryCodes, err := model.GenerateRecoveryCodes(10)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to generate recovery codes")
	}

	// Encrypt the secret for storage
	encryptedSecret, err := model.EncryptMFASecret(secret)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to encrypt MFA secret")
	}

	// Create or update MFA settings
	if existingSettings != nil && existingSettings.Status == model.MFAStatusPending {
		// Update existing pending settings
		err = updatePendingMFASettings(user.ID, encryptedSecret)
	} else {
		// Create new MFA settings
		_, err = model.CreateMFASettings(user.ID, encryptedSecret, model.MFAMethodTOTP)
	}
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to save MFA settings")
	}

	// Generate the TOTP URL for authenticator apps
	totpURL := auth.GenerateTOTPUrl(secret, user.Email, "Vault Hub")

	slog.Info("MFA setup initiated", "userID", user.ID)

	response := MFASetupResponse{
		Secret:        auth.FormatMFASecretForDisplay(secret),
		QRCodeURL:     totpURL,
		RecoveryCodes: recoveryCodes,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// EnableMFA enables MFA after verification
func (s Server) EnableMFA(c *fiber.Ctx) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Parse request body
	var req MFAVerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "invalid request body")
	}

	if req.Code == "" {
		return handler.SendError(c, fiber.StatusBadRequest, "verification code is required")
	}

	// Get MFA settings
	settings, err := model.GetMFASettings(user.ID)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "MFA not set up")
	}

	if settings.Status == model.MFAStatusEnabled {
		return handler.SendError(c, fiber.StatusBadRequest, "MFA is already enabled")
	}

	// Decrypt the secret
	secret, err := model.DecryptMFASecret(settings.Secret)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to decrypt MFA secret")
	}

	// Verify the code
	if !auth.ValidateTOTP(secret, req.Code) {
		return handler.SendError(c, fiber.StatusBadRequest, "invalid verification code")
	}

	// Generate and hash recovery codes
	recoveryCodes, err := model.GenerateRecoveryCodes(10)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to generate recovery codes")
	}

	// Enable MFA
	if err := model.EnableMFA(user.ID, recoveryCodes); err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to enable MFA")
	}

	slog.Info("MFA enabled", "userID", user.ID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "MFA enabled successfully",
	})
}

// DisableMFA disables MFA for the current user
func (s Server) DisableMFA(c *fiber.Ctx) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Parse request body
	var req MFAVerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "invalid request body")
	}

	if req.Code == "" {
		return handler.SendError(c, fiber.StatusBadRequest, "verification code is required")
	}

	// Get MFA settings
	settings, err := model.GetMFASettings(user.ID)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "MFA not set up")
	}

	if settings.Status != model.MFAStatusEnabled {
		return handler.SendError(c, fiber.StatusBadRequest, "MFA is not enabled")
	}

	// Unpack recovery codes
	var hashedCodes []string
	if err := json.Unmarshal([]byte(settings.RecoveryCodes), &hashedCodes); err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to parse recovery codes")
	}

	// Verify with a recovery code first
	var isValid bool
	if len(req.Code) == 16 { // Recovery codes are 16 characters
		isValid = model.VerifyRecoveryCode(hashedCodes, req.Code)
	} else {
		// Verify with TOTP
		secret, err := model.DecryptMFASecret(settings.Secret)
		if err == nil {
			isValid = auth.ValidateTOTP(secret, req.Code)
		}
	}

	if !isValid {
		return handler.SendError(c, fiber.StatusBadRequest, "invalid verification code")
	}

	// Disable MFA
	if err := model.DisableMFA(user.ID); err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to disable MFA")
	}

	slog.Info("MFA disabled", "userID", user.ID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "MFA disabled successfully",
	})
}

// VerifyMFA verifies an MFA code (for use during login)
func (s Server) VerifyMFA(c *fiber.Ctx) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Parse request body
	var req MFAVerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "invalid request body")
	}

	if req.Code == "" {
		return handler.SendError(c, fiber.StatusBadRequest, "verification code is required")
	}

	// Get MFA settings
	settings, err := model.GetMFASettings(user.ID)
	if err != nil {
		return handler.SendError(c, fiber.StatusBadRequest, "MFA not set up")
	}

	if settings.Status != model.MFAStatusEnabled {
		return handler.SendError(c, fiber.StatusBadRequest, "MFA is not enabled")
	}

	// Unpack recovery codes
	var hashedCodes []string
	if err := json.Unmarshal([]byte(settings.RecoveryCodes), &hashedCodes); err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, "failed to parse recovery codes")
	}

	var isValid bool
	if len(req.Code) == 16 { // Recovery codes are 16 characters
		// Try to use recovery code
		isValid = model.VerifyRecoveryCode(hashedCodes, req.Code)
		if isValid {
			// Mark the recovery code as used
			_ = model.UseRecoveryCode(user.ID, req.Code)
		}
	} else {
		// Verify with TOTP
		secret, err := model.DecryptMFASecret(settings.Secret)
		if err == nil {
			isValid = auth.ValidateTOTP(secret, req.Code)
		}
	}

	if !isValid {
		return handler.SendError(c, fiber.StatusBadRequest, "invalid verification code")
	}

	// Update last used timestamp
	_ = model.UpdateMFALastUsed(user.ID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "MFA verification successful",
	})
}

// GetMFAStatus returns the MFA status for the current user
func (s Server) GetMFAStatus(c *fiber.Ctx) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	settings, err := model.GetMFASettings(user.ID)
	if err != nil {
		// MFA not set up
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"enabled": false,
			"status":  "disabled",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"enabled": settings.Status == model.MFAStatusEnabled,
		"status":  string(settings.Status),
		"method":  string(settings.Method),
	})
}

// updatePendingMFASettings updates the secret for pending MFA settings
func updatePendingMFASettings(userID uint, encryptedSecret string) error {
	return model.DB.Model(&model.MFASettings{}).
		Where("user_id = ? AND status = ?", userID, model.MFAStatusPending).
		Update("secret", encryptedSecret).Error
}
