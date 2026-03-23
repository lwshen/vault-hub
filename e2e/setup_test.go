package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestServer manages a test instance of vault-hub server
type TestServer struct {
	URL       string
	APIKey    string
	VaultName string
	VaultID   string
	UserID    uint
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

	// Build server binary if not exists
	serverBinary := "/tmp/vault-hub-test-server"
	if _, err := os.Stat(serverBinary); os.IsNotExist(err) {
		cmd := exec.Command("go", "build", "-o", serverBinary, "./apps/server/main.go")
		cmd.Dir = "/home/openclaw/vault-hub"
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

	// Start server with SQLite
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, serverBinary)
	cmd.Env = append(os.Environ(),
		"APP_PORT=0", // Let OS assign port
		"DATABASE_TYPE=sqlite",
		"DATABASE_URL="+dbPath,
		"JWT_SECRET="+jwtSecret,
		"ENCRYPTION_KEY="+encryptionKey,
		"DEMO_ENABLED=true", // Enable demo mode to auto-create demo user
	)

	// Capture server output for port detection
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

	// Combine stdout and stderr for port detection
	var outputBuf bytes.Buffer
	go io.Copy(&outputBuf, stdoutPipe)
	go io.Copy(&outputBuf, stderrPipe)

	// Wait for server to start and extract port
	var port string
	for i := 0; i < 100; i++ {
		time.Sleep(100 * time.Millisecond)
		output := outputBuf.String()

		// Look for "addr=" in the log output (Fiber format)
		if strings.Contains(output, "addr=") {
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				if idx := strings.Index(line, "addr="); idx != -1 {
					addrPart := line[idx+5:]
					// Extract port from addr=:3000 or addr=0.0.0.0:3000
					if colonIdx := strings.LastIndex(addrPart, ":"); colonIdx != -1 {
						port = strings.TrimSpace(addrPart[colonIdx+1:])
						// Remove any trailing characters
						if spaceIdx := strings.IndexAny(port, " \t\"'"); spaceIdx != -1 {
							port = port[:spaceIdx]
						}
						break
					}
				}
			}
		}
		if port != "" {
			break
		}
	}

	if port == "" {
		cmd.Process.Kill()
		cancel()
		t.Fatalf("Server failed to start within timeout. Output:\n%s", outputBuf.String())
	}

	server := &TestServer{
		URL: "http://localhost:" + port,
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

	userID, ok := loginResp["userId"].(float64)
	if !ok {
		s.Stop()
		t.Fatalf("No userId in login response")
	}
	s.UserID = uint(userID)

	// Create API key
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

	plainKey, ok := apiKeyResp["plainKey"].(string)
	if !ok {
		s.Stop()
		t.Fatalf("No plainKey in API key response")
	}
	s.APIKey = plainKey
}

// setupTestVault creates a test vault via API
func (s *TestServer) setupTestVault(t *testing.T) {
	t.Helper()

	vaultName := "test-vault-" + generateRandomString(8)
	s.VaultName = vaultName

	// Create vault via API using API key
	vaultReq := map[string]interface{}{
		"uniqueId": uuid.New().String(),
		"name":     vaultName,
		"value":    "initial-test-value",
	}
	vaultBody, _ := json.Marshal(vaultReq)

	req, _ := http.NewRequest("POST", s.URL+"/api/cli/vaults", bytes.NewBuffer(vaultBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

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

	// Build CLI binary if not exists
	cliBinary := "/tmp/vault-hub-test-cli"
	if _, err := os.Stat(cliBinary); os.IsNotExist(err) {
		cmd := exec.Command("go", "build", "-o", cliBinary, "./apps/cli/main.go")
		cmd.Dir = "/home/openclaw/vault-hub"
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

// generateRandomString generates a random string for test data
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}