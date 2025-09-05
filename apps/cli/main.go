package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	openapi "github.com/lwshen/vault-hub-go-client"
	"github.com/spf13/cobra"

	"github.com/lwshen/vault-hub/internal/version"
)

var (
	apiKey  string
	baseURL string
	debug   bool
	client  *openapi.APIClient
)

// debugLog prints debug messages to stderr when debug mode is enabled
func debugLog(format string, args ...any) {
	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
	}
}

// getStringFlag is a helper function to extract string flag values with error handling
func getStringFlag(cmd *cobra.Command, flag string) (string, error) {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		debugLog("Failed to get %s flag: %v", flag, err)
		return "", fmt.Errorf("failed to get %s flag: %v", flag, err)
	}
	return value, nil
}

// mustGetStringFlag extracts string flag values and exits on error
func mustGetStringFlag(cmd *cobra.Command, flag string) string {
	value, err := getStringFlag(cmd, flag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return value
}

var rootCmd = &cobra.Command{
	Use:   "vault-hub",
	Short: "VaultHub CLI - Secure environment variable and API key management",
	Long: `VaultHub CLI is a command-line interface for managing your secure
environment variables and API keys stored in VaultHub.

This CLI allows you to list and retrieve vaults from your VaultHub instance.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debugLog("Initializing VaultHub CLI")
		// Remove trailing slash from baseURL if present
		baseURL = strings.TrimSuffix(baseURL, "/")
		debugLog("Base URL: %s", baseURL)
		debugLog("Debug mode: %v", debug)

		if apiKey == "" {
			fmt.Fprintf(os.Stderr, "Error: --api-key is required\n")
			os.Exit(1)
		}
		if baseURL == "" {
			fmt.Fprintf(os.Stderr, "Error: --base-url is required\n")
			os.Exit(1)
		}

		cfg := openapi.NewConfiguration()
		cfg.Debug = debug
		cfg.Servers = openapi.ServerConfigurations{
			{
				URL: baseURL,
			},
		}
		debugLog("Creating API client with configuration")
		client = openapi.NewAPIClient(cfg)
		client.GetConfig().DefaultHeader["Authorization"] = "Bearer " + apiKey
		client.GetConfig().UserAgent = fmt.Sprintf("VaultHub-CLI/%s (%s)", version.Version, version.Commit)
		debugLog("API client initialized successfully with User-Agent: %s", client.GetConfig().UserAgent)
	},
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key for authentication")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "Base URL of VaultHub server")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")

	// Add subcommands to root
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(versionCmd)

	listCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	getCmd.Flags().StringP("name", "n", "", "Vault name")
	getCmd.Flags().StringP("id", "i", "", "Vault Unique ID")
	getCmd.Flags().StringP("output", "o", "", "Output to file instead of stdout")
	getCmd.Flags().StringP("exec", "e", "", "Command to execute if vault has been updated")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all accessible vaults",
	Long: `List all vaults that you have access to.
This command will display basic information about each vault including
name, unique ID, and description.`,
	Run: func(cmd *cobra.Command, args []string) {
		debugLog("Executing list command")

		ctx := context.Background()
		debugLog("Making API request to get vaults by API key")
		vaults, _, err := client.CliAPI.GetVaultsByAPIKey(ctx).Execute()
		if err != nil {
			debugLog("API request failed: %v", err)
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		debugLog("API request successful, received %d vault(s)", len(vaults))

		jsonOutput, _ := cmd.Flags().GetBool("json")
		debugLog("JSON output mode: %v", jsonOutput)

		if jsonOutput {
			debugLog("Marshaling vaults to JSON")
			output, err := json.MarshalIndent(vaults, "", "  ")
			if err != nil {
				debugLog("JSON marshaling failed: %v", err)
				fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
				os.Exit(1)
			}
			debugLog("JSON marshaling successful")
			fmt.Println(string(output))
			return
		}

		if len(vaults) == 0 {
			debugLog("No vaults found, displaying empty message")
			fmt.Println("No vaults found.")
			return
		}

		debugLog("Displaying %d vault(s) in formatted output", len(vaults))
		fmt.Printf("Found %d vault(s):\n\n", len(vaults))
		for i, vault := range vaults {
			fmt.Printf("  %d. 📦 %s\n", i+1, vault.GetName())
			fmt.Printf("     ID: %s\n", vault.GetUniqueId())
			if vault.Category != nil && *vault.Category != "" {
				fmt.Printf("     Category: %s\n", *vault.Category)
			}
			if vault.Description != nil && *vault.Description != "" {
				fmt.Printf("     Description: %s\n", *vault.Description)
			}
			if i < len(vaults)-1 {
				fmt.Println()
			}
		}
		debugLog("List command completed successfully")
	},
}

var getCmd = &cobra.Command{
	Use:   "get <vault-name-or-id>",
	Short: "Get a specific vault by name or unique ID",
	Long: `Get detailed information about a specific vault including its encrypted value.
You can specify the vault by either its name or unique ID.

Examples:
  vault-hub get --name my-api-keys
  vault-hub get --id abc123-def456-ghi789
  vault-hub get --name my-api-keys --output ./secrets.txt
  vault-hub get --name my-api-keys --output ./secrets.txt --exec "source ./secrets.txt && npm start"`,
	Run: func(cmd *cobra.Command, args []string) {
		debugLog("Executing get command")

		name := mustGetStringFlag(cmd, "name")
		id := mustGetStringFlag(cmd, "id")
		outputFile := mustGetStringFlag(cmd, "output")
		followUpCommand := mustGetStringFlag(cmd, "exec")

		debugLog("Parameters - name: '%s', id: '%s', output: '%s', exec: '%s'", name, id, outputFile, followUpCommand)

		if name == "" && id == "" {
			debugLog("Validation failed: neither name nor id provided")
			fmt.Fprintf(os.Stderr, "Error: either name or id must be provided\n")
			os.Exit(1)
		}

		ctx := context.Background()
		var vault *openapi.Vault
		var err error
		if name != "" {
			debugLog("Making API request to get vault by name: %s", name)
			vault, _, err = client.CliAPI.GetVaultByNameAPIKey(ctx, name).Execute()
		} else {
			debugLog("Making API request to get vault by ID: %s", id)
			vault, _, err = client.CliAPI.GetVaultByAPIKey(ctx, id).Execute()
		}

		if err != nil {
			debugLog("API request failed: %v", err)
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		debugLog("API request successful, vault retrieved")

		if outputFile != "" {
			debugLog("Writing vault value to file: %s", outputFile)
			
			// Check if file exists and get its modification time
			var fileModTime time.Time
			var fileExists bool
			if fileInfo, err := os.Stat(outputFile); err == nil {
				fileModTime = fileInfo.ModTime()
				fileExists = true
				debugLog("File exists, last modified: %s", fileModTime.Format(time.RFC3339))
			} else {
				debugLog("File does not exist")
			}
			
			// Parse vault's updated timestamp
			var vaultUpdatedAt time.Time
			if vault.UpdatedAt != nil {
				vaultUpdatedAt = *vault.UpdatedAt
			} else {
				vaultUpdatedAt = time.Now() // Fallback to current time if no timestamp
			}
			debugLog("Vault last updated: %s", vaultUpdatedAt.Format(time.RFC3339))
			
			// Determine if vault has been updated
			vaultHasUpdates := !fileExists || vaultUpdatedAt.After(fileModTime)
			debugLog("Vault has updates: %v", vaultHasUpdates)
			
			err = os.WriteFile(outputFile, []byte(vault.Value), 0600)
			if err != nil {
				debugLog("File write failed: %v", err)
				fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
				os.Exit(1)
			}
			debugLog("File write successful")
			
			if vaultHasUpdates {
				fmt.Printf("Vault value written to %s (vault was updated)\n", outputFile)
				
				// Execute follow-up command if specified
				if followUpCommand != "" {
					debugLog("Executing follow-up command: %s", followUpCommand)
					fmt.Printf("Executing follow-up command: %s\n", followUpCommand)
					
					cmdParts := strings.Fields(followUpCommand)
					if len(cmdParts) > 0 {
						cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						
						if err := cmd.Run(); err != nil {
							debugLog("Follow-up command failed: %v", err)
							fmt.Fprintf(os.Stderr, "Warning: Follow-up command failed: %v\n", err)
						} else {
							debugLog("Follow-up command completed successfully")
						}
					}
				}
			} else {
				fmt.Printf("Vault value written to %s (no updates since last fetch)\n", outputFile)
			}
		} else {
			debugLog("Outputting vault value to stdout")
			fmt.Println(vault.Value)
		}
		debugLog("Get command completed successfully")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display version and commit hash information for VaultHub CLI.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Skip authentication for version command
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("VaultHub CLI\n")
		fmt.Printf("Version: %s\n", version.Version)
		fmt.Printf("Commit:  %s\n", version.Commit)
	},
}
