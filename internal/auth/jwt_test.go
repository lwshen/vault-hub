package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lwshen/vault-hub/internal/config"
)

func TestGenerateToken(t *testing.T) {
	originalJwtSecret := config.JwtSecret
	defer func() {
		config.JwtSecret = originalJwtSecret
	}()

	config.JwtSecret = "test-secret-key"

	tests := []struct {
		name   string
		userId uint
	}{
		{
			name:   "valid user ID",
			userId: 1,
		},
		{
			name:   "zero user ID",
			userId: 0,
		},
		{
			name:   "large user ID",
			userId: 999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := GenerateToken(tt.userId)
			if err != nil {
				t.Fatalf("GenerateToken() error = %v", err)
			}

			if tokenString == "" {
				t.Fatal("GenerateToken() returned empty token")
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					t.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(config.JwtSecret), nil
			})

			if err != nil {
				t.Fatalf("Failed to parse token: %v", err)
			}

			if !token.Valid {
				t.Fatal("Token is not valid")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				t.Fatal("Failed to get token claims")
			}

			subClaim, ok := claims["sub"]
			if !ok {
				t.Fatal("Missing 'sub' claim")
			}

			subFloat, ok := subClaim.(float64)
			if !ok {
				t.Fatal("'sub' claim is not a number")
			}

			if uint(subFloat) != tt.userId {
				t.Errorf("Expected userId %d, got %d", tt.userId, uint(subFloat))
			}

			expClaim, ok := claims["exp"]
			if !ok {
				t.Fatal("Missing 'exp' claim")
			}

			expFloat, ok := expClaim.(float64)
			if !ok {
				t.Fatal("'exp' claim is not a number")
			}

			expectedExp := time.Now().Add(time.Hour * expirationHour).Unix()
			actualExp := int64(expFloat)

			if actualExp < time.Now().Unix() {
				t.Fatal("Token is already expired")
			}

			if actualExp-expectedExp > 2 || expectedExp-actualExp > 2 {
				t.Errorf("Expected expiration around %d, got %d", expectedExp, actualExp)
			}

			iatClaim, ok := claims["iat"]
			if !ok {
				t.Fatal("Missing 'iat' claim")
			}

			iatFloat, ok := iatClaim.(float64)
			if !ok {
				t.Fatal("'iat' claim is not a number")
			}

			now := time.Now().Unix()
			actualIat := int64(iatFloat)

			if actualIat-now > 2 || now-actualIat > 2 {
				t.Errorf("Expected iat around %d, got %d", now, actualIat)
			}

			nbfClaim, ok := claims["nbf"]
			if !ok {
				t.Fatal("Missing 'nbf' claim")
			}

			nbfFloat, ok := nbfClaim.(float64)
			if !ok {
				t.Fatal("'nbf' claim is not a number")
			}

			actualNbf := int64(nbfFloat)

			if actualNbf-now > 2 || now-actualNbf > 2 {
				t.Errorf("Expected nbf around %d, got %d", now, actualNbf)
			}
		})
	}
}

func TestGenerateTokenWithEmptySecret(t *testing.T) {
	originalJwtSecret := config.JwtSecret
	defer func() {
		config.JwtSecret = originalJwtSecret
	}()

	config.JwtSecret = ""

	tokenString, err := GenerateToken(1)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	if tokenString == "" {
		t.Fatal("GenerateToken() returned empty token")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if !token.Valid {
		t.Fatal("Token is not valid")
	}
}

func TestGenerateTokenConsistency(t *testing.T) {
	originalJwtSecret := config.JwtSecret
	defer func() {
		config.JwtSecret = originalJwtSecret
	}()

	config.JwtSecret = "consistent-secret"

	userId := uint(42)

	token1, err := GenerateToken(userId)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	time.Sleep(time.Second * 1)

	token2, err := GenerateToken(userId)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	parseToken := func(tokenString string) jwt.MapClaims {
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JwtSecret), nil
		})
		if err != nil {
			t.Fatalf("Failed to parse token: %v", err)
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			t.Fatal("Failed to get token claims")
		}
		return claims
	}

	claims1 := parseToken(token1)
	claims2 := parseToken(token2)

	if claims1["sub"] != claims2["sub"] {
		t.Error("Expected same user ID in both tokens")
	}

	exp1 := int64(claims1["exp"].(float64))
	exp2 := int64(claims2["exp"].(float64))
	
	if exp2 <= exp1 {
		t.Error("Second token should have later expiration time")
	}

	iat1 := int64(claims1["iat"].(float64))
	iat2 := int64(claims2["iat"].(float64))
	
	if iat2 <= iat1 {
		t.Error("Second token should have later issued at time")
	}
}