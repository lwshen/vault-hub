package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lwshen/vault-hub/handler"
	clientencryption "github.com/lwshen/vault-hub/internal/cli/encryption"
	"github.com/lwshen/vault-hub/internal/config"
	"github.com/lwshen/vault-hub/internal/constants"
	"github.com/lwshen/vault-hub/model"
)

type apiErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func TestUpdateVaultByAPIKeySuccess(t *testing.T) {
	app := setupCLIUpdateTestApp(t)
	user := createTestUser(t)
	vault := createTestVault(t, user.ID, "prod-secrets", "initial-secret")
	apiKey, plainAPIKey := createTestAPIKey(t, user.ID, "cli-updater", []uint{vault.ID})

	requestBody := map[string]any{
		"value":       "updated-secret",
		"description": "rotated production token",
		"favourite":   true,
	}
	response := makeCLIRequest(
		t,
		app,
		http.MethodPut,
		fmt.Sprintf("/api/cli/vault/%s", vault.UniqueID),
		plainAPIKey,
		requestBody,
		nil,
	)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.StatusCode, readBody(t, response))
	}

	var payload Vault
	decodeJSON(t, response, &payload)
	if payload.Value != "updated-secret" {
		t.Fatalf("expected updated value in response, got %q", payload.Value)
	}
	if payload.Description == nil || *payload.Description != "rotated production token" {
		t.Fatalf("expected updated description, got %#v", payload.Description)
	}
	if payload.Favourite == nil || !*payload.Favourite {
		t.Fatalf("expected favourite=true, got %#v", payload.Favourite)
	}

	var updated model.Vault
	if err := updated.GetByUniqueID(vault.UniqueID, user.ID); err != nil {
		t.Fatalf("failed to reload updated vault: %v", err)
	}
	if updated.Value != "updated-secret" {
		t.Fatalf("expected database value to be updated, got %q", updated.Value)
	}
	if updated.Description != "rotated production token" {
		t.Fatalf("expected database description to be updated, got %q", updated.Description)
	}
	if !updated.Favourite {
		t.Fatalf("expected database favourite=true")
	}

	var auditLog model.AuditLog
	err := model.DB.
		Where("vault_id = ? AND action = ? AND source = ?", vault.ID, model.ActionUpdateVault, model.SourceCLI).
		Order("id DESC").
		First(&auditLog).Error
	if err != nil {
		t.Fatalf("expected update audit log entry, got error: %v", err)
	}
	if auditLog.APIKeyID == nil || *auditLog.APIKeyID != apiKey.ID {
		t.Fatalf("expected API key ID %d in audit log, got %#v", apiKey.ID, auditLog.APIKeyID)
	}
}

func TestUpdateVaultByNameAPIKeySuccessWithRename(t *testing.T) {
	app := setupCLIUpdateTestApp(t)
	user := createTestUser(t)
	_ = createTestVault(t, user.ID, "old-name", "old-secret")
	_, plainAPIKey := createTestAPIKey(t, user.ID, "cli-updater", nil)

	requestBody := map[string]any{
		"name":     "new-name",
		"value":    "new-secret",
		"category": "infra",
	}
	response := makeCLIRequest(
		t,
		app,
		http.MethodPut,
		"/api/cli/vault/name/old-name",
		plainAPIKey,
		requestBody,
		nil,
	)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.StatusCode, readBody(t, response))
	}

	var payload Vault
	decodeJSON(t, response, &payload)
	if payload.Name != "new-name" {
		t.Fatalf("expected renamed vault in response, got %q", payload.Name)
	}
	if payload.Value != "new-secret" {
		t.Fatalf("expected updated value in response, got %q", payload.Value)
	}
	if payload.Category == nil || *payload.Category != "infra" {
		t.Fatalf("expected updated category, got %#v", payload.Category)
	}

	var updated model.Vault
	if err := updated.GetByName("new-name", user.ID); err != nil {
		t.Fatalf("expected renamed vault in database, got error: %v", err)
	}
	if updated.Value != "new-secret" {
		t.Fatalf("expected updated vault value in database, got %q", updated.Value)
	}
}

