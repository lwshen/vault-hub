package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lwshen/vault-hub/internal/version"
)

// NewVersionCommand creates the version command
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
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
}
