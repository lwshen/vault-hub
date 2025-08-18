# VaultHub CLI

A command-line interface for managing secure environment variables and API keys stored in VaultHub using the `github.com/lwshen/vault-hub-go-client` package.

## Features

- **List Command**: Display all accessible vaults with their unique IDs, names, and descriptions
- **Get Command**: Retrieve specific vault details by name or unique ID
- **Flexible Configuration**: Support for both command-line flags and environment variables
- **Proper Error Handling**: User-friendly error messages and validation

## Installation

```bash
go build -o vault-hub-cli ./apps/cli/
```

## Configuration

The CLI requires two configuration parameters:

### Option 1: Command-line flags
```bash
./vault-hub-cli --api-key="your-api-key" --base-url="https://your-vault-hub.com" [command]
```

### Option 2: Environment variables
```bash
export VAULT_HUB_API_KEY="your-api-key"
export VAULT_HUB_BASE_URL="https://your-vault-hub.com"
./vault-hub-cli [command]
```

### Option 3: Mixed (flags override environment variables)
```bash
export VAULT_HUB_BASE_URL="https://your-vault-hub.com"
./vault-hub-cli --api-key="your-api-key" [command]
```

## Usage

### List all vaults
```bash
./vault-hub-cli list
# or
./vault-hub-cli ls
```

### Get a specific vault
```bash
# By unique ID
./vault-hub-cli get abc123-def456-ghi789

# By name (case-insensitive)
./vault-hub-cli get my-api-keys
```

### Help
```bash
./vault-hub-cli --help
./vault-hub-cli list --help
./vault-hub-cli get --help
```

## Example Output

### List Command
```
✓ Initialized VaultHub client
  Base URL: https://vault-hub.example.com
  API Key:  test...-123
Fetching vaults...
UNIQUE ID                            NAME                 DESCRIPTION
------------------------------------ -------------------- --------------------
abc123-def456-ghi789                 my-api-keys          Production API keys
xyz789-abc123-def456                 database-creds       Database credentials

Total: 2 vault(s)
```

### Get Command
```
✓ Initialized VaultHub client
  Base URL: https://vault-hub.example.com
  API Key:  test...-123
Vault Details:
  Unique ID:   abc123-def456-ghi789
  Name:        my-api-keys
  Description: Production API keys
  Category:    secrets
  Value:       {"api_key": "secret-value", "token": "another-secret"}
  Created:     2024-01-15 10:30:45
  Updated:     2024-01-20 14:22:10
```

## Important Note

**Current Status**: The `github.com/lwshen/vault-hub-go-client` package has compilation issues due to conflicting type definitions (APIKey type is defined in both `model_api_key.go` and `configuration.go`). This implementation uses mock structures to demonstrate the correct usage pattern.

**To Fix**: Once the package maintainer resolves the compilation issues, replace the mock structures in `main.go` with:

```go
import openapi "github.com/lwshen/vault-hub-go-client"

// In initializeClient():
config := openapi.NewConfiguration()
client := openapi.NewAPIClient(config)
```

## Implementation Details

The CLI application demonstrates proper usage of the VaultHub Go client with:

1. **Configuration Management**: Handles API key and base URL via flags or environment variables
2. **Client Initialization**: Sets up the OpenAPI client with proper authentication headers
3. **Command Structure**: Uses Cobra for clean CLI command organization
4. **Error Handling**: Provides user-friendly error messages and proper HTTP status code handling
5. **Smart Vault Lookup**: Attempts direct lookup by unique ID, falls back to name-based search
6. **Formatted Output**: Clean, readable output formatting for both list and detailed views

## API Endpoints Used

- `GET /api/vaults` - List all accessible vaults (returns VaultLite objects)
- `GET /api/vault/{uniqueId}` - Get specific vault details (returns full Vault object)

## Authentication

The CLI uses Bearer token authentication, setting the `Authorization: Bearer {api-key}` header for all API requests.