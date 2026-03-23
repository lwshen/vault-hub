package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	openapi "github.com/lwshen/vault-hub-go-client"
	"github.com/spf13/cobra"

	"github.com/lwshen/vault-hub/internal/cli/encryption"
	"github.com/lwshen/vault-hub/internal/constants"
)

// NewUpdateCommand creates the update command
func NewUpdateCommand(ctx *CommandContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update --name/--id <vault-name-or-id>",
		Short: "Update a vault value by name or unique ID",
		Long: `Update an existing vault's value by name or unique ID.

Client-side encryption is ENABLED BY DEFAULT for enhanced security when updating values.
The vault value is encrypted with a per-vault key derived from your API key.
Use --no-client-encryption to disable this feature if needed.

Examples:
  vault-hub update --name my-api-keys --value "new-secret-value"
  vault-hub update --id abc123-def456-ghi789 --value "new-value"
  vault-hub update --name my-api-keys --value-file ./secret.txt
  vault-hub update --id abc123 --value "plain-value" --no-client-encryption`,
		Run: func(cmd *cobra.Command, args []string) {
			runUpdateCommand(cmd, args, ctx)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Vault name")
	cmd.Flags().StringP("id", "i", "", "Vault Unique ID")
	cmd.Flags().StringP("value", "v", "", "New vault value")
	cmd.Flags().String("value-file", "", "Read value from file (takes precedence over --value)")
	cmd.Flags().StringP("output", "o", "text", "Output format: text|json")
	cmd.Flags().Bool("no-client-encryption", false, "Disable client-side encryption (less secure)")

	return cmd
}

// updateCommandParams holds the parsed command parameters
type updateCommandParams struct {
	name               string
	id                 string
	value              string
	valueFile          string
	output             string
	noClientEncryption bool
}

// runUpdateCommand executes the vault update operation
func runUpdateCommand(cmd *cobra.Command, _ []string, ctx *CommandContext) {
	ctx.DebugLog("Executing update command")

	// Parse command flags
	params := parseUpdateCommandFlags(cmd, ctx)
	ctx.DebugLog("Parameters - name: '%s', id: '%s', output: '%s'",
		params.name, params.id, params.output)

	// Validate required parameters
	if err := validateUpdateParams(params, cmd); err != nil {
		ctx.DebugLog("Validation failed: %v", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Read value from file if specified
	if params.valueFile != "" {
		value, err := os.ReadFile(params.valueFile)
		if err != nil {
			ctx.DebugLog("Failed to read value file: %v", err)
			fmt.Fprintf(os.Stderr, "Error: Failed to read value file: %v\n", err)
			os.Exit(1)
		}
		if len(value) == 0 {
			ctx.DebugLog("Value file is empty")
			fmt.Fprintf(os.Stderr, "Error: Value file is empty\n")
			os.Exit(1)
		}
		params.value = string(value)
		ctx.DebugLog("Value read from file: %s", params.valueFile)
	}

	// Build update request
	updateReq := buildUpdateRequest(params)

	// Apply client-side encryption if enabled and value is being updated
	if !params.noClientEncryption && params.value != "" {
		ctx.DebugLog("Encrypting vault value with client-side encryption")
		ctx.DebugLog("Original value length: %d bytes", len(params.value))

		// Determine the salt (vault identifier used for key derivation)
		salt := params.name
		if salt == "" {
			salt = params.id
		}
		ctx.DebugLog("Using salt for key derivation: %s", salt)

		// Encrypt the vault value
		encryptedValue, err := encryption.EncryptForClient(params.value, ctx.GetAPIKey(), salt)
		if err != nil {
			ctx.DebugLog("Encryption failed: %v", err)
			fmt.Fprintf(os.Stderr, "Error: Failed to encrypt vault value: %v\n", err)
			os.Exit(1)
		}

		updateReq.SetValue(encryptedValue)
		ctx.DebugLog("Vault value encrypted successfully")
		ctx.DebugLog("Encrypted value length: %d bytes", len(encryptedValue))

		// Set client encryption header
		ctx.GetClient().GetConfig().DefaultHeader[constants.HeaderClientEncryption] = "true"
		defer delete(ctx.GetClient().GetConfig().DefaultHeader, constants.HeaderClientEncryption)
	} else {
		ctx.DebugLog("Client-side encryption disabled or no value to encrypt")
	}

	// Update vault via API
	vault, err := updateVault(params, updateReq, ctx)
	if err != nil {
		ctx.DebugLog("API request failed: %v", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if vault == nil {
		ctx.DebugLog("API returned nil vault")
		fmt.Fprintf(os.Stderr, "Error: API returned nil vault\n")
		os.Exit(1)
	}
	ctx.DebugLog("API request successful, vault updated")

	// Handle output
	handleUpdateOutput(vault, params, ctx)
	ctx.DebugLog("Update command completed successfully")
}

// parseUpdateCommandFlags extracts and returns command flags
func parseUpdateCommandFlags(cmd *cobra.Command, ctx *CommandContext) updateCommandParams {
	noClientEncryption, _ := cmd.Flags().GetBool("no-client-encryption")
	output, _ := cmd.Flags().GetString("output")

	return updateCommandParams{
		name:               ctx.MustGetStringFlag(cmd, "name"),
		id:                 ctx.MustGetStringFlag(cmd, "id"),
		value:              ctx.MustGetStringFlag(cmd, "value"),
		valueFile:          ctx.MustGetStringFlag(cmd, "value-file"),
		output:             output,
		noClientEncryption: noClientEncryption,
	}
}

// validateUpdateParams ensures required parameters are provided
func validateUpdateParams(params updateCommandParams, cmd *cobra.Command) error {
	// Must specify either name or id, but not both
	if params.name == "" && params.id == "" {
		return fmt.Errorf("either --name or --id must be provided")
	}
	if params.name != "" && params.id != "" {
		return fmt.Errorf("cannot specify both --name and --id, please use only one")
	}

	// Check if value-file flag was explicitly set
	valueFileSet := cmd.Flags().Changed("value-file")

	// At least one update field must be provided
	// Note: check valueFileSet instead of params.valueFile since file is read after validation
	if params.value == "" && !valueFileSet {
		return fmt.Errorf("either --value or --value-file must be provided")
	}

	// Validate output format
	if params.output != "text" && params.output != "json" {
		return fmt.Errorf("invalid output format '%s', must be 'text' or 'json'", params.output)
	}

	return nil
}

// buildUpdateRequest creates the update request from parameters
func buildUpdateRequest(params updateCommandParams) openapi.UpdateVaultRequest {
	req := openapi.NewUpdateVaultRequest()

	if params.value != "" {
		req.SetValue(params.value)
	}

	return *req
}

// updateVault updates a vault via the API
func updateVault(params updateCommandParams, updateReq openapi.UpdateVaultRequest, ctx *CommandContext) (*openapi.Vault, error) {
	apiCtx := context.Background()

	var vault *openapi.Vault
	var err error

	if params.name != "" {
		ctx.DebugLog("Making API request to update vault by name: %s", params.name)
		vault, _, err = ctx.GetClient().CliAPI.UpdateVaultByNameAPIKey(apiCtx, params.name).UpdateVaultRequest(updateReq).Execute()
	} else {
		ctx.DebugLog("Making API request to update vault by ID: %s", params.id)
		vault, _, err = ctx.GetClient().CliAPI.UpdateVaultByAPIKey(apiCtx, params.id).UpdateVaultRequest(updateReq).Execute()
	}

	return vault, err
}

// handleUpdateOutput manages the output of vault data
func handleUpdateOutput(vault *openapi.Vault, params updateCommandParams, ctx *CommandContext) {
	if params.output == "json" {
		printUpdateJSONOutput(vault, params, ctx)
	} else {
		printUpdateTextOutput(vault)
	}
}

// printUpdateJSONOutput marshals and prints vault in JSON format
// If client-side encryption was used, decrypts the value before output
func printUpdateJSONOutput(vault *openapi.Vault, params updateCommandParams, ctx *CommandContext) {
	ctx.DebugLog("Marshaling vault to JSON")

	// Decrypt the value if client-side encryption was used
	if !params.noClientEncryption && vault.Value != "" {
		ctx.DebugLog("Decrypting vault value for JSON output")

		// Determine the salt (same as used for encryption)
		salt := params.name
		if salt == "" {
			salt = params.id
		}

		decryptedValue, err := encryption.DecryptForClient(vault.Value, ctx.GetAPIKey(), salt)
		if err != nil {
			ctx.DebugLog("Decryption failed: %v", err)
			fmt.Fprintf(os.Stderr, "Error: Failed to decrypt vault value: %v\n", err)
			os.Exit(1)
		}

		vault.Value = decryptedValue
		ctx.DebugLog("Vault value decrypted successfully for JSON output")
	}

	output, err := json.MarshalIndent(vault, "", "  ")
	if err != nil {
		ctx.DebugLog("JSON marshaling failed: %v", err)
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	ctx.DebugLog("JSON marshaling successful")
	fmt.Println(string(output))
}

// printUpdateTextOutput displays vault in a human-readable format
func printUpdateTextOutput(vault *openapi.Vault) {
	fmt.Printf("✅ Vault updated successfully\n\n")
	fmt.Printf("📦 %s\n", vault.GetName())
	fmt.Printf("   ID: %s\n", vault.GetUniqueId())

	if vault.UpdatedAt != nil {
		fmt.Printf("   Updated: %s\n", vault.UpdatedAt.Format("2006-01-02 %H:%M:%S"))
	}
}
