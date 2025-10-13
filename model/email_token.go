package model

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// TokenType represents the type of email token
type TokenType string

const (
	TokenTypeEmailVerification TokenType = "email_verification"
	TokenTypePasswordReset     TokenType = "password_reset"
	TokenTypeMagicLink         TokenType = "magic_link"
)

// EmailToken represents a token for email-based operations
type EmailToken struct {
	gorm.Model
	UserID    uint      `gorm:"index;not null"`
	Token     string    `gorm:"uniqueIndex;not null;size:255"`
	TokenType TokenType `gorm:"type:varchar(50);not null;index"`
	ExpiresAt time.Time `gorm:"not null;index"`
	UsedAt    *time.Time
	User      User `gorm:"foreignKey:UserID"`
}

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// CreateEmailVerificationToken creates a new email verification token for the user
func CreateEmailVerificationToken(userID uint) (*EmailToken, error) {
	token, err := GenerateSecureToken()
	if err != nil {
		return nil, err
	}

	emailToken := &EmailToken{
		UserID:    userID,
		Token:     token,
		TokenType: TokenTypeEmailVerification,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours expiry
	}

	if err := DB.Create(emailToken).Error; err != nil {
		return nil, fmt.Errorf("failed to create email verification token: %w", err)
	}

	return emailToken, nil
}

// CreatePasswordResetToken creates a new password reset token for the user
func CreatePasswordResetToken(userID uint) (*EmailToken, error) {
	token, err := GenerateSecureToken()
	if err != nil {
		return nil, err
	}

	emailToken := &EmailToken{
		UserID:    userID,
		Token:     token,
		TokenType: TokenTypePasswordReset,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
	}

	if err := DB.Create(emailToken).Error; err != nil {
		return nil, fmt.Errorf("failed to create password reset token: %w", err)
	}

	return emailToken, nil
}

// CreateMagicLinkToken creates a new magic link token for the user
func CreateMagicLinkToken(userID uint) (*EmailToken, error) {
	token, err := GenerateSecureToken()
	if err != nil {
		return nil, err
	}

	emailToken := &EmailToken{
		UserID:    userID,
		Token:     token,
		TokenType: TokenTypeMagicLink,
		ExpiresAt: time.Now().Add(15 * time.Minute), // 15 minutes expiry
	}

	if err := DB.Create(emailToken).Error; err != nil {
		return nil, fmt.Errorf("failed to create magic link token: %w", err)
	}

	return emailToken, nil
}

// FindValidToken finds a valid (not used, not expired) token
func FindValidToken(token string, tokenType TokenType) (*EmailToken, error) {
	var emailToken EmailToken
	err := DB.Where("token = ? AND token_type = ? AND expires_at > ? AND used_at IS NULL",
		token, tokenType, time.Now()).
		Preload("User").
		First(&emailToken).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid or expired token")
		}
		return nil, fmt.Errorf("failed to find token: %w", err)
	}

	return &emailToken, nil
}

// MarkAsUsed marks the token as used
func (t *EmailToken) MarkAsUsed() error {
	now := time.Now()
	t.UsedAt = &now

	if err := DB.Save(t).Error; err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	return nil
}

// IsValid checks if the token is still valid (not used and not expired)
func (t *EmailToken) IsValid() bool {
	return t.UsedAt == nil && t.ExpiresAt.After(time.Now())
}

// DeleteExpiredTokens removes all expired tokens from the database
// This should be called periodically via a cron job
func DeleteExpiredTokens() error {
	result := DB.Where("expires_at < ?", time.Now()).Delete(&EmailToken{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete expired tokens: %w", result.Error)
	}
	return nil
}

// DeleteUserTokensByType deletes all tokens of a specific type for a user
func DeleteUserTokensByType(userID uint, tokenType TokenType) error {
	result := DB.Where("user_id = ? AND token_type = ?", userID, tokenType).Delete(&EmailToken{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user tokens: %w", result.Error)
	}
	return nil
}
