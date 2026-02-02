package auth

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"strings"
	"time"
)

// TOTP constants
const (
	TOTPdigits    = 6
	TOTPperiod    = 30 // seconds
	TOTPalgorithm = "SHA1"
)

// base32Encoding is a base32 encoding without padding
var base32Encoding = base32.StdEncoding.WithPadding(base32.NoPadding)

// GenerateTOTPSecret generates a new TOTP secret
func GenerateTOTPSecret() (string, error) {
	// Generate 20 bytes of random data for a good secret
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base32Encoding.EncodeToString(b), nil
}

// GenerateTOTPUrl generates the TOTP URL for authenticator apps
func GenerateTOTPUrl(secret, accountName, issuer string) string {
	// URL format: otpauth://totp/ISSUER:ACCOUNT?secret=SECRET&issuer=ISSUER
	encodedSecret := base32Encoding.EncodeToString([]byte(secret))
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		issuer, accountName, encodedSecret, issuer)
}

// ValidateTOTP validates a TOTP code
func ValidateTOTP(secret, code string) bool {
	// Get current and adjacent time windows
	now := time.Now().Unix()
	currentWindow := now / TOTPperiod

	// Check current window and adjacent windows (±1) for clock drift
	for i := -1; i <= 1; i++ {
		window := currentWindow + int64(i)
		expectedCode := generateTOTP(secret, window)
		if code == expectedCode {
			return true
		}
	}
	return false
}

// HMAC computes HMAC using the provided key and message
func HMAC(key, message []byte) []byte {
	// Simple HMAC-SHA1 implementation
	blockSize := 64

	// If key is longer than block size, hash it first
	if len(key) > blockSize {
		// #nosec G401
		h := sha1.New()
		h.Write(key)
		key = h.Sum(nil)
	}

	// Pad key to block size
	paddedKey := make([]byte, blockSize)
	copy(paddedKey, key)

	// Inner padding
	innerPad := make([]byte, blockSize)
	for i := 0; i < blockSize; i++ {
		innerPad[i] = 0x36
	}

	// Outer padding
	outerPad := make([]byte, blockSize)
	for i := 0; i < blockSize; i++ {
		outerPad[i] = 0x5c
	}

	// Inner hash - #nosec G401
	innerHash := sha1.New()
	innerHash.Write(innerPad)
	innerHash.Write(message)
	innerResult := innerHash.Sum(nil)

	// Outer hash - #nosec G401
	outerHash := sha1.New()
	outerHash.Write(outerPad)
	outerHash.Write(innerResult)
	return outerHash.Sum(nil)
}

// generateTOTP generates a TOTP code for a specific time window
func generateTOTP(secret string, window int64) string {
	// Decode the secret
	secretBytes, err := base32Encoding.DecodeString(secret)
	if err != nil {
		return ""
	}

	// Convert window to 8-byte big-endian
	b := make([]byte, 8)
	b[0] = byte(window >> 56)
	b[1] = byte(window >> 48)
	b[2] = byte(window >> 40)
	b[3] = byte(window >> 32)
	b[4] = byte(window >> 24)
	b[5] = byte(window >> 16)
	b[6] = byte(window >> 8)
	b[7] = byte(window)

	// HMAC-SHA1
	mac := HMAC(secretBytes, b)
	offset := mac[len(mac)-1] & 0x0f
	truncated := mac[offset : offset+4]

	// Take the last 31 bits
	code := int32(truncated[0]&0x7f)<<24 |
		int32(truncated[1])<<16 |
		int32(truncated[2])<<8 |
		int32(truncated[3])

	// Generate the code
	result := code % 1000000
	return fmt.Sprintf("%06d", result)
}

// FormatMFASecretForDisplay formats a TOTP secret for display
func FormatMFASecretForDisplay(secret string) string {
	// Format as groups of 4 characters for readability
	chunks := make([]string, 0, len(secret)/4+1)
	for i := 0; i < len(secret); i += 4 {
		end := i + 4
		if end > len(secret) {
			end = len(secret)
		}
		chunks = append(chunks, secret[i:end])
	}
	return strings.Join(chunks, " ")
}
