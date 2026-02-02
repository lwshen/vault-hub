package model

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"gorm.io/gorm"
)

// RefreshToken represents a JWT refresh token stored in the database
type RefreshToken struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint   `gorm:"index;not null"`
	User      User   `gorm:"foreignKey:UserID"`
	TokenHash string `gorm:"size:64;index;not null"`
	ExpiresAt time.Time
	RevokedAt *time.Time
}

// RefreshTokenTTL is the default refresh token expiration (7 days)
const RefreshTokenTTL = 7 * 24 * time.Hour

// generateRefreshToken creates a new refresh token and its hash
func generateRefreshToken() (string, string, error) {
	// 32 random bytes, base64-url encode (no padding)
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	token := base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(token))
	hash := base64.RawURLEncoding.EncodeToString(sum[:])
	return token, hash, nil
}

// CreateRefreshToken creates a new refresh token for a user
func CreateRefreshToken(userID uint) (string, *RefreshToken, error) {
	token, hash, err := generateRefreshToken()
	if err != nil {
		return "", nil, err
	}
	rt := &RefreshToken{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(RefreshTokenTTL),
	}
	if err := DB.Create(rt).Error; err != nil {
		return "", nil, err
	}
	return token, rt, nil
}

// ValidateRefreshToken validates a refresh token and returns the user ID
// Returns (userID, nil) if valid, (0, error) if invalid
func ValidateRefreshToken(plaintextToken string) (uint, error) {
	sum := sha256.Sum256([]byte(plaintextToken))
	hash := base64.RawURLEncoding.EncodeToString(sum[:])

	var rt RefreshToken
	err := DB.
		Where("token_hash = ? AND revoked_at IS NULL AND expires_at >= ?", hash, time.Now()).
		First(&rt).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, errors.New("invalid or expired refresh token")
	}
	if err != nil {
		return 0, err
	}
	return rt.UserID, nil
}

// RevokeRefreshToken revokes a refresh token
func RevokeRefreshToken(plaintextToken string) error {
	sum := sha256.Sum256([]byte(plaintextToken))
	hash := base64.RawURLEncoding.EncodeToString(sum[:])

	now := time.Now()
	result := DB.Model(&RefreshToken{}).
		Where("token_hash = ? AND revoked_at IS NULL", hash).
		Updates(map[string]interface{}{
			"revoked_at": now,
			"updated_at": now,
		})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// RevokeAllUserRefreshTokens revokes all refresh tokens for a user (logout all devices)
func RevokeAllUserRefreshTokens(userID uint) error {
	now := time.Now()
	return DB.Model(&RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Updates(map[string]interface{}{
			"revoked_at": now,
			"updated_at": now,
		}).Error
}

// CleanupExpiredRefreshTokens removes expired refresh tokens older than the given duration
func CleanupExpiredRefreshTokens(olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	return DB.Where("expires_at < ?", cutoff).Unscoped().Delete(&RefreshToken{}).Error
}
