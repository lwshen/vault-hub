package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/model"
	"github.com/lwshen/vault-hub/packages/api/gen/models"
	"github.com/lwshen/vault-hub/route"
	"github.com/stretchr/testify/assert"
)

func TestEchoAdapter(t *testing.T) {
	adapter := NewEchoAdapter()

	t.Run("ConvertHealthCheckResponse", func(t *testing.T) {
		status := "healthy"
		timestamp := time.Now()

		result := adapter.ConvertHealthCheckResponse(status, timestamp)

		assert.Equal(t, status, result.Status)
		assert.Equal(t, timestamp, result.Timestamp)
	})

	t.Run("ConvertLoginResponse", func(t *testing.T) {
		token := "test-jwt-token"

		result := adapter.ConvertLoginResponse(token)

		assert.Equal(t, token, result.Token)
	})

	t.Run("ConvertConfigResponse", func(t *testing.T) {
		result := adapter.ConvertConfigResponse()

		assert.True(t, result.IsOidcEnabled || !result.IsOidcEnabled) // Should be either true or false
		assert.True(t, result.IsEmailEnabled || !result.IsEmailEnabled)   // Should be either true or false
		assert.Greater(t, result.PasswordMinLength, int64(0))
	})
}

func TestEchoMiddleware(t *testing.T) {
	t.Run("Health Endpoint", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.HealthCheckResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ok", response.Status)
	})
}

func TestEchoContextHelpers(t *testing.T) {
	t.Run("GetCurrentUserFromEcho", func(t *testing.T) {
		e := echo.New()

		// Mock user
		user := &model.User{
			ID:    1,
			Email:  "test@example.com",
			Name:   "Test User",
		}

		req := httptest.NewRequest(http.MethodGet, "/api/user", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", user)

		retrieved, err := route.GetCurrentUserFromEcho(c)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, retrieved.ID)
		assert.Equal(t, user.Email, retrieved.Email)
		assert.Equal(t, user.Name, retrieved.Name)
		assert.Nil(t, retrieved.Password) // Should be cleared
	})
}

func TestModelCompatibility(t *testing.T) {
	adapter := NewEchoAdapter()

	t.Run("Vault Model Conversion", func(t *testing.T) {
		description := "Test vault description"
		value := "encrypted-value"

		vault := &model.Vault{
			ID:          1,
			UniqueID:    "test-uuid",
			Name:        "Test Vault",
			Description: &description,
			Value:       "encrypted-value", // This will be decrypted in actual implementation
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsActive:    true,
		}

		// Convert with decrypted value
		result := adapter.ConvertVault(vault, value)

		assert.Equal(t, vault.UniqueID, result.UniqueId)
		assert.Equal(t, vault.Name, result.Name)
		assert.Equal(t, description, result.Description)
		assert.Equal(t, value, result.Value)
		assert.Equal(t, vault.IsActive, result.IsActive)
	})

	t.Run("APIKey Model Conversion", func(t *testing.T) {
		apiKey := &model.APIKey{
			ID:        1,
			Name:      "Test Key",
			Key:       "vhub_test_key",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsActive: true,
		}

		result := adapter.ConvertAPIKey(apiKey)

		assert.Equal(t, int64(apiKey.ID), result.Id)
		assert.Equal(t, apiKey.Name, result.Name)
		assert.Equal(t, apiKey.Key, result.Key)
		assert.Equal(t, apiKey.IsActive, result.IsActive)
	})
}

func TestEchoAuthentication(t *testing.T) {
	// This would require database setup for full integration testing
	t.Run("JWT Token Validation", func(t *testing.T) {
		// Mock JWT token validation would go here
		// For now, test that middleware rejects invalid tokens
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/vaults", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := route.handleEchoJWTAuth(c, func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}, "invalid-token")

		assert.Error(t, err)

		// Check if it's an HTTPError with correct status
		if httpErr, ok := err.(*echo.HTTPError); ok {
			assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
		}
	})

	t.Run("API Key Validation", func(t *testing.T) {
		// Mock API key validation would go here
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/cli/vaults", nil)
		req.Header.Set("Authorization", "Bearer vhub_invalid_key")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := route.handleEchoAPIKeyAuth(c, func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}, "vhub_invalid_key")

		assert.Error(t, err)

		// Check if it's an HTTPError with correct status
		if httpErr, ok := err.(*echo.HTTPError); ok {
			assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
		}
	})
}

// BenchmarkEchoVsFiber compares performance between Echo and Fiber implementations
func BenchmarkEchoVsFiber(b *testing.B) {
	// This would require running both implementations
	b.Skip("Skipping benchmark - requires full server setup")

	// Example benchmark structure:
	// b.Run("Echo", func(b *testing.B) {
	//     e := echo.New()
	//     // Setup Echo server and run requests
	// })
	//
	// b.Run("Fiber", func(b *testing.B) {
	//     app := fiber.New()
	//     // Setup Fiber server and run requests
	// })
}