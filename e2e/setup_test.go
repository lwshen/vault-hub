package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestServer manages a test instance of vault-hub server
type TestServer struct {
	URL       string
	APIKey    string
	JWTToken  string // JWT token for authenticated requests
	VaultName string
	VaultID   string
	cmd       *exec.Cmd
	cancel    context.CancelFunc
}

// CLIResult holds the result of a CLI execution
type CLIResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

// StartTestServer starts a test server with SQLite in-memory database
func StartTestServer(t *testing.T) *TestServer {
	t.Helper()

	// Get project root directory (parent of e2e/)
	projectRoot := os.Getenv("PROJECT_ROOT")
	if projectRoot == "" {
		// Try to find project root from current directory
		wd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Failed to get working directory: %v", err)
		}
		projectRoot = wd
	}

	// Build server binary if not exists
	serverBinary := "/tmp/vault-hub-test-server"
	if _, err := os.Stat(serverBinary); os.IsNotExist(err) {
		cmd := exec.Command("go", "build", "-o", serverBinary, "./apps/server/main.go")
		cmd.Dir = projectRoot
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Failed to build server: %v\n%s", err, output)
		}
	}

	// Generate test secrets
	jwtSecret := generateRandomString(32)
	encryptionKey := generateRandomString(32)

	// Create temp directory for test data (SQLite needs a file, not :memory: for concurrent access)
	tempDir, err := os.MkdirTemp("", "vault-hub-e2e-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	dbPath := tempDir + "/test.db"

	// Find an available port
	port := findAvailablePort(t)

	// Start server with SQLite
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, serverBinary)
	cmd.Env = append(os.Environ(),
		"APP_PORT="+port, // Use fixed port
		"DATABASE_TYPE=sqlite",
		"DATABASE_URL="+dbPath,
		"JWT_SECRET="+jwtSecret,
		"ENCRYPTION_KEY="+encryptionKey,
		"DEMO_ENABLED=true", // Enable demo mode to auto-create demo user
	)

	// Capture server output
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		t.Fatalf("Failed to create stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		cancel()
		t.Fatalf("Failed to start server: %v", err)
	}

	// Discard output (we don't need to parse port anymore)
	go io.Copy(io.Discard, stdoutPipe)
	go io.Copy(io.Discard, stderrPipe)

	// Wait for server to be ready (check health endpoint)
	serverURL := "http://localhost:" + port
	var serverReady bool
	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		resp, err := http.Get(serverURL + "/api/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			serverReady = true
			break
		}
	}

	if !serverReady {
		cmd.Process.Kill()
		cancel()
		t.Fatalf("Server failed to start within timeout on port %s", port)
	}

	server := &TestServer{
		URL: serverURL,
		cmd: cmd,
		cancel: cancel,
	}

	// Wait a bit for server to be fully ready
	time.Sleep(200 * time.Millisecond)

	// Setup test user and API key
	server.setupTestUserAndAPIKey(t)

	// Create test vault
	server.setupTestVault(t)

	// Cleanup on test end
	t.Cleanup(func() {
		server.Stop()
		os.RemoveAll(tempDir)
	})

	return server
}

