package model

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// MFAStatus represents the MFA status for a user
type MFAStatus string

const (
	MFAStatusDisabled MFAStatus = "disabled"
	MFAStatusPending  MFAStatus = "pending" // MFA setup initiated but not verified
	MFAStatusEnabled  MFAStatus = "enabled"
)

// MFAMethod represents the MFA method
type MFAMethod string

const (
	MFAMethodTOTP MFAMethod = "totp" // Time-based One-Time Password
)

// MFASettings represents the MFA settings for a user
type MFASettings struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	UserID        uint       `gorm:"uniqueIndex;not null"`
	User          User       `gorm:"foreignKey:UserID"`
	Secret        string     `gorm:"type:text"`        // Encrypted TOTP secret
	Method        MFAMethod  `gorm:"size:16;not null"` // MFA method (totp)
	Status        MFAStatus  `gorm:"size:16;not null;default:'disabled'"`
	RecoveryCodes string     `gorm:"type:json"` // JSON array of hashed recovery codes
	LastUsedAt    *time.Time // Last time MFA was used
}

// GenerateMFASecret generates a new random TOTP secret
func GenerateMFASecret() (string, error) {
	// Generate 20 bytes of random data
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base32Encoding.EncodeToString(b), nil
}

// base32Encoding is a base32 encoding without padding
var base32Encoding = base32.StdEncoding.WithPadding(base32.NoPadding)

// GenerateRecoveryCodes generates a set of recovery codes
func GenerateRecoveryCodes(count int) ([]string, error) {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		// Generate 8 random bytes
		b := make([]byte, 8)
		if _, err := rand.Read(b); err != nil {
			return nil, err
		}
		codes[i] = base32Encoding.EncodeToString(b)
	}
	return codes, nil
}

// HashRecoveryCodes hashes recovery codes for storage
func HashRecoveryCodes(codes []string) ([]string, error) {
	hashed := make([]string, len(codes))
	for i, code := range codes {
		hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashed[i] = string(hash)
	}
	return hashed, nil
}

// VerifyRecoveryCode checks if a recovery code is valid
func VerifyRecoveryCode(hashedCodes []string, code string) bool {
	for _, hashed := range hashedCodes {
		if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(code)); err == nil {
			return true
		}
	}
	return false
}

// GetMFASettings retrieves MFA settings for a user
func GetMFASettings(userID uint) (*MFASettings, error) {
	var settings MFASettings
	err := DB.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// CreateMFASettings creates MFA settings for a user
func CreateMFASettings(userID uint, secret string, method MFAMethod) (*MFASettings, error) {
	settings := &MFASettings{
		UserID: userID,
		Secret: secret,
		Method: method,
		Status: MFAStatusPending,
	}
	err := DB.Create(settings).Error
	return settings, err
}

// EnableMFA enables MFA after verification
func EnableMFA(userID uint, recoveryCodes []string) error {
	hashedCodes, err := HashRecoveryCodes(recoveryCodes)
	if err != nil {
		return err
	}

	hashedJSON, _ := json.Marshal(hashedCodes)

	return DB.Model(&MFASettings{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"status":         MFAStatusEnabled,
			"recovery_codes": string(hashedJSON),
			"last_used_at":   nil, // Reset last used on enable
		}).Error
}

// DisableMFA disables MFA for a user
func DisableMFA(userID uint) error {
	return DB.Model(&MFASettings{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"status":         MFAStatusDisabled,
			"secret":         "",
			"recovery_codes": "[]",
		}).Error
}

// UseRecoveryCode marks a recovery code as used by removing it from the list
func UseRecoveryCode(userID uint, code string) error {
	var settings MFASettings
	if err := DB.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		return err
	}

	var codes []string
	if err := json.Unmarshal([]byte(settings.RecoveryCodes), &codes); err != nil {
		return err
	}

	// Remove the used code
	newCodes := make([]string, 0, len(codes))
	for _, c := range codes {
		if c != code {
			newCodes = append(newCodes, c)
		}
	}

	newJSON, _ := json.Marshal(newCodes)

	return DB.Model(&settings).
		Where("user_id = ?", userID).
		Update("recovery_codes", string(newJSON)).Error
}

// UpdateMFALastUsed updates the last used timestamp
func UpdateMFALastUsed(userID uint) error {
	now := time.Now()
	return DB.Model(&MFASettings{}).
		Where("user_id = ?", userID).
		Update("last_used_at", now).Error
}

// EncryptMFASecret encrypts an MFA secret
func EncryptMFASecret(secret string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(secret)), nil
}

// DecryptMFASecret decrypts an MFA secret
func DecryptMFASecret(encryptedSecret string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedSecret)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
