package commands

import (
	openapi "github.com/lwshen/vault-hub-go-client"
	"github.com/spf13/cobra"
)

// CommandContext holds dependencies for commands
type CommandContext struct {
	GetClient         func() *openapi.APIClient
	GetAPIKey         func() string
	DebugLog          func(string, ...any)
	MustGetStringFlag func(*cobra.Command, string) string
}