// setupTestUserAndAPIKey creates a test user and API key via API
func (s *TestServer) setupTestUserAndAPIKey(t *testing.T) {
	t.Helper()

	// Demo mode creates a demo user automatically, we need to create an API key for it
	// First, login as demo user to get a token
	loginReq := map[string]string{
		"email":    "mock@demo.com",
		"password": "Test1234!",
	}

	loginBody, _ := json.Marshal(loginReq)
	resp, err := http.Post(s.URL+"/api/auth/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		s.Stop()
		t.Fatalf("Failed to login as demo user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.Stop()
		t.Fatalf("Login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var loginResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		s.Stop()
		t.Fatalf("Failed to decode login response: %v", err)
	}

	token, ok := loginResp["token"].(string)
	if !ok {
		s.Stop()
		t.Fatalf("No token in login response")
	}
	s.JWTToken = token // Save JWT token for later use

	// Create API key using the token
	apiKeyReq := map[string]interface{}{
		"name": "e2e-test-key",
	}
	apiKeyBody, _ := json.Marshal(apiKeyReq)

	req, _ := http.NewRequest("POST", s.URL+"/api/api-keys", bytes.NewBuffer(apiKeyBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	apiResp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.Stop()
		t.Fatalf("Failed to create API key: %v", err)
	}
	defer apiResp.Body.Close()

	if apiResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(apiResp.Body)
		s.Stop()
		t.Fatalf("API key creation failed with status %d: %s", apiResp.StatusCode, string(body))
	}

	var apiKeyResp map[string]interface{}
	if err := json.NewDecoder(apiResp.Body).Decode(&apiKeyResp); err != nil {
		s.Stop()
		t.Fatalf("Failed to decode API key response: %v", err)
	}

	// Try "key" first (CreateAPIKeyResponse uses "key"), then "plainKey"
	var apiKey string
	if key, ok := apiKeyResp["key"].(string); ok {
		apiKey = key
	} else if key, ok := apiKeyResp["plainKey"].(string); ok {
		apiKey = key
	} else {
		s.Stop()
		t.Fatalf("No key or plainKey in API key response: %v", apiKeyResp)
	}
	s.APIKey = apiKey
}

// setupTestVault creates a test vault via API using JWT authentication
func (s *TestServer) setupTestVault(t *testing.T) {
	t.Helper()

	vaultName := "test-vault-" + generateRandomString(8)
	s.VaultName = vaultName

	// Create vault via API using JWT token (not API key)
	// POST /api/vaults requires JWT authentication
	vaultReq := map[string]interface{}{
		"uniqueId": uuid.New().String(),
		"name":     vaultName,
		"value":    "initial-test-value",
	}
	vaultBody, _ := json.Marshal(vaultReq)

	req, _ := http.NewRequest("POST", s.URL+"/api/vaults", bytes.NewBuffer(vaultBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.JWTToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.Stop()
		t.Fatalf("Failed to create vault: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		s.Stop()
		t.Fatalf("Vault creation failed with status %d: %s", resp.StatusCode, string(body))
	}

	var vaultResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&vaultResp); err != nil {
		s.Stop()
		t.Fatalf("Failed to decode vault response: %v", err)
	}

	vaultID, ok := vaultResp["uniqueId"].(string)
	if !ok {
		s.Stop()
		t.Fatalf("No uniqueId in vault response")
	}
	s.VaultID = vaultID
}

// Stop stops the test server
func (s *TestServer) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	if s.cmd != nil && s.cmd.Process != nil {
		s.cmd.Process.Kill()
		s.cmd.Wait()
	}
}

// RunCLI executes the CLI with given arguments
func RunCLI(t *testing.T, args ...string) *CLIResult {
	t.Helper()

	// Get project root directory
	projectRoot := os.Getenv("PROJECT_ROOT")
	if projectRoot == "" {
		wd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Failed to get working directory: %v", err)
		}
		projectRoot = wd
	}

	// Build CLI binary if not exists
	cliBinary := "/tmp/vault-hub-test-cli"
	if _, err := os.Stat(cliBinary); os.IsNotExist(err) {
		cmd := exec.Command("go", "build", "-o", cliBinary, "./apps/cli/main.go")
		cmd.Dir = projectRoot
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Failed to build CLI: %v\n%s", err, output)
		}
	}

	cmd := exec.Command(cliBinary, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return &CLIResult{
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}
}

// MustSucceed fails the test if CLI result indicates failure
func (r *CLIResult) MustSucceed(t *testing.T) {
	t.Helper()
	if r.ExitCode != 0 {
		t.Fatalf("Expected success but got exit code %d\nStdout: %s\nStderr: %s",
			r.ExitCode, r.Stdout, r.Stderr)
	}
}

// MustFail fails the test if CLI result indicates success
func (r *CLIResult) MustFail(t *testing.T, expectedExitCode int) {
	t.Helper()
	if r.ExitCode != expectedExitCode {
		t.Fatalf("Expected exit code %d but got %d\nStdout: %s\nStderr: %s",
			expectedExitCode, r.ExitCode, r.Stdout, r.Stderr)
	}
}

// ContainsStdout checks if stdout contains expected string
func (r *CLIResult) ContainsStdout(t *testing.T, expected string) bool {
	t.Helper()
	return strings.Contains(r.Stdout, expected)
}

// ContainsStderr checks if stderr contains expected string
func (r *CLIResult) ContainsStderr(t *testing.T, expected string) bool {
	t.Helper()
	return strings.Contains(r.Stderr, expected)
}

// findAvailablePort finds an available TCP port for the test server
func findAvailablePort(t *testing.T) string {
	t.Helper()

	// Try to find an available port
	for i := 0; i < 10; i++ {
		// Start from a high port range to avoid conflicts
		port := 30000 + (time.Now().Nanosecond() % 30000)
		addr := fmt.Sprintf(":%d", port)

		// Try to listen on the port
		ln, err := net.Listen("tcp", addr)
		if err == nil {
			ln.Close()
			return fmt.Sprintf("%d", port)
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatal("Could not find an available port")
	return ""
}

// generateRandomString generates a random string for test data
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}