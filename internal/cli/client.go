package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	openapi "github.com/lwshen/vault-hub/packages/api/openapi/client"
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

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// getEnvBool gets an environment variable as a boolean with a fallback value
func getEnvBool(key string, fallback bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	// Parse boolean from string - accepts "true", "1", "yes", etc.
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return boolVal
}

// getFlagOrEnv gets a flag value or falls back to environment variable
func getFlagOrEnv(cmd *cobra.Command, flagName, envVar, fallback string) string {
	// First check if the flag was explicitly set by user
	if cmd.Flags().Changed(flagName) {
		value, _ := cmd.Flags().GetString(flagName)
		DebugLog("Using --%s flag value", flagName)
		return value
	}

	// If flag not set, check environment variable
	value := getEnv(envVar, fallback)
	if value != fallback {
		DebugLog("Using %s environment variable", envVar)
	} else if fallback != "" {
		DebugLog("Using default value for %s", flagName)
	}
	return value
}

// getBoolFlagOrEnv gets a boolean flag value or falls back to environment variable
func getBoolFlagOrEnv(cmd *cobra.Command, flagName, envVar string, fallback bool) bool {
	// First check if the flag was explicitly set by user
	if cmd.Flags().Changed(flagName) {
		value, _ := cmd.Flags().GetBool(flagName)
		DebugLog("Using --%s flag value", flagName)
		return value
	}

	// If flag not set, check environment variable
	value := getEnvBool(envVar, fallback)
	if value != fallback {
		DebugLog("Using %s environment variable", envVar)
	} else {
		DebugLog("Using default value for %s", flagName)
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
		fmt.Fprintf(os.Stderr, "Error: API key is required. Set via --api-key flag or VAULT_HUB_API_KEY environment variable\n")
		os.Exit(1)
	}
	if BaseURL == "" {
		fmt.Fprintf(os.Stderr, "Error: Base URL is required. Set via --base-url flag or VAULT_HUB_BASE_URL environment variable\n")
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
