package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	openapi "github.com/lwshen/vault-hub-go-client"
	"github.com/spf13/cobra"
)

// NewGetCommand creates the get command
func NewGetCommand(ctx *CommandContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get --name/--id <vault-name-or-id> --output <file> --exec <command>",
		Short: "Get a specific vault by name or unique ID",
		Long: `Get detailed information about a specific vault including its encrypted value.
You can specify the vault by either its name or unique ID.

Examples:
  vault-hub get --name my-api-keys
  vault-hub get --id abc123-def456-ghi789
  vault-hub get --name my-api-keys --output ./secrets.txt
  vault-hub get --name my-api-keys --output .env --exec "source .env && echo 'Environment loaded'"`,
		Run: func(cmd *cobra.Command, args []string) {
			runGetCommand(cmd, args, ctx)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Vault name")
	cmd.Flags().StringP("id", "i", "", "Vault Unique ID")
	cmd.Flags().StringP("output", "o", "", "Output to file instead of stdout")
	cmd.Flags().StringP("exec", "e", "", "Command to execute if vault has been updated")

	return cmd
}

// runGetCommand executes the vault retrieval operation
// It validates parameters, fetches the vault, and handles output
func runGetCommand(cmd *cobra.Command, _ []string, ctx *CommandContext) {
	ctx.DebugLog("Executing get command")

	// Parse command flags
	params := parseGetCommandFlags(cmd, ctx)
	ctx.DebugLog("Parameters - name: '%s', id: '%s', output: '%s', exec: '%s'",
		params.name, params.id, params.outputFile, params.followUpCommand)

	// Validate required parameters
	if err := validateGetParams(params); err != nil {
		ctx.DebugLog("Validation failed: %v", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Fetch vault from API
	vault, err := fetchVault(params, ctx)
	if err != nil {
		ctx.DebugLog("API request failed: %v", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	ctx.DebugLog("API request successful, vault retrieved")

	// Handle output
	handleVaultOutput(vault, params, ctx)
	ctx.DebugLog("Get command completed successfully")
}

// getCommandParams holds the parsed command parameters
type getCommandParams struct {
	name            string
	id              string
	outputFile      string
	followUpCommand string
}

// parseGetCommandFlags extracts and returns command flags
func parseGetCommandFlags(cmd *cobra.Command, ctx *CommandContext) getCommandParams {
	return getCommandParams{
		name:            ctx.MustGetStringFlag(cmd, "name"),
		id:              ctx.MustGetStringFlag(cmd, "id"),
		outputFile:      ctx.MustGetStringFlag(cmd, "output"),
		followUpCommand: ctx.MustGetStringFlag(cmd, "exec"),
	}
}

// validateGetParams ensures required parameters are provided
func validateGetParams(params getCommandParams) error {
	if params.name == "" && params.id == "" {
		return fmt.Errorf("either name or id must be provided")
	}
	return nil
}

// fetchVault retrieves a vault from the API by name or ID
func fetchVault(params getCommandParams, ctx *CommandContext) (*openapi.Vault, error) {
	apiCtx := context.Background()

	if params.name != "" {
		ctx.DebugLog("Making API request to get vault by name: %s", params.name)
		vault, _, err := ctx.GetClient().CliAPI.GetVaultByNameAPIKey(apiCtx, params.name).Execute()
		return vault, err
	}

	ctx.DebugLog("Making API request to get vault by ID: %s", params.id)
	vault, _, err := ctx.GetClient().CliAPI.GetVaultByAPIKey(apiCtx, params.id).Execute()
	return vault, err
}

// handleVaultOutput manages the output of vault data
func handleVaultOutput(vault *openapi.Vault, params getCommandParams, ctx *CommandContext) {
	if params.outputFile != "" {
		handleFileOutput(vault, params.outputFile, params.followUpCommand, ctx.DebugLog)
	} else {
		ctx.DebugLog("Outputting vault value to stdout")
		fmt.Println(vault.Value)
	}
}

type UpdateResult struct {
	HasUpdates       bool
	TimestampChanged bool
	ContentChanged   bool
	FileExists       bool
	Reason           string
}

// isVaultUpdated determines if a vault has been updated since the last file write
// It compares timestamps and content to decide if the file should be updated
func isVaultUpdated(vault *openapi.Vault, filePath string, debugLog func(string, ...any)) (UpdateResult, error) {
	var result UpdateResult

	// Check file existence and get modification info
	fileModTime, fileExists := getFileModTime(filePath, debugLog)
	result.FileExists = fileExists

	// Compare file content with vault value
	result.ContentChanged = checkContentChanged(vault.Value, filePath, fileExists, debugLog)

	// Check timestamp changes
	vaultUpdatedAt, hasValidTimestamp := getVaultTimestamp(vault, debugLog)
	result.TimestampChanged = hasValidTimestamp && fileExists && vaultUpdatedAt.After(fileModTime)

	// Determine overall update status
	result.HasUpdates = !fileExists || result.TimestampChanged || result.ContentChanged
	result.Reason = determineUpdateReason(result)

	debugLog("Vault has updates: %v (timestamp: %v, content: %v)",
		result.HasUpdates, result.TimestampChanged, result.ContentChanged)
	return result, nil
}

// getFileModTime returns file modification time and existence status
func getFileModTime(filePath string, debugLog func(string, ...any)) (time.Time, bool) {
	if fileInfo, err := os.Stat(filePath); err == nil {
		fileModTime := fileInfo.ModTime()
		debugLog("File exists, last modified: %s", fileModTime.Format(time.RFC3339))
		return fileModTime, true
	}
	debugLog("File does not exist")
	return time.Time{}, false
}

// checkContentChanged compares file content with vault value
func checkContentChanged(vaultValue, filePath string, fileExists bool, debugLog func(string, ...any)) bool {
	if !fileExists {
		return true // New file, so content is "changed"
	}

	existingContent, readErr := os.ReadFile(filePath)
	if readErr != nil {
		debugLog("Could not read existing file for content comparison: %v", readErr)
		return true // Assume content changed if we can't read the file
	}

	contentChanged := string(existingContent) != vaultValue
	debugLog("Content comparison - file differs from vault: %v", contentChanged)
	return contentChanged
}

// getVaultTimestamp extracts and validates vault timestamp
func getVaultTimestamp(vault *openapi.Vault, debugLog func(string, ...any)) (time.Time, bool) {
	if vault.UpdatedAt != nil {
		vaultUpdatedAt := *vault.UpdatedAt
		debugLog("Vault last updated: %s", vaultUpdatedAt.Format(time.RFC3339))
		return vaultUpdatedAt, true
	}
	debugLog("Vault has no timestamp - treating as new vault")
	return time.Time{}, false
}

// determineUpdateReason provides a human-readable reason for the update decision
func determineUpdateReason(result UpdateResult) string {
	if !result.FileExists {
		return "new file"
	}
	if result.TimestampChanged && result.ContentChanged {
		return "vault updated and content differs"
	}
	if result.TimestampChanged {
		return "vault was updated"
	}
	if result.ContentChanged {
		return "content differs from vault"
	}
	return "no updates - content matches vault"
}

func handleFileOutput(vault *openapi.Vault, outputFile, followUpCommand string, debugLog func(string, ...any)) {
	debugLog("Writing vault value to file: %s", outputFile)

	updateResult, err := isVaultUpdated(vault, outputFile, debugLog)
	if err != nil {
		debugLog("Error checking vault updates: %v", err)
		fmt.Fprintf(os.Stderr, "Error checking vault updates: %v\n", err)
		os.Exit(1)
	}

	if updateResult.HasUpdates {
		err = os.WriteFile(outputFile, []byte(vault.Value), 0600)
		if err != nil {
			debugLog("File write failed: %v", err)
			fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
			os.Exit(1)
		}
		debugLog("File write successful")
		fmt.Printf("Vault value written to %s (%s)\n", outputFile, updateResult.Reason)

		// Execute follow-up command if specified
		if followUpCommand != "" {
			executeFollowUpCommand(followUpCommand, debugLog)
		}
	} else {
		debugLog("No updates detected, file not modified")
	}
}

func executeFollowUpCommand(followUpCommand string, debugLog func(string, ...any)) {
	debugLog("Executing follow-up command: %s", followUpCommand)
	fmt.Printf("Executing follow-up command: %s\n", followUpCommand)

	// Use shell to handle complex commands properly
	cmd := exec.Command("sh", "-c", followUpCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		debugLog("Follow-up command failed: %v", err)
		fmt.Fprintf(os.Stderr, "Warning: Follow-up command failed: %v\n", err)
	} else {
		debugLog("Follow-up command completed successfully")
	}
}
