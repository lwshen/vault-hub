package auth

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"github.com/lwshen/vault-hub/internal/config"
	"golang.org/x/oauth2"
)

func TestGenerateState(t *testing.T) {
	state1 := generateState()
	state2 := generateState()

	if state1 == "" {
		t.Fatal("generateState() returned empty string")
	}

	if state2 == "" {
		t.Fatal("generateState() returned empty string")
	}

	if state1 == state2 {
		t.Error("generateState() should return different states each time")
	}

	_, err := uuid.Parse(state1)
	if err != nil {
		t.Errorf("generateState() returned invalid UUID: %v", err)
	}

	_, err = uuid.Parse(state2)
	if err != nil {
		t.Errorf("generateState() returned invalid UUID: %v", err)
	}
}

func TestSessionOperations(t *testing.T) {
	app := fiber.New()
	
	sessionStore = session.New(session.Config{
		KeyLookup:  "cookie:test_session",
		Expiration: time.Hour * 1,
	})

	t.Run("store and retrieve session value", func(t *testing.T) {
		app.Get("/test", func(c *fiber.Ctx) error {
			err := storeInSession(c, "test_key", "test_value")
			if err != nil {
				return c.Status(500).SendString(err.Error())
			}

			value, err := getFromSession(c, "test_key")
			if err != nil {
				return c.Status(500).SendString(err.Error())
			}

			return c.SendString(value)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		buf := make([]byte, 1024)
		n, _ := resp.Body.Read(buf)
		body := string(buf[:n])

		if body != "test_value" {
			t.Errorf("Expected 'test_value', got '%s'", body)
		}
	})

	t.Run("get non-existent session key", func(t *testing.T) {
		app.Get("/test-missing", func(c *fiber.Ctx) error {
			_, err := getFromSession(c, "non_existent_key")
			if err != nil {
				return c.Status(400).SendString("expected error")
			}
			return c.SendString("unexpected success")
		})

		req := httptest.NewRequest("GET", "/test-missing", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 400 {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestAuthCodeURL(t *testing.T) {
	if !config.OidcEnabled {
		originalConfig := []string{config.OidcClientId, config.OidcClientSecret, config.OidcIssuer}
		defer func() {
			config.OidcClientId = originalConfig[0]
			config.OidcClientSecret = originalConfig[1]
			config.OidcIssuer = originalConfig[2]
			config.OidcEnabled = false
		}()

		config.OidcClientId = "test-client-id"
		config.OidcClientSecret = "test-client-secret"
		config.OidcIssuer = "https://example.com"
		config.OidcEnabled = true

		oauthConfig = &oauth2.Config{
			ClientID:     config.OidcClientId,
			ClientSecret: config.OidcClientSecret,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://example.com/auth",
				TokenURL: "https://example.com/token",
			},
		}

		sessionStore = session.New(session.Config{
			KeyLookup:  "cookie:auth_session",
			Expiration: time.Hour * 1,
		})
	}

	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		url, err := AuthCodeURL(c, "https://localhost:3000")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.SendString(url)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	buf := make([]byte, 1024)
	n, _ := resp.Body.Read(buf)
	url := string(buf[:n])

	if url == "" {
		t.Fatal("AuthCodeURL() returned empty URL")
	}

	expectedPrefix := "https://example.com/auth"
	if len(url) < len(expectedPrefix) || url[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Expected URL to start with %s, got %s", expectedPrefix, url)
	}

	if !contains(url, "client_id=test-client-id") {
		t.Error("URL should contain client_id parameter")
	}

	if !contains(url, "response_type=code") {
		t.Error("URL should contain response_type=code")
	}

	if !contains(url, "scope=openid+email+profile") {
		t.Error("URL should contain correct scope")
	}

	if !contains(url, "state=") {
		t.Error("URL should contain state parameter")
	}

	if !contains(url, "redirect_uri=https%3A%2F%2Flocalhost%3A3000%2Fapi%2Fauth%2Fcallback%2Foidc") {
		t.Error("URL should contain correct redirect_uri")
	}
}

func TestVerifyState(t *testing.T) {
	app := fiber.New()

	sessionStore = session.New(session.Config{
		KeyLookup:  "cookie:auth_session",
		Expiration: time.Hour * 1,
	})

	t.Run("valid state verification", func(t *testing.T) {
		app.Get("/test-valid", func(c *fiber.Ctx) error {
			state := "test-state-123"
			
			err := storeInSession(c, "oauth", state)
			if err != nil {
				return c.Status(500).SendString(err.Error())
			}

			err = VerifyState(c, state)
			if err != nil {
				return c.Status(400).SendString(err.Error())
			}

			return c.SendString("success")
		})

		req := httptest.NewRequest("GET", "/test-valid", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("invalid state verification", func(t *testing.T) {
		app.Get("/test-invalid", func(c *fiber.Ctx) error {
			err := storeInSession(c, "oauth", "stored-state")
			if err != nil {
				return c.Status(500).SendString(err.Error())
			}

			err = VerifyState(c, "different-state")
			if err != nil {
				return c.Status(400).SendString("expected error")
			}

			return c.SendString("unexpected success")
		})

		req := httptest.NewRequest("GET", "/test-invalid", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 400 {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("missing state in session", func(t *testing.T) {
		app.Get("/test-missing", func(c *fiber.Ctx) error {
			err := VerifyState(c, "any-state")
			if err != nil {
				return c.Status(400).SendString("expected error")
			}

			return c.SendString("unexpected success")
		})

		req := httptest.NewRequest("GET", "/test-missing", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 400 {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestSetupOIDCError(t *testing.T) {
	originalConfig := []string{config.OidcClientId, config.OidcClientSecret, config.OidcIssuer}
	defer func() {
		config.OidcClientId = originalConfig[0]
		config.OidcClientSecret = originalConfig[1]
		config.OidcIssuer = originalConfig[2]
	}()

	config.OidcClientId = "test-client"
	config.OidcClientSecret = "test-secret"
	config.OidcIssuer = "invalid-url"

	err := SetupOIDC()
	if err == nil {
		t.Error("Expected SetupOIDC() to return error for invalid issuer URL")
	}
}

func TestVerifyWithInvalidToken(t *testing.T) {
	if oauthConfig == nil {
		oauthConfig = &oauth2.Config{
			ClientID:     "test-client",
			ClientSecret: "test-secret",
		}
	}

	ctx := context.Background()
	
	_, err := Verify(ctx, "invalid-code")
	if err == nil {
		t.Error("Expected Verify() to return error for invalid code")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) != -1
}

func findSubstring(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestVerifyTokenError(t *testing.T) {
	ctx := context.Background()
	
	token := &oauth2.Token{
		AccessToken: "test-token",
		TokenType:   "Bearer",
	}

	_, err := verifyToken(ctx, token)
	if err == nil {
		t.Error("Expected verifyToken() to return error for token without id_token")
	}

	expectedError := "missing ID token"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestUserInfoError(t *testing.T) {
	if provider == nil {
		t.Skip("Skipping UserInfo test as provider is not initialized")
	}

	ctx := context.Background()
	token := &oauth2.Token{
		AccessToken: "invalid-token",
		TokenType:   "Bearer",
	}

	_, err := UserInfo(ctx, token)
	if err == nil {
		t.Error("Expected UserInfo() to return error for invalid token")
	}
}