# VaultHub CLI

A command-line interface for VaultHub - secure environment variable and API key management system.

## Overview

The VaultHub CLI provides a convenient way to interact with your VaultHub server from the command line. It allows you to list and retrieve encrypted vaults using API key authentication, making it perfect for scripts, CI/CD pipelines, and development workflows.

## Features

- 🔐 **Secure API key authentication**
- 📋 **List all accessible vaults** with detailed information
- 🔍 **Retrieve specific vaults** by name or unique ID
- 💾 **Save vault contents to files** for easy integration
- 🐛 **Debug mode** for troubleshooting
- 📄 **JSON output support** for programmatic usage

## Installation

### Build from Source

```bash
# From the vault-hub project root
go build -o vault-hub-cli ./apps/cli/main.go

# Make it executable and optionally move to PATH
chmod +x vault-hub-cli
sudo mv vault-hub-cli /usr/local/bin/  # Optional: for global access
```

### Prerequisites

- Go 1.24+ (for building from source)
- Access to a running VaultHub server
- Valid API key from your VaultHub instance

## Configuration

The CLI requires two essential parameters that can be provided via command-line flags:

- `--api-key`: Your VaultHub API key for authentication
- `--base-url`: Base URL of your VaultHub server (e.g., `https://vault.example.com`)

### Environment Variables (Optional)

While not currently supported, you can create wrapper scripts to avoid repetitive flag usage:

```bash
#!/bin/bash
# ~/.local/bin/vh (example wrapper script)
vault-hub-cli --api-key="$VAULT_HUB_API_KEY" --base-url="$VAULT_HUB_URL" "$@"
```

## Usage

### Basic Command Structure

```bash
vault-hub-cli [global-flags] <command> [command-flags] [arguments]
```

### Global Flags

- `--api-key <key>`: API key for authentication (required)
- `--base-url <url>`: Base URL of VaultHub server (required)
- `--debug`: Enable debug mode for detailed logging

### Commands

#### `list` (alias: `ls`)

List all vaults you have access to.

```bash
# Basic list
vault-hub-cli --api-key="your-key" --base-url="https://vault.example.com" list

# JSON output for scripts
vault-hub-cli --api-key="your-key" --base-url="https://vault.example.com" list --json
```

**Flags:**
- `-j, --json`: Output in JSON format

**Example output:**
```
Found 3 vault(s):

  1. 📦 production-api-keys
     ID: abc123-def456-ghi789
     Category: API Keys
     Description: Production environment API keys

  2. 📦 database-credentials
     ID: xyz789-uvw456-rst123
     Category: Database
     Description: Database connection strings

  3. 📦 third-party-tokens
     ID: mno345-pqr678-stu901
     Description: External service authentication tokens
```

#### `get`

Retrieve a specific vault by name or unique ID.

```bash
# Get by name
vault-hub-cli --api-key="your-key" --base-url="https://vault.example.com" get --name "production-api-keys"

# Get by unique ID
vault-hub-cli --api-key="your-key" --base-url="https://vault.example.com" get --id "abc123-def456-ghi789"

# Save to file
vault-hub-cli --api-key="your-key" --base-url="https://vault.example.com" get --name "production-api-keys" --output "./secrets.env"
```

**Flags:**
- `-n, --name <name>`: Vault name
- `-i, --id <id>`: Vault unique ID
- `-o, --output <file>`: Save output to file instead of stdout

**Note:** Either `--name` or `--id` must be provided, but not both.

## Examples

### Development Workflow

```bash
# Set up your environment
export VAULT_API_KEY="your-api-key-here"
export VAULT_URL="https://vault.company.com"

# List available vaults
vault-hub-cli --api-key="$VAULT_API_KEY" --base-url="$VAULT_URL" list

# Get development environment variables
vault-hub-cli --api-key="$VAULT_API_KEY" --base-url="$VAULT_URL" get --name "dev-env" --output ".env.local"

# Source the environment file
source .env.local
```

### CI/CD Pipeline

```bash
#!/bin/bash
# ci-script.sh - Example CI/CD integration

set -e

# Retrieve production secrets
vault-hub-cli \
  --api-key="$CI_VAULT_API_KEY" \
  --base-url="$CI_VAULT_URL" \
  get --name "production-deploy-keys" \
  --output "./deploy-keys.env"

# Source secrets and deploy
source ./deploy-keys.env
./deploy.sh

# Clean up sensitive files
rm -f ./deploy-keys.env
```

### Debugging Connection Issues

```bash
# Enable debug mode to see detailed logs
vault-hub-cli \
  --api-key="your-key" \
  --base-url="https://vault.example.com" \
  --debug \
  list
```

Debug output will show:
- API client initialization
- HTTP request details  
- Response processing
- Error details

## Error Handling

The CLI provides clear error messages for common issues:

- **Authentication errors**: Invalid API key or expired token
- **Network errors**: Connection timeouts or unreachable server
- **Vault not found**: Specified vault name or ID doesn't exist
- **Permission errors**: API key lacks access to specific vault
- **File write errors**: Permission issues when using `--output`

## Integration Tips

### Shell Scripts

```bash
#!/bin/bash
# get-secret.sh - Helper script for retrieving secrets

VAULT_NAME="$1"
if [ -z "$VAULT_NAME" ]; then
    echo "Usage: $0 <vault-name>"
    exit 1
fi

vault-hub-cli \
  --api-key="$VAULT_HUB_API_KEY" \
  --base-url="$VAULT_HUB_URL" \
  get --name "$VAULT_NAME"
```

### JSON Processing

```bash
# Get vault list as JSON and process with jq
vault-hub-cli --api-key="$API_KEY" --base-url="$URL" list --json | \
  jq -r '.[] | select(.category == "API Keys") | .name'
```

### Docker Integration

```dockerfile
# Dockerfile example
FROM golang:1.24-alpine AS builder
COPY . /app
WORKDIR /app
RUN go build -o vault-hub-cli ./apps/cli/main.go

FROM alpine:3.22
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/vault-hub-cli /usr/local/bin/
ENTRYPOINT ["vault-hub-cli"]
```

## Security Considerations

- **API Key Protection**: Never hardcode API keys in scripts. Use environment variables or secure CI/CD variables.
- **File Permissions**: Output files are created with `0600` permissions (owner read/write only).
- **Debug Mode**: Avoid debug mode in production as it may log sensitive information.
- **Network Security**: Always use HTTPS URLs for production VaultHub servers.

## Troubleshooting

### Common Issues

1. **"Error: either name or id must be provided"**
   - Solution: Use either `--name` or `--id` flag with the `get` command

2. **"Error: connection refused"**
   - Check if the VaultHub server is running
   - Verify the `--base-url` is correct
   - Ensure network connectivity

3. **"Error: 401 Unauthorized"**
   - Verify your API key is correct and not expired
   - Check if the API key has the required permissions

4. **"Error: vault not found"**
   - Confirm the vault name or ID is correct
   - Ensure your API key has access to the vault

### Getting Help

```bash
# General help
vault-hub-cli --help

# Command-specific help
vault-hub-cli list --help
vault-hub-cli get --help
```

## Development

This CLI is part of the larger VaultHub project. For development and contribution:

1. See the main [VaultHub README](../../README.md) for project setup
2. The CLI source code is in `./main.go`
3. Uses [Cobra](https://github.com/spf13/cobra) for command structure
4. Integrates with the VaultHub Go client library

## License

Apache License 2.0 - see [LICENSE](../../LICENSE) file.