func TestUpdateVaultByAPIKeyForbiddenWithoutAccess(t *testing.T) {
	app := setupCLIUpdateTestApp(t)
	user := createTestUser(t)
	allowedVault := createTestVault(t, user.ID, "allowed-vault", "allowed-secret")
	deniedVault := createTestVault(t, user.ID, "denied-vault", "denied-secret")
	_, plainAPIKey := createTestAPIKey(t, user.ID, "restricted-key", []uint{allowedVault.ID})

	requestBody := map[string]any{
		"value": "attempted-update",
	}
	response := makeCLIRequest(
		t,
		app,
		http.MethodPut,
		fmt.Sprintf("/api/cli/vault/%s", deniedVault.UniqueID),
		plainAPIKey,
		requestBody,
		nil,
	)

	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", response.StatusCode, readBody(t, response))
	}

	var errResponse apiErrorResponse
	decodeJSON(t, response, &errResponse)
	if !strings.Contains(errResponse.Error.Message, "does not have access") {
		t.Fatalf("expected forbidden access message, got %q", errResponse.Error.Message)
	}

	var unchanged model.Vault
	if err := unchanged.GetByUniqueID(deniedVault.UniqueID, user.ID); err != nil {
		t.Fatalf("failed to reload denied vault: %v", err)
	}
	if unchanged.Value != "denied-secret" {
		t.Fatalf("expected denied vault value to remain unchanged, got %q", unchanged.Value)
	}
}

func TestUpdateVaultByAPIKeyNotFound(t *testing.T) {
	app := setupCLIUpdateTestApp(t)
	user := createTestUser(t)
	_, plainAPIKey := createTestAPIKey(t, user.ID, "cli-updater", nil)

	requestBody := map[string]any{
		"value": "new-value",
	}

	testCases := []struct {
		name string
		path string
	}{
		{
			name: "by unique id",
			path: "/api/cli/vault/not-a-real-unique-id",
		},
		{
			name: "by name",
			path: "/api/cli/vault/name/not-a-real-name",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			response := makeCLIRequest(t, app, http.MethodPut, testCase.path, plainAPIKey, requestBody, nil)
			if response.StatusCode != http.StatusNotFound {
				t.Fatalf("expected status 404, got %d: %s", response.StatusCode, readBody(t, response))
			}

			var errResponse apiErrorResponse
			decodeJSON(t, response, &errResponse)
			if errResponse.Error.Message != "vault not found" {
				t.Fatalf("expected not found message, got %q", errResponse.Error.Message)
			}
		})
	}
}

func TestUpdateVaultByAPIKeyBadRequestForInvalidPayload(t *testing.T) {
	app := setupCLIUpdateTestApp(t)
	user := createTestUser(t)
	vault := createTestVault(t, user.ID, "test-vault", "test-secret")
	_, plainAPIKey := createTestAPIKey(t, user.ID, "cli-updater", []uint{vault.ID})

	requestBody := map[string]any{
		"value": "",
	}
	response := makeCLIRequest(
		t,
		app,
		http.MethodPut,
		fmt.Sprintf("/api/cli/vault/%s", vault.UniqueID),
		plainAPIKey,
		requestBody,
		nil,
	)

	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.StatusCode, readBody(t, response))
	}

	var errResponse apiErrorResponse
	decodeJSON(t, response, &errResponse)
	if !strings.Contains(errResponse.Error.Message, "value cannot be empty") {
		t.Fatalf("expected validation error for empty value, got %q", errResponse.Error.Message)
	}
}

