package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lwshen/vault-hub/internal/config"
)

const expirationHour = 24

func GenerateToken(userId uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Hour * expirationHour).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
	})
	return token.SignedString([]byte(config.JwtSecret))
}
