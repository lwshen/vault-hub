package model

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"gorm.io/gorm"
)

type TokenPurpose string

const (
	TokenPurposeVerifyEmail   TokenPurpose = "verify_email"
	TokenPurposeResetPassword TokenPurpose = "reset_password"
	TokenPurposeMagicLink     TokenPurpose = "magic_link"
)

type EmailToken struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserID     uint         `gorm:"index;constraint:OnDelete:CASCADE"`
	User       User         `gorm:"foreignKey:UserID"`
	TokenHash  string       `gorm:"size:64;index"`
	Purpose    TokenPurpose `gorm:"size:32;index"`
	ExpiresAt  time.Time    `gorm:"index"`
	ConsumedAt *time.Time   `gorm:"index"`
}

func generateToken() (string, string, error) {
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

func CreateEmailToken(userID uint, purpose TokenPurpose, ttl time.Duration) (string, *EmailToken, error) {
	token, hash, err := generateToken()
	if err != nil {
		return "", nil, err
	}
	t := &EmailToken{
		UserID:    userID,
		TokenHash: hash,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(ttl),
	}
	if err := DB.Create(t).Error; err != nil {
		return "", nil, err
	}
	return token, t, nil
}

func VerifyAndConsumeEmailToken(plaintextToken string, purpose TokenPurpose) (*EmailToken, error) {
	sum := sha256.Sum256([]byte(plaintextToken))
	hash := base64.RawURLEncoding.EncodeToString(sum[:])

	now := time.Now()
	update := DB.Model(&EmailToken{}).
		Where("token_hash = ? AND purpose = ? AND consumed_at IS NULL AND expires_at >= ?", hash, purpose, now).
		Updates(map[string]interface{}{
			"consumed_at": now,
			"updated_at":  now,
		})
	if update.Error != nil {
		return nil, update.Error
	}

	var t EmailToken
	if update.RowsAffected == 0 {
		if err := DB.Where("token_hash = ? AND purpose = ?", hash, purpose).First(&t).Error; err != nil {
			return nil, err
		}
		if t.ConsumedAt != nil {
			return nil, errors.New("token already used")
		}
		if now.After(t.ExpiresAt) {
			return nil, errors.New("token expired")
		}
		return nil, errors.New("unable to consume token")
	}

	if err := DB.Where("token_hash = ? AND purpose = ?", hash, purpose).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

// EmailTokenRateLimited reports whether a user has recently requested a token for the given purpose.
// If the most recent token was created within the provided window, the call returns true along with
// the remaining cooldown. A zero or negative window disables rate limiting.
func EmailTokenRateLimited(userID uint, purpose TokenPurpose, window time.Duration) (bool, time.Duration, error) {
	if window <= 0 {
		return false, 0, nil
	}

	var token EmailToken
	err := DB.
		Where("user_id = ? AND purpose = ?", userID, purpose).
		Order("created_at DESC").
		First(&token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}

	elapsed := time.Since(token.CreatedAt)
	if elapsed < window {
		return true, window - elapsed, nil
	}
	return false, 0, nil
}