func TestUpdateVaultByNameAPIKeyClientEncryptionResponse(t *testing.T) {
	app := setupCLIUpdateTestApp(t)
	user := createTestUser(t)
	_ = createTestVault(t, user.ID, "salt-name", "before-update")
	_, plainAPIKey := createTestAPIKey(t, user.ID, "cli-updater", nil)

	requestBody := map[string]any{
		"value": "encrypted-response-secret",
	}
	response := makeCLIRequest(
		t,
		app,
		http.MethodPut,
		"/api/cli/vault/name/salt-name",
		plainAPIKey,
		requestBody,
		map[string]string{
			constants.HeaderClientEncryption: "true",
		},
	)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.StatusCode, readBody(t, response))
	}

	var payload Vault
	decodeJSON(t, response, &payload)
	if payload.Value == "encrypted-response-secret" {
		t.Fatalf("expected encrypted response value when %s=true", constants.HeaderClientEncryption)
	}

	decryptedValue, err := clientencryption.DecryptForClient(payload.Value, plainAPIKey, "salt-name")
	if err != nil {
		t.Fatalf("failed to decrypt client-encrypted response: %v", err)
	}
	if decryptedValue != "encrypted-response-secret" {
		t.Fatalf("unexpected decrypted value, got %q", decryptedValue)
	}
}

func setupCLIUpdateTestApp(t *testing.T) *fiber.App {
	t.Helper()
	setupCLIUpdateTestDatabase(t)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		if !strings.HasPrefix(c.Path(), "/api/cli/") {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return handler.SendError(c, fiber.StatusUnauthorized, "API key required for this endpoint")
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return handler.SendError(c, fiber.StatusUnauthorized, "invalid authorization header")
		}

		tokenString := tokenParts[1]
		if !strings.HasPrefix(tokenString, "vhub_") {
			return handler.SendError(c, fiber.StatusUnauthorized, "API key required for this endpoint")
		}

		apiKey, err := model.ValidateAPIKey(tokenString)
		if err != nil {
			return handler.SendError(c, fiber.StatusUnauthorized, "invalid API key")
		}

		c.Locals("user_id", &apiKey.UserID)
		c.Locals("api_key", apiKey)

		return c.Next()
	})

	server := NewServer()
	RegisterHandlers(app, server)
	return app
}

func setupCLIUpdateTestDatabase(t *testing.T) {
	t.Helper()

	config.DatabaseType = config.DatabaseTypeSQLite
	config.DatabaseUrl = filepath.Join(t.TempDir(), "cli-vault-update-test.db")

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	if err := model.Open(logger); err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
}

func createTestUser(t *testing.T) *model.User {
	t.Helper()

	user := &model.User{
		Email: fmt.Sprintf("cli-update-%s@example.com", uuid.NewString()),
	}
	if err := model.DB.Create(user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return user
}

func createTestVault(t *testing.T, userID uint, name, value string) *model.Vault {
	t.Helper()

	createParams := model.CreateVaultParams{
		UniqueID: uuid.NewString(),
		UserID:   userID,
		Name:     name,
		Value:    value,
	}
	vault, err := createParams.Create()
	if err != nil {
		t.Fatalf("failed to create test vault: %v", err)
	}

	return vault
}

func createTestAPIKey(t *testing.T, userID uint, name string, vaultIDs []uint) (*model.APIKey, string) {
	t.Helper()

	createParams := model.CreateAPIKeyParams{
		UserID:   userID,
		Name:     name,
		VaultIDs: vaultIDs,
	}
	apiKey, plainKey, err := createParams.Create()
	if err != nil {
		t.Fatalf("failed to create test API key: %v", err)
	}

	return apiKey, plainKey
}

func makeCLIRequest(
	t *testing.T,
	app *fiber.App,
	method string,
	path string,
	apiKey string,
	body any,
	extraHeaders map[string]string,
) *http.Response {
	t.Helper()

	var requestBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		requestBody = bytes.NewReader(jsonBody)
	}

	request := httptest.NewRequest(method, path, requestBody)
	request.Header.Set("Authorization", "Bearer "+apiKey)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	for key, value := range extraHeaders {
		request.Header.Set(key, value)
	}

	response, err := app.Test(request, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	return response
}

func readBody(t *testing.T, response *http.Response) string {
	t.Helper()
	defer response.Body.Close()

	rawBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	return string(rawBody)
}

func decodeJSON(t *testing.T, response *http.Response, target any) {
	t.Helper()
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
}
