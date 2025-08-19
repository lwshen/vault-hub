package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	openapi "github.com/lwshen/vault-hub-go-client"
	"github.com/spf13/cobra"
)

var (
	apiKey  string
	baseURL string
	debug   bool
	client  *openapi.APIClient
)

// debugLog prints debug messages to stderr when debug mode is enabled
func debugLog(format string, args ...interface{}) {
	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
	}
}

var rootCmd = &cobra.Command{
	Use:   "vault-hub",
	Short: "VaultHub CLI - Secure environment variable and API key management",
	Long: `VaultHub CLI is a command-line interface for managing your secure
environment variables and API keys stored in VaultHub.

This CLI allows you to list and retrieve vaults from your VaultHub instance.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debugLog("Initializing VaultHub CLI")
		debugLog("Base URL: %s", baseURL)
		debugLog("Debug mode: %v", debug)
		
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
		debugLog("API client initialized successfully")
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

	listCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	getCmd.Flags().StringP("name", "n", "", "Vault name")
	getCmd.Flags().StringP("id", "i", "", "Vault Unique ID")
	getCmd.Flags().StringP("output", "o", "", "Output to file instead of stdout")
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
  vault-hub get --name my-api-keys --output ./secrets.txt`,
	Run: func(cmd *cobra.Command, args []string) {
		debugLog("Executing get command")
		
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			debugLog("Failed to get name flag: %v", err)
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			debugLog("Failed to get id flag: %v", err)
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		outputFile, err := cmd.Flags().GetString("output")
		if err != nil {
			debugLog("Failed to get output flag: %v", err)
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		debugLog("Parameters - name: '%s', id: '%s', output: '%s'", name, id, outputFile)

		if name == "" && id == "" {
			debugLog("Validation failed: neither name nor id provided")
			fmt.Fprintf(os.Stderr, "Error: either name or id must be provided\n")
			os.Exit(1)
		}

		ctx := context.Background()
		var vault *openapi.Vault
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
			err = os.WriteFile(outputFile, []byte(vault.Value), 0600)
			if err != nil {
				debugLog("File write failed: %v", err)
				fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
				os.Exit(1)
			}
			debugLog("File write successful")
			fmt.Printf("Vault value written to %s\n", outputFile)
		} else {
			debugLog("Outputting vault value to stdout")
			fmt.Println(vault.Value)
		}
		debugLog("Get command completed successfully")
	},
}
