# CLI Guide

The VaultHub CLI is the primary way to interact with VaultHub. It provides secure, programmatic access to your vaults and is perfect for development workflows, CI/CD pipelines, and automation. Start here to quickly access your secrets without needing to set up a server.

## CLI Installation

### Download Pre-built Binaries

Download the latest CLI binary for your platform from the [GitHub releases page](https://github.com/lwshen/vault-hub/releases/latest).

**Latest Release**

Download the latest version of VaultHub CLI from our GitHub releases page:

[👉 Download Latest Release](https://github.com/lwshen/vault-hub/releases/latest)

**Supported Platforms:**
- **Linux** - amd64, arm64
- **Windows** - amd64  
- **macOS** - amd64, arm64

### Build from Source

```bash
# Clone the repository
git clone https://github.com/lwshen/vault-hub.git
cd vault-hub

# Build the CLI
go build -o vault-hub-cli ./apps/cli/main.go

# Make it executable and move to PATH (Linux/macOS)
chmod +x vault-hub-cli
sudo mv vault-hub-cli /usr/local/bin/vault-hub
```

## Authentication

The CLI uses API keys for authentication. First, create an API key in the web interface:

1. Log into the VaultHub web interface
2. Navigate to Dashboard → API Keys
3. Click "Create API Key" and give it a name
4. Copy the generated API key (starts with `vhub_`)

### Setting Up Authentication

```bash
# Set the API key as an environment variable
export VAULT_HUB_API_KEY=vhub_your_api_key_here

# Set the server URL (if different from default)
export VAULT_HUB_BASE_URL=https://your-vaulthub-server.com

# Enable debug mode for troubleshooting
export VAULT_HUB_DEBUG=true
# or
export DEBUG=true

# Or pass flags directly to commands
vault-hub --api-key vhub_your_api_key_here --base-url https://your-server.com list
```

## Commands

### List Vaults

```bash
# List all accessible vaults
vault-hub list

# Short form
vault-hub ls

# Output in JSON format for scripting
vault-hub list --json
vault-hub ls -j
```

### Get Vault Contents

```bash
# Get vault by name
vault-hub get --name production-secrets
vault-hub get -n production-secrets

# Get vault by ID
vault-hub get --id vault-uuid-here
vault-hub get -i vault-uuid-here

# Export to file (creates file with 0600 permissions for security)
vault-hub get --name production-secrets --output .env
vault-hub get -n production-secrets -o secrets.txt

# Execute command only if vault has been updated since last file write
vault-hub get --name production-secrets --output .env --exec "source .env && npm start"
vault-hub get -n production-secrets -o .env -e "docker build -t myapp ."

# The CLI intelligently detects updates by comparing:
# - Vault modification timestamp vs file modification time
# - Vault content vs existing file content
# Files are only updated and commands only executed when changes are detected
```

### Version Information

```bash
# Show version and build information
vault-hub version
```

## Example Workflows

### Development Workflow

```bash
# Get development secrets and start your app
vault-hub get --name dev-secrets --exec "npm run dev"

# Export secrets to .env file for local development
vault-hub get --name dev-secrets --output .env
```

### CI/CD Pipeline

```bash
# In your CI/CD pipeline
export VAULT_HUB_API_KEY=${{ secrets.VAULT_HUB_API_KEY }}
export VAULT_HUB_BASE_URL=https://vault.company.com
vault-hub get --name production-secrets --exec "docker build -t myapp ."

# GitHub Actions example
- name: Deploy with secrets
  env:
    VAULT_HUB_API_KEY: ${{ secrets.VAULT_HUB_API_KEY }}
    VAULT_HUB_BASE_URL: ${{ secrets.VAULT_HUB_BASE_URL }}
  run: |
    vault-hub get --name prod-env --output .env --exec "source .env && ./deploy.sh"
```

## Docker Usage

### One-shot Execution

```bash
# Run CLI in Docker container
docker run --rm \
  -e VAULT_HUB_API_KEY=vhub_your_api_key_here \
  -e VAULT_HUB_BASE_URL=https://your-server.com \
  vaulthub/cli:latest list

# Get vault and save to mounted volume
docker run --rm \
  -v $(pwd):/output \
  -e VAULT_HUB_API_KEY=vhub_your_api_key_here \
  -e VAULT_HUB_BASE_URL=https://your-server.com \
  -e VAULT_HUB_CLI_ARGS="get --name prod-secrets --output /output/.env" \
  vaulthub/cli:latest
```

### Scheduled Execution with Cron

```bash
# Run CLI on a schedule (default: every hour)
docker run -d \
  --name vault-hub-sync \
  -v $(pwd)/logs:/var/log/cron \
  -v $(pwd)/secrets:/output \
  -e RUN_MODE=cron \
  -e CRON_SCHEDULE="*/30 * * * *" \
  -e VAULT_HUB_API_KEY=vhub_your_api_key_here \
  -e VAULT_HUB_BASE_URL=https://your-server.com \
  -e VAULT_HUB_CLI_ARGS="get --name prod-secrets --output /output/.env" \
  vaulthub/cli:latest

# Check logs
docker logs vault-hub-sync
tail -f ./logs/vault-hub.log
```

## Environment Variables

The CLI supports the following environment variables:

| Variable | Flag Equivalent | Description |
|----------|-----------------|-------------|
| `VAULT_HUB_API_KEY` | `--api-key` | API key for authentication |
| `VAULT_HUB_BASE_URL` | `--base-url` | VaultHub server URL |
| `VAULT_HUB_DEBUG` or `DEBUG` | `--debug` | Enable debug logging |

## Troubleshooting

### Enable Debug Mode

```bash
# Enable debug logging to see detailed request/response info
export VAULT_HUB_DEBUG=true
vault-hub list

# Or use the flag
vault-hub --debug list
```

### Common Issues

**Authentication Failed**
```bash
# Verify your API key
vault-hub --debug list
# Check the debug output for authentication errors
```

**Connection Issues**
```bash
# Test connectivity to your server
curl -H "Authorization: Bearer vhub_your_api_key" https://your-server.com/api/status

# Verify base URL is correct
export VAULT_HUB_DEBUG=true
vault-hub --base-url https://your-server.com list
```

**File Permissions**
```bash
# CLI creates output files with 0600 permissions (owner read/write only)
# If you need different permissions, change them after file creation
vault-hub get --name secrets --output .env
chmod 644 .env  # If needed for your use case
```