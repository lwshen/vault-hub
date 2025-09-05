package cli

import (
	"fmt"
	"os"
	"strings"

	openapi "github.com/lwshen/vault-hub-go-client"
	"github.com/spf13/cobra"

	"github.com/lwshen/vault-hub/internal/version"
)

var (
	APIKey  string
	BaseURL string
	Debug   bool
	Client  *openapi.APIClient
)

// DebugLog prints debug messages to stderr when debug mode is enabled
func DebugLog(format string, args ...any) {
	if Debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
	}
}

// GetStringFlag is a helper function to extract string flag values with error handling
func GetStringFlag(cmd *cobra.Command, flag string) (string, error) {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		DebugLog("Failed to get %s flag: %v", flag, err)
		return "", fmt.Errorf("failed to get %s flag: %v", flag, err)
	}
	return value, nil
}

// MustGetStringFlag extracts string flag values and exits on error
func MustGetStringFlag(cmd *cobra.Command, flag string) string {
	value, err := GetStringFlag(cmd, flag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return value
}

// InitializeClient sets up the OpenAPI client with authentication
func InitializeClient() {
	DebugLog("Initializing VaultHub CLI")
	// Remove trailing slash from baseURL if present
	BaseURL = strings.TrimSuffix(BaseURL, "/")
	DebugLog("Base URL: %s", BaseURL)
	DebugLog("Debug mode: %v", Debug)

	if APIKey == "" {
		fmt.Fprintf(os.Stderr, "Error: --api-key is required\n")
		os.Exit(1)
	}
	if BaseURL == "" {
		fmt.Fprintf(os.Stderr, "Error: --base-url is required\n")
		os.Exit(1)
	}

	cfg := openapi.NewConfiguration()
	cfg.Debug = Debug
	cfg.Servers = openapi.ServerConfigurations{
		{
			URL: BaseURL,
		},
	}
	DebugLog("Creating API client with configuration")
	Client = openapi.NewAPIClient(cfg)
	Client.GetConfig().DefaultHeader["Authorization"] = "Bearer " + APIKey
	Client.GetConfig().UserAgent = fmt.Sprintf("VaultHub-CLI/%s (%s)", version.Version, version.Commit)
	DebugLog("API client initialized successfully with User-Agent: %s", Client.GetConfig().UserAgent)
}
