# VaultHub CLI

A command-line interface for managing secure environment variables and API keys stored in VaultHub.

## Features

- **List Vaults**: View all accessible vaults with basic information
- **Get Vault**: Retrieve specific vaults by name or unique ID
- **Health Check**: Verify server connectivity and health
- **Flexible Configuration**: Support for command-line flags, environment variables, and config files

## Installation

The CLI is part of the VaultHub project. To build it:

```bash
go build -o vault-hub apps/cli/main.go
```

## Configuration

The CLI can be configured using multiple methods (in order of precedence):

### 1. Command Line Flags

```bash
vault-hub --api-key "your-key" --base-url "https://server.com" list
```

### 2. Environment Variables

```bash
export VAULT_HUB_API_KEY="your-api-key"
export VAULT_HUB_BASE_URL="https://your-server.com"
vault-hub list
```

### 3. Configuration File

Create `~/.vault-hub/config.yaml`:

```yaml
api_key: "your-api-key-here"
base_url: "https://your-vault-hub-server.com"
timeout: "30s"
```

## Usage

### List All Vaults

```bash
vault-hub list
# or
vault-hub ls
```

### Get a Specific Vault

```bash
# By name
vault-hub get my-api-keys

# By unique ID
vault-hub get abc123-def456-ghi789
```

### Health Check

```bash
vault-hub health
```

### Global Options

- `--api-key`: API key for authentication
- `--base-url`: Base URL of VaultHub server
- `--timeout`: Request timeout (default: 30s)

## Examples

```bash
# List vaults with custom configuration
vault-hub --api-key "my-key" --base-url "https://vault.example.com" list

# Get a vault with custom timeout
vault-hub --timeout "60s" get production-db-credentials

# Use environment variables
export VAULT_HUB_API_KEY="my-key"
export VAULT_HUB_BASE_URL="https://vault.example.com"
vault-hub list
```

## Error Handling

The CLI provides clear error messages for common issues:

- Missing API key or base URL
- Invalid server URLs
- Network timeouts
- Authentication failures
- Vault not found errors

## Security

- API keys are never logged or displayed
- All communication uses HTTPS by default
- Sensitive vault values are displayed as-is (encrypted)
- Timeout protection prevents hanging requests