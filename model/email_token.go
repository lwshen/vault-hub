package model

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"
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

	var t EmailToken
	if err := DB.Where("token_hash = ? AND purpose = ?", hash, purpose).First(&t).Error; err != nil {
		return nil, err
	}
	if t.ConsumedAt != nil {
		return nil, errors.New("token already used")
	}
	if time.Now().After(t.ExpiresAt) {
		return nil, errors.New("token expired")
	}
	now := time.Now()
	t.ConsumedAt = &now
	if err := DB.Save(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
