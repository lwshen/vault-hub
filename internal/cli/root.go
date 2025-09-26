package cli

import (
	"fmt"
	"os"

	openapi "github.com/lwshen/vault-hub-go-client"
	"github.com/spf13/cobra"

	"github.com/lwshen/vault-hub/internal/cli/commands"
)

// NewRootCommand creates the root command with all subcommands
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "vault-hub",
		Short: "VaultHub CLI - Secure environment variable and API key management",
		Long: `VaultHub CLI is a command-line interface for managing your secure
environment variables and API keys stored in VaultHub.

This CLI allows you to list and retrieve vaults from your VaultHub instance.

Global flags can be set via environment variables:
  --api-key     VAULT_HUB_API_KEY
  --base-url    VAULT_HUB_BASE_URL  
  --debug       DEBUG`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Update global variables with flag values or environment variables
			APIKey = getFlagOrEnv(cmd, "api-key", "VAULT_HUB_API_KEY", "")
			BaseURL = getFlagOrEnv(cmd, "base-url", "VAULT_HUB_BASE_URL", "")
			Debug = getBoolFlagOrEnv(cmd, "debug", "DEBUG", false)

			InitializeClient()
		},
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&APIKey, "api-key", "", "API key for authentication (env: VAULT_HUB_API_KEY)")
	rootCmd.PersistentFlags().StringVar(&BaseURL, "base-url", "", "Base URL of VaultHub server (env: VAULT_HUB_BASE_URL)")
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Enable debug mode (env: DEBUG)")

	// Create a context to pass dependencies to commands
	ctx := &commands.CommandContext{
		GetClient:         func() *openapi.APIClient { return Client },
		DebugLog:          DebugLog,
		MustGetStringFlag: MustGetStringFlag,
	}

	// Add subcommands
	rootCmd.AddCommand(commands.NewListCommand(ctx))
	rootCmd.AddCommand(commands.NewGetCommand(ctx))
	rootCmd.AddCommand(commands.NewVersionCommand())

	return rootCmd
}

// Execute runs the root command
func Execute() {
	rootCmd := NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
